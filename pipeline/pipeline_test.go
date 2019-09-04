package pipeline

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"testing"
)

type httpClientMock struct {
	merchantCounter int
	categoryCounter int
}

func (client *httpClientMock) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	if url == "http://golang_test_server:8091/enrich/merchant" {
		client.merchantCounter++
		randMerchants := []string{
			"Netflix",
			"Mark & Spencer",
			"Amazon",
		}

		return &http.Response{
			Body: ioutil.NopCloser(strings.NewReader(randMerchants[rand.Intn(len(randMerchants))])),
		}, nil
	}

	client.categoryCounter++
	randCats := []string{
		"Entertainment",
		"Leisure",
		"Shopping",
	}

	return &http.Response{
		Body: ioutil.NopCloser(strings.NewReader(randCats[rand.Intn(len(randCats))])),
	}, nil
}

func TestPipeline(t *testing.T) {
	testCases, err := ioutil.ReadFile("./pipeline_test/mock.json")
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	sut := &Pipeline{}
	mockedHttpClient := &httpClientMock{0, 0}
	sut.HttpClient = mockedHttpClient

	actual, err := sut.Enrich(testCases)

	t.Run("should_finish_enrichment_process_without_an_error", func(t *testing.T) {
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}
	})

	t.Run("should_return_three_enriched_transactions", func(t *testing.T) {
		if len(actual) != 3 {
			t.Log(actual)
			t.Log(fmt.Sprintf("Length is: %d", len(actual)))
			t.Fail()
		}
	})

	t.Run("should_have_transactions_enriched_with_merchant_and_category", func(t *testing.T) {
		for _, act := range actual {
			t.Log(act)
			if len(act.Enrichment.Merchant) < 1 {
				t.Log(act.Enrichment.Merchant)
				t.Fail()
			}

			if len(act.Enrichment.Category) < 1 {
				t.Log(act.Enrichment.Category)
				t.Fail()
			}
		}
	})

	t.Run("should_have_called_category_enrichment_microservice_three_times", func(t *testing.T) {
		if mockedHttpClient.categoryCounter != 3 {
			t.Log(mockedHttpClient.categoryCounter)
			t.Fail()
		}
	})

	t.Run("should_have_called_merchant_enrichment_microservice_three_times", func(t *testing.T) {
		if mockedHttpClient.merchantCounter != 3 {
			t.Log(mockedHttpClient.merchantCounter)
			t.Fail()
		}
	})
}

type httpClientErrorMock struct {
	merchantCounter int
	categoryCounter int
}

func (client *httpClientErrorMock) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	if url == "http://golang_test_server:8091/enrich/merchant" {
		client.merchantCounter++
		randMerchants := []string{
			"Netflix",
			"Mark & Spencer",
			"Amazon",
		}

		switch client.merchantCounter {
		case 1:
			return &http.Response{
					Body:       ioutil.NopCloser(strings.NewReader(randMerchants[rand.Intn(len(randMerchants))])),
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
				},
				errors.New(fmt.Sprintf("made up error call number: %d", 1))
		case 2:
			return &http.Response{
					Body:       ioutil.NopCloser(strings.NewReader(randMerchants[rand.Intn(len(randMerchants))])),
					StatusCode: http.StatusNotFound,
					Status:     http.StatusText(http.StatusNotFound),
				},
				errors.New(fmt.Sprintf("made up error call number: %d", 2))
		default:
			return &http.Response{
				Body: ioutil.NopCloser(strings.NewReader(randMerchants[rand.Intn(len(randMerchants))])),
			}, nil
		}
	}

	client.categoryCounter++
	randCats := []string{
		"Entertainment",
		"Leisure",
		"Shopping",
	}

	switch client.categoryCounter {
	case 1:
		return &http.Response{
				Body:       ioutil.NopCloser(strings.NewReader(randCats[rand.Intn(len(randCats))])),
				StatusCode: http.StatusNotFound,
			},
			errors.New(fmt.Sprintf("made up error call number: %d", 1))
	case 2:
		return &http.Response{
				Body:       ioutil.NopCloser(strings.NewReader(randCats[rand.Intn(len(randCats))])),
				StatusCode: http.StatusNotFound,
			},
			errors.New(fmt.Sprintf("made up error call number: %d", 2))
	default:
		return &http.Response{
			Body: ioutil.NopCloser(strings.NewReader(randCats[rand.Intn(len(randCats))])),
		}, nil
	}
}

func TestPipelineFailureCases(t *testing.T) {
	sut := &Pipeline{}
	mockedHttpClient := &httpClientErrorMock{0, 0}
	sut.HttpClient = mockedHttpClient

	t.Run("should_return_an_error_when_single_transaction_is_passed_in", func(t *testing.T) {
		testCase := `{
					  "AccountId": "3234672871",
					  "TransactionId": "3ed19f1c-096b-412d-b27b-9bc2b757a572",
					  "Amount": {
						"Amount": "48.95",
						"Currency": "GBP"
					  },
					  "TransactionInformation": "Tesco SuperStore Kensington London"
					}`

		actual, err := sut.Enrich([]byte(testCase))

		if len(actual) > 0 {
			t.Fail()
		}

		if err == nil {
			t.Log(err)
		}
	})

	t.Run("should_return_an_error_when_one_of_fields_is_incorrect_type", func(t *testing.T) {
		testCase := `[{
					  "AccountId": "3234672871",
					  "TransactionId": "3ed19f1c-096b-412d-b27b-9bc2b757a572",
					  "Amount": {
						"Amount": 48.95,
						"Currency": "GBP"
					  },
					  "TransactionInformation": "Tesco SuperStore Kensington London"
					}]`

		actual, err := sut.Enrich([]byte(testCase))

		if len(actual) > 0 {
			t.Fail()
		}

		if err == nil {
			t.Fail()
			t.Log(err)
		}
		t.Log(err)
	})

	t.Run("should_return_an_error_identifying_merchant", func(t *testing.T) {
		testCase := `[{
					  "AccountId": "3234672871",
					  "TransactionId": "3ed19f1c-096b-412d-b27b-9bc2b757a572",
					  "Amount": {
						"Amount": "48.95",
						"Currency": "GBP"
					  },
					  "TransactionInformation": "Tesco SuperStore Kensington London"
					}]`

		actual, err := sut.Enrich([]byte(testCase))

		if len(actual) > 0 {
			t.Fail()
		}

		if err == nil {
			t.Fail()
			t.Log(err)
		}
		t.Log(err)
	})
}
