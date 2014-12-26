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

func (this *injector) GetTagged(from interface{}, tag string) (interface{}, error) {
	return this.get(newTaggedBindingKey(reflect.TypeOf(from), tag))
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
	returnValues := reflect.ValueOf(taggedFunction).Call([]reflect.Value{*populateTaggedFuncStruct(taggedFuncReflectType.In(0), reflectValues)})
	return reflectValuesToValues(returnValues), nil
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

func reflectValuesToValues(reflectValues []reflect.Value) []interface{} {
	lenReflectValues := len(reflectValues)
	values := make([]interface{}, lenReflectValues)
	for i := 0; i < lenReflectValues; i++ {
		values[i] = reflectValues[i].Interface()
	}
	return values
}
