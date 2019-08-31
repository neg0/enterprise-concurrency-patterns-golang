package future_test

import (
	Future "concurrency-patterns-go/future"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestFuture(t *testing.T) {
	sut := &Future.Promise{}

	t.Run("when_closure_used_as_parameter_with_success", func(t *testing.T) {
		sut.
			Future(setContext("http://golang_test_server:8091")).
			Catch(func(e error) {
				t.Log(e.Error())
				t.Fail()
			}).
			Then(func(s string) {
				if len(s) < 1 {
					t.Fail()
				}
			}).
			Execute()
	})

	t.Run("when_closure_used_as_parameter_with_error", func(t *testing.T) {
		sut.
			Future(setContext("http://golang_test_server:666")).
			Catch(func(e error) {
				if e == nil {
					t.Fail()
				}
			}).
			Then(func(s string) {
				if len(s) > 0 {
					t.Log(s)
					t.Fail()
				}
			}).
			Execute()
	})
}

func setContext(url string) func() (string, error) {
	return func() (string, error) {
		httpClient := &http.Client{}
		resp, err := httpClient.Get(url)
		if err != nil {
			return "", err
		}

		bodyResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		return string(bodyResp), nil
	}
}
