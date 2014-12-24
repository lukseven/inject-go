package inject

import (
	"reflect"
	"sync"
	"sync/atomic"
)

const (
	taggedConstructorStructFieldTag = "inject"
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
	for _, bindingKey := range this.getBindingKeys() {
		_, err := injector.getBinding(bindingKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *constructorBinding) get(injector *injector) (interface{}, error) {
	bindingKeys := this.getBindingKeys()
	numIn := len(bindingKeys)
	parameterValues := make([]reflect.Value, numIn)
	for i := 0; i < numIn; i++ {
		parameter, err := injector.get(bindingKeys[i])
		if err != nil {
			return nil, err
		}
		parameterValues[i] = reflect.ValueOf(parameter)
	}
	return callConstructor(this.constructor, parameterValues)
}

func (this *constructorBinding) resolvedBinding(module *module) (resolvedBinding, error) {
	return this, nil
}

func (this *constructorBinding) getBindingKeys() []bindingKey {
	constructorReflectType := reflect.TypeOf(this.constructor)
	numIn := constructorReflectType.NumIn()
	bindingKeys := make([]bindingKey, numIn)
	for i := 0; i < numIn; i++ {
		inReflectType := constructorReflectType.In(i)
		// TODO(pedge): this is really specific logic, and there wil need to be more
		// of this if more types are allowed for binding - this should be abstracted
		if inReflectType.Kind() == reflect.Interface {
			inReflectType = reflect.PtrTo(inReflectType)
		}
		bindingKeys[i] = newBindingKey(inReflectType)
	}
	return bindingKeys
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

func (this *singletonConstructorBinding) resolvedBinding(module *module) (resolvedBinding, error) {
	return this, nil
}

type taggedConstructorBinding struct {
	constructor interface{}
}

func newTaggedConstructorBinding(constructor interface{}) binding {
	return &taggedConstructorBinding{constructor}
}

func (this *taggedConstructorBinding) validate(injector *injector) error {
	for _, bindingKey := range this.getBindingKeys() {
		_, err := injector.getBinding(bindingKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *taggedConstructorBinding) get(injector *injector) (interface{}, error) {
	bindingKeys := this.getBindingKeys()
	constructorReflectType := reflect.TypeOf(this.constructor)
	inReflectType := constructorReflectType.In(0)
	numFields := inReflectType.NumField()
	value := reflect.Zero(inReflectType)
	for i := 0; i < numFields; i++ {
		field, err := injector.get(bindingKeys[i])
		if err != nil {
			return nil, err
		}
		value.Field(i).Set(reflect.ValueOf(field))
	}
	return callConstructor(this.constructor, []reflect.Value{value})
}

func (this *taggedConstructorBinding) resolvedBinding(module *module) (resolvedBinding, error) {
	return this, nil
}

func (this *taggedConstructorBinding) getBindingKeys() []bindingKey {
	constructorReflectType := reflect.TypeOf(this.constructor)
	inReflectType := constructorReflectType.In(0)
	numFields := inReflectType.NumField()
	bindingKeys := make([]bindingKey, numFields)
	for i := 0; i < numFields; i++ {
		structField := inReflectType.Field(i)
		structFieldReflectType := structField.Type
		// TODO(pedge): this is really specific logic, and there wil need to be more
		// of this if more types are allowed for binding - this should be abstracted
		if structFieldReflectType.Kind() == reflect.Interface {
			structFieldReflectType = reflect.PtrTo(structFieldReflectType)
		}
		tag := structField.Tag.Get(taggedConstructorStructFieldTag)
		if tag != "" {
			bindingKeys[i] = newTaggedBindingKey(structFieldReflectType, tag)
		} else {
			bindingKeys[i] = newBindingKey(structFieldReflectType)
		}
	}
	return bindingKeys
}

type taggedSingletonConstructorBinding struct {
	taggedConstructorBinding
	// TODO(pedge): is atomic.Value the equivalent of a volatile variable in Java?
	value atomic.Value
	once  sync.Once
}

func newTaggedSingletonConstructorBinding(constructor interface{}) binding {
	return &taggedSingletonConstructorBinding{taggedConstructorBinding{constructor}, atomic.Value{}, sync.Once{}}
}

func (this *taggedSingletonConstructorBinding) get(injector *injector) (interface{}, error) {
	this.once.Do(func() {
		value, err := this.taggedConstructorBinding.get(injector)
		this.value.Store(&valueErr{value, err})
	})
	valueErr := this.value.Load().(*valueErr)
	return valueErr.value, valueErr.err
}

func (this *taggedSingletonConstructorBinding) resolvedBinding(module *module) (resolvedBinding, error) {
	return this, nil
}

func callConstructor(constructor interface{}, reflectValues []reflect.Value) (interface{}, error) {
	returnValues := reflect.ValueOf(constructor).Call(reflectValues)
	return1 := returnValues[0].Interface()
	return2 := returnValues[1].Interface()
	switch {
	case return2 != nil:
		return nil, return2.(error)
	default:
		return return1, nil
	}
}

type valueErr struct {
	value interface{}
	err   error
}
