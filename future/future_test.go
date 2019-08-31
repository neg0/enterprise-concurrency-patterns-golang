package future

import (
	"errors"
	"strings"
	"testing"
)

func TestFuture(t *testing.T) {
	sut := &Promise{}

	t.Run("when_result_is_successful", func(t *testing.T) {
		sut.Then(func(s string) {
			if !strings.Contains(s, "Golang Concurrency Ninja!") {
				t.Fail()
			}
		}).Catch(func(e error) {
			t.Fail()
		}).Future(func() (string, error) {
			return "Golang Concurrency Ninja!", nil
		})

	})

	t.Run("when_result_is_failed", func(t *testing.T) {
		successCallsCounter := 0
		sut.Then(func(s string) {
			successCallsCounter++
		}).Catch(func(e error) {
			if strings.Compare(e.Error(), "error occurred") == -1 {
				t.Fail()
			}
		}).Future(func() (string, error) {
			return "", errors.New("error occurred")
		})

		if successCallsCounter > 0 {
			t.Fail()
		}
	})
}
