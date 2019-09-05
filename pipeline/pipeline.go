package pipeline

import (
	"bytes"
	Bank "concurrency-patterns-go/pipeline/bank_transaction"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Pipeline struct {
	sync.WaitGroup
	HttpClient
}

var once sync.Once
var instance *Pipeline
var DefaultTimeOutInSecond = 1

func init() {
	once.Do(func() {
		instance = &Pipeline{}
		instance.HttpClient = &http.Client{
			Timeout: time.Duration(int64(DefaultTimeOutInSecond)),
		}
	})
}

func Instance() *Pipeline {
	return instance
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
		errCh <- err
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
	numOfWorkers := positiveLength(len(in))

	out := make(chan Bank.Transaction, numOfWorkers)
	errCh := make(chan error, numOfWorkers)
	done := make(chan bool)
	var err error

	p.Add(1)
	go func(hasError error) {
		for v := range in {
			payload, _ := json.Marshal(v.TransactionID)
			resp, err := p.HttpClient.Post(
				"http://golang_test_server:8091/enrich/merchant",
				"text/plain",
				bytes.NewBuffer(payload),
			)
			if err != nil {
				errCh <- err
				continue
			}

			bodyResp, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errCh <- err
				continue
			}

			v.Enrichment.Merchant = string(bodyResp)
			out <- v
		}
		close(out)
		err = <-errCh
		close(errCh)
		done<- true
	}(err)

	select {
		case <-done:
			return out, err
		case <-time.After(time.Millisecond * time.Duration(int64(DefaultTimeOutInSecond))):
			return out, err
	}
}

func (p *Pipeline) identifyCategory(in <-chan Bank.Transaction) (<-chan Bank.Transaction, error) {
	numOfWorkers := positiveLength(len(in))

	out := make(chan Bank.Transaction, numOfWorkers)
	errOut := make(chan error, numOfWorkers)
	done := make(chan bool)
	var err error

	go func(hasError error) {
		for v := range in {
			payload, err := json.Marshal(v.TransactionID)
			if err != nil {
				errOut <- err
				break
			}

			resp, err := p.HttpClient.Post(
				"http://golang_test_server:8091/enrich/category",
				"text/plain",
				bytes.NewBuffer(payload),
			)
			if err != nil {
				errOut <- err
				break
			}

			bodyResp, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				errOut <- err
				break
			}

			v.Enrichment.Category = string(bodyResp)
			out <- v
		}
		close(out)
		err = <-errOut
		close(errOut)
		done<- true
	}(err)
	time.Sleep(time.Millisecond * 1)

	select {
		case <-done:
			return out, err
		case <-time.After(time.Millisecond * time.Duration(int64(DefaultTimeOutInSecond))):
			return out, err
	}
}

func positiveLength(numOfWorkers int) int {
	if numOfWorkers == 0 {
		numOfWorkers = 1
	}

	return numOfWorkers
}