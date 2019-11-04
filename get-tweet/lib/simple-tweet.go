package lib

import (
	"encoding/json"
	"errors"

	"github.com/dghubble/go-twitter/twitter"
)

// SimpleTweetDTO creates a tweet with
// only the message and its ID.
type SimpleTweetDTO struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}

// TweetToJSON takes a simple tweet
// and returns a JSON version of it
func (t *SimpleTweetDTO) TweetToJSON() (string, error) {
	tweetsJSON, err := json.Marshal(t)
	return string(tweetsJSON), err
}

// CreateSimpleTweetDTO creates a SimpleTweetDTO struct
// from an array of tweets
func CreateSimpleTweetDTO(tweets *([]twitter.Tweet)) ([]SimpleTweetDTO, error) {

	if len(*tweets) == 0 {
		e := errors.New("Cannot create simple tweet from empty array")
		emptyTweet := []SimpleTweetDTO{SimpleTweetDTO{}}
		return emptyTweet, e
	}

	var message string
	var simpleTweets []SimpleTweetDTO
	for i := 0; i < len(*tweets); i++ {
		firstTweet := (*tweets)[i]

		if len(firstTweet.FullText) == 0 {
			message = firstTweet.Text
		} else {
			message = firstTweet.FullText
		}

		t := SimpleTweetDTO{
			Message: message,
			ID:      firstTweet.IDStr,
		}

		simpleTweets = append(simpleTweets, t)
	}

	return simpleTweets, nil
}
