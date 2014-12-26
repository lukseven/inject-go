package inject

import (
	"fmt"
)

type Module interface {
	fmt.Stringer
	Bind(from interface{}) Builder
	BindTagged(from interface{}, tag string) Builder
}

func CreateModule() Module { return createModule() }

type Builder interface {
	To(to interface{})
	ToSingleton(singleton interface{})
	ToConstructor(constructor interface{})
	ToSingletonConstructor(constructor interface{})
	ToTaggedConstructor(constructor interface{})
	ToTaggedSingletonConstructor(constructor interface{})
}

type Injector interface {
	fmt.Stringer
	Get(from interface{}) (interface{}, error)
	GetTagged(from interface{}, tag string) (interface{}, error)
	Call(function interface{}) ([]interface{}, error)
	CallTagged(taggedFunction interface{}) ([]interface{}, error)
}

func CreateInjector(modules ...Module) (Injector, error) { return createInjector(modules) }
