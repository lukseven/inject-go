package inject

import (
	"reflect"
)

type noOpBuilder struct{}

func newNoOpBuilder() Builder {
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

func newBuilder(module *module, bindingKeys []bindingKey) Builder {
	return &baseBuilder{module, bindingKeys}
}

func (this *baseBuilder) To(to interface{}) {
	this.to(to, verifyToReflectType, newIntermediateBinding)
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

func verifyToReflectType(bindingKeyReflectType reflect.Type, toReflectType reflect.Type) error {
	// TODO(pedge): is this restriction necessary/warranted? how about structs with anonymous fields?
	if !isInterfacePtr(bindingKeyReflectType) {
		eb := newErrorBuilder(injectErrorTypeNotInterfacePtr)
		eb = eb.addTag("bindingKeyReflectType", bindingKeyReflectType)
		return eb.build()
	}
	if !toReflectType.Implements(bindingKeyReflectType.Elem()) {
		eb := newErrorBuilder(injectErrorTypeDoesNotImplement)
		eb = eb.addTag("toReflectType", toReflectType)
		eb = eb.addTag("bindingKeyReflectType", bindingKeyReflectType)
		return eb.build()
	}
	return nil
}

func verifyBindingReflectType(bindingKeyReflectType reflect.Type, bindingReflectType reflect.Type) error {
	switch {
	case isInterfacePtr(bindingKeyReflectType):
		if !bindingReflectType.Implements(bindingKeyReflectType.Elem()) {
			eb := newErrorBuilder(injectErrorTypeDoesNotImplement)
			eb = eb.addTag("bindingKeyReflectType", bindingKeyReflectType)
			eb = eb.addTag("bindingReflectType", bindingReflectType)
			return eb.build()
		}
	case isStructPtr(bindingKeyReflectType), isStruct(bindingKeyReflectType):
		// TODO(pedge): is this correct?
		if !bindingReflectType.AssignableTo(bindingKeyReflectType) {
			eb := newErrorBuilder(injectErrorTypeNotAssignable)
			eb = eb.addTag("bindingKeyReflectType", bindingKeyReflectType)
			eb = eb.addTag("bindingReflectType", bindingReflectType)
			return eb.build()
		}
	// nothing else is supported for now
	// TODO(pedge): at least support primitives with tags
	default:
		eb := newErrorBuilder(injectErrorTypeNotSupportedYet)
		eb = eb.addTag("bindingKeyReflectType", bindingKeyReflectType)
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
	// TODO(pedge): can this be simplified?
	if !constructorReflectType.Out(1).AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
		eb := newErrorBuilder(injectErrorTypeConstructorReturnValuesInvalid)
		eb = eb.addTag("constructorReflectType", constructorReflectType)
		return eb.build()
	}
	return nil
}
