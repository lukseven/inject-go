package inject // import "go.pedge.io/inject"

import (
	"fmt"
)

type Module interface {
	fmt.Stringer
	Bind(from ...interface{}) Builder
	BindTagged(tag string, from ...interface{}) Builder
	BindInterface(fromInterface ...interface{}) InterfaceBuilder
	BindTaggedInterface(tag string, fromInterface ...interface{}) InterfaceBuilder
	BindTaggedBool(tag string) Builder
	BindTaggedInt(tag string) Builder
	BindTaggedInt8(tag string) Builder
	BindTaggedInt16(tag string) Builder
	BindTaggedInt32(tag string) Builder
	BindTaggedInt64(tag string) Builder
	BindTaggedUint(tag string) Builder
	BindTaggedUint8(tag string) Builder
	BindTaggedUint16(tag string) Builder
	BindTaggedUint32(tag string) Builder
	BindTaggedUint64(tag string) Builder
	BindTaggedFloat32(tag string) Builder
	BindTaggedFloat64(tag string) Builder
	BindTaggedComplex64(tag string) Builder
	BindTaggedComplex128(tag string) Builder
	BindTaggedString(tag string) Builder
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
	GetTaggedBool(tag string) (bool, error)
	GetTaggedInt(tag string) (int, error)
	GetTaggedInt8(tag string) (int8, error)
	GetTaggedInt16(tag string) (int16, error)
	GetTaggedInt32(tag string) (int32, error)
	GetTaggedInt64(tag string) (int64, error)
	GetTaggedUint(tag string) (uint, error)
	GetTaggedUint8(tag string) (uint8, error)
	GetTaggedUint16(tag string) (uint16, error)
	GetTaggedUint32(tag string) (uint32, error)
	GetTaggedUint64(tag string) (uint64, error)
	GetTaggedFloat32(tag string) (float32, error)
	GetTaggedFloat64(tag string) (float64, error)
	GetTaggedComplex64(tag string) (complex64, error)
	GetTaggedComplex128(tag string) (complex128, error)
	GetTaggedString(tag string) (string, error)
	Call(function interface{}) ([]interface{}, error)
	CallTagged(taggedFunction interface{}) ([]interface{}, error)
	Populate(populateStruct interface{}) error
}

func CreateInjector(modules ...Module) (Injector, error) { return createInjector(modules) }
