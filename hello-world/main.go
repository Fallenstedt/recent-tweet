package main

import (
	"encoding/json"
	"hello-sam/hello-world/lib"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	t := lib.Tweet{Message: "hello", ID: "5"}
	jsonData, err := json.Marshal(t.Message)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(jsonData),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
