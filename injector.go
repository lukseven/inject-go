package inject

import (
	"reflect"
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
	for bindingKey, binding := range module.bindings {
		if foundBinding, ok := injector.bindings[bindingKey]; ok {
			eb := newErrorBuilder(InjectErrorTypeAlreadyBound)
			eb.addTag("bindingKey", bindingKey)
			eb.addTag("foundBinding", foundBinding)
			return eb.build()
		}
		resolvedBinding, err := binding.resolvedBinding(module)
		if err != nil {
			return err
		}
		injector.bindings[bindingKey] = resolvedBinding
	}
	return nil
}

func validate(injector *injector) error {
	for _, resolvedBinding := range injector.bindings {
		err := resolvedBinding.validate(injector)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *injector) Get(from interface{}) (interface{}, error) {
	return this.get(reflect.TypeOf(from))
}

func (this *injector) get(bindingType reflect.Type) (interface{}, error) {
	binding, err := this.getBinding(bindingType)
	if err != nil {
		return nil, err
	}
	return binding.get(this)
}

func (this *injector) GetTagged(from interface{}, tag string) (interface{}, error) {
	return this.getTagged(reflect.TypeOf(from), tag)
}

func (this *injector) getTagged(bindingType reflect.Type, tag string) (interface{}, error) {
	binding, err := this.getTaggedBinding(bindingType, tag)
	if err != nil {
		return nil, err
	}
	return binding.get(this)
}

func (this *injector) getBinding(bindingType reflect.Type) (resolvedBinding, error) {
	if bindingType == nil {
		return nil, newErrorBuilder(InjectErrorTypeReflectTypeNil).build()
	}
	binding, ok := this.bindings[newBindingKey(bindingType)]
	if !ok {
		eb := newErrorBuilder(InjectErrorTypeNoBinding)
		eb.addTag("bindingType", bindingType)
		return nil, eb.build()
	}
	return binding, nil
}

func (this *injector) getTaggedBinding(bindingType reflect.Type, tag string) (resolvedBinding, error) {
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
	return binding, nil
}
