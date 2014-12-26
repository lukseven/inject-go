package inject

import (
	"fmt"
	"reflect"
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
	bindingKeys := getParameterBindingKeysForFunc(reflect.TypeOf(constructor))
	return &constructorBindingCache{len(bindingKeys), bindingKeys}
}

func (this *constructorBinding) String() string {
	return fmt.Sprintf("%v", this.constructor)
}

func (this *constructorBinding) validate() error {
	return this.injector.validateBindingKeys(this.cache.bindingKeys)
}

func (this *constructorBinding) get() (interface{}, error) {
	reflectValues, err := this.injector.getReflectValues(this.cache.bindingKeys)
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
	bindingKeys := getParameterBindingKeysForTaggedFunc(constructorReflectType)
	return &taggedConstructorBindingCache{constructorReflectType.In(0), len(bindingKeys), bindingKeys}
}

func (this *taggedConstructorBinding) String() string {
	return fmt.Sprintf("%v", this.constructor)
}

func (this *taggedConstructorBinding) validate() error {
	return this.injector.validateBindingKeys(this.cache.bindingKeys)
}

func (this *taggedConstructorBinding) get() (interface{}, error) {
	reflectValues, err := this.injector.getReflectValues(this.cache.bindingKeys)
	if err != nil {
		return nil, err
	}
	return callConstructor(this.constructor, []reflect.Value{*populateTaggedFuncStruct(this.cache.inReflectType, reflectValues)})
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
