package bank_transaction

import (
	"concurrency-patterns-go/pipeline/enrichment"
	"fmt"
	"strconv"
)

type ValueObjectFloat interface {
	Value() (float64, error)
}

type Transaction struct {
	AccountID              string                `json:"AccountId"`
	TransactionID          string                `json:"TransactionId"`
	Amount                 Amount                `json:"Amount"`
	TransactionInformation string                `json:"TransactionInformation"`
	Enrichment             enrichment.Enrichment `json:"-"`
}

type Amount struct {
	Amount   string `json:"amount"`
	Currency string `json:"currency"`
}

func (a Amount) Value() (float64, error) {
	val, err := strconv.ParseFloat(a.Amount, 2)
	if err != nil {
		return 0.00, err
	}

	return strconv.ParseFloat(fmt.Sprintf("%.2f", val), 2)
}
