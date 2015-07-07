package inject

import (
	"fmt"
	"strings"
)

const (
	injectErrorPrefix                             = "inject: "
	injectErrorTypeNil                            = "Parameter is nil"
	injectErrorTypeReflectTypeNil                 = "reflect.TypeOf() returns nil"
	injectErrorTypeNotSupportedYet                = "Binding type not supported yet, feel free to help!"
	injectErrorTypeNotAssignable                  = "Binding not assignable"
	injectErrorTypeConstructorReturnValuesInvalid = "Constructor can only have two return values, the first providing the value, the second being an error"
	injectErrorTypeIntermediateBinding            = "Trying to get for an intermediate binding"
	injectErrorTypeFinalBinding                   = "Trying to get bindingKey for a final binding"
	injectErrorTypeCannotCastModule               = "Cannot cast Module to internal module type"
	injectErrorTypeNoBinding                      = "No binding for binding key"
	injectErrorTypeNoFinalBinding                 = "No final binding for binding key"
	injectErrorTypeAlreadyBound                   = "Already found a binding for this binding key"
	injectErrorTypeTagEmpty                       = "Tag empty"
	injectErrorTypeTaggedParametersInvalid        = "Tagged function must have one anonymous struct parameter"
	injectErrorTypeNotFunction                    = "Argument is not a function"
	injectErrorTypeNotInterfacePtr                = "Value is not an interface pointer"
	injectErrorTypeNotStructPtr                   = "Value is not a struct pointer"
	injectErrorTypeNotSupportedBindType           = "Type is not supported for this binding method"
	injectErrorTypeBindingErrors                  = "Errors with bindings"
)

type injectError struct {
	errorType string
	tags      map[string]interface{}
	// TODO(pedge): there has to be a better way to do this
	tagOrder []string
}

func (i *injectError) Error() string {
	value := fmt.Sprintf("%s%s", injectErrorPrefix, i.errorType)
	tagStrings := i.tagStrings()
	if len(tagStrings) > 0 {
		value = fmt.Sprintf("%s tags{%s}", value, strings.Join(tagStrings, " "))
	}
	return value
}

func (i *injectError) Type() string {
	return i.errorType
}

func (i *injectError) GetTag(key string) (interface{}, bool) {
	value, ok := i.tags[key]
	return value, ok
}

func (i *injectError) tagStrings() []string {
	strings := make([]string, len(i.tags))
	ii := 0
	for _, key := range i.tagOrder {
		value := i.tags[key]
		var valueString string
		if stringer, ok := value.(fmt.Stringer); ok {
			valueString = fmt.Sprintf("%v", stringer.String())
		} else {
			valueString = fmt.Sprintf("%v", value)
		}
		strings[ii] = fmt.Sprintf("%s:%s", key, valueString)
		ii++
	}
	return strings
}

type injectErrorBuilder struct {
	errorType string
	tags      map[string]interface{}
	tagOrder  []string
}

func newErrorBuilder(errorType string) *injectErrorBuilder {
	return &injectErrorBuilder{errorType, make(map[string]interface{}), make([]string, 0)}
}

func (i *injectErrorBuilder) addTag(key string, value interface{}) *injectErrorBuilder {
	i.tags[key] = value
	i.tagOrder = append(i.tagOrder, key)
	return i
}

func (i *injectErrorBuilder) build() *injectError {
	return &injectError{i.errorType, i.tags, i.tagOrder}
}
