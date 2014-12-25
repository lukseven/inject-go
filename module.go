package inject

import (
	"bytes"
	"reflect"
	"strings"
)

type module struct {
	bindings map[bindingKey]binding
}

func createModule() *module {
	return &module{make(map[bindingKey]binding)}
}

func (this *module) Bind(from interface{}) Builder {
	fromReflectType := reflect.TypeOf(from)
	if fromReflectType == nil {
		return newPropogatedErrorBuilder(newErrorBuilder(InjectErrorTypeNil).build())
	}
	return newBuilder(this, newBindingKey(fromReflectType))
}

func (this *module) BindTagged(from interface{}, tag string) Builder {
	fromReflectType := reflect.TypeOf(from)
	if fromReflectType == nil {
		return newPropogatedErrorBuilder(newErrorBuilder(InjectErrorTypeNil).build())
	}
	if tag == "" {
		return newPropogatedErrorBuilder(newErrorBuilder(InjectErrorTypeTagEmpty).build())
	}
	return newBuilder(this, newTaggedBindingKey(fromReflectType, tag))
}

func (this *module) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("module{")
	buffer.WriteString(strings.Join(this.keyValueStrings(), " "))
	buffer.WriteString("}")
	return buffer.String()
}

func (this *module) keyValueStrings() []string {
	strings := make([]string, len(this.bindings))
	i := 0
	for bindingKey, binding := range this.bindings {
		var buffer bytes.Buffer
		buffer.WriteString(bindingKey.String())
		buffer.WriteString(":")
		buffer.WriteString(binding.String())
		strings[i] = buffer.String()
		i++
	}
	return strings
}

func (this *module) binding(bindingKey bindingKey) (binding, bool) {
	binding, ok := this.bindings[bindingKey]
	return binding, ok
}

func (this *module) setBinding(bindingKey bindingKey, binding binding) error {
	foundBinding, ok := this.bindings[bindingKey]
	if ok {
		eb := newErrorBuilder(InjectErrorTypeAlreadyBound)
		eb.addTag("bindingKey", bindingKey)
		eb.addTag("foundBinding", foundBinding)
		return eb.build()
	}
	this.bindings[bindingKey] = binding
	return nil
}
