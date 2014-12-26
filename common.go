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

func verifyIsStructPtr(reflectType reflect.Type) error {
	if !isStructPtr(reflectType) {
		eb := newErrorBuilder(injectErrorTypeNotStructPtr)
		eb.addTag("reflectType", reflectType)
		return eb.build()
	}
	return nil
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
	return getStructFieldBindingKeys(funcReflectType.In(0))
}

func getStructFieldBindingKeys(structReflectType reflect.Type) []bindingKey {
	numFields := structReflectType.NumField()
	bindingKeys := make([]bindingKey, numFields)
	for i := 0; i < numFields; i++ {
		structField := structReflectType.Field(i)
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

func getTaggedFuncStructReflectValue(structReflectType reflect.Type, reflectValues []reflect.Value) *reflect.Value {
	structReflectValue := reflect.Indirect(reflect.New(structReflectType))
	populateStructReflectValue(&structReflectValue, reflectValues)
	return &structReflectValue
}

func newStructReflectValue(structReflectType reflect.Type) reflect.Value {
	return reflect.Indirect(reflect.New(structReflectType))
}

func populateStructReflectValue(structReflectValue *reflect.Value, reflectValues []reflect.Value) {
	numReflectValues := len(reflectValues)
	for i := 0; i < numReflectValues; i++ {
		structReflectValue.Field(i).Set(reflectValues[i])
	}
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
