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

func (l *loader) load(f func() (interface{}, error)) (interface{}, error) {
	l.once.Do(func() {
		value, err := f()
		l.value.Store(&valueErr{value, err})
	})
	valueErr := l.value.Load().(*valueErr)
	return valueErr.value, valueErr.err
}

type valueErr struct {
	value interface{}
	err   error
}
