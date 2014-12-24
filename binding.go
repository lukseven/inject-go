package inject

import (
	"reflect"
	"sync"
	"sync/atomic"
)

type binding interface {
	resolvedBinding(*module) (resolvedBinding, error)
}

type resolvedBinding interface {
	validate(*injector) error
	get(*injector) (interface{}, error)
}

type intermediateBinding struct {
	bindingKey bindingKey
}

func newIntermediateBinding(bindingKey bindingKey) binding {
	return &intermediateBinding{bindingKey}
}

func (this *intermediateBinding) resolvedBinding(module *module) (resolvedBinding, error) {
	binding, ok := module.binding(this.bindingKey)
	if !ok {
		eb := newErrorBuilder(InjectErrorTypeNoFinalBinding)
		eb.addTag("bindingKey", this.bindingKey)
		return nil, eb.build()
	}
	return binding.resolvedBinding(module)
}

type singletonBinding struct {
	singleton interface{}
}

func newSingletonBinding(singleton interface{}) binding {
	return &singletonBinding{singleton}
}

func (this *singletonBinding) validate(injector *injector) error {
	return nil
}

func (this *singletonBinding) get(injector *injector) (interface{}, error) {
	return this.singleton, nil
}

func (this *singletonBinding) resolvedBinding(module *module) (resolvedBinding, error) {
	return this, nil
}

type constructorBinding struct {
	constructor interface{}
}

func newConstructorBinding(constructor interface{}) binding {
	return &constructorBinding{constructor}
}

func (this *constructorBinding) validate(injector *injector) error {
	for _, inReflectType := range this.getInReflectTypes() {
		_, err := injector.getBinding(inReflectType)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *constructorBinding) get(injector *injector) (interface{}, error) {
	inReflectTypes := this.getInReflectTypes()
	numIn := len(inReflectTypes)
	parameterValues := make([]reflect.Value, numIn)
	for i := 0; i < numIn; i++ {
		parameter, err := injector.get(inReflectTypes[i])
		if err != nil {
			return nil, err
		}
		parameterValues[i] = reflect.ValueOf(parameter)
	}
	returnValues := reflect.ValueOf(this.constructor).Call(parameterValues)
	return1 := returnValues[0].Interface()
	return2 := returnValues[1].Interface()
	switch {
	case return2 != nil:
		return nil, return2.(error)
	default:
		return return1, nil
	}
}

func (this *constructorBinding) getInReflectTypes() []reflect.Type {
	constructorReflectType := reflect.TypeOf(this.constructor)
	numIn := constructorReflectType.NumIn()
	inReflectTypes := make([]reflect.Type, numIn)
	for i := 0; i < numIn; i++ {
		inReflectType := constructorReflectType.In(i)
		// TODO(pedge): this is really specific logic, and there wil need to be more
		// of this if more types are allowed for binding - this should be abstracted
		if inReflectType.Kind() == reflect.Interface {
			inReflectType = reflect.PtrTo(inReflectType)
		}
		inReflectTypes[i] = inReflectType
	}
	return inReflectTypes
}

func (this *constructorBinding) resolvedBinding(module *module) (resolvedBinding, error) {
	return this, nil
}

type singletonConstructorBinding struct {
	constructorBinding
	// TODO(pedge): is atomic.Value the equivalent of a volatile variable in Java?
	value atomic.Value
	once  sync.Once
}

func newSingletonConstructorBinding(constructor interface{}) binding {
	return &singletonConstructorBinding{constructorBinding{constructor}, atomic.Value{}, sync.Once{}}
}

func (this *singletonConstructorBinding) get(injector *injector) (interface{}, error) {
	this.once.Do(func() {
		value, err := this.constructorBinding.get(injector)
		this.value.Store(&valueErr{value, err})
	})
	valueErr := this.value.Load().(*valueErr)
	return valueErr.value, valueErr.err
}

func (this *singletonConstructorBinding) resolvedBinding(moduke *module) (resolvedBinding, error) {
	return this, nil
}

type valueErr struct {
	value interface{}
	err   error
}
