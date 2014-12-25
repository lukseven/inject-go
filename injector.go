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
			return nil, newErrorBuilder(InjectErrorTypeCannotCastModule).build()
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
		eb := newErrorBuilder(InjectErrorTypeBindingErrors)
		for i := 0; i < numBindingErrors; i++ {
			eb.addTag(strconv.Itoa(i+1), module.bindingErrors[i].Error())
		}
		return eb.build()
	}
	for bindingKey, binding := range module.bindings {
		if foundBinding, ok := injector.bindings[bindingKey]; ok {
			eb := newErrorBuilder(InjectErrorTypeAlreadyBound)
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
		eb := newErrorBuilder(InjectErrorTypeNoBinding)
		eb.addTag("bindingKey", bindingKey)
		return nil, eb.build()
	}
	return binding, nil
}
