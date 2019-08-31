package future

import "sync"

type resolve func(string)
type reject func(error)
type execution func() (string, error)

type Promise struct {
	resolve
	reject
	sync.WaitGroup
}

func (p *Promise) Then(f resolve) *Promise {
	p.resolve = f

	return p
}

func (p *Promise) Catch(f reject) *Promise {
	p.reject = f

	return p
}

func (p *Promise) Execute() {
	p.Wait()
}

func (p *Promise) Future(f execution) *Promise {
	p.Add(1)
	go func(p *Promise) {
		str, err := f()
		if err != nil {
			p.reject(err)
		} else {
			p.resolve(str)
		}
		p.Done()
	}(p)

	return p
}
