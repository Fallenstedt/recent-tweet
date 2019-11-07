package lib

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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
func (s *SimpleTweetTable) QueryTweetFromDynamo(t *SimpleTweetDTO) (*dynamodb.GetItemOutput, error) {
	result, err := s.session.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(s.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(t.ID),
			},
		},
	})

	return result, err
}
