package inject

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNil                         = errors.New("inject: Parameter is nil")
	ErrReflectTypeNil              = errors.New("inject: reflect.TypeOf() returns nil")
	ErrUnknownBinderType           = errors.New("inject: Unknown binder type")
	ErrUnknownValidBinderType      = errors.New("inject: Unknown valid binder type")
	ErrNotInterfacePtr             = errors.New("inject: Binding with Binder.ToType() and from is not an interface pointer")
	ErrDoesNotImplement            = errors.New("inject: to binding does not implement from binding")
	ErrNotSupportedYet             = errors.New("inject.: Binding type not supported yet, feel free to help!")
	ErrNotAssignable               = errors.New("inject: Binding not assignable")
	ErrProviderNotFunction         = errors.New("inject: Provider is not a function")
	ErrProviderReturnValuesInvalid = errors.New("inject: Provider can only have two return values, the first providing the value, the second being an error")
	ErrInvalidReturnFromProvider   = errors.New("inject: Invalid return values from provider")
	ErrBindingTypeIntermediate     = errors.New("inject: Binding type is intermediate")
	ErrBindingTypeUnknown          = errors.New("inject: Binding type is unknown")
	ErrFinalBindingTypeUnknown     = errors.New("inject: Final binding type is unknown")
)

func CreateInjector() Injector { return createInjector() }

type Injector interface {
	Bind(from interface{}) Binder
	BindTagged(from interface{}, tag interface{}) Binder
	CreateContainer() (Container, error)
}

type Binder interface {
	ToType(to interface{}) error
	ToSingleton(singleton interface{}) error
	ToProvider(provider interface{}) error
	ToProviderAsSingleton(provider interface{}) error
}

type Container interface {
	Get(from interface{}) (interface{}, error)
	GetTagged(from interface{}, tag interface{}) (interface{}, error)
}

// private

const (
	noBindingMsg                      = "has no binding"
	noBindingToSingletonOrProviderMsg = "has no binding to a singleton or provider"
)

type taggedBoundType struct {
	boundType reflect.Type
	tag       interface{}
}

type injector struct {
	boundTypeToBinding       map[reflect.Type]*binding
	taggedBoundTypeToBinding map[taggedBoundType]*binding
}

func createInjector() *injector {
	return &injector{
		make(map[reflect.Type]*binding),
		make(map[taggedBoundType]*binding),
	}
}

func (this *injector) Bind(from interface{}) Binder {
	if from == nil {
		return newBinderErr(ErrNil)
	}
	fromReflectType := reflect.TypeOf(from)
	if fromReflectType == nil {
		return newBinderErr(ErrReflectTypeNil)
	}
	return newBinderBoundType(this, fromReflectType)
}

func (this *injector) BindTagged(from interface{}, tag interface{}) Binder {
	if from == nil {
		return newBinderErr(ErrNil)
	}
	if tag == nil {
		return newBinderErr(ErrNil)
	}
	fromReflectType := reflect.TypeOf(from)
	if fromReflectType == nil {
		return newBinderErr(ErrReflectTypeNil)
	}
	return newBinderTaggedBoundType(this, taggedBoundType{fromReflectType, tag})
}

func (this *injector) CreateContainer() (Container, error) {
	container := container{
		make(map[reflect.Type]*finalBinding),
		make(map[taggedBoundType]*finalBinding),
	}
	for taggedBoundType, binding := range this.taggedBoundTypeToBinding {
		finalBinding, ok := this.getFinalBinding(binding)
		if !ok {
			return nil, fmt.Errorf("inject: %v %v", taggedBoundType, noBindingToSingletonOrProviderMsg)
		}
		container.taggedBoundTypeToBinding[taggedBoundType] = finalBinding
	}
	for boundType, binding := range this.boundTypeToBinding {
		finalBinding, ok := this.getFinalBinding(binding)
		if !ok {
			return nil, fmt.Errorf("inject: %v %v", boundType, noBindingToSingletonOrProviderMsg)
		}
		container.boundTypeToBinding[boundType] = finalBinding
	}
	return &container, nil
}

func (this *injector) getFinalBinding(b *binding) (*finalBinding, bool) {
	var ok bool
	for b.bindingType == bindingTypeIntermediate {
		b, ok = this.boundTypeToBinding[b.intermediateBinding]
		if !ok {
			return nil, false
		}
	}
	return b.finalBinding, true
}

type container struct {
	boundTypeToBinding       map[reflect.Type]*finalBinding
	taggedBoundTypeToBinding map[taggedBoundType]*finalBinding
}

func (this *container) Get(from interface{}) (interface{}, error) {
	if from == nil {
		return nil, ErrNil
	}
	return this.get(reflect.TypeOf(from))
}

func (this *container) get(fromReflectType reflect.Type) (interface{}, error) {
	if fromReflectType == nil {
		return nil, ErrReflectTypeNil
	}
	finalBinding, ok := this.boundTypeToBinding[fromReflectType]
	if !ok {
		return nil, fmt.Errorf("inject: %v %v", fromReflectType, noBindingMsg)
	}
	return finalBinding.get(this)
}

func (this *container) GetTagged(from interface{}, tag interface{}) (interface{}, error) {
	if from == nil {
		return nil, ErrNil
	}
	return this.getTagged(reflect.TypeOf(from), tag)
}

func (this *container) getTagged(fromReflectType reflect.Type, tag interface{}) (interface{}, error) {
	if tag == nil {
		return nil, ErrNil
	}
	if fromReflectType == nil {
		return nil, ErrReflectTypeNil
	}
	taggedBoundType := taggedBoundType{fromReflectType, tag}
	finalBinding, ok := this.taggedBoundTypeToBinding[taggedBoundType]
	if !ok {
		return nil, fmt.Errorf("inject: %v with tag %v %v", fromReflectType, tag, noBindingMsg)
	}
	return finalBinding.get(this)
}
