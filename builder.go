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

func (this *noOpBuilder) To(to interface{}) {}

func (this *noOpBuilder) ToSingleton(singleton interface{}) {}

func (this *noOpBuilder) ToConstructor(constructor interface{}) {}

func (this *noOpBuilder) ToSingletonConstructor(construtor interface{}) {}

func (this *noOpBuilder) ToTaggedConstructor(constructor interface{}) {}

func (this *noOpBuilder) ToTaggedSingletonConstructor(constructor interface{}) {}

type baseBuilder struct {
	module      *module
	bindingKeys []bindingKey
}

func newBuilder(module *module, bindingKeys []bindingKey) InterfaceBuilder {
	return &baseBuilder{module, bindingKeys}
}

func (this *baseBuilder) To(to interface{}) {
	this.to(to, verifyBindingReflectType, newIntermediateBinding)
}

func (this *baseBuilder) ToSingleton(singleton interface{}) {
	this.to(singleton, verifyBindingReflectType, newSingletonBinding)
}

func (this *baseBuilder) ToConstructor(constructor interface{}) {
	this.to(constructor, verifyConstructorReflectType, newConstructorBinding)
}

func (this *baseBuilder) ToSingletonConstructor(constructor interface{}) {
	this.to(constructor, verifyConstructorReflectType, newSingletonConstructorBinding)
}

func (this *baseBuilder) ToTaggedConstructor(constructor interface{}) {
	this.to(constructor, verifyTaggedConstructorReflectType, newTaggedConstructorBinding)
}

func (this *baseBuilder) ToTaggedSingletonConstructor(constructor interface{}) {
	this.to(constructor, verifyTaggedConstructorReflectType, newTaggedSingletonConstructorBinding)
}

func (this *baseBuilder) to(object interface{}, verifyFunc func(reflect.Type, reflect.Type) error, newBindingFunc func(interface{}) binding) {
	objectReflectType := reflect.TypeOf(object)
	for _, bindingKey := range this.bindingKeys {
		err := verifyFunc(bindingKey.reflectType(), objectReflectType)
		if err != nil {
			this.module.addBindingError(err)
			return
		}
	}
	binding := newBindingFunc(object)
	for _, bindingKey := range this.bindingKeys {
		this.setBinding(bindingKey, binding)
	}
}

func (this *baseBuilder) setBinding(bindingKey bindingKey, binding binding) {
	this.module.setBinding(bindingKey, binding)
}

func verifyBindingReflectType(bindingKeyReflectType reflect.Type, bindingReflectType reflect.Type) error {
	if !isSupportedBindingKeyReflectType(bindingKeyReflectType) {
		eb := newErrorBuilder(injectErrorTypeNotSupportedYet)
		eb = eb.addTag("bindingKeyReflectType", bindingKeyReflectType)
		return eb.build()
	}
	if isInterfacePtr(bindingKeyReflectType) {
		bindingKeyReflectType = bindingKeyReflectType.Elem()
	}
	if !bindingReflectType.AssignableTo(bindingKeyReflectType) {
		eb := newErrorBuilder(injectErrorTypeNotAssignable)
		eb = eb.addTag("bindingKeyReflectType", bindingKeyReflectType)
		eb = eb.addTag("bindingReflectType", bindingReflectType)
		return eb.build()
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
		eb := newErrorBuilder(injectErrorTypeConstructorReturnValuesInvalid)
		eb = eb.addTag("constructorReflectType", constructorReflectType)
		return eb.build()
	}
	err := verifyBindingReflectType(bindingKeyReflectType, constructorReflectType.Out(0))
	if err != nil {
		return err
	}
	if !constructorReflectType.Out(1).AssignableTo(errorReflectType) {
		eb := newErrorBuilder(injectErrorTypeConstructorReturnValuesInvalid)
		eb = eb.addTag("constructorReflectType", constructorReflectType)
		return eb.build()
	}
	return nil
}
