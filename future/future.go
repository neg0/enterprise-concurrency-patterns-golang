package future

import "sync"

type successful func(string)
type failure func(error)
type execution func() (string, error)

type Promise struct {
	sync.WaitGroup
	successful
	failure
}

func (p *Promise) Success(f successful) *Promise {
	p.successful = f

	return p
}

func (p *Promise) Fail(f failure) *Promise {
	p.failure = f

	return p
}

func (p *Promise) Execute(f execution) {
	p.Add(1)
	go func(p *Promise) {
		str, err := f()
		if err != nil {
			p.failure(err)
		} else {
			p.successful(str)
		}
		p.Done()
	}(p)
	p.Wait()
}

func (p *Promise) Finally() {
	p.Wait()
}