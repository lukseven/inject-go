[![Codeship Status](http://img.shields.io/codeship/34b974b0-6dfa-0132-51b4-66f2bf861e14/master.svg?style=flat-square)](https://codeship.com/projects/54288)
[![API Documentation](http://img.shields.io/badge/api-Godoc-blue.svg?style=flat-square)](https://godoc.org/github.com/peter-edge/inject)
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)](https://github.com/peter-edge/inject/blob/master/LICENSE)

### Installation
```bash
go get -u gopkg.in/peter-edge/inject.v1
```

### Import
```go
import (
    "gopkg.in/peter-edge/inject.v1"
)
```

# Godoc

Dependency injection for Go

## Usage

```go
const (
	InjectErrorTypeNil                                = "Parameter is nil"
	InjectErrorTypeReflectTypeNil                     = "reflect.TypeOf() returns nil"
	InjectErrorTypeNotInterfacePtr                    = "Binding with Binder.ToType() and from is not an interface pointer"
	InjectErrorTypeDoesNotImplement                   = "to binding does not implement from binding"
	InjectErrorTypeNotSupportedYet                    = "Binding type not supported yet, feel free to help!"
	InjectErrorTypeNotAssignable                      = "Binding not assignable"
	InjectErrorTypeConstructorNotFunction             = "Constructor is not a function"
	InjectErrorTypeConstructorReturnValuesInvalid     = "Constructor can only have two return values, the first providing the value, the second being an error"
	InjectErrorTypeIntermediateBinding                = "Trying to get for an intermediate binding"
	InjectErrorTypeFinalBinding                       = "Trying to get bindingKey for a final binding"
	InjectErrorTypeCannotCastModule                   = "Cannot cast Module to internal module type"
	InjectErrorTypeNoBinding                          = "No binding for binding key"
	InjectErrorTypeNoFinalBinding                     = "No final binding for binding key"
	InjectErrorTypeAlreadyBound                       = "Already found a binding for this binding key"
	InjectErrorTypeTagEmpty                           = "Tag empty"
	InjectErrorTypeTaggedConstructorParametersInvalid = "Tagged constructor must have one anonymous struct parameter"
)
```

#### type Builder

```go
type Builder interface {
	To(to interface{}) error
	ToSingleton(singleton interface{}) error
	ToConstructor(constructor interface{}) error
	ToSingletonConstructor(constructor interface{}) error
	ToTaggedConstructor(constructor interface{}) error
	ToTaggedSingletonConstructor(constructor interface{}) error
}
```


#### type InjectError

```go
type InjectError struct {
}
```


#### func (*InjectError) Error

```go
func (this *InjectError) Error() string
```

#### func (*InjectError) GetTag

```go
func (this *InjectError) GetTag(key string) (interface{}, bool)
```

#### func (*InjectError) Type

```go
func (this *InjectError) Type() string
```

#### type Injector

```go
type Injector interface {
	Get(from interface{}) (interface{}, error)
	GetTagged(from interface{}, tag string) (interface{}, error)
}
```


#### func  CreateInjector

```go
func CreateInjector(modules ...Module) (Injector, error)
```

#### type Module

```go
type Module interface {
	Bind(from interface{}) Builder
	BindTagged(from interface{}, tag string) Builder
}
```


#### func  CreateModule

```go
func CreateModule() Module
```
