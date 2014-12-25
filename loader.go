package inject

import (
	"sync"
	"sync/atomic"
)

type loader struct {
	once  sync.Once
	value atomic.Value
}

func newLoader() *loader {
	return &loader{sync.Once{}, atomic.Value{}}
}

func (this *loader) load(f func() (interface{}, error)) (interface{}, error) {
	this.once.Do(func() {
		value, err := f()
		this.value.Store(&valueErr{value, err})
	})
	valueErr := this.value.Load().(*valueErr)
	return valueErr.value, valueErr.err
}

type valueErr struct {
	value interface{}
	err   error
}
