package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/fallenstedt/latest-tweet/get-tweet/lib"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	twitterFactory := lib.Credentials{
		ConsumerKey:       os.Getenv("CONSUMER_KEY"),
		ConsumerSecret:    os.Getenv("CONSUMER_SECRET"),
		AccessToken:       os.Getenv("ACCESS_KEY"),
		AccessTokenSecret: os.Getenv("ACCESS_SECRET"),
	}

	client, err := twitterFactory.GetClient()

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		Count: 1,
	})

	tweetsJSON, err := json.Marshal(tweets)

	if err == nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       string(tweetsJSON),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
