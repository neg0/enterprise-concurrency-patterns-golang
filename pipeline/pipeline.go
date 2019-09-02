package pipeline

import (
	"bytes"
	Bank "concurrency-patterns-go/pipeline/bank_transaction"
	"encoding/json"
	"io/ioutil"
	"sync"
)

type Pipeline struct {
	sync.WaitGroup
	HttpClient
}

func (p *Pipeline) Enrich(transactions []byte) ([]Bank.Transaction, error) {
	parsedTrx, err := p.process(transactions)
	if err != nil {
		return nil, err
	}

	transactionsForEnrichment := p.generatePipelines(parsedTrx)

	enrichedTransactionsWithMerchant, err := p.identifyMerchant(transactionsForEnrichment)
	if err != nil {
		return nil, err
	}
	enrichedTransactionsWithCategory, err := p.identifyCategory(enrichedTransactionsWithMerchant)
	if err != nil {
		return nil, err
	}

	var trxCollection []Bank.Transaction
	for i := 0; i < len(parsedTrx); i++ {
		trxCollection = append(trxCollection, <-enrichedTransactionsWithCategory)
	}

	return trxCollection, nil
}

func (p *Pipeline) process(transactions []byte) ([]Bank.Transaction, error) {
	var parsedTrx []Bank.Transaction

	errCh := make(chan error, 1)
	go func(trx []byte) {
		err := json.Unmarshal(transactions, &parsedTrx)
		errCh<- err
		close(errCh)
	}(transactions)

	err := <-errCh
	if err != nil {
		return nil, err
	}

	return parsedTrx, nil
}

func (p *Pipeline) generatePipelines(transactions []Bank.Transaction) <-chan Bank.Transaction {
	outChString := make(chan Bank.Transaction, len(transactions))

	go func() {
		for _, trx := range transactions {
			outChString <- trx
		}
		close(outChString)
	}()

	return outChString
}

func (p *Pipeline) identifyMerchant(in <-chan Bank.Transaction) (<-chan Bank.Transaction, error) {
	out := make(chan Bank.Transaction, len(in))
	errCh := make(chan error, 1)
	var err error

	go func(hasError error) {
		for v := range in {
			payload, _ := json.Marshal(v.TransactionID)
			resp, err := p.HttpClient.Post(
				"http://golang_test_server:8091/enrich/merchant",
				"text/plain",
				bytes.NewBuffer(payload),
			)
			if err != nil {
				errCh<- err
				close(errCh)
				continue
			}
			bodyResp, _ := ioutil.ReadAll(resp.Body)

			v.Enrichment.Merchant = string(bodyResp)
			out <- v
		}
		close(out)
		err = <-errCh
	}(err)

	return out, err
}

func (p *Pipeline) identifyCategory(in <-chan Bank.Transaction) (<-chan Bank.Transaction, error) {
	out := make(chan Bank.Transaction, len(in))

	go func() {
		for v := range in {
			payload, _ := json.Marshal(v.TransactionID)

			resp, _ := p.HttpClient.Post(
				"http://golang_test_server:8091/enrich/category",
				"text/plain",
				bytes.NewBuffer(payload),
			)
			bodyResp, _ := ioutil.ReadAll(resp.Body)

			v.Enrichment.Category = string(bodyResp)
			out <- v
		}
		close(out)
	}()

	return out, nil
}
