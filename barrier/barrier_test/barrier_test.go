package barrier_test

import (
	"concurrency-patterns-go/barrier"
	"strings"
	"testing"
)

func TestBarrier(t *testing.T) {
	t.Run("with_correct_endpoints", func(t *testing.T) {
		endpoints := []string{
			"http://golang_test_server:8091/oauth/account",
			"http://golang_test_server:8091/kms",
		}
		result, _ := barrier.NewBarrier().Barrier(endpoints...)

		if len(result) != 2 {
			t.Fail()
		}
	})

	t.Run("with_one_incorrect_endpoint", func(t *testing.T) {
		endpoints := []string{
			"http://malformed-url",
			"http://golang_test_server:8091/oauth/account",
		}
		_, err := barrier.NewBarrier().Barrier(endpoints...)

		if !strings.Contains(err.Error(), "no such host") {
			t.Fail()
		}
	})

	t.Run("with_short_timeout", func(t *testing.T) {
		endpoints := []string{
			"http://golang_test_server:8091/oauth/account",
			"http://golang_test_server:8091/kms",
		}
		barrier.TimeoutMilliseconds = 1

		res, err := barrier.NewBarrier().Barrier(endpoints...)

		if len(res) > 0 {
			t.Fail()
		}
		if !strings.Contains(err.Error(), "Client.Timeout") {
			t.Fail()
		}
	})
}
