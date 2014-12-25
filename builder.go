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
	module     *module
	bindingKey bindingKey
}

func newBuilder(module *module, bindingKey bindingKey) Builder {
	return &baseBuilder{module, bindingKey}
}

func (this *baseBuilder) To(to interface{}) {
	this.to(to, this.verifyToReflectType, newIntermediateBinding)
}

func (this *baseBuilder) ToSingleton(singleton interface{}) {
	this.to(singleton, this.verifyBindingReflectType, newSingletonBinding)
}

func (this *baseBuilder) ToConstructor(constructor interface{}) {
	this.to(constructor, this.verifyConstructorReflectType, newConstructorBinding)
}

func (this *baseBuilder) ToSingletonConstructor(constructor interface{}) {
	this.to(constructor, this.verifyConstructorReflectType, newSingletonConstructorBinding)
}

func (this *baseBuilder) ToTaggedConstructor(constructor interface{}) {
	this.to(constructor, this.verifyTaggedConstructorReflectType, newTaggedConstructorBinding)
}

func (this *baseBuilder) ToTaggedSingletonConstructor(constructor interface{}) {
	this.to(constructor, this.verifyTaggedConstructorReflectType, newTaggedSingletonConstructorBinding)
}

func (this *baseBuilder) to(object interface{}, verifyFunc func(reflect.Type) error, newBindingFunc func(interface{}) binding) {
	objectReflectType := reflect.TypeOf(object)
	err := verifyFunc(objectReflectType)
	if err != nil {
		this.module.addBindingError(err)
		return
	}
	this.setBinding(newBindingFunc(object))
}

func (this *baseBuilder) verifyToReflectType(toReflectType reflect.Type) error {
	bindingKeyReflectType := this.bindingKey.reflectType()
	// TODO(pedge): is this restriction necessary/warranted? how about structs with anonymous fields?
	if !(bindingKeyReflectType.Kind() == reflect.Ptr && bindingKeyReflectType.Elem().Kind() == reflect.Interface) {
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

func (this *baseBuilder) verifyBindingReflectType(bindingReflectType reflect.Type) error {
	bindingKeyReflectType := this.bindingKey.reflectType()
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

func (this *baseBuilder) verifyConstructorReflectType(constructorReflectType reflect.Type) error {
	if !isFunc(constructorReflectType) {
		eb := newErrorBuilder(injectErrorTypeConstructorNotFunction)
		eb = eb.addTag("constructorReflectType", constructorReflectType)
		return eb.build()
	}
	if constructorReflectType.NumOut() != 2 {
		eb := newErrorBuilder(injectErrorTypeConstructorReturnValuesInvalid)
		eb = eb.addTag("constructorReflectType", constructorReflectType)
		return eb.build()
	}
	err := this.verifyBindingReflectType(constructorReflectType.Out(0))
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

func (this *baseBuilder) verifyTaggedConstructorReflectType(constructorReflectType reflect.Type) error {
	err := this.verifyConstructorReflectType(constructorReflectType)
	if err != nil {
		return err
	}
	if constructorReflectType.NumIn() != 1 {
		eb := newErrorBuilder(injectErrorTypeTaggedConstructorParametersInvalid)
		eb = eb.addTag("constructorReflectType", constructorReflectType)
		return eb.build()
	}
	inReflectType := constructorReflectType.In(0)
	if !isStruct(inReflectType) {
		eb := newErrorBuilder(injectErrorTypeTaggedConstructorParametersInvalid)
		eb = eb.addTag("constructorReflectType", constructorReflectType)
		return eb.build()
	}
	if inReflectType.Name() != "" {
		eb := newErrorBuilder(injectErrorTypeTaggedConstructorParametersInvalid)
		eb = eb.addTag("constructorReflectType", constructorReflectType)
		return eb.build()
	}
	return nil
}

func (this *baseBuilder) setBinding(binding binding) {
	this.module.setBinding(this.bindingKey, binding)
}

func isInterfacePtr(reflectType reflect.Type) bool {
	return isPtr(reflectType) && isInterface(reflectType.Elem())
}

func isStructPtr(reflectType reflect.Type) bool {
	return isPtr(reflectType) && isStruct(reflectType.Elem())
}

func isInterface(reflectType reflect.Type) bool {
	return reflectType.Kind() == reflect.Interface
}

func isStruct(reflectType reflect.Type) bool {
	return reflectType.Kind() == reflect.Struct
}

func isPtr(reflectType reflect.Type) bool {
	return reflectType.Kind() == reflect.Ptr
}

func isFunc(reflectType reflect.Type) bool {
	return reflectType.Kind() == reflect.Func
}
