package inject

import (
	"bytes"
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
	var buffer bytes.Buffer
	buffer.WriteString("{type:")
	buffer.WriteString(b.reflectType().String())
	buffer.WriteString("}")
	return buffer.String()
}

type taggedBindingKey struct {
	baseBindingKey
	tag string
}

func newTaggedBindingKey(reflectType reflect.Type, tag string) bindingKey {
	return taggedBindingKey{baseBindingKey{reflectType}, tag}
}

func (b taggedBindingKey) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("{type:")
	buffer.WriteString(b.reflectType().String())
	buffer.WriteString(" tag:")
	buffer.WriteString(b.tag)
	buffer.WriteString("}")
	return buffer.String()
}
