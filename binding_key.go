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

func (this baseBindingKey) reflectType() reflect.Type {
	return this.rt
}

func (this baseBindingKey) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("baseBindingKey{reflectType:")
	buffer.WriteString(this.reflectType().String())
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

func (this taggedBindingKey) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("baseBindingKey{reflectType:")
	buffer.WriteString(this.reflectType().String())
	buffer.WriteString(" tag:")
	buffer.WriteString(this.tag)
	buffer.WriteString("}")
	return buffer.String()
}
