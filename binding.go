package inject

import (
	"fmt"
	"reflect"
)

const (
	taggedConstructorStructFieldTag = "inject"
)

type binding interface {
	fmt.Stringer
	// has to be a copy constructor
	// https://github.com/peter-edge/inject/commit/e525825afc80f0de819f35a6afc26a4bf3d3a192
	// this could be designed better
	resolvedBinding(*module, *injector) (resolvedBinding, error)
}

type resolvedBinding interface {
	fmt.Stringer
	validate() error
	get() (interface{}, error)
}

type intermediateBinding struct {
	bindingKey bindingKey
}

func newIntermediateBinding(to interface{}) binding {
	return &intermediateBinding{newBindingKey(reflect.TypeOf(to))}
}

func (this *intermediateBinding) String() string {
	return this.bindingKey.String()
}

func (this *intermediateBinding) resolvedBinding(module *module, injector *injector) (resolvedBinding, error) {
	binding, ok := module.binding(this.bindingKey)
	if !ok {
		eb := newErrorBuilder(injectErrorTypeNoFinalBinding)
		eb.addTag("bindingKey", this.bindingKey)
		return nil, eb.build()
	}
	return binding.resolvedBinding(module, injector)
}

type singletonBinding struct {
	singleton interface{}
	injector  *injector
}

func newSingletonBinding(singleton interface{}) binding {
	return &singletonBinding{singleton, nil}
}

func (this *singletonBinding) String() string {
	return fmt.Sprintf("%v", this.singleton)
}

func (this *singletonBinding) validate() error {
	return nil
}

func (this *singletonBinding) get() (interface{}, error) {
	return this.singleton, nil
}

func (this *singletonBinding) resolvedBinding(module *module, injector *injector) (resolvedBinding, error) {
	return &singletonBinding{this.singleton, injector}, nil
}

type constructorBinding struct {
	constructor interface{}
	cache       *constructorBindingCache
	injector    *injector
}

type constructorBindingCache struct {
	numIn       int
	bindingKeys []bindingKey
}

func newConstructorBinding(constructor interface{}) binding {
	return &constructorBinding{constructor, newConstructorBindingCache(constructor), nil}
}

func newConstructorBindingCache(constructor interface{}) *constructorBindingCache {
	constructorReflectType := reflect.TypeOf(constructor)
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
	return &constructorBindingCache{
		numIn,
		bindingKeys,
	}
}

func (this *constructorBinding) String() string {
	return fmt.Sprintf("%v", this.constructor)
}

func (this *constructorBinding) validate() error {
	return validateBindingKeys(this.cache.bindingKeys, this.injector)
}

func (this *constructorBinding) get() (interface{}, error) {
	reflectValues, err := getReflectValues(this.injector, this.cache.bindingKeys)
	if err != nil {
		return nil, err
	}
	return callConstructor(this.constructor, reflectValues)
}

func (this *constructorBinding) resolvedBinding(module *module, injector *injector) (resolvedBinding, error) {
	return &constructorBinding{this.constructor, this.cache, injector}, nil
}

func (this *singletonConstructorBinding) String() string {
	return fmt.Sprintf("%v", this.constructor)
}

type singletonConstructorBinding struct {
	constructorBinding
	loader *loader
}

func newSingletonConstructorBinding(constructor interface{}) binding {
	return &singletonConstructorBinding{constructorBinding{constructor, newConstructorBindingCache(constructor), nil}, nil}
}

func (this *singletonConstructorBinding) get() (interface{}, error) {
	return this.loader.load(this.constructorBinding.get)
}

func (this *singletonConstructorBinding) resolvedBinding(module *module, injector *injector) (resolvedBinding, error) {
	return &singletonConstructorBinding{constructorBinding{this.constructorBinding.constructor, this.constructorBinding.cache, injector}, newLoader()}, nil
}

type taggedConstructorBinding struct {
	constructor interface{}
	cache       *taggedConstructorBindingCache
	injector    *injector
}

type taggedConstructorBindingCache struct {
	inReflectType reflect.Type
	numFields     int
	bindingKeys   []bindingKey
}

func newTaggedConstructorBinding(constructor interface{}) binding {
	return &taggedConstructorBinding{constructor, newTaggedConstructorBindingCache(constructor), nil}
}

func newTaggedConstructorBindingCache(constructor interface{}) *taggedConstructorBindingCache {
	constructorReflectType := reflect.TypeOf(constructor)
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
	return &taggedConstructorBindingCache{
		inReflectType,
		numFields,
		bindingKeys,
	}
}

func (this *taggedConstructorBinding) String() string {
	return fmt.Sprintf("%v", this.constructor)
}

func (this *taggedConstructorBinding) validate() error {
	return validateBindingKeys(this.cache.bindingKeys, this.injector)
}

func (this *taggedConstructorBinding) get() (interface{}, error) {
	reflectValues, err := getReflectValues(this.injector, this.cache.bindingKeys)
	if err != nil {
		return nil, err
	}
	numReflectValues := len(reflectValues)
	value := reflect.Indirect(reflect.New(this.cache.inReflectType))
	for i := 0; i < numReflectValues; i++ {
		value.Field(i).Set(reflectValues[i])
	}
	return callConstructor(this.constructor, []reflect.Value{value})
}

func (this *taggedConstructorBinding) resolvedBinding(module *module, injector *injector) (resolvedBinding, error) {
	return &taggedConstructorBinding{this.constructor, this.cache, injector}, nil
}

type taggedSingletonConstructorBinding struct {
	taggedConstructorBinding
	loader *loader
}

func newTaggedSingletonConstructorBinding(constructor interface{}) binding {
	return &taggedSingletonConstructorBinding{taggedConstructorBinding{constructor, newTaggedConstructorBindingCache(constructor), nil}, nil}
}

func (this *taggedSingletonConstructorBinding) String() string {
	return fmt.Sprintf("%v", this.constructor)
}

func (this *taggedSingletonConstructorBinding) get() (interface{}, error) {
	return this.loader.load(this.taggedConstructorBinding.get)
}

func (this *taggedSingletonConstructorBinding) resolvedBinding(module *module, injector *injector) (resolvedBinding, error) {
	return &taggedSingletonConstructorBinding{taggedConstructorBinding{this.taggedConstructorBinding.constructor, this.taggedConstructorBinding.cache, injector}, newLoader()}, nil
}

func getReflectValues(injector *injector, bindingKeys []bindingKey) ([]reflect.Value, error) {
	numBindingKeys := len(bindingKeys)
	reflectValues := make([]reflect.Value, numBindingKeys)
	for i := 0; i < numBindingKeys; i++ {
		value, err := injector.get(bindingKeys[i])
		if err != nil {
			return nil, err
		}
		reflectValues[i] = reflect.ValueOf(value)
	}
	return reflectValues, nil
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

func validateBindingKeys(bindingKeys []bindingKey, injector *injector) error {
	for _, bindingKey := range bindingKeys {
		_, err := injector.getBinding(bindingKey)
		if err != nil {
			return err
		}
	}
	return nil
}
