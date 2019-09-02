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
	p.Add(1)
	var parsedTrx []Bank.Transaction

	var err error
	go func(trx []byte, w *sync.WaitGroup, err error) {
		err = json.Unmarshal(transactions, &parsedTrx)
		if err != nil {
			return
		}
		w.Done()
	}(transactions, &p.WaitGroup, err)
	p.WaitGroup.Wait()

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
		defer close(outChString)
	}()

	return outChString
}

func (p *Pipeline) identifyMerchant(in <-chan Bank.Transaction) (<-chan Bank.Transaction, error) {
	out := make(chan Bank.Transaction, len(in))

	go func() {
		for v := range in {
			payload, err := json.Marshal(v.TransactionID)
			if err != nil {
				return
			}
			resp, err := p.HttpClient.Post(
				"http://golang_test_server:8091/enrich/merchant",
				"text/plain",
				bytes.NewBuffer(payload),
			)
			if err != nil {
				return
			}

			bodyResp, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return
			}

			v.Enrichment.Merchant = string(bodyResp)
			out <- v
		}
		defer close(out)
	}()

	return out, nil
}

func (p *Pipeline) identifyCategory(in <-chan Bank.Transaction) (<-chan Bank.Transaction, error) {
	var err error
	out := make(chan Bank.Transaction, len(in))

	go func(err error) {
		for v := range in {
			payload, err := json.Marshal(v.TransactionID)
			if err != nil {
				return
			}
			resp, err := p.HttpClient.Post(
				"http://golang_test_server:8091/enrich/category",
				"text/plain",
				bytes.NewBuffer(payload),
			)
			if err != nil {
				return
			}

			bodyResp, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return
			}

			v.Enrichment.Category = string(bodyResp)
			out <- v
		}
		defer close(out)
	}(err)

	if err != nil {
		return nil, err
	}

	return out, err
}
