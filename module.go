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

func (m *module) Bind(froms ...interface{}) Builder {
	if !m.verifySupportedTypes(froms, isSupportedBindReflectType) {
		return newNoOpBuilder()
	}
	return m.bind(newBindingKey, froms)
}

func (m *module) BindTagged(tag string, froms ...interface{}) Builder {
	if !m.verifyTag(tag) {
		return newNoOpBuilder()
	}
	if !m.verifySupportedTypes(froms, isSupportedBindReflectType) {
		return newNoOpBuilder()
	}
	return m.bind(func(fromReflectType reflect.Type) bindingKey { return newTaggedBindingKey(fromReflectType, tag) }, froms)
}

func (m *module) BindInterface(fromInterfaces ...interface{}) InterfaceBuilder {
	if !m.verifySupportedTypes(fromInterfaces, isSupportedBindInterfaceReflectType) {
		return newNoOpBuilder()
	}
	return m.bind(newBindingKey, fromInterfaces)
}

func (m *module) BindTaggedInterface(tag string, fromInterfaces ...interface{}) InterfaceBuilder {
	if !m.verifyTag(tag) {
		return newNoOpBuilder()
	}
	if !m.verifySupportedTypes(fromInterfaces, isSupportedBindInterfaceReflectType) {
		return newNoOpBuilder()
	}
	return m.bind(func(fromReflectType reflect.Type) bindingKey { return newTaggedBindingKey(fromReflectType, tag) }, fromInterfaces)
}

func (m *module) BindTaggedBool(tag string) Builder {
	return m.bindTaggedConstant(tag, boolConstantKind)
}

func (m *module) BindTaggedInt(tag string) Builder {
	return m.bindTaggedConstant(tag, intConstantKind)
}

func (m *module) BindTaggedInt8(tag string) Builder {
	return m.bindTaggedConstant(tag, int8ConstantKind)
}

func (m *module) BindTaggedInt16(tag string) Builder {
	return m.bindTaggedConstant(tag, int16ConstantKind)
}

func (m *module) BindTaggedInt32(tag string) Builder {
	return m.bindTaggedConstant(tag, int32ConstantKind)
}

func (m *module) BindTaggedInt64(tag string) Builder {
	return m.bindTaggedConstant(tag, int64ConstantKind)
}

func (m *module) BindTaggedUint(tag string) Builder {
	return m.bindTaggedConstant(tag, uintConstantKind)
}

func (m *module) BindTaggedUint8(tag string) Builder {
	return m.bindTaggedConstant(tag, uint8ConstantKind)
}

func (m *module) BindTaggedUint16(tag string) Builder {
	return m.bindTaggedConstant(tag, uint16ConstantKind)
}

func (m *module) BindTaggedUint32(tag string) Builder {
	return m.bindTaggedConstant(tag, uint32ConstantKind)
}

func (m *module) BindTaggedUint64(tag string) Builder {
	return m.bindTaggedConstant(tag, uint64ConstantKind)
}

func (m *module) BindTaggedFloat32(tag string) Builder {
	return m.bindTaggedConstant(tag, float32ConstantKind)
}

func (m *module) BindTaggedFloat64(tag string) Builder {
	return m.bindTaggedConstant(tag, float64ConstantKind)
}

func (m *module) BindTaggedComplex64(tag string) Builder {
	return m.bindTaggedConstant(tag, complex64ConstantKind)
}

func (m *module) BindTaggedComplex128(tag string) Builder {
	return m.bindTaggedConstant(tag, complex128ConstantKind)
}

func (m *module) BindTaggedString(tag string) Builder {
	return m.bindTaggedConstant(tag, stringConstantKind)
}

func (m *module) bindTaggedConstant(tag string, constantKind constantKind) Builder {
	if !m.verifyTag(tag) {
		return newNoOpBuilder()
	}
	if !m.verifySupportedTypes([]interface{}{constantKind.constant()}, isSupportedBindConstantReflectType) {
		return newNoOpBuilder()
	}
	return m.bind(func(fromReflectType reflect.Type) bindingKey { return newTaggedBindingKey(fromReflectType, tag) }, []interface{}{constantKind.constant()})
}

func (m *module) bind(newBindingKeyFunc func(reflect.Type) bindingKey, from []interface{}) InterfaceBuilder {
	lenFrom := len(from)
	if lenFrom == 0 {
		m.addBindingError(newErrorBuilder(injectErrorTypeNil).build())
		return newNoOpBuilder()
	}
	bindingKeys := make([]bindingKey, lenFrom)
	for i := 0; i < lenFrom; i++ {
		fromReflectType := reflect.TypeOf(from[i])
		if fromReflectType == nil {
			m.addBindingError(newErrorBuilder(injectErrorTypeNil).build())
			return newNoOpBuilder()
		}
		bindingKeys[i] = newBindingKeyFunc(fromReflectType)
	}
	return newBuilder(m, bindingKeys)
}

func (m *module) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("module{")
	buffer.WriteString(strings.Join(m.keyValueStrings(), " "))
	buffer.WriteString("}")
	return buffer.String()
}

func (m *module) keyValueStrings() []string {
	strings := make([]string, len(m.bindings))
	i := 0
	for bindingKey, binding := range m.bindings {
		var buffer bytes.Buffer
		buffer.WriteString(bindingKey.String())
		buffer.WriteString(":")
		buffer.WriteString(binding.String())
		strings[i] = buffer.String()
		i++
	}
	return strings
}

func (m *module) addBindingError(err error) {
	m.bindingErrors = append(m.bindingErrors, err)
}

func (m *module) binding(bindingKey bindingKey) (binding, bool) {
	binding, ok := m.bindings[bindingKey]
	return binding, ok
}

func (m *module) setBinding(bindingKey bindingKey, binding binding) {
	foundBinding, ok := m.bindings[bindingKey]
	if ok {
		eb := newErrorBuilder(injectErrorTypeAlreadyBound)
		eb.addTag("bindingKey", bindingKey)
		eb.addTag("foundBinding", foundBinding)
		m.addBindingError(eb.build())
		return
	}
	m.bindings[bindingKey] = binding
}

func (m *module) verifyTag(tag string) bool {
	if tag == "" {
		m.addBindingError(newErrorBuilder(injectErrorTypeTagEmpty).build())
		return false
	}
	return true
}

func (m *module) verifySupportedTypes(froms []interface{}, isSupportedFunc func(reflect.Type) bool) bool {
	var ok bool
	for _, from := range froms {
		if !m.verifySupportedType(reflect.TypeOf(from), isSupportedFunc) {
			ok = false
		}
	}
	return ok
}

func (m *module) verifySupportedType(reflectType reflect.Type, isSupportedFunc func(reflect.Type) bool) bool {
	if !isSupportedFunc(reflectType) {
		m.addNotSupportedBindTypeError(reflectType)
		return false
	}
	return true
}

func (m *module) addNotSupportedBindTypeError(reflectType reflect.Type) {
	eb := newErrorBuilder(injectErrorTypeNotSupportedBindType)
	eb.addTag("reflectType", reflectType)
	m.addBindingError(eb.build())
}
