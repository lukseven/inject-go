package inject

import (
	"errors"
	"fmt"
	"reflect"
)

const (
	binderTypeTo = iota
	binderTypeTaggedTo

	bindingTypeTo
	bindingTypeToSingleton
	bindingTypeToProvider
	bindingTypeToProviderAsSingleton

	noBindingMsg = "has no binding to a singleton or provider"
)

var (
	ErrNil                = errors.New("inject: Parameter is nil")
	ErrReflectTypeNil     = errors.New("inject: reflect.TypeOf() returns nil")
	ErrTagValueInvalid    = errors.New("inject: Tag value is invalid")
	ErrUnknownBinderType  = errors.New("inject: Unknown binder type")
	ErrUnknownBindingType = errors.New("inject: Unknown binder type")
	ErrNotInterfacePtr    = errors.New("inject: Binding with Binder.To() and from is not an interface pointer")
	ErrDoesNotImplement   = errors.New("inject: to binding does not implement from binding")
	ErrNotSupportedYet    = errors.New("inject.: Binding type not supported yet, feel free to help!")
	ErrNotAssignable      = errors.New("inject: Binding not assignable")
)

func CreateInjector() Injector {
	return &injector{
		make(map[boundType]binding),
		make(map[taggedBoundType]binding),
	}
}

type Injector interface {
	Bind(from interface{}) Binder
	BindTagged(from interface{}, tag interface{}) Binder

	CreateContainer() (Container, error)
}

type Binder interface {
	To(to interface{}) error
	ToSingleton(singleton interface{}) error
	ToProvider(provider interface{}) error
	ToProviderAsSingleton(provider interface{}) error
}

type Container interface {
	Get(from interface{}) (interface{}, error)
	GetTagged(from interface{}, tag interface{}) (interface{}, error)
}

// private

type boundType reflect.Type

type taggedBoundType struct {
	boundType
	tag interface{}
}

type binding struct {
	bindingType int

	to                    reflect.Type
	toSingleton           interface{}
	toProvider            interface{}
	toProviderAsSingleton interface{}
}

type injector struct {
	boundTypeToBinding       map[boundType]binding
	taggedBoundTypeToBinding map[taggedBoundType]binding
}

func (this *injector) Bind(from interface{}) Binder {
	if from == nil {
		return &binder{binderTypeTo, nil, nil, nil, ErrNil}
	}
	fromReflectType := reflect.TypeOf(from)
	if fromReflectType == nil {
		return &binder{binderTypeTo, nil, nil, nil, ErrReflectTypeNil}
	}
	return &binder{binderTypeTo, this, fromReflectType, nil, nil}
}

func (this *injector) BindTagged(from interface{}, tag interface{}) Binder {
	if from == nil {
		return &binder{binderTypeTaggedTo, nil, nil, nil, ErrNil}
	}
	if tag == nil {
		return &binder{binderTypeTaggedTo, nil, nil, nil, ErrNil}
	}
	fromReflectType := reflect.TypeOf(from)
	if fromReflectType == nil {
		return &binder{binderTypeTaggedTo, nil, nil, nil, ErrReflectTypeNil}
	}
	tagReflectValue := reflect.ValueOf(tag)
	if !tagReflectValue.IsValid() {
		return &binder{binderTypeTaggedTo, nil, nil, nil, ErrTagValueInvalid}
	}
	return &binder{binderTypeTaggedTo, this, fromReflectType, tag, nil}
}

func (this *injector) CreateContainer() (Container, error) {
	container := container{
		make(map[boundType]binding),
		make(map[taggedBoundType]binding),
	}
	for taggedBoundType, binding := range this.taggedBoundTypeToBinding {
		finalBinding, ok := this.getFinalBinding(binding)
		if !ok {
			return nil, fmt.Errorf("inject: %v %v", taggedBoundType, noBindingMsg)
		}
		container.taggedBoundTypeToBinding[taggedBoundType] = finalBinding
	}
	for boundType, binding := range this.boundTypeToBinding {
		finalBinding, ok := this.getFinalBinding(binding)
		if !ok {
			return nil, fmt.Errorf("inject: %v %v", boundType, noBindingMsg)
		}
		container.boundTypeToBinding[boundType] = finalBinding
	}
	return &container, nil
}

func (this *injector) getFinalBinding(b binding) (binding, bool) {
	var ok bool
	for b.bindingType == bindingTypeTo {
		b, ok = this.boundTypeToBinding[b.to]
		if !ok {
			return binding{}, false
		}
	}
	return b, true
}

type binder struct {
	binderType int

	injector        *injector
	fromReflectType reflect.Type
	tag             interface{}
	err             error
}

func (this *binder) To(to interface{}) error {
	if this.err != nil {
		return this.err
	}
	toReflectType := reflect.TypeOf(to)
	if !(this.fromReflectType.Kind() == reflect.Ptr && this.fromReflectType.Elem().Kind() == reflect.Interface) {
		return ErrNotInterfacePtr
	}
	if !toReflectType.Implements(this.fromReflectType.Elem()) {
		return ErrDoesNotImplement
	}
	switch this.binderType {
	case binderTypeTo:
		this.injector.boundTypeToBinding[this.fromReflectType] = binding{
			bindingType: bindingTypeTo,
			to:          toReflectType,
		}
		return nil
	case binderTypeTaggedTo:
		this.injector.taggedBoundTypeToBinding[taggedBoundType{this.fromReflectType, this.tag}] = binding{
			bindingType: bindingTypeTo,
			to:          toReflectType,
		}
		return nil
	default:
		return ErrUnknownBinderType
	}
}

func (this *binder) ToSingleton(singleton interface{}) error {
	if this.err != nil {
		return this.err
	}
	singletonReflectType := reflect.TypeOf(singleton)
	switch {
	// from is an interface
	case this.fromReflectType.Kind() == reflect.Ptr && this.fromReflectType.Elem().Kind() == reflect.Interface:
		if !singletonReflectType.Implements(this.fromReflectType.Elem()) {
			return ErrDoesNotImplement
		}
	// from is a struct pointer
	case this.fromReflectType.Kind() == reflect.Ptr && this.fromReflectType.Elem().Kind() == reflect.Struct:
		// TODO(pedge): is this correct?
		if !singletonReflectType.AssignableTo(this.fromReflectType) {
			return ErrNotAssignable
		}
	// from is a struct
	case this.fromReflectType.Kind() == reflect.Struct:
		// TODO(pedge): is this correct?
		if !singletonReflectType.AssignableTo(this.fromReflectType) {
			return ErrNotAssignable
		}
	// nothing else is supported for now
	// TODO(pedge): at least support primitives with tags
	default:
		return ErrNotSupportedYet
	}
	switch this.binderType {
	case binderTypeTo:
		this.injector.boundTypeToBinding[this.fromReflectType] = binding{
			bindingType: bindingTypeToSingleton,
			toSingleton: singleton,
		}
		return nil
	case binderTypeTaggedTo:
		this.injector.taggedBoundTypeToBinding[taggedBoundType{this.fromReflectType, this.tag}] = binding{
			bindingType: bindingTypeToSingleton,
			toSingleton: singleton,
		}
		return nil
	default:
		return ErrUnknownBinderType
	}
}

func (this *binder) ToProvider(provider interface{}) error {
	if this.err != nil {
		return this.err
	}
	//providerReflectType := reflect.TypeOf(provider)
	return nil
}

func (this *binder) ToProviderAsSingleton(provider interface{}) error {
	if this.err != nil {
		return this.err
	}
	//providerReflectType := reflect.TypeOf(provider)
	return nil
}

type container struct {
	boundTypeToBinding       map[boundType]binding
	taggedBoundTypeToBinding map[taggedBoundType]binding
}

func (this *container) Get(from interface{}) (interface{}, error) {
	if from == nil {
		return nil, ErrNil
	}
	fromReflectType := reflect.TypeOf(from)
	if fromReflectType == nil {
		return nil, ErrReflectTypeNil
	}
	binding, ok := this.boundTypeToBinding[fromReflectType]
	if !ok {
		return nil, fmt.Errorf("inject: No binding for %v", fromReflectType)
	}
	switch binding.bindingType {
	case bindingTypeToSingleton:
		return binding.toSingleton, nil
	default:
		return nil, ErrNotSupportedYet
	}
}

func (this *container) GetTagged(from interface{}, tag interface{}) (interface{}, error) {
	if from == nil {
		return nil, ErrNil
	}
	if tag == nil {
		return nil, ErrNil
	}
	fromReflectType := reflect.TypeOf(from)
	if fromReflectType == nil {
		return nil, ErrReflectTypeNil
	}
	tagReflectValue := reflect.ValueOf(tag)
	if !tagReflectValue.IsValid() {
		return nil, ErrTagValueInvalid
	}
	taggedBoundType := taggedBoundType{fromReflectType, tag}
	binding, ok := this.taggedBoundTypeToBinding[taggedBoundType]
	if !ok {
		return nil, fmt.Errorf("inject: No binding for %v with tag %v", fromReflectType, tag)
	}
	switch binding.bindingType {
	case bindingTypeToSingleton:
		return binding.toSingleton, nil
	default:
		return nil, ErrNotSupportedYet
	}
}
