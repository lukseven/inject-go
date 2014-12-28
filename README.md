[![Codeship Status](http://img.shields.io/codeship/34b974b0-6dfa-0132-51b4-66f2bf861e14/master.svg?style=flat-square)](https://codeship.com/projects/54288)
[![API Documentation](http://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](https://godoc.org/github.com/peter-edge/inject)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/peter-edge/inject/blob/master/LICENSE)

Guice-inspired dependency injection for Go

## Installation
```bash
go get -u gopkg.in/peter-edge/inject.v1
```

## Import
```go
import (
    "gopkg.in/peter-edge/inject.v1"
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


#### type ConstantKind

```go
type ConstantKind uint
```


```go
const (
	Bool ConstantKind = iota
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	String
)
```

#### func  ConstantKindOf

```go
func ConstantKindOf(reflectKind reflect.Kind) ConstantKind
```

#### func (ConstantKind) ReflectKind

```go
func (this ConstantKind) ReflectKind() reflect.Kind
```

#### type Injector

```go
type Injector interface {
	fmt.Stringer
	Get(from interface{}) (interface{}, error)
	GetTagged(tag string, from interface{}) (interface{}, error)
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
}
```


#### func  CreateModule

```go
func CreateModule() Module
```
