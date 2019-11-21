package lib

import "log"

type Handler interface {
	Request(s *DynamoDbInstance, q *SimpleTweetDTO, t *SimpleTweetDTO) *SimpleTweetDTO
}

type QueriedTweetEmpty struct {
	Next Handler
}

func (h *QueriedTweetEmpty) Request(session *DynamoDbInstance, queriedTweet *SimpleTweetDTO, recentTweet *SimpleTweetDTO) *SimpleTweetDTO {
	if queriedTweet.ID == "" {
		log.Print("Updating latest tweet in Dynamo")
		updatedTweet := performOperationOnDynamo(session, recentTweet, UpdateTweet{})
		return updatedTweet
	} else {
		return h.Next.Request(session, queriedTweet, recentTweet)
	}
}

type QueriedTweetIdDoesNotMatchRecentTweetId struct {
	Next Handler
}

func (h *QueriedTweetIdDoesNotMatchRecentTweetId) Request(session *DynamoDbInstance, queriedTweet *SimpleTweetDTO, recentTweet *SimpleTweetDTO) *SimpleTweetDTO {
	if queriedTweet.ID != "" && queriedTweet.ID != recentTweet.ID {
		log.Printf("I have compred to the two tweets, %s, %s", queriedTweet.ID, recentTweet.ID)
		log.Print("Deleting tweet in Dynamo")
		_ = performOperationOnDynamo(session, queriedTweet, DeleteTweet{})

		log.Print("Updating latest tweet in Dynamo")
		updatedTweet := performOperationOnDynamo(session, recentTweet, UpdateTweet{})

		return updatedTweet
	} else {
		return h.Next.Request(session, queriedTweet, recentTweet)
	}
}

type QueriedTweetMatchesRecentTweet struct{}

func (h *QueriedTweetMatchesRecentTweet) Request(session *DynamoDbInstance, queriedTweet *SimpleTweetDTO, recentTweet *SimpleTweetDTO) *SimpleTweetDTO {
	if queriedTweet.ID == recentTweet.ID {
		return queriedTweet
	} else {
		panic("Unable to resolve operation! Could not determine if user has the latest tweet")
	}
}

func performOperationOnDynamo(s *DynamoDbInstance, t *SimpleTweetDTO, o DynamoOperator) *SimpleTweetDTO {

	operation := DynamoOperation{
		DynamoOperator: o,
	}
	operationResult := operation.ExecuteOperation(s, t)
	return operationResult
}
