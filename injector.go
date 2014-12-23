package inject

import (
	"reflect"
)

type injector struct {
	bindings map[bindingKey]binding
}

func createInjector(modules []Module) (Injector, error) {
	injector := injector{make(map[bindingKey]binding)}
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
	return &injector, nil
}

func installModuleToInjector(injector *injector, module *module) error {
	for bindingKey, binding := range module.bindings {
		if foundBinding, ok := injector.bindings[bindingKey]; ok {
			eb := newErrorBuilder(InjectErrorTypeAlreadyBound)
			eb.addTag("bindingKey", bindingKey)
			eb.addTag("foundBinding", foundBinding)
			return eb.build()
		}
		finalBinding, ok := getFinalBinding(injector, binding)
		if !ok {
			eb := newErrorBuilder(InjectErrorTypeNoFinalBinding)
			eb.addTag("bindingKey", bindingKey)
			return eb.build()
		}
		injector.bindings[bindingKey] = finalBinding
	}
	return nil
}

func getFinalBinding(injector *injector, binding binding) (binding, bool) {
	var ok bool
	for bindingKey, err := binding.bindingKey(); err != nil; bindingKey, err = binding.bindingKey() {
		binding, ok = injector.bindings[bindingKey]
		if !ok {
			return nil, false
		}
	}
	return binding, true
}

func (this *injector) Get(from interface{}) (interface{}, error) {
	return this.get(reflect.TypeOf(from))
}

func (this *injector) get(bindingType reflect.Type) (interface{}, error) {
	if bindingType == nil {
		return nil, newErrorBuilder(InjectErrorTypeReflectTypeNil).build()
	}
	binding, ok := this.bindings[newBindingKey(bindingType)]
	if !ok {
		eb := newErrorBuilder(InjectErrorTypeNoBinding)
		eb.addTag("bindingType", bindingType)
		return nil, eb.build()
	}
	return binding.get(this)
}

func (this *injector) GetTagged(from interface{}, tag string) (interface{}, error) {
	return this.getTagged(reflect.TypeOf(from), tag)
}

func (this *injector) getTagged(bindingType reflect.Type, tag string) (interface{}, error) {
	if bindingType == nil {
		return nil, newErrorBuilder(InjectErrorTypeReflectTypeNil).build()
	}
	if tag == "" {
		return nil, newErrorBuilder(InjectErrorTypeTagEmpty).build()
	}
	binding, ok := this.bindings[newTaggedBindingKey(bindingType, tag)]
	if !ok {
		eb := newErrorBuilder(InjectErrorTypeNoBinding)
		eb.addTag("bindingType", bindingType)
		eb.addTag("tag", tag)
		return nil, eb.build()
	}
	return binding.get(this)
}
