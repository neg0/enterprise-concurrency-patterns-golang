package bank_transaction

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestCreationOfBankTransactionFromJSONByte(t *testing.T) {
	var sut Transaction
	err := json.Unmarshal([]byte(testCase), &sut)

	t.Run("should_successfully_create_transaction_from_json_string", func(t *testing.T) {
		if err != nil {
			t.Log(err.Error())
			t.Fail()
		}
	})

	t.Run("should_have_value_of_float_with_two_decimal_and_currency_code", func(t *testing.T) {
		actual, _ := sut.Amount.Value()

		if reflect.TypeOf(actual).String() != "float64" {
			t.Log(reflect.TypeOf(actual).String())
			t.Fail()
		}

		if actual != 11.99 {
			t.Log(actual)
			t.Fail()
		}
	})
}

const testCase = `{
					"AccountId": "3234672871",
					"TransactionId": "1b20d4cc-29a5-4c51-a8f3-f8b08a1a7661",
					"Amount": {
					  "Amount": "11.9923",
					  "Currency": "GBP"
					},
					"TransactionInformation": "Netflix Subscription Intl."
				  }`