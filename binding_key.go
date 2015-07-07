package inject

import (
	"fmt"
	"reflect"
)

type bindingKey interface {
	fmt.Stringer
	reflectType() reflect.Type
}

type baseBindingKey struct {
	rt reflect.Type
}

func newBindingKey(reflectType reflect.Type) bindingKey {
	return baseBindingKey{reflectType}
}

func (b baseBindingKey) reflectType() reflect.Type {
	return b.rt
}

func (b baseBindingKey) String() string {
	return fmt.Sprintf("{type:%s}", b.reflectType().String())
}

type taggedBindingKey struct {
	baseBindingKey
	tag string
}

func newTaggedBindingKey(reflectType reflect.Type, tag string) bindingKey {
	return taggedBindingKey{baseBindingKey{reflectType}, tag}
}

func (t taggedBindingKey) String() string {
	return fmt.Sprintf("{type:%s tag:%s}", t.reflectType().String(), t.tag)
}
