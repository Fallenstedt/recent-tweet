package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/fallenstedt/latest-tweet/lib"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	session := lib.CreateSimpleTweetTableSession(os.Getenv("TABLE_NAME"))

	operation := lib.DynamoOperation{
		DynamoOperator: lib.GetLatestTweet{},
	}
	tweet := operation.ExecuteOperation(session, nil)

	return events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		Body:       tweet.TweetToJSON(),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
