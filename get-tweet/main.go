package main

import (
	"fmt"
	"log"
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

	log.Print("Getting Client")
	client, err := twitterFactory.GetClient()
	if err != nil {
		log.Print(err)
		return events.APIGatewayProxyResponse{}, err
	}
	log.Print("Got Client")

	log.Print("Getting Tweets")
	tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		TweetMode: "extended",
		Count:     1,
	})

	if err != nil {
		log.Print(err)
		return events.APIGatewayProxyResponse{}, err
	}
	log.Print("Got Tweets")
	log.Print("Create Simple Tweet")
	simpleTweets, err := lib.CreateSimpleTweetDTO(&tweets)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       string("User has no tweets"),
		}, nil
	}
	log.Print("Got Simple Tweet")
	log.Print("Create Session and Query")
	session := lib.CreateSimpleTweetTableSession(os.Getenv("TABLE_NAME"))
	queryResult, err := session.QueryTweetFromDynamo(&simpleTweets[0])

	fmt.Printf("%+v\n", queryResult)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	log.Print("Got Session and Query")

	return events.APIGatewayProxyResponse{
		Body:       string(queryResult.String()),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
