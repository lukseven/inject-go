package inject

import (
	"reflect"
	"sync"
	"sync/atomic"
)

const (
	bindingTypeIntermediate = iota
	bindingTypeFinal

	finalBindingTypeSingleton
	finalBindingTypeProvider
	finalBindingTypeSingletonProvider
)

type binding struct {
	bindingType int

	intermediateBinding reflect.Type
	finalBinding        *finalBinding
}

func newBindingIntermediate(intermediate reflect.Type) *binding {
	return &binding{bindingTypeIntermediate, intermediate, nil}
}

func newBindingFinal(final *finalBinding) *binding {
	return &binding{bindingTypeFinal, nil, final}
}

type finalBinding struct {
	finalBindingType int

	singleton         interface{}
	provider          interface{}
	singletonProvider *singletonProvider
}

func newFinalBindingSingleton(singleton interface{}) *finalBinding {
	return &finalBinding{finalBindingTypeSingleton, singleton, nil, nil}
}

func newFinalBindingProvider(provider interface{}) *finalBinding {
	return &finalBinding{finalBindingTypeProvider, nil, provider, nil}
}

func newFinalBindingSingletonProvider(provider interface{}) *finalBinding {
	return &finalBinding{finalBindingTypeSingletonProvider, nil, nil, newSingletonProvider(provider)}
}

func (this *finalBinding) get(c *container) (interface{}, error) {
	switch this.finalBindingType {
	case finalBindingTypeSingleton:
		return this.singleton, nil
	case finalBindingTypeProvider:
		return getFromProvider(c, this.provider)
	case finalBindingTypeSingletonProvider:
		return this.singletonProvider.get(c)
	default:
		return nil, ErrNotSupportedYet
	}
}

type singletonProvider struct {
	provider interface{}
	// TODO(pedge): is atomic.Value the equivalent of a volatile variable in Java?
	value atomic.Value
	once  sync.Once
}

func newSingletonProvider(provider interface{}) *singletonProvider {
	return &singletonProvider{provider, atomic.Value{}, sync.Once{}}
}

func (this *singletonProvider) get(c *container) (interface{}, error) {
	this.once.Do(func() {
		value, err := getFromProvider(c, this.provider)
		this.value.Store(&valueErr{value, err})
	})
	valueErr := this.value.Load().(*valueErr)
	return valueErr.value, valueErr.err
}

type valueErr struct {
	value interface{}
	err   error
}

// TODO(pedge): this is really hacky, and probably slow, clean this up
func getFromProvider(c *container, provider interface{}) (interface{}, error) {
	// assuming this is a valid provider/that this is already checked
	providerReflectType := reflect.TypeOf(provider)
	numIn := providerReflectType.NumIn()
	parameterValues := make([]reflect.Value, numIn)
	if numIn == 1 && providerReflectType.In(0).AssignableTo(reflect.TypeOf((*Container)(nil)).Elem()) {
		parameterValues[0] = reflect.ValueOf(c)
	} else {
		for i := 0; i < numIn; i++ {
			inReflectType := providerReflectType.In(i)
			// TODO(pedge): this is really specific logic, and there wil need to be more
			// of this if more types are allowed for binding - this should be abstracted
			if inReflectType.Kind() == reflect.Interface {
				inReflectType = reflect.PtrTo(inReflectType)
			}
			parameter, err := c.get(inReflectType)
			if err != nil {
				return nil, err
			}
			parameterValues[i] = reflect.ValueOf(parameter)
		}
	}
	returnValues := reflect.ValueOf(provider).Call(parameterValues)
	return1 := returnValues[0].Interface()
	return2 := returnValues[1].Interface()
	switch {
	case return1 != nil && return2 != nil:
		return nil, ErrInvalidReturnFromProvider
	case return2 != nil:
		return nil, return2.(error)
	default:
		return return1, nil
	}
}
