package lib

import (
	"fmt"
	"os"

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

	return &tweet
}

// UpdateTweet is a DynamoOperator that allows us to Update the latest tweet
type UpdateTweet struct{}

// Execute from UpdateTweet allows us to Update the latest tweet in DynamoDB
func (UpdateTweet) Execute(s *DynamoDbInstance, t *SimpleTweetDTO) *SimpleTweetDTO {
	av, err := dynamodbattribute.MarshalMap(t)

	fmt.Printf("%+v\n", av)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal Record, %v", err))
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	_, err = s.Session.PutItem(input)
	if err != nil {
		panic(fmt.Sprintf("Failed to Put Item, %v", err))
	}

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	fmt.Println("Successfully added tweet " + t.ID)
	return t

}
