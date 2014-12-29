package inject

import (
	"bytes"
	"reflect"
	"strconv"
	"strings"
)

type injector struct {
	bindings map[bindingKey]resolvedBinding
}

func createInjector(modules []Module) (Injector, error) {
	injector := injector{make(map[bindingKey]resolvedBinding)}
	for _, m := range modules {
		castModule, ok := m.(*module)
		if !ok {
			return nil, newErrorBuilder(injectErrorTypeCannotCastModule).build()
		}
		err := installModuleToInjector(&injector, castModule)
		if err != nil {
			return nil, err
		}
	}
	err := validate(&injector)
	if err != nil {
		return nil, err
	}
	return &injector, nil
}

func installModuleToInjector(injector *injector, module *module) error {
	numBindingErrors := len(module.bindingErrors)
	if numBindingErrors > 0 {
		eb := newErrorBuilder(injectErrorTypeBindingErrors)
		for i := 0; i < numBindingErrors; i++ {
			eb.addTag(strconv.Itoa(i+1), module.bindingErrors[i].Error())
		}
		return eb.build()
	}
	for bindingKey, binding := range module.bindings {
		if foundBinding, ok := injector.bindings[bindingKey]; ok {
			eb := newErrorBuilder(injectErrorTypeAlreadyBound)
			eb.addTag("bindingKey", bindingKey)
			eb.addTag("foundBinding", foundBinding)
			return eb.build()
		}
		resolvedBinding, err := binding.resolvedBinding(module, injector)
		if err != nil {
			return err
		}
		injector.bindings[bindingKey] = resolvedBinding
	}
	return nil
}

func validate(injector *injector) error {
	for _, resolvedBinding := range injector.bindings {
		err := resolvedBinding.validate()
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *injector) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("injector{")
	buffer.WriteString(strings.Join(this.keyValueStrings(), " "))
	buffer.WriteString("}")
	return buffer.String()
}

func (this *injector) keyValueStrings() []string {
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

func (this *injector) Get(from interface{}) (interface{}, error) {
	return this.get(newBindingKey(reflect.TypeOf(from)))
}

func (this *injector) GetTagged(tag string, from interface{}) (interface{}, error) {
	return this.get(newTaggedBindingKey(reflect.TypeOf(from), tag))
}

func (this *injector) GetTaggedBool(tag string) (bool, error) {
	obj, err := this.getTaggedConstant(tag, boolConstantKind)
	if err != nil {
		return boolConstant, err
	}
	return obj.(bool), nil
}

func (this *injector) GetTaggedInt(tag string) (int, error) {
	obj, err := this.getTaggedConstant(tag, intConstantKind)
	if err != nil {
		return intConstant, err
	}
	return obj.(int), nil
}

func (this *injector) GetTaggedInt8(tag string) (int8, error) {
	obj, err := this.getTaggedConstant(tag, int8ConstantKind)
	if err != nil {
		return int8Constant, err
	}
	return obj.(int8), nil
}

func (this *injector) GetTaggedInt16(tag string) (int16, error) {
	obj, err := this.getTaggedConstant(tag, int16ConstantKind)
	if err != nil {
		return int16Constant, err
	}
	return obj.(int16), nil
}

func (this *injector) GetTaggedInt32(tag string) (int32, error) {
	obj, err := this.getTaggedConstant(tag, int32ConstantKind)
	if err != nil {
		return int32Constant, err
	}
	return obj.(int32), nil
}

func (this *injector) GetTaggedInt64(tag string) (int64, error) {
	obj, err := this.getTaggedConstant(tag, int64ConstantKind)
	if err != nil {
		return int64Constant, err
	}
	return obj.(int64), nil
}

func (this *injector) GetTaggedUint(tag string) (uint, error) {
	obj, err := this.getTaggedConstant(tag, uintConstantKind)
	if err != nil {
		return uintConstant, err
	}
	return obj.(uint), nil
}

func (this *injector) GetTaggedUint8(tag string) (uint8, error) {
	obj, err := this.getTaggedConstant(tag, uint8ConstantKind)
	if err != nil {
		return uint8Constant, err
	}
	return obj.(uint8), nil
}

func (this *injector) GetTaggedUint16(tag string) (uint16, error) {
	obj, err := this.getTaggedConstant(tag, uint16ConstantKind)
	if err != nil {
		return uint16Constant, err
	}
	return obj.(uint16), nil
}

func (this *injector) GetTaggedUint32(tag string) (uint32, error) {
	obj, err := this.getTaggedConstant(tag, uint32ConstantKind)
	if err != nil {
		return uint32Constant, err
	}
	return obj.(uint32), nil
}

func (this *injector) GetTaggedUint64(tag string) (uint64, error) {
	obj, err := this.getTaggedConstant(tag, uint64ConstantKind)
	if err != nil {
		return uint64Constant, err
	}
	return obj.(uint64), nil
}

func (this *injector) GetTaggedFloat32(tag string) (float32, error) {
	obj, err := this.getTaggedConstant(tag, float32ConstantKind)
	if err != nil {
		return float32Constant, err
	}
	return obj.(float32), nil
}

func (this *injector) GetTaggedFloat64(tag string) (float64, error) {
	obj, err := this.getTaggedConstant(tag, float64ConstantKind)
	if err != nil {
		return float64Constant, err
	}
	return obj.(float64), nil
}

func (this *injector) GetTaggedComplex64(tag string) (complex64, error) {
	obj, err := this.getTaggedConstant(tag, complex64ConstantKind)
	if err != nil {
		return complex64Constant, err
	}
	return obj.(complex64), nil
}

func (this *injector) GetTaggedComplex128(tag string) (complex128, error) {
	obj, err := this.getTaggedConstant(tag, complex128ConstantKind)
	if err != nil {
		return complex128Constant, err
	}
	return obj.(complex128), nil
}

func (this *injector) GetTaggedString(tag string) (string, error) {
	obj, err := this.getTaggedConstant(tag, stringConstantKind)
	if err != nil {
		return stringConstant, err
	}
	return obj.(string), nil
}

func (this *injector) getTaggedConstant(tag string, constantKind constantKind) (interface{}, error) {
	return this.get(newTaggedBindingKey(constantKind.reflectType(), tag))
}

func (this *injector) Call(function interface{}) ([]interface{}, error) {
	funcReflectType := reflect.TypeOf(function)
	err := verifyIsFunc(funcReflectType)
	if err != nil {
		return nil, err
	}
	bindingKeys := getParameterBindingKeysForFunc(funcReflectType)
	err = this.validateBindingKeys(bindingKeys)
	if err != nil {
		return nil, err
	}
	reflectValues, err := this.getReflectValues(bindingKeys)
	if err != nil {
		return nil, err
	}
	returnValues := reflect.ValueOf(function).Call(reflectValues)
	return reflectValuesToValues(returnValues), nil
}

func (this *injector) CallTagged(taggedFunction interface{}) ([]interface{}, error) {
	taggedFuncReflectType := reflect.TypeOf(taggedFunction)
	err := verifyIsTaggedFunc(taggedFuncReflectType)
	if err != nil {
		return nil, err
	}
	bindingKeys := getParameterBindingKeysForTaggedFunc(taggedFuncReflectType)
	err = this.validateBindingKeys(bindingKeys)
	if err != nil {
		return nil, err
	}
	reflectValues, err := this.getReflectValues(bindingKeys)
	if err != nil {
		return nil, err
	}
	structReflectValue := newStructReflectValue(taggedFuncReflectType.In(0))
	populateStructReflectValue(&structReflectValue, reflectValues)
	returnValues := reflect.ValueOf(taggedFunction).Call([]reflect.Value{structReflectValue})
	return reflectValuesToValues(returnValues), nil
}

func (this *injector) Populate(populateStructPtr interface{}) error {
	populateStructPtrReflectType := reflect.TypeOf(populateStructPtr)
	err := verifyIsStructPtr(populateStructPtrReflectType)
	if err != nil {
		return err
	}
	populateStructValue := reflect.Indirect(reflect.ValueOf(populateStructPtr))
	bindingKeys := getStructFieldBindingKeys(populateStructValue.Type())
	err = this.validateBindingKeys(bindingKeys)
	if err != nil {
		return err
	}
	reflectValues, err := this.getReflectValues(bindingKeys)
	if err != nil {
		return err
	}
	populateStructReflectValue(&populateStructValue, reflectValues)
	return nil
}

func (this *injector) get(bindingKey bindingKey) (interface{}, error) {
	binding, err := this.getBinding(bindingKey)
	if err != nil {
		return nil, err
	}
	return binding.get()
}

func (this *injector) getBinding(bindingKey bindingKey) (resolvedBinding, error) {
	binding, ok := this.bindings[bindingKey]
	if !ok {
		eb := newErrorBuilder(injectErrorTypeNoBinding)
		eb.addTag("bindingKey", bindingKey)
		return nil, eb.build()
	}
	return binding, nil
}

func (this *injector) getReflectValues(bindingKeys []bindingKey) ([]reflect.Value, error) {
	numBindingKeys := len(bindingKeys)
	reflectValues := make([]reflect.Value, numBindingKeys)
	for i := 0; i < numBindingKeys; i++ {
		value, err := this.get(bindingKeys[i])
		if err != nil {
			return nil, err
		}
		reflectValues[i] = reflect.ValueOf(value)
	}
	return reflectValues, nil
}

func (this *injector) validateBindingKeys(bindingKeys []bindingKey) error {
	for _, bindingKey := range bindingKeys {
		_, err := this.getBinding(bindingKey)
		if err != nil {
			return err
		}
	}
	return nil
}

func verifyIsStructPtr(reflectType reflect.Type) error {
	if !isStructPtr(reflectType) {
		eb := newErrorBuilder(injectErrorTypeNotStructPtr)
		eb.addTag("reflectType", reflectType)
		return eb.build()
	}
	return nil
}

func reflectValuesToValues(reflectValues []reflect.Value) []interface{} {
	lenReflectValues := len(reflectValues)
	values := make([]interface{}, lenReflectValues)
	for i := 0; i < lenReflectValues; i++ {
		values[i] = reflectValues[i].Interface()
	}
	return values
}
