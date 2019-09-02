package enrichment

import (
	"reflect"
	"testing"
)

func TestCreatingEnrichment(t *testing.T) {
	sut := &Enrichment{
		Merchant: "Sony Inc",
		Category: "Entertainment",
	}

	t.Run("should_be_enrichment_struct_type", func(t *testing.T) {
		if reflect.TypeOf(sut).String() != "*enrichment.Enrichment" {
			t.Log(reflect.TypeOf(sut).String())
			t.Fail()
		}
	})

	t.Run("should_have_merchant_name", func(t *testing.T) {
		if len(sut.Merchant) < 1 {
			t.Log(sut.Merchant)
			t.Fail()
		}
	})

	t.Run("should_have_category_name", func(t *testing.T) {
		if len(sut.Category) < 1 {
			t.Log(sut.Category)
			t.Fail()
		}
	})
}