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
	recentTweet := &simpleTweets[0]
	queriedTweet := performOperationOnDynamo(session, nil, lib.GetLatestTweet{})

	if queriedTweet.ID == "" {
		log.Print("Updating latest tweet in Dynamo")
		updatedTweet := performOperationOnDynamo(session, recentTweet, lib.UpdateTweet{})
		returnedTweet = updatedTweet
	}

	if queriedTweet.ID != "" && queriedTweet.ID != recentTweet.ID {
		log.Printf("I have compred to the two tweets, %s, %s", queriedTweet.ID, recentTweet.ID)
		log.Print("Deleting tweet in Dynamo")
		_ = performOperationOnDynamo(session, queriedTweet, lib.DeleteTweet{})

		log.Print("Updating latest tweet in Dynamo")
		updatedTweet := performOperationOnDynamo(session, recentTweet, lib.UpdateTweet{})

		returnedTweet = updatedTweet
	}

	if queriedTweet.ID == recentTweet.ID {
		log.Printf("Returning queried tweet, %s", queriedTweet.ID)
		returnedTweet = queriedTweet
	}

	resp := buildResponse(returnedTweet)
	return resp, nil
}

func performOperationOnDynamo(s *lib.DynamoDbInstance, t *lib.SimpleTweetDTO, o lib.DynamoOperator) *lib.SimpleTweetDTO {

	operation := lib.DynamoOperation{
		DynamoOperator: o,
	}
	operationResult := operation.ExecuteOperation(s, t)
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
