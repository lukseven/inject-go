package inject

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
)

var (
	ErrNil                         = errors.New("inject: Parameter is nil")
	ErrReflectTypeNil              = errors.New("inject: reflect.TypeOf() returns nil")
	ErrUnknownBinderType           = errors.New("inject: Unknown binder type")
	ErrUnknownBindingType          = errors.New("inject: Unknown binding type")
	ErrNotInterfacePtr             = errors.New("inject: Binding with Binder.ToType() and from is not an interface pointer")
	ErrDoesNotImplement            = errors.New("inject: to binding does not implement from binding")
	ErrNotSupportedYet             = errors.New("inject.: Binding type not supported yet, feel free to help!")
	ErrNotAssignable               = errors.New("inject: Binding not assignable")
	ErrProviderNotFunction         = errors.New("inject: Provider is not a function")
	ErrProviderReturnValuesInvalid = errors.New("inject: Provider can only have two return values, the first providing the value, the second being an error")
	ErrInvalidReturnFromProvider   = errors.New("inject: Invalid return values from provider")
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
	binderTypeTo = iota
	binderTypeTaggedTo

	// TODO(pedge): rename to bindingTypeToType
	bindingTypeTo
	bindingTypeToSingleton
	bindingTypeToProvider
	bindingTypeToProviderAsSingleton

	noBindingMsg                      = "has no binding"
	noBindingToSingletonOrProviderMsg = "has no binding to a singleton or provider"
)

type boundType reflect.Type

type taggedBoundType struct {
	boundType
	tag interface{}
}

type binding struct {
	bindingType int

	// TODO(pedge): rename to toType
	to                    reflect.Type
	toSingleton           interface{}
	toProvider            interface{}
	toProviderAsSingleton interface{}

	// TODO(pedge): is atomic.Value the equivalent of a volatile variable in Java?
	providerSingletonValue  atomic.Value
	providerSingletonLoader sync.Once
}

type valueErr struct {
	value interface{}
	err   error
}

func createInjector() *injector {
	return &injector{
		make(map[boundType]*binding),
		make(map[taggedBoundType]*binding),
	}
}

type injector struct {
	boundTypeToBinding       map[boundType]*binding
	taggedBoundTypeToBinding map[taggedBoundType]*binding
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
	return &binder{binderTypeTaggedTo, this, fromReflectType, tag, nil}
}

func (this *injector) CreateContainer() (Container, error) {
	container := container{
		make(map[boundType]*binding),
		make(map[taggedBoundType]*binding),
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

func (this *injector) getFinalBinding(b *binding) (*binding, bool) {
	var ok bool
	for b.bindingType == bindingTypeTo {
		b, ok = this.boundTypeToBinding[b.to]
		if !ok {
			return nil, false
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

func (this *binder) ToType(to interface{}) error {
	if this.err != nil {
		return this.err
	}
	if to == nil {
		return ErrNil
	}
	toReflectType := reflect.TypeOf(to)
	// TODO(pedge): is this restriction necessary/warranted? how about structs with anonymous fields?
	if !(this.fromReflectType.Kind() == reflect.Ptr && this.fromReflectType.Elem().Kind() == reflect.Interface) {
		return ErrNotInterfacePtr
	}
	if !toReflectType.Implements(this.fromReflectType.Elem()) {
		return ErrDoesNotImplement
	}
	switch this.binderType {
	case binderTypeTo:
		this.injector.boundTypeToBinding[this.fromReflectType] = &binding{
			bindingType: bindingTypeTo,
			to:          toReflectType,
		}
		return nil
	case binderTypeTaggedTo:
		this.injector.taggedBoundTypeToBinding[taggedBoundType{this.fromReflectType, this.tag}] = &binding{
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
	if singleton == nil {
		return ErrNil
	}
	err := isValidBinding(this.fromReflectType, reflect.TypeOf(singleton))
	if err != nil {
		return err
	}
	switch this.binderType {
	case binderTypeTo:
		this.injector.boundTypeToBinding[this.fromReflectType] = &binding{
			bindingType: bindingTypeToSingleton,
			toSingleton: singleton,
		}
		return nil
	case binderTypeTaggedTo:
		this.injector.taggedBoundTypeToBinding[taggedBoundType{this.fromReflectType, this.tag}] = &binding{
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
	if provider == nil {
		return ErrNil
	}
	err := isValidProvider(this.fromReflectType, provider)
	if err != nil {
		return err
	}
	switch this.binderType {
	case binderTypeTo:
		this.injector.boundTypeToBinding[this.fromReflectType] = &binding{
			bindingType: bindingTypeToProvider,
			toProvider:  provider,
		}
		return nil
	case binderTypeTaggedTo:
		this.injector.taggedBoundTypeToBinding[taggedBoundType{this.fromReflectType, this.tag}] = &binding{
			bindingType: bindingTypeToProvider,
			toProvider:  provider,
		}
		return nil
	default:
		return ErrUnknownBinderType
	}
}

func (this *binder) ToProviderAsSingleton(provider interface{}) error {
	if this.err != nil {
		return this.err
	}
	if this.err != nil {
		return this.err
	}
	if provider == nil {
		return ErrNil
	}
	err := isValidProvider(this.fromReflectType, provider)
	if err != nil {
		return err
	}
	switch this.binderType {
	case binderTypeTo:
		this.injector.boundTypeToBinding[this.fromReflectType] = &binding{
			bindingType:             bindingTypeToProviderAsSingleton,
			toProviderAsSingleton:   provider,
			providerSingletonValue:  atomic.Value{},
			providerSingletonLoader: sync.Once{},
		}
		return nil
	case binderTypeTaggedTo:
		this.injector.taggedBoundTypeToBinding[taggedBoundType{this.fromReflectType, this.tag}] = &binding{
			bindingType:             bindingTypeToProviderAsSingleton,
			toProviderAsSingleton:   provider,
			providerSingletonValue:  atomic.Value{},
			providerSingletonLoader: sync.Once{},
		}
		return nil
	default:
		return ErrUnknownBinderType
	}
}

func isValidProvider(fromReflectType reflect.Type, provider interface{}) error {
	providerReflectType := reflect.TypeOf(provider)
	if providerReflectType.Kind() != reflect.Func {
		return ErrProviderNotFunction
	}
	if providerReflectType.NumOut() != 2 {
		return ErrProviderReturnValuesInvalid
	}
	err := isValidBinding(fromReflectType, providerReflectType.Out(0))
	if err != nil {
		return err
	}
	// TODO(pedge): can this be simplified?
	if !providerReflectType.Out(1).AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
		return ErrProviderReturnValuesInvalid
	}
	return nil
}

func isValidBinding(fromReflectType reflect.Type, toReflectType reflect.Type) error {
	switch {
	// from is an interface
	case fromReflectType.Kind() == reflect.Ptr && fromReflectType.Elem().Kind() == reflect.Interface:
		if !toReflectType.Implements(fromReflectType.Elem()) {
			return ErrDoesNotImplement
		}
	// from is a struct pointer
	case fromReflectType.Kind() == reflect.Ptr && fromReflectType.Elem().Kind() == reflect.Struct:
		fallthrough
	// from is a struct
	case fromReflectType.Kind() == reflect.Struct:
		// TODO(pedge): is this correct?
		if !toReflectType.AssignableTo(fromReflectType) {
			return ErrNotAssignable
		}
	// nothing else is supported for now
	// TODO(pedge): at least support primitives with tags
	default:
		return ErrNotSupportedYet
	}
	return nil
}

type container struct {
	boundTypeToBinding       map[boundType]*binding
	taggedBoundTypeToBinding map[taggedBoundType]*binding
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
	binding, ok := this.boundTypeToBinding[fromReflectType]
	if !ok {
		return nil, fmt.Errorf("inject: %v %v", fromReflectType, noBindingMsg)
	}
	return this.getFromBinding(binding)
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
	binding, ok := this.taggedBoundTypeToBinding[taggedBoundType]
	if !ok {
		return nil, fmt.Errorf("inject: %v with tag %v %v", fromReflectType, tag, noBindingMsg)
	}
	return this.getFromBinding(binding)
}

func (this *container) getFromBinding(binding *binding) (interface{}, error) {
	switch binding.bindingType {
	case bindingTypeToSingleton:
		return binding.toSingleton, nil
	case bindingTypeToProvider:
		return this.getFromProvider(binding.toProvider)
	case bindingTypeToProviderAsSingleton:
		binding.providerSingletonLoader.Do(func() {
			value, err := this.getFromProvider(binding.toProviderAsSingleton)
			binding.providerSingletonValue.Store(&valueErr{value, err})
		})
		valueErr := binding.providerSingletonValue.Load().(*valueErr)
		return valueErr.value, valueErr.err
	default:
		return nil, ErrNotSupportedYet
	}
}

// TODO(pedge): this is really hacky, and probably slow, clean this up
func (this *container) getFromProvider(provider interface{}) (interface{}, error) {
	// assuming this is a valid provider/that this is already checked
	providerReflectType := reflect.TypeOf(provider)
	numIn := providerReflectType.NumIn()
	parameterValues := make([]reflect.Value, numIn)
	if numIn == 1 && providerReflectType.In(0).AssignableTo(reflect.TypeOf((*Container)(nil)).Elem()) {
		parameterValues[0] = reflect.ValueOf(this)
	} else {
		for i := 0; i < numIn; i++ {
			inReflectType := providerReflectType.In(i)
			// TODO(pedge): this is really specific logic, and there wil need to be more
			// of this if more types are allowed for binding - this should be abstracted
			if inReflectType.Kind() == reflect.Interface {
				inReflectType = reflect.PtrTo(inReflectType)
			}
			parameter, err := this.get(inReflectType)
			if err != nil {
				return nil, err
			}
			parameterValues[i] = reflect.ValueOf(parameter)
		}
	}
	returnValues := reflect.ValueOf(provider).Call(parameterValues)
	return1 := returnValues[0].Interface()
	return2 := returnValues[1].Interface()
	switch {
	case return1 != nil && return2 != nil:
		return nil, ErrInvalidReturnFromProvider
	case return2 != nil:
		return nil, return2.(error)
	default:
		return return1, nil
	}
}
