package barrier

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

type httpClientMock struct {
	counter int8
}

func (c *httpClientMock) Get(url string) (resp *http.Response, err error)  {
	c.counter++
	switch url {
	case "https://jpmorganchase.com":
		return &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("Welcome to JP Morgan & Chase")),
		}, nil
	case "https://home.barclays":
		return &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("Welcome to Barclays Home")),
		}, nil
	default:
		panic("unexpected call received")
	}
}

func TestBarrierWithSuccessfulResponses(t *testing.T) {
	httpClientMock := &httpClientMock{0}
	sut := &Barrier{httpClientMock}

	t.Run("when two requests completed successfully", func(t *testing.T) {
		mockedEndpoints := []string{
			"https://jpmorganchase.com",
			"https://home.barclays",
		}

		actual, err := sut.Barrier(mockedEndpoints...)
		if err != nil {
			t.Fail()
		}
		if len(actual) != 2 {
			t.Fail()
		}
		if httpClientMock.counter != 2 {
			t.Fail()
		}
	})
}

type httpClientMockWithError struct {
	counter int8
}

func (c *httpClientMockWithError) Get(url string) (resp *http.Response, err error)  {
	c.counter++
	if c.counter == 1 {
		return &http.Response{
			Body: ioutil.NopCloser(strings.NewReader("Welcome to JP Morgan & Chase")),
		}, nil
	}
	return  nil, errors.New("error establishing http connection")
}

func TestBarrierWithErrorResponses(t *testing.T) {
	httpClientMock := &httpClientMockWithError{0}
	sut := &Barrier{httpClientMock}

	t.Run("when second request fails", func(t *testing.T) {
		mockedEndpoints := []string{
			"https://jpmorganchase.com",
			"https://home.barclays",
		}

		actual, err := sut.Barrier(mockedEndpoints...)
		if err == nil {
			t.Fail()
		}
		if len(actual) > 0 {
			t.Fail()
		}
	})
}