package future

import "sync"

type resolve func([]byte) Promise
type reject func(error) Promise
type execution func() ([]byte, error)

type Promise interface {
	Promise(f execution) Promise
	Then(f resolve) Promise
	Catch(f reject) Promise
	Await()
}

type promise struct {
	resolve
	reject
	sync.WaitGroup
}

func NewPromise() Promise {
	return &promise{}
}

func (p *promise) Then(f resolve) Promise {
	p.resolve = f

	return p
}

func (p *promise) Catch(f reject) Promise {
	p.reject = f

	return p
}

func (p *promise) Await() {
	p.Wait()
}

func (p *promise) Promise(f execution) Promise {
	p.Add(1)
	go func(p *promise) {
		bytes, err := f()
		if err != nil {
			p.reject(err)
		} else {
			p.resolve(bytes)
		}
		p.Done()
	}(p)

	return p
}
