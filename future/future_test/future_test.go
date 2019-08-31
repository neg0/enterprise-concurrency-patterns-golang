package future_test

import (
	Future "concurrency-patterns-go/future"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"testing"
)

func TestFuture(t *testing.T) {
	sut := &Future.Promise{}

	httpClient := &http.Client{}

	t.Run("when_result_is_successful", func(t *testing.T) {
		wg := sync.WaitGroup{}

		sut.Execute(func() (string, error) {
			resp, _ := httpClient.Get("https://httpbin.org/get")
			_, _ = ioutil.ReadAll(resp.Body)

			return "sdsdsdsds", nil
		}, &wg)

		sut.Success(func(s string) {
			if !strings.Contains(s, "Golang Concurrency Ninja!") {
				t.Fail()
			}
		}).Fail(func(e error) {
			t.Fail()
		})
	})

	//t.Run("when_result_is_failed", func(t *testing.T) {
	//	sut.Execute(func() (string, error) {
	//		resp, err := httpClient.Get("http://golang_test_server:8091/error404")
	//		t.Log(resp)
	//		t.Log(err)
	//		if err != nil {
	//			return "", err
	//		}
	//
	//		body, err := ioutil.ReadAll(resp.Body)
	//
	//		return string(body), err
	//	})
	//
	//	successCallsCounter := 0
	//	sut.Success(func(s string) {
	//		successCallsCounter++
	//	}).Fail(func(e error) {
	//		if strings.Contains(strings.ToLower(e.Error()), "error") {
	//			t.Fail()
	//		}
	//	})
	//
	//	if successCallsCounter > 0 {
	//		t.Fail()
	//	}
	//})
}
