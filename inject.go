package inject

import (
	"fmt"
)

type Module interface {
	fmt.Stringer
	Bind(from ...interface{}) Builder
	BindTagged(tag string, from ...interface{}) Builder
	BindInterface(fromInterface ...interface{}) InterfaceBuilder
	BindTaggedInterface(tag string, fromInterface ...interface{}) InterfaceBuilder
	BindTaggedConstant(tag string, constantKind ConstantKind) Builder
}

func CreateModule() Module { return createModule() }

type Builder interface {
	ToSingleton(singleton interface{})
	ToConstructor(constructor interface{})
	ToSingletonConstructor(constructor interface{})
	ToTaggedConstructor(constructor interface{})
	ToTaggedSingletonConstructor(constructor interface{})
}

type InterfaceBuilder interface {
	Builder
	To(to interface{})
}

type Injector interface {
	fmt.Stringer
	Get(from interface{}) (interface{}, error)
	GetTagged(tag string, from interface{}) (interface{}, error)
	Call(function interface{}) ([]interface{}, error)
	CallTagged(taggedFunction interface{}) ([]interface{}, error)
	Populate(populateStruct interface{}) error
}

func CreateInjector(modules ...Module) (Injector, error) { return createInjector(modules) }
