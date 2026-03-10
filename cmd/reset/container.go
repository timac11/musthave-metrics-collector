package reset

import (
	"sync"
)

type Resettable interface {
	Reset()
}

type Pool[T Resettable] struct {
	pool sync.Pool
}

func New[T Resettable](new func() T) *Pool[T] {
	p := &Pool[T]{}
	p.pool.New = func() any {
		return new()
	}
	return p
}

func (p *Pool[T]) Get() T {
	obj := p.pool.Get().(T)
	obj.Reset()
	return obj
}

func (p *Pool[T]) Put(obj T) {
	p.pool.Put(obj)
}
