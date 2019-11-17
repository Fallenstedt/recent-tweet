package lib

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// DynamoDbInstance returns a DynamoDB session
// to the table it talks to.
type DynamoDbInstance struct {
	Session   *dynamodb.DynamoDB
	TableName string
}

// CreateSimpleTweetTableSession returns a dynamodb session
func CreateSimpleTweetTableSession(tableName string) *DynamoDbInstance {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	obj := DynamoDbInstance{
		Session:   svc,
		TableName: tableName,
	}
	return &obj
}

// DynamoOperator defines a single operation that
// can be operated on our database
type DynamoOperator interface {
	Execute(*DynamoDbInstance, *SimpleTweetDTO) *SimpleTweetDTO
}

// DynamoOperation defines a single Operator that can be executed
type DynamoOperation struct {
	DynamoOperator DynamoOperator
}

// ExecuteOperation allows us to call the chosen DynamoOperator's `Execute` function
func (o *DynamoOperation) ExecuteOperation(s *DynamoDbInstance, t *SimpleTweetDTO) *SimpleTweetDTO {
	return o.DynamoOperator.Execute(s, t)
}

// QueryTweet is a DynamoOperator that allows us to Query a Tweet
type QueryTweet struct{}

// Execute from QueryTweet allows us to fetch the latest tweet from DynamoDB
func (QueryTweet) Execute(s *DynamoDbInstance, t *SimpleTweetDTO) *SimpleTweetDTO {
	log.Printf("Attempting to query tweet %v", t.ID)
	result, err := s.Session.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(s.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(t.ID),
			},
		},
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to GetItem from dynamo, %v", err))
	}

	tweet := SimpleTweetDTO{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &tweet)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	log.Printf("Successfully queried tweet %v", tweet.ID)
	return &tweet
}

// UpdateTweet is a DynamoOperator that allows us to Update the latest tweet
type UpdateTweet struct{}

// Execute from UpdateTweet allows us to Update the latest tweet in DynamoDB
func (UpdateTweet) Execute(s *DynamoDbInstance, t *SimpleTweetDTO) *SimpleTweetDTO {
	av, err := dynamodbattribute.MarshalMap(t)

	log.Printf("%+v\n", av)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal Record, %v", err))
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(s.TableName),
		Item:      av,
	}

	_, err = s.Session.PutItem(input)

	if err != nil {
		panic(fmt.Sprintf("Failed to Put Item, %v", err))
	}

	log.Println("Successfully added tweet " + t.ID)
	return t
}

// DeleteTweet is a DynamoOperator that allows us to Delete the latest tweet from Dynamo
type DeleteTweet struct{}

// Execute from DeleteTweet allows us to delete the latest tweet in DynamoDB
func (DeleteTweet) Execute(s *DynamoDbInstance, t *SimpleTweetDTO) *SimpleTweetDTO {
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(s.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(t.ID),
			},
		},
	}

	_, err := s.Session.DeleteItem(input)
	if err != nil {
		panic(fmt.Sprintf("Failed to DeleteItem from dynamo, %v", err))
	}

	log.Println("Successfully deleted tweet " + t.ID)
	return t
}

// GetLatestTweet is a DynamoOperator that allows us to get the latest tweet from dynamo
type GetLatestTweet struct{}

// Execute from GetLatest allows us to scan our table for a tweet.
func (GetLatestTweet) Execute(s *DynamoDbInstance, t *SimpleTweetDTO) *SimpleTweetDTO {
	input := &dynamodb.ScanInput{
		TableName: aws.String(s.TableName),
		Limit:     aws.Int64(1),
	}

	result, err := s.Session.Scan(input)
	if err != nil {
		panic(fmt.Sprintf("Failed to query the latest tweet from dynamo, %v", err))
	}

	tweet := SimpleTweetDTO{}

	if len(result.Items) == 0 {
		return &tweet
	}

	err = dynamodbattribute.UnmarshalMap(result.Items[0], &tweet)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	log.Println("Successfully got the latest tweet " + tweet.ID)

	return &tweet
}
