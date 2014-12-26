[![Codeship Status](http://img.shields.io/codeship/34b974b0-6dfa-0132-51b4-66f2bf861e14/master.svg?style=flat-square)](https://codeship.com/projects/54288)
[![API Documentation](http://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](https://godoc.org/github.com/peter-edge/inject)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/peter-edge/inject/blob/master/LICENSE)

Dependency injection for Go

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
	To(to interface{})
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
	GetTagged(from interface{}, tag string) (interface{}, error)
	Call(function interface{}) ([]interface{}, error)
	CallTagged(taggedFunction interface{}) ([]interface{}, error)
}
```


#### func  CreateInjector

```go
func CreateInjector(modules ...Module) (Injector, error)
```

#### type Module

```go
type Module interface {
	fmt.Stringer
	Bind(from interface{}) Builder
	BindTagged(from interface{}, tag string) Builder
}
```


#### func  CreateModule

```go
func CreateModule() Module
```
