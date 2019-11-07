package main

import (
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

	client, err := twitterFactory.GetClient()
	if err != nil {
		log.Print(err)
		return events.APIGatewayProxyResponse{}, err
	}

	tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		TweetMode: "extended",
		Count:     1,
	})

	if err != nil {
		log.Print(err)
		return events.APIGatewayProxyResponse{}, err
	}

	simpleTweets, err := lib.CreateSimpleTweetDTO(&tweets)

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       string("User has no tweets"),
		}, nil
	}

	session := lib.CreateSimpleTweetTableSession(os.Getenv("TABLE_NAME"))
	queryResult := session.QueryTweetFromDynamo(&simpleTweets[0])
	isTweetFromDynamoNotMyLatestTweet := queryResult.ID == "" || &queryResult.ID != &simpleTweets[0].ID

	if isTweetFromDynamoNotMyLatestTweet {
		log.Print("Updating latest tweet in Dynamo")
		session.UpdateLatestTweetInDynamo(&simpleTweets[0])
		return events.APIGatewayProxyResponse{
			Body:       string("bar"),
			StatusCode: 200,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       string("foo"),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
