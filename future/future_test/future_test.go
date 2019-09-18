package future_test

import (
	Future "concurrency-patterns-go/future"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestFuture(t *testing.T) {
	sut := Future.NewPromise()

	t.Run("when_closure_used_as_parameter_with_success", func(t *testing.T) {
		sut.
			Promise(setContext("http://golang_test_server:8091")).
			Catch(func(e error) {
				t.Log(e.Error())
				t.Fail()
			}).
			Then(func(b []byte) {
				if len(b) < 1 {
					t.Fail()
				}
			}).
			Await()
	})

	t.Run("when_closure_used_as_parameter_with_error", func(t *testing.T) {
		sut.
			Promise(setContext("http://golang_test_server:666")).
			Catch(func(e error) {
				if e == nil {
					t.Fail()
				}
			}).
			Then(func(b []byte) {
				if len(b) > 0 {
					t.Log(b)
					t.Fail()
				}
			}).
			Await()
	})
}

func setContext(url string) func() ([]byte, error) {
	return func() ([]byte, error) {
		httpClient := &http.Client{}
		resp, err := httpClient.Get(url)
		if err != nil {
			return nil, err
		}

		bodyResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		httpStatusCode := resp.StatusCode
		if httpStatusCode > 299 || httpStatusCode < 200 {
			errMsg := fmt.Sprintf(
				"status code %d %s is returned. %s",
				httpStatusCode,
				http.StatusText(httpStatusCode),
				string(bodyResp),
			)

			return nil, errors.New(errMsg)
		}

		return bodyResp, nil
	}
}
