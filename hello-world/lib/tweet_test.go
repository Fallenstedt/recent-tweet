package lib_test

import (
	"hello-sam/hello-world/lib"
	"testing"
)

func TestTweet(t *testing.T) {
	t.Run("should exist", func(t *testing.T) {
		result := lib.Tweet{Message: "Hello", ID: "123"}

		if len(result.Message) == 0 {
			t.Fatal("Tweet struct should have a message property")
		}

		if len(result.ID) == 0 {
			t.Fatal("Tweet struct should have an ID property")
		}
	})
}
