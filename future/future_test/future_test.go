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
		sut.Success(func(s string) {
			t.Log(s)
		}).Fail(func(e error) {
			t.Log(e.Error())
			t.Fail()
		}).Execute(setContext("http://127.0.0.1:8091"))
	})

	t.Run("when_closure_used_as_parameter_with_Error", func(t *testing.T) {
		sut.Success(func(s string) {
			t.Log(s)
			t.Fail()
		}).Fail(func(e error) {
			t.Log(e.Error())
		}).Execute(setContext("http://127.0.0.1:9999"))
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
