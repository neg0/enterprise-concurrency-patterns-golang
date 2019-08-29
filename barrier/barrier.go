package barrier

import (
	"io/ioutil"
	"net/http"
	"time"
)

type Barrier struct {
	HttpClient
}

type barrierResp struct {
	Err  error
	Resp string
}

var TimeoutMilliseconds = 5000

func NewBarrier() *Barrier {
	return &Barrier{
		&http.Client{
			Timeout: time.Duration(TimeoutMilliseconds) * time.Millisecond,
		},
	}
}

func (b *Barrier) Barrier(endpoints ...string) ([]string, error) {
	numberOfRequests := len(endpoints)

	in := make(chan barrierResp, len(endpoints))
	for _, endpoint := range endpoints {
		go b.doRequest(in, endpoint)
	}

	responses := make([]string, numberOfRequests)
	for i := 0; i < numberOfRequests; i++ {
		resp := <-in
		if resp.Err != nil {
			return nil, resp.Err
		}
		responses[i] = resp.Resp
	}
	defer close(in)

	return responses, nil
}

func (b *Barrier) doRequest(out chan<- barrierResp, url string) {
	res := barrierResp{}

	resp, err := b.HttpClient.Get(url)
	if err != nil {
		res.Err = err
		out <- res
		return
	}

	byt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		res.Err = err
		out <- res
		return
	}

	res.Resp = string(byt)
	out <- res
}
