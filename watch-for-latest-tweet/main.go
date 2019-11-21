package main

import (
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/fallenstedt/latest-tweet/lib"
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

	var returnedTweet *lib.SimpleTweetDTO
	session := lib.CreateSimpleTweetTableSession(os.Getenv("TABLE_NAME"))
	queriedTweet := getLatestTweetFromDynamo(session, lib.GetLatestTweet{})
	recentTweet := &simpleTweets[0]
	chain := buildTweetHandleChain()

	returnedTweet = chain.Request(session, queriedTweet, recentTweet)

	resp := buildResponse(returnedTweet)
	return resp, nil
}

func buildTweetHandleChain() lib.Handler {
	quriedTweetEmpty := lib.QueriedTweetEmpty{}
	queriedTweetIDDoesNotMatchRecentTweetID := lib.QueriedTweetIdDoesNotMatchRecentTweetId{}
	queriedTweetMatchesRecentTweet := lib.QueriedTweetMatchesRecentTweet{}

	quriedTweetEmpty.Next = &queriedTweetIDDoesNotMatchRecentTweetID
	queriedTweetIDDoesNotMatchRecentTweetID.Next = &queriedTweetMatchesRecentTweet

	return &quriedTweetEmpty
}

func getLatestTweetFromDynamo(s *lib.DynamoDbInstance, o lib.DynamoOperator) *lib.SimpleTweetDTO {
	operation := lib.DynamoOperation{
		DynamoOperator: o,
	}
	operationResult := operation.ExecuteOperation(s, nil)
	return operationResult
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
