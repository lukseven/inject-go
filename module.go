package inject

import (
	"bytes"
	"reflect"
	"strings"
)

type module struct {
	bindings      map[bindingKey]binding
	bindingErrors []error
}

func createModule() *module {
	return &module{make(map[bindingKey]binding), make([]error, 0)}
}

func (this *module) Bind(from ...interface{}) Builder {
	lenFrom := len(from)
	if lenFrom == 0 {
		this.addBindingError(newErrorBuilder(injectErrorTypeNil).build())
		return newNoOpBuilder()
	}
	bindingKeys := make([]bindingKey, lenFrom)
	for i := 0; i < lenFrom; i++ {
		fromReflectType := reflect.TypeOf(from[i])
		if fromReflectType == nil {
			this.addBindingError(newErrorBuilder(injectErrorTypeNil).build())
			return newNoOpBuilder()
		}
		bindingKeys[i] = newBindingKey(fromReflectType)
	}
	return newBuilder(this, bindingKeys)
}

func (this *module) BindTagged(tag string, from ...interface{}) Builder {
	if tag == "" {
		this.addBindingError(newErrorBuilder(injectErrorTypeTagEmpty).build())
		return newNoOpBuilder()
	}
	lenFrom := len(from)
	if lenFrom == 0 {
		this.addBindingError(newErrorBuilder(injectErrorTypeNil).build())
		return newNoOpBuilder()
	}
	bindingKeys := make([]bindingKey, lenFrom)
	for i := 0; i < lenFrom; i++ {
		fromReflectType := reflect.TypeOf(from[i])
		if fromReflectType == nil {
			this.addBindingError(newErrorBuilder(injectErrorTypeNil).build())
			return newNoOpBuilder()
		}
		bindingKeys[i] = newTaggedBindingKey(fromReflectType, tag)
	}
	return newBuilder(this, bindingKeys)
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

func (this *module) addBindingError(err error) {
	this.bindingErrors = append(this.bindingErrors, err)
}

func (this *module) binding(bindingKey bindingKey) (binding, bool) {
	binding, ok := this.bindings[bindingKey]
	return binding, ok
}

func (this *module) setBinding(bindingKey bindingKey, binding binding) {
	foundBinding, ok := this.bindings[bindingKey]
	if ok {
		eb := newErrorBuilder(injectErrorTypeAlreadyBound)
		eb.addTag("bindingKey", bindingKey)
		eb.addTag("foundBinding", foundBinding)
		this.addBindingError(eb.build())
		return
	}
	this.bindings[bindingKey] = binding
}
