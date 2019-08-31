package future

import (
	"errors"
	"strings"
	"testing"
)

func TestFuture(t *testing.T) {
	sut := &Promise{}

	t.Run("when_result_is_successful", func(t *testing.T) {
		sut.
			Future(func() (string, error) {
				return "Golang Concurrency Ninja!", nil
			}).Then(func(s string) {
				if !strings.Contains(s, "Golang Concurrency Ninja!") {
					t.Fail()
				}
			}).Catch(func(e error) {
				t.Fail()
			}).
			Execute()
	})

	t.Run("when_result_is_failed", func(t *testing.T) {
		successCallsCounter := 0
		sut.
			Future(func() (string, error) {
				return "", errors.New("error occurred")
			}).Then(func(s string) {
				successCallsCounter++
			}).Catch(func(e error) {
				if strings.Compare(e.Error(), "error occurred") == -1 {
					t.Fail()
				}
			}).Execute()

		if successCallsCounter > 0 {
			t.Fail()
		}
	})
}
