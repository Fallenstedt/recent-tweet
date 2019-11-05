package lib

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// CreateSimpleTweetTableSession returns a dynamodb session
func CreateSimpleTweetTableSession() *dynamodb.DynamoDB {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)

	return svc
}

func GetLatestTweetFromDynamo(s *dynamodb.DynamoDB) string {
	return "Hello"
}
