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
	return switchBinderType(this, to, func(validBinder *validBinder, object interface{}) error {
		return validBinder.toType(object)
	})
}

func (this *binder) ToSingleton(singleton interface{}) error {
	return switchBinderType(this, singleton, func(validBinder *validBinder, object interface{}) error {
		return validBinder.toSingleton(object)
	})
}

func (this *binder) ToProvider(provider interface{}) error {
	return switchBinderType(this, provider, func(validBinder *validBinder, object interface{}) error {
		return validBinder.toProvider(object)
	})
}

func (this *binder) ToProviderAsSingleton(provider interface{}) error {
	return switchBinderType(this, provider, func(validBinder *validBinder, object interface{}) error {
		return validBinder.toProviderAsSingleton(object)
	})
}

func switchBinderType(binder *binder, object interface{}, validBinderFunc func(*validBinder, interface{}) error) error {
	switch binder.binderType {
	case binderTypeErr:
		return binder.err
	case binderTypeValidBinder:
		return validBinderFunc(binder.validBinder, object)
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
	err := verifyValidBinderInput(this, to, verifyToType)
	if err != nil {
		return err
	}
	return assignValidBinderBinding(this, newBindingIntermediate(reflect.TypeOf(to)))
}

func (this *validBinder) toSingleton(singleton interface{}) error {
	err := verifyValidBinderInput(this, singleton, verifyBinding)
	if err != nil {
		return err
	}
	return assignValidBinderBinding(this, newBindingFinal(newFinalBindingSingleton(singleton)))
}

func (this *validBinder) toProvider(provider interface{}) error {
	err := verifyValidBinderInput(this, provider, verifyProvider)
	if err != nil {
		return err
	}
	return assignValidBinderBinding(this, newBindingFinal(newFinalBindingProvider(provider)))
}

func (this *validBinder) toProviderAsSingleton(provider interface{}) error {
	err := verifyValidBinderInput(this, provider, verifyProvider)
	if err != nil {
		return err
	}
	return assignValidBinderBinding(this, newBindingFinal(newFinalBindingSingletonProvider(provider)))
}

func verifyValidBinderInput(validBinder *validBinder, object interface{}, verifyFunc func(reflect.Type, interface{}) error) error {
	switch validBinder.validBinderType {
	case validBinderTypeBoundType:
		err := verifyFunc(validBinder.boundType, object)
		if err != nil {
			return err
		}
		return nil
	case validBinderTypeTaggedBoundType:
		err := verifyToType(validBinder.taggedBoundType.boundType, object)
		if err != nil {
			return err
		}
		return nil
	default:
		return ErrUnknownValidBinderType
	}
}

func assignValidBinderBinding(validBinder *validBinder, binding *binding) error {
	switch validBinder.validBinderType {
	case validBinderTypeBoundType:
		validBinder.injector.boundTypeToBinding[validBinder.boundType] = binding
		return nil
	case validBinderTypeTaggedBoundType:
		validBinder.injector.taggedBoundTypeToBinding[validBinder.taggedBoundType] = binding
		return nil
	default:
		return ErrUnknownValidBinderType
	}
}

func verifyToType(fromReflectType reflect.Type, to interface{}) error {
	toReflectType := reflect.TypeOf(to)
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
	err := verifyBindingReflectType(fromReflectType, providerReflectType.Out(0))
	if err != nil {
		return err
	}
	// TODO(pedge): can this be simplified?
	if !providerReflectType.Out(1).AssignableTo(reflect.TypeOf((*error)(nil)).Elem()) {
		return ErrProviderReturnValuesInvalid
	}
	return nil
}

func verifyBinding(fromReflectType reflect.Type, to interface{}) error {
	return verifyBindingReflectType(fromReflectType, reflect.TypeOf(to))
}

func verifyBindingReflectType(fromReflectType reflect.Type, toReflectType reflect.Type) error {
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
