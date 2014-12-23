package inject

import (
	"reflect"
)

const (
	binderTypeValidBinder = iota
	binderTypeErr

	validBinderTypeBoundType
	validBinderTypeTaggedBoundType
)

type binder struct {
	binderType int

	validBinder *validBinder
	err         error
}

func newBinderBoundType(injector *injector, boundType reflect.Type) *binder {
	return &binder{binderTypeValidBinder, newValidBinderBoundType(injector, boundType), nil}
}

func newBinderTaggedBoundType(injector *injector, taggedBoundType taggedBoundType) *binder {
	return &binder{binderTypeValidBinder, newValidBinderTaggedBoundType(injector, taggedBoundType), nil}
}

func newBinderErr(err error) *binder {
	return &binder{binderTypeErr, nil, err}
}

func (this *binder) ToType(to interface{}) error {
	switch this.binderType {
	case binderTypeErr:
		return this.err
	case binderTypeValidBinder:
		return this.validBinder.toType(to)
	default:
		return ErrUnknownBinderType
	}
}

func (this *binder) ToSingleton(singleton interface{}) error {
	switch this.binderType {
	case binderTypeErr:
		return this.err
	case binderTypeValidBinder:
		return this.validBinder.toSingleton(singleton)
	default:
		return ErrUnknownBinderType
	}
}

func (this *binder) ToProvider(provider interface{}) error {
	switch this.binderType {
	case binderTypeErr:
		return this.err
	case binderTypeValidBinder:
		return this.validBinder.toProvider(provider)
	default:
		return ErrUnknownBinderType
	}
}

func (this *binder) ToProviderAsSingleton(provider interface{}) error {
	switch this.binderType {
	case binderTypeErr:
		return this.err
	case binderTypeValidBinder:
		return this.validBinder.toProviderAsSingleton(provider)
	default:
		return ErrUnknownBinderType
	}
}

type validBinder struct {
	injector *injector

	validBinderType int

	boundType       reflect.Type
	taggedBoundType taggedBoundType
}

func newValidBinderBoundType(injector *injector, boundType reflect.Type) *validBinder {
	return &validBinder{injector, validBinderTypeBoundType, boundType, taggedBoundType{}}
}

func newValidBinderTaggedBoundType(injector *injector, taggedBoundType taggedBoundType) *validBinder {
	return &validBinder{injector, validBinderTypeTaggedBoundType, nil, taggedBoundType}
}

func (this *validBinder) toType(to interface{}) error {
	if to == nil {
		return ErrNil
	}
	toReflectType := reflect.TypeOf(to)
	switch this.validBinderType {
	case validBinderTypeBoundType:
		err := verifyToType(this.boundType, toReflectType)
		if err != nil {
			return err
		}
	case validBinderTypeTaggedBoundType:
		err := verifyToType(this.taggedBoundType.boundType, toReflectType)
		if err != nil {
			return err
		}
	default:
		return ErrUnknownValidBinderType
	}
	return this.assignBinding(newBindingIntermediate(toReflectType))
}

func (this *validBinder) toSingleton(singleton interface{}) error {
	if singleton == nil {
		return ErrNil
	}
	singletonReflectType := reflect.TypeOf(singleton)
	switch this.validBinderType {
	case validBinderTypeBoundType:
		err := verifyBinding(this.boundType, singletonReflectType)
		if err != nil {
			return err
		}
	case validBinderTypeTaggedBoundType:
		err := verifyBinding(this.taggedBoundType.boundType, singletonReflectType)
		if err != nil {
			return err
		}
	default:
		return ErrUnknownValidBinderType
	}
	return this.assignBinding(newBindingFinal(newFinalBindingSingleton(singleton)))
}

func (this *validBinder) toProvider(provider interface{}) error {
	if provider == nil {
		return ErrNil
	}
	switch this.validBinderType {
	case validBinderTypeBoundType:
		err := verifyProvider(this.boundType, provider)
		if err != nil {
			return err
		}
	case validBinderTypeTaggedBoundType:
		err := verifyProvider(this.taggedBoundType.boundType, provider)
		if err != nil {
			return err
		}
	default:
		return ErrUnknownValidBinderType
	}
	return this.assignBinding(newBindingFinal(newFinalBindingProvider(provider)))
}

func (this *validBinder) toProviderAsSingleton(provider interface{}) error {
	if provider == nil {
		return ErrNil
	}
	switch this.validBinderType {
	case validBinderTypeBoundType:
		err := verifyProvider(this.boundType, provider)
		if err != nil {
			return err
		}
	case validBinderTypeTaggedBoundType:
		err := verifyProvider(this.taggedBoundType.boundType, provider)
		if err != nil {
			return err
		}
	default:
		return ErrUnknownValidBinderType
	}
	return this.assignBinding(newBindingFinal(newFinalBindingSingletonProvider(provider)))
}

func (this *validBinder) assignBinding(binding *binding) error {
	switch this.validBinderType {
	case validBinderTypeBoundType:
		this.injector.boundTypeToBinding[this.boundType] = binding
		return nil
	case validBinderTypeTaggedBoundType:
		this.injector.taggedBoundTypeToBinding[this.taggedBoundType] = binding
		return nil
	default:
		return ErrUnknownValidBinderType
	}
}

func verifyToType(fromReflectType reflect.Type, toReflectType reflect.Type) error {
	// TODO(pedge): is this restriction necessary/warranted? how about structs with anonymous fields?
	if !(fromReflectType.Kind() == reflect.Ptr && fromReflectType.Elem().Kind() == reflect.Interface) {
		return ErrNotInterfacePtr
	}
	if !toReflectType.Implements(fromReflectType.Elem()) {
		return ErrDoesNotImplement
	}
	return nil
}

func verifyProvider(fromReflectType reflect.Type, provider interface{}) error {
	providerReflectType := reflect.TypeOf(provider)
	if providerReflectType.Kind() != reflect.Func {
		return ErrProviderNotFunction
	}
	if providerReflectType.NumOut() != 2 {
		return ErrProviderReturnValuesInvalid
	}
	err := verifyBinding(fromReflectType, providerReflectType.Out(0))
	if err != nil {
		return err
	}
	// TODO(pedge): can this be simplified?
	if !providerReflectType.Out(1).AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
		return ErrProviderReturnValuesInvalid
	}
	return nil
}

func verifyBinding(fromReflectType reflect.Type, toReflectType reflect.Type) error {
	switch {
	// from is an interface
	case fromReflectType.Kind() == reflect.Ptr && fromReflectType.Elem().Kind() == reflect.Interface:
		if !toReflectType.Implements(fromReflectType.Elem()) {
			return ErrDoesNotImplement
		}
	// from is a struct pointer
	case fromReflectType.Kind() == reflect.Ptr && fromReflectType.Elem().Kind() == reflect.Struct:
		fallthrough
	// from is a struct
	case fromReflectType.Kind() == reflect.Struct:
		// TODO(pedge): is this correct?
		if !toReflectType.AssignableTo(fromReflectType) {
			return ErrNotAssignable
		}
	// nothing else is supported for now
	// TODO(pedge): at least support primitives with tags
	default:
		return ErrNotSupportedYet
	}
	return nil
}
