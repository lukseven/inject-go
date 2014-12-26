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

func (this *module) Bind(froms ...interface{}) Builder {
	return this.bind(newBindingKey, froms)
}

func (this *module) BindTagged(tag string, froms ...interface{}) Builder {
	ok := this.verifyTag(tag)
	if !ok {
		return newNoOpBuilder()
	}
	return this.bind(func(fromReflectType reflect.Type) bindingKey { return newTaggedBindingKey(fromReflectType, tag) }, froms)
}

func (this *module) BindInterface(fromInterfaces ...interface{}) InterfaceBuilder {
	ok := this.verifyAreInterfacePtrs(fromInterfaces)
	if !ok {
		return newNoOpBuilder()
	}
	return this.bind(newBindingKey, fromInterfaces)
}

func (this *module) BindTaggedInterface(tag string, fromInterfaces ...interface{}) InterfaceBuilder {
	ok := this.verifyTag(tag)
	if !ok {
		return newNoOpBuilder()
	}
	ok = this.verifyAreInterfacePtrs(fromInterfaces)
	if !ok {
		return newNoOpBuilder()
	}
	return this.bind(func(fromReflectType reflect.Type) bindingKey { return newTaggedBindingKey(fromReflectType, tag) }, fromInterfaces)
}

func (this *module) bind(newBindingKeyFunc func(reflect.Type) bindingKey, from []interface{}) InterfaceBuilder {
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
		bindingKeys[i] = newBindingKeyFunc(fromReflectType)
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

func (this *module) verifyTag(tag string) bool {
	if tag == "" {
		this.addBindingError(newErrorBuilder(injectErrorTypeTagEmpty).build())
		return false
	}
	return true
}

func (this *module) verifyAreInterfacePtrs(fromInterfaces []interface{}) bool {
	var ok bool = true
	for _, fromInterface := range fromInterfaces {
		if !this.verifyIsInterfacePtr(reflect.TypeOf(fromInterface)) {
			ok = false
		}
	}
	return ok
}

func (this *module) verifyIsInterfacePtr(reflectType reflect.Type) bool {
	if !isInterfacePtr(reflectType) {
		eb := newErrorBuilder(injectErrorTypeNotInterfacePtr)
		eb.addTag("reflectType", reflectType)
		this.addBindingError(eb.build())
		return false
	}
	return true
}
