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

func (n *noOpBuilder) ToSingletonConstructor(construtor interface{}) {}

func (n *noOpBuilder) ToTaggedConstructor(constructor interface{}) {}

func (n *noOpBuilder) ToTaggedSingletonConstructor(constructor interface{}) {}

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

func (b *baseBuilder) ToSingletonConstructor(constructor interface{}) {
	b.to(constructor, verifyConstructorReflectType, newSingletonConstructorBinding)
}

func (b *baseBuilder) ToTaggedConstructor(constructor interface{}) {
	b.to(constructor, verifyTaggedConstructorReflectType, newTaggedConstructorBinding)
}

func (b *baseBuilder) ToTaggedSingletonConstructor(constructor interface{}) {
	b.to(constructor, verifyTaggedConstructorReflectType, newTaggedSingletonConstructorBinding)
}

func (b *baseBuilder) to(object interface{}, verifyFunc func(reflect.Type, reflect.Type) error, newBindingFunc func(interface{}) binding) {
	objectReflectType := reflect.TypeOf(object)
	for _, bindingKey := range b.bindingKeys {
		err := verifyFunc(bindingKey.reflectType(), objectReflectType)
		if err != nil {
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
	err := verifyIsFunc(constructorReflectType)
	if err != nil {
		return err
	}
	err = verifyConstructorReturnValues(bindingKeyReflectType, constructorReflectType)
	if err != nil {
		return err
	}
	return nil
}

func verifyTaggedConstructorReflectType(bindingKeyReflectType reflect.Type, constructorReflectType reflect.Type) error {
	err := verifyIsTaggedFunc(constructorReflectType)
	if err != nil {
		return err
	}
	err = verifyConstructorReturnValues(bindingKeyReflectType, constructorReflectType)
	if err != nil {
		return err
	}
	return nil
}

func verifyConstructorReturnValues(bindingKeyReflectType reflect.Type, constructorReflectType reflect.Type) error {
	if constructorReflectType.NumOut() != 2 {
		return errConstructorReturnValuesInvalid.withTag("constructorReflectType", constructorReflectType)
	}
	err := verifyBindingReflectType(bindingKeyReflectType, constructorReflectType.Out(0))
	if err != nil {
		return err
	}
	if !constructorReflectType.Out(1).AssignableTo(errorReflectType) {
		return errConstructorReturnValuesInvalid.withTag("constructorReflectType", constructorReflectType)
	}
	return nil
}
