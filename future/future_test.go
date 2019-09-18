package future

import (
	"errors"
	"strings"
	"testing"
)

func TestFuture(t *testing.T) {
	sut := NewPromise()

	t.Run("when_result_is_successful", func(t *testing.T) {
		sut.
			Promise(func() ([]byte, error) {
				return []byte("Golang2 Concurrency Ninja!"), nil
			}).
			Then(func(b []byte) Promise {
				if !strings.Contains(string(b), "Golang2 Concurrency Ninja!") {
					t.Fail()
				}

				return  NewPromise().Promise(func() ([]byte, error) {
					return []byte("Golang Concurrency Ninja!"), nil
				}).Then(func(b []byte) Promise {
					if !strings.Contains(string(b), "Golang Concurrency Ninja!") {
						t.Fail()
					}
					return sut
				}).Catch(func(e error) Promise {
					t.Fail()
					return sut
				})
			}).
			Catch(func(e error) Promise {
				t.Fail()
				return sut
			}).
			Await()
	})

	t.Run("when_result_is_failed", func(t *testing.T) {
		successCallsCounter := 0
		sut.
			Promise(func() ([]byte, error) {
				return nil, errors.New("error occurred")
			}).Then(func(b []byte) Promise {
			successCallsCounter++
			return sut
		}).Catch(func(e error) Promise {
			if strings.Compare(e.Error(), "error occurred") == -1 {
				t.Fail()
			}
			return sut
		}).Await()

		if successCallsCounter > 0 {
			t.Fail()
		}
	})
}
