package pipeline_test

import (
	"concurrency-patterns-go/pipeline"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestPipeline(t *testing.T) {
	testCases, err := ioutil.ReadFile("./mock.json")
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}

	sut := &pipeline.Pipeline{}
	sut.HttpClient = &http.Client{}

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
}
