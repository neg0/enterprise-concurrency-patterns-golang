package future

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestFuture(t *testing.T) {
	sut := &Promise{}

	t.Run("when_result_is_successful", func(t *testing.T) {
		sut.Success(func(s string) {
			if !strings.Contains(s, "Golang Concurrency Ninja!") {
				t.Fail()
			}
		}).Fail(func(e error) {
			t.Fail()
		}).Execute(func() (string, error) {
			return "Golang Concurrency Ninja!", nil
		})

	})

	t.Run("when_result_is_failed", func(t *testing.T) {
		successCallsCounter := 0
		sut.Success(func(s string) {
			successCallsCounter++
		}).Fail(func(e error) {
			if strings.Compare(e.Error(), "error occurred") == -1 {
				t.Fail()
			}
		}).Execute(func() (string, error) {
			return "", errors.New("error occurred")
		})

		if successCallsCounter > 0 {
			t.Fail()
		}
	})

	t.Run("when_closure_used_as_parameter_with_Success result", func(t *testing.T) {
		sut.Success(func(s string) {
			t.Log(s)
		}).Fail(func(e error) {
			t.Log(e.Error())
			t.Fail()
		}).Execute(setContext("Hello"))
	})
}

func setContext(msg string) execution {
	return func() (string, error) {
		httpClient := &http.Client{}
		resp, err := httpClient.Get("http://127.0.0.1:8091/")
		if err != nil {
			return "", err
		}

		bodyResp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}

		msg = fmt.Sprintf("%s Closure!\n", string(bodyResp))

		return msg, nil
	}
}
