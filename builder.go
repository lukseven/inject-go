package inject

import (
	"reflect"
)

var (
	errorReflectType = reflect.TypeOf((*error)(nil)).Elem()
)

type noOpBuilder struct{}

func newNoOpBuilder() InterfaceBuilder {
	return &noOpBuilder{}
}

func (n *noOpBuilder) To(to interface{}) {}

func (n *noOpBuilder) ToSingleton(singleton interface{}) {}

func (n *noOpBuilder) ToConstructor(constructor interface{}) {}

func (n *noOpBuilder) ToSingletonConstructor(construtor interface{}) SingletonBuilder {
	return nil
}

func (n *noOpBuilder) ToTaggedConstructor(constructor interface{}) {}

func (n *noOpBuilder) ToTaggedSingletonConstructor(constructor interface{}) SingletonBuilder {
	return nil
}

type baseBuilder struct {
	module      *module
	bindingKeys []bindingKey
}

func newBuilder(module *module, bindingKeys []bindingKey) InterfaceBuilder {
	return &baseBuilder{module, bindingKeys}
}

func (b *baseBuilder) To(to interface{}) {
	b.to(to, verifyBindingReflectType, newIntermediateBinding)
}

func (b *baseBuilder) ToSingleton(singleton interface{}) {
	b.to(singleton, verifyBindingReflectType, newSingletonBinding)
}

func (b *baseBuilder) ToConstructor(constructor interface{}) {
	b.to(constructor, verifyConstructorReflectType, newConstructorBinding)
}

func (b *baseBuilder) ToSingletonConstructor(constructor interface{}) SingletonBuilder {
	b.to(constructor, verifyConstructorReflectType, newSingletonConstructorBinding)
	return newSingletonBuilder(b.module, b.bindingKeys[0].reflectType())
}

func (b *baseBuilder) ToTaggedConstructor(constructor interface{}) {
	b.to(constructor, verifyTaggedConstructorReflectType, newTaggedConstructorBinding)
}

func (b *baseBuilder) ToTaggedSingletonConstructor(constructor interface{}) SingletonBuilder {
	b.to(constructor, verifyTaggedConstructorReflectType, newTaggedSingletonConstructorBinding)
	return newSingletonBuilder(b.module, b.bindingKeys[0].reflectType())
}

func (b *baseBuilder) to(object interface{}, verifyFunc func(reflect.Type, reflect.Type) error, newBindingFunc func(interface{}) binding) {
	objectReflectType := reflect.TypeOf(object)
	for _, bindingKey := range b.bindingKeys {
		if err := verifyFunc(bindingKey.reflectType(), objectReflectType); err != nil {
			b.module.addBindingError(err)
			return
		}
	}
	binding := newBindingFunc(object)
	for _, bindingKey := range b.bindingKeys {
		b.setBinding(bindingKey, binding)
	}
}

func (b *baseBuilder) setBinding(bindingKey bindingKey, binding binding) {
	b.module.setBinding(bindingKey, binding)
}

type singletonBuilder struct {
	module *module
	t      reflect.Type
	fn     interface{}
}

func (b *singletonBuilder) Eagerly() {
	if b == nil {
		return
	}
	b.module.eager = append(b.module.eager, b)
}

func (b *singletonBuilder) EagerlyAndCall(function interface{}) {
	if b == nil {
		return
	}
	b.fn = function
	b.module.eager = append(b.module.eager, b)
}

func newSingletonBuilder(module *module, t reflect.Type) SingletonBuilder {
	return &singletonBuilder{module: module, t: t}
}

func verifyBindingReflectType(bindingKeyReflectType reflect.Type, bindingReflectType reflect.Type) error {
	if !isSupportedBindingKeyReflectType(bindingKeyReflectType) {
		return errNotSupportedYet.withTag("bindingKeyReflectType", bindingReflectType)
	}
	if isInterfacePtr(bindingKeyReflectType) {
		bindingKeyReflectType = bindingKeyReflectType.Elem()
	}
	if !bindingReflectType.AssignableTo(bindingKeyReflectType) {
		return errNotAssignable.withTag("bindingKeyReflectType", bindingKeyReflectType).withTag("bindingReflectType", bindingReflectType)
	}
	return nil
}

func verifyConstructorReflectType(bindingKeyReflectType reflect.Type, constructorReflectType reflect.Type) error {
	if err := verifyIsFunc(constructorReflectType); err != nil {
		return err
	}
	return verifyConstructorReturnValues(bindingKeyReflectType, constructorReflectType)
}

func verifyTaggedConstructorReflectType(bindingKeyReflectType reflect.Type, constructorReflectType reflect.Type) error {
	if err := verifyIsTaggedFunc(constructorReflectType); err != nil {
		return err
	}
	return verifyConstructorReturnValues(bindingKeyReflectType, constructorReflectType)
}

func verifyConstructorReturnValues(bindingKeyReflectType reflect.Type, constructorReflectType reflect.Type) error {
	numOut := constructorReflectType.NumOut()
	if numOut < 1 || numOut > 2 {
		return errConstructorReturnValuesInvalid.withTag("constructorReflectType", constructorReflectType)
	}
	if bindingKeyReflectType != nil {
		if err := verifyBindingReflectType(bindingKeyReflectType, constructorReflectType.Out(0)); err != nil {
			return err
		}
	}
	if numOut == 2 && !constructorReflectType.Out(1).AssignableTo(errorReflectType) {
		return errConstructorReturnValuesInvalid.withTag("constructorReflectType", constructorReflectType)
	}
	return nil
}
