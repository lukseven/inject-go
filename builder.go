package inject

import (
	"reflect"
)

type propogatedErrorBuilder struct {
	err error
}

func newPropogatedErrorBuilder(err error) Builder {
	return &propogatedErrorBuilder{err}
}

func (this *propogatedErrorBuilder) To(to interface{}) error {
	return this.err
}

func (this *propogatedErrorBuilder) ToSingleton(singleton interface{}) error {
	return this.err
}

func (this *propogatedErrorBuilder) ToConstructor(constructor interface{}) error {
	return this.err
}

func (this *propogatedErrorBuilder) ToSingletonConstructor(construtor interface{}) error {
	return this.err
}

func (this *propogatedErrorBuilder) ToTaggedConstructor(constructor interface{}) error {
	return this.err
}

func (this *propogatedErrorBuilder) ToTaggedSingletonConstructor(constructor interface{}) error {
	return this.err
}

type baseBuilder struct {
	module     *module
	bindingKey bindingKey
}

func newBuilder(module *module, bindingKey bindingKey) Builder {
	return &baseBuilder{module, bindingKey}
}

func (this *baseBuilder) To(to interface{}) error {
	return this.to(to, this.verifyToReflectType, newIntermediateBinding)
}

func (this *baseBuilder) ToSingleton(singleton interface{}) error {
	return this.to(singleton, this.verifyBindingReflectType, newSingletonBinding)
}

func (this *baseBuilder) ToConstructor(constructor interface{}) error {
	return this.to(constructor, this.verifyConstructorReflectType, newConstructorBinding)
}

func (this *baseBuilder) ToSingletonConstructor(constructor interface{}) error {
	return this.to(constructor, this.verifyConstructorReflectType, newSingletonConstructorBinding)
}

func (this *baseBuilder) ToTaggedConstructor(constructor interface{}) error {
	return this.to(constructor, this.verifyTaggedConstructorReflectType, newTaggedConstructorBinding)
}

func (this *baseBuilder) ToTaggedSingletonConstructor(constructor interface{}) error {
	return this.to(constructor, this.verifyTaggedConstructorReflectType, newTaggedSingletonConstructorBinding)
}

func (this *baseBuilder) to(object interface{}, verifyFunc func(reflect.Type) error, newBindingFunc func(interface{}) binding) error {
	objectReflectType := reflect.TypeOf(object)
	err := verifyFunc(objectReflectType)
	if err != nil {
		return err
	}
	return this.setBinding(newBindingFunc(object))
}

func (this *baseBuilder) verifyToReflectType(toReflectType reflect.Type) error {
	bindingKeyReflectType := this.bindingKey.reflectType()
	// TODO(pedge): is this restriction necessary/warranted? how about structs with anonymous fields?
	if !(bindingKeyReflectType.Kind() == reflect.Ptr && bindingKeyReflectType.Elem().Kind() == reflect.Interface) {
		eb := newErrorBuilder(InjectErrorTypeNotInterfacePtr)
		eb = eb.addTag("bindingKeyReflectType", bindingKeyReflectType)
		return eb.build()
	}
	if !toReflectType.Implements(bindingKeyReflectType.Elem()) {
		eb := newErrorBuilder(InjectErrorTypeDoesNotImplement)
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
			eb := newErrorBuilder(InjectErrorTypeDoesNotImplement)
			eb = eb.addTag("bindingKeyReflectType", bindingKeyReflectType)
			eb = eb.addTag("bindingReflectType", bindingReflectType)
			return eb.build()
		}
	case isStructPtr(bindingKeyReflectType), isStruct(bindingKeyReflectType):
		// TODO(pedge): is this correct?
		if !bindingReflectType.AssignableTo(bindingKeyReflectType) {
			eb := newErrorBuilder(InjectErrorTypeNotAssignable)
			eb = eb.addTag("bindingKeyReflectType", bindingKeyReflectType)
			eb = eb.addTag("bindingReflectType", bindingReflectType)
			return eb.build()
		}
	// nothing else is supported for now
	// TODO(pedge): at least support primitives with tags
	default:
		eb := newErrorBuilder(InjectErrorTypeNotSupportedYet)
		eb = eb.addTag("bindingKeyReflectType", bindingKeyReflectType)
		return eb.build()
	}
	return nil
}

func (this *baseBuilder) verifyConstructorReflectType(constructorReflectType reflect.Type) error {
	if !isFunc(constructorReflectType) {
		eb := newErrorBuilder(InjectErrorTypeConstructorNotFunction)
		eb = eb.addTag("constructorReflectType", constructorReflectType)
		return eb.build()
	}
	if constructorReflectType.NumOut() != 2 {
		eb := newErrorBuilder(InjectErrorTypeConstructorReturnValuesInvalid)
		eb = eb.addTag("constructorReflectType", constructorReflectType)
		return eb.build()
	}
	err := this.verifyBindingReflectType(constructorReflectType.Out(0))
	if err != nil {
		return err
	}
	// TODO(pedge): can this be simplified?
	if !constructorReflectType.Out(1).AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
		eb := newErrorBuilder(InjectErrorTypeConstructorReturnValuesInvalid)
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
		eb := newErrorBuilder(InjectErrorTypeTaggedConstructorParametersInvalid)
		eb = eb.addTag("constructorReflectType", constructorReflectType)
		return eb.build()
	}
	inReflectType := constructorReflectType.In(0)
	if !isStruct(inReflectType) {
		eb := newErrorBuilder(InjectErrorTypeTaggedConstructorParametersInvalid)
		eb = eb.addTag("constructorReflectType", constructorReflectType)
		return eb.build()
	}
	if inReflectType.Name() != "" {
		eb := newErrorBuilder(InjectErrorTypeTaggedConstructorParametersInvalid)
		eb = eb.addTag("constructorReflectType", constructorReflectType)
		return eb.build()
	}
	return nil
}

func (this *baseBuilder) setBinding(binding binding) error {
	return this.module.setBinding(this.bindingKey, binding)
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
