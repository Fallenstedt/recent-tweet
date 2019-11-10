package lib

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// SimpleTweetTable returns a DynamoDB session
// to the table it talks to.
type SimpleTweetTable struct {
	session   *dynamodb.DynamoDB
	TableName string
}

// CreateSimpleTweetTableSession returns a dynamodb session
func CreateSimpleTweetTableSession(tableName string) *SimpleTweetTable {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	obj := SimpleTweetTable{
		session:   svc,
		TableName: tableName,
	}
	return &obj
}

//QueryTweetFromDynamo gets the tweet from Dynamo if it exists
func (s *SimpleTweetTable) QueryTweetFromDynamo(t *SimpleTweetDTO) *SimpleTweetDTO {
	result, err := s.session.GetItem(&dynamodb.GetItemInput{
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

func (s *SimpleTweetTable) UpdateLatestTweetInDynamo(t *SimpleTweetDTO) {
	av, err := dynamodbattribute.MarshalMap(t)

	fmt.Printf("%+v\n", av)
	if err != nil {
		panic(fmt.Sprintf("Failed to marshal Record, %v", err))
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	_, err = s.session.PutItem(input)
	if err != nil {
		panic(fmt.Sprintf("Failed to Put Item, %v", err))
	}

	fmt.Println("Successfully added tweet " + t.ID)
}
