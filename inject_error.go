package inject

import (
	"fmt"
	"strings"
)

const (
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

var (
	errNil                            = newInjectError(injectErrorTypeNil)
	errReflectTypeNil                 = newInjectError(injectErrorTypeReflectTypeNil)
	errNotSupportedYet                = newInjectError(injectErrorTypeNotSupportedYet)
	errNotAssignable                  = newInjectError(injectErrorTypeNotAssignable)
	errConstructorReturnValuesInvalid = newInjectError(injectErrorTypeConstructorReturnValuesInvalid)
	errIntermediateBinding            = newInjectError(injectErrorTypeIntermediateBinding)
	errFinalBinding                   = newInjectError(injectErrorTypeFinalBinding)
	errCannotCastModule               = newInjectError(injectErrorTypeCannotCastModule)
	errNoBinding                      = newInjectError(injectErrorTypeNoBinding)
	errNoFinalBinding                 = newInjectError(injectErrorTypeNoFinalBinding)
	errAlreadyBound                   = newInjectError(injectErrorTypeAlreadyBound)
	errTagEmpty                       = newInjectError(injectErrorTypeTagEmpty)
	errTaggedParametersInvalid        = newInjectError(injectErrorTypeTaggedParametersInvalid)
	errNotFunction                    = newInjectError(injectErrorTypeNotFunction)
	errNotInterfacePtr                = newInjectError(injectErrorTypeNotInterfacePtr)
	errNotStructPtr                   = newInjectError(injectErrorTypeNotStructPtr)
	errNotSupportedBindType           = newInjectError(injectErrorTypeNotSupportedBindType)
	errBindingErrors                  = newInjectError(injectErrorTypeBindingErrors)
)

type injectError struct {
	errorType string
	tags      injectErrorTags
}

func newInjectError(errorType string) *injectError {
	return &injectError{errorType, make([]*injectErrorTag, 0)}
}

func (i *injectError) Error() string {
	value := fmt.Sprintf("inject: %s", i.errorType)
	if len(i.tags) == 0 {
		return value
	}
	return fmt.Sprintf("%s %s", value, i.tags.String())
}

func (i *injectError) withTag(key string, value interface{}) *injectError {
	return &injectError{i.errorType, append(i.tags, newInjectErrorTag(key, value))}
}

type injectErrorTag struct {
	key   string
	value interface{}
}

func newInjectErrorTag(key string, value interface{}) *injectErrorTag {
	return &injectErrorTag{key, value}
}

func (t *injectErrorTag) String() string {
	if stringer, ok := t.value.(fmt.Stringer); ok {
		return fmt.Sprintf("%s:%s", t.key, stringer.String())
	}
	return fmt.Sprintf("%s:%s", t.key, t.value)
}

type injectErrorTags []*injectErrorTag

func (ts injectErrorTags) String() string {
	if len(ts) == 0 {
		return ""
	}
	s := make([]string, len(ts))
	for i, tag := range ts {
		s[i] = tag.String()
	}
	return fmt.Sprintf("tags{%s}", strings.Join(s, " "))
}
