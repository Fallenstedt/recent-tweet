package lib_test

import (
	"testing"

	"github.com/fallenstedt/latest-tweet/watch-for-latest-tweet/lib"
)

func TestSimpleTweetDTO(t *testing.T) {

	t.Run("should exist", func(t *testing.T) {
		result := givenASimpleTweetDTO()

		if len(result.Message) == 0 {
			t.Fatal("Tweet struct should have a message property")
		}

		if len(result.ID) == 0 {
			t.Fatal("Tweet struct should have an ID property")
		}
	})

	t.Run("TweetToJSON should return a string version of the struct", func(t *testing.T) {
		tweet := givenASimpleTweetDTO()
		result, err := tweet.TweetToJSON()
		expected := `{"message":"Hello","id":"123"}`

		if err != nil {
			t.Fatal("Could not create JSON from struct")
		}

		if result != expected {
			t.Fatal("Could not parse JSON correctly")
		}
	})
}

// func TestCreateSimpleTweetDTO() {}

func givenASimpleTweetDTO() lib.SimpleTweetDTO {
	result := lib.SimpleTweetDTO{
		Message: "Hello",
		ID:      "123",
	}
	return result
}
