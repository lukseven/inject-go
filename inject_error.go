package inject

import (
	"bytes"
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

func (this *injectError) Error() string {
	var buffer bytes.Buffer
	buffer.WriteString(injectErrorPrefix)
	buffer.WriteString(this.errorType)
	if len(this.tags) > 0 {
		buffer.WriteString(" tags{")
		buffer.WriteString(strings.Join(this.tagStrings(), " "))
		buffer.WriteString("}")
	}
	return buffer.String()
}

func (this *injectError) Type() string {
	return this.errorType
}

func (this *injectError) GetTag(key string) (interface{}, bool) {
	value, ok := this.tags[key]
	return value, ok
}

func (this *injectError) tagStrings() []string {
	strings := make([]string, len(this.tags))
	i := 0
	for _, key := range this.tagOrder {
		value := this.tags[key]
		var buffer bytes.Buffer
		buffer.WriteString(key)
		buffer.WriteString(":")
		if stringer, ok := value.(fmt.Stringer); ok {
			buffer.WriteString(fmt.Sprintf("%v", stringer.String()))
		} else {
			buffer.WriteString(fmt.Sprintf("%v", value))
		}
		strings[i] = buffer.String()
		i++
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

func (this *injectErrorBuilder) addTag(key string, value interface{}) *injectErrorBuilder {
	this.tags[key] = value
	this.tagOrder = append(this.tagOrder, key)
	return this
}

func (this *injectErrorBuilder) build() *injectError {
	return &injectError{this.errorType, this.tags, this.tagOrder}
}
