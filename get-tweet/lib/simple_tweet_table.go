package lib

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// CreateSimpleTweetTableSession returns a dynamodb session
func CreateSimpleTweetTableSession() (*dynamodb.DynamoDB, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	if err != nil {
		return nil, err
	}

	svc := dynamodb.New(sess)

	return svc, nil
}
