[![Codeship Status](http://img.shields.io/codeship/34b974b0-6dfa-0132-51b4-66f2bf861e14/master.svg?style=flat-square)](https://codeship.com/projects/54288)
[![API Documentation](http://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](https://godoc.org/github.com/peter-edge/inject)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/peter-edge/inject/blob/master/LICENSE)

Guice-inspired dependency injection for Go

## Installation
```bash
go get -u github.com/peter-edge/inject
```

## Import
```go
import (
    "github.com/peter-edge/inject"
)
```

## Usage

#### type Builder

```go
type Builder interface {
	ToSingleton(singleton interface{})
	ToConstructor(constructor interface{})
	ToSingletonConstructor(constructor interface{})
	ToTaggedConstructor(constructor interface{})
	ToTaggedSingletonConstructor(constructor interface{})
}
```


#### type Injector

```go
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
```


#### func  CreateInjector

```go
func CreateInjector(modules ...Module) (Injector, error)
```

#### type InterfaceBuilder

```go
type InterfaceBuilder interface {
	Builder
	To(to interface{})
}
```


#### type Module

```go
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
```


#### func  CreateModule

```go
func CreateModule() Module
```
