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

	firstTweet := &simpleTweets[0]
	session := lib.CreateSimpleTweetTableSession(os.Getenv("TABLE_NAME"))
	queryResult := session.QueryTweetFromDynamo(firstTweet)
	//TODO Implement a strategy design here.
	isTweetFromDynamoNotMyLatestTweet := queryResult.ID == "" || queryResult.ID != firstTweet.ID
	if isTweetFromDynamoNotMyLatestTweet {
		log.Print("Updating latest tweet in Dynamo")
		session.UpdateLatestTweetInDynamo(firstTweet)
	}

	resp := buildResponse(queryResult)
	return resp, nil
}

func buildResponse(t *lib.SimpleTweetDTO) events.APIGatewayProxyResponse {

	r := events.APIGatewayProxyResponse{
		Body:       t.TweetToJSON(),
		StatusCode: 200,
	}
	return r
}

func main() {
	lambda.Start(handler)
}
