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

func newBuilder(module *module, bindingType reflect.Type) Builder {
	return &baseBuilder{module, newBindingKey(bindingType)}
}

func newTaggedBuilder(module *module, bindingType reflect.Type, tag string) Builder {
	return &baseBuilder{module, newTaggedBindingKey(bindingType, tag)}
}

func (this *baseBuilder) To(to interface{}) error {
	toReflectType := reflect.TypeOf(to)
	err := this.verifyToReflectType(toReflectType)
	if err != nil {
		return nil
	}
	return this.assignBinding(newIntermediateBinding(newBindingKey(toReflectType)))
}

func (this *baseBuilder) ToSingleton(singleton interface{}) error {
	singletonReflectType := reflect.TypeOf(singleton)
	err := this.verifyBindingReflectType(singletonReflectType)
	if err != nil {
		return nil
	}
	return this.assignBinding(newSingletonBinding(singleton))
}

func (this *baseBuilder) ToConstructor(constructor interface{}) error {
	constructorReflectType := reflect.TypeOf(constructor)
	err := this.verifyConstructorReflectType(constructorReflectType)
	if err != nil {
		return nil
	}
	return this.assignBinding(newConstructorBinding(constructor))
}

func (this *baseBuilder) ToSingletonConstructor(constructor interface{}) error {
	constructorReflectType := reflect.TypeOf(constructor)
	err := this.verifyConstructorReflectType(constructorReflectType)
	if err != nil {
		return nil
	}
	return this.assignBinding(newSingletonConstructorBinding(constructor))
}

func (this *baseBuilder) ToTaggedConstructor(constructor interface{}) error {
	return newErrorBuilder("NOT IMPLEMENTED").build()
}

func (this *baseBuilder) ToTaggedSingletonConstructor(constructor interface{}) error {
	return newErrorBuilder("NOT IMPLEMENTED").build()
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

func (this *baseBuilder) verifyConstructorReflectType(constructorReflectType reflect.Type) error {
	bindingKeyReflectType := this.bindingKey.reflectType()
	if constructorReflectType.Kind() != reflect.Func {
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

func (this *baseBuilder) verifyBindingReflectType(bindingReflectType reflect.Type) error {
	bindingKeyReflectType := this.bindingKey.reflectType()
	switch {
	// from is an interface
	case bindingKeyReflectType.Kind() == reflect.Ptr && bindingKeyReflectType.Elem().Kind() == reflect.Interface:
		if !bindingReflectType.Implements(bindingKeyReflectType.Elem()) {
			eb := newErrorBuilder(InjectErrorTypeDoesNotImplement)
			eb = eb.addTag("bindingKeyReflectType", bindingKeyReflectType)
			eb = eb.addTag("bindingReflectType", bindingReflectType)
			return eb.build()
		}
	// from is a struct pointer
	case bindingKeyReflectType.Kind() == reflect.Ptr && bindingKeyReflectType.Elem().Kind() == reflect.Struct:
		fallthrough
	// from is a struct
	case bindingKeyReflectType.Kind() == reflect.Struct:
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

func (this *baseBuilder) assignBinding(binding binding) error {
	return this.module.assignBinding(this.bindingKey, binding)
}
