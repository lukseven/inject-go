package inject

import (
	"reflect"
)

const (
	taggedFuncStructFieldTag = "injectTag"
)

func verifyIsFunc(funcReflectType reflect.Type) error {
	if !isFunc(funcReflectType) {
		eb := newErrorBuilder(injectErrorTypeNotFunction)
		eb = eb.addTag("funcReflectType", funcReflectType)
		return eb.build()
	}
	numIn := funcReflectType.NumIn()
	for i := 0; i < numIn; i++ {
		err := verifyParameterCanBeInjected(funcReflectType.In(i))
		if err != nil {
			return err
		}
	}
	return nil
}

func verifyIsTaggedFunc(funcReflectType reflect.Type) error {
	if !isFunc(funcReflectType) {
		eb := newErrorBuilder(injectErrorTypeNotFunction)
		eb = eb.addTag("funcReflectType", funcReflectType)
		return eb.build()
	}
	if funcReflectType.NumIn() != 1 {
		eb := newErrorBuilder(injectErrorTypeTaggedParametersInvalid)
		eb = eb.addTag("funcReflectType", funcReflectType)
		return eb.build()
	}
	inReflectType := funcReflectType.In(0)
	if !isStruct(inReflectType) {
		eb := newErrorBuilder(injectErrorTypeTaggedParametersInvalid)
		eb = eb.addTag("funcReflectType", funcReflectType)
		return eb.build()
	}
	if inReflectType.Name() != "" {
		eb := newErrorBuilder(injectErrorTypeTaggedParametersInvalid)
		eb = eb.addTag("funcReflectType", funcReflectType)
		return eb.build()
	}
	numFields := inReflectType.NumField()
	for i := 0; i < numFields; i++ {
		err := verifyParameterCanBeInjected(inReflectType.Field(i).Type)
		if err != nil {
			return err
		}
	}
	return nil
}

func verifyParameterCanBeInjected(parameterReflectType reflect.Type) error {
	switch {
	case isInterface(parameterReflectType), isStructPtr(parameterReflectType), isStruct(parameterReflectType):
		return nil
	default:
		eb := newErrorBuilder(injectErrorTypeNotSupportedYet)
		eb.addTag("parameterReflectType", parameterReflectType)
		return eb.build()
	}
}

func getParameterBindingKeysForFunc(funcReflectType reflect.Type) []bindingKey {
	numIn := funcReflectType.NumIn()
	bindingKeys := make([]bindingKey, numIn)
	for i := 0; i < numIn; i++ {
		inReflectType := funcReflectType.In(i)
		// TODO(pedge): this is really specific logic, and there wil need to be more
		// of this if more types are allowed for binding - this should be abstracted
		if inReflectType.Kind() == reflect.Interface {
			inReflectType = reflect.PtrTo(inReflectType)
		}
		bindingKeys[i] = newBindingKey(inReflectType)
	}
	return bindingKeys
}

func getParameterBindingKeysForTaggedFunc(funcReflectType reflect.Type) []bindingKey {
	inReflectType := funcReflectType.In(0)
	numFields := inReflectType.NumField()
	bindingKeys := make([]bindingKey, numFields)
	for i := 0; i < numFields; i++ {
		structField := inReflectType.Field(i)
		structFieldReflectType := structField.Type
		// TODO(pedge): this is really specific logic, and there wil need to be more
		// of this if more types are allowed for binding - this should be abstracted
		if structFieldReflectType.Kind() == reflect.Interface {
			structFieldReflectType = reflect.PtrTo(structFieldReflectType)
		}
		tag := structField.Tag.Get(taggedFuncStructFieldTag)
		if tag != "" {
			bindingKeys[i] = newTaggedBindingKey(structFieldReflectType, tag)
		} else {
			bindingKeys[i] = newBindingKey(structFieldReflectType)
		}
	}
	return bindingKeys
}

func populateTaggedFuncStruct(structReflectType reflect.Type, reflectValues []reflect.Value) *reflect.Value {
	numReflectValues := len(reflectValues)
	value := reflect.Indirect(reflect.New(structReflectType))
	for i := 0; i < numReflectValues; i++ {
		value.Field(i).Set(reflectValues[i])
	}
	return &value
}

func isInterfacePtr(reflectType reflect.Type) bool {
	return isPtr(reflectType) && isInterface(reflectType.Elem())
}

func isStructPtr(reflectType reflect.Type) bool {
	return isPtr(reflectType) && isStruct(reflectType.Elem())
}

func isInterface(reflectType reflect.Type) bool {
	return reflectType.Kind() == reflect.Interface
}

func isStruct(reflectType reflect.Type) bool {
	return reflectType.Kind() == reflect.Struct
}

func isPtr(reflectType reflect.Type) bool {
	return reflectType.Kind() == reflect.Ptr
}

func isFunc(reflectType reflect.Type) bool {
	return reflectType.Kind() == reflect.Func
}
