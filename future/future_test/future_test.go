package future_test

import (
	Future "concurrency-patterns-go/future"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestFuture(t *testing.T) {
	sut := &Future.Promise{}

	t.Run("when_closure_used_as_parameter_with_Success", func(t *testing.T) {
		sut.Then(func(s string) {
			t.Log(s)
		}).Catch(func(e error) {
			t.Log(e.Error())
			t.Fail()
		}).Future(setContext("http://golang_test_server:8091"))
	})

	t.Run("when_closure_used_as_parameter_with_Error", func(t *testing.T) {
		sut.Then(func(s string) {
			t.Log(s)
			t.Fail()
		}).Catch(func(e error) {
			t.Log(e.Error())
		}).Future(setContext("http://golang_test_server:9999"))
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
