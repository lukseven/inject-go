package inject

import (
	"reflect"
)

type bindingKey interface {
	reflectType() reflect.Type
}

type baseBindingKey struct {
	rt reflect.Type
}

func newBindingKey(reflectType reflect.Type) bindingKey {
	return &baseBindingKey{reflectType}
}

func (this *baseBindingKey) reflectType() reflect.Type {
	return this.rt
}

type taggedBindingKey struct {
	baseBindingKey
	tag string
}

func newTaggedBindingKey(reflectType reflect.Type, tag string) bindingKey {
	return &taggedBindingKey{baseBindingKey{reflectType}, tag}
}
