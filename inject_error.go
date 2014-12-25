package inject

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	injectErrorPrefix                                 = "inject: "
	InjectErrorTypeNil                                = "Parameter is nil"
	InjectErrorTypeReflectTypeNil                     = "reflect.TypeOf() returns nil"
	InjectErrorTypeNotInterfacePtr                    = "Binding with Binder.ToType() and from is not an interface pointer"
	InjectErrorTypeDoesNotImplement                   = "to binding does not implement from binding"
	InjectErrorTypeNotSupportedYet                    = "Binding type not supported yet, feel free to help!"
	InjectErrorTypeNotAssignable                      = "Binding not assignable"
	InjectErrorTypeConstructorNotFunction             = "Constructor is not a function"
	InjectErrorTypeConstructorReturnValuesInvalid     = "Constructor can only have two return values, the first providing the value, the second being an error"
	InjectErrorTypeIntermediateBinding                = "Trying to get for an intermediate binding"
	InjectErrorTypeFinalBinding                       = "Trying to get bindingKey for a final binding"
	InjectErrorTypeCannotCastModule                   = "Cannot cast Module to internal module type"
	InjectErrorTypeNoBinding                          = "No binding for binding key"
	InjectErrorTypeNoFinalBinding                     = "No final binding for binding key"
	InjectErrorTypeAlreadyBound                       = "Already found a binding for this binding key"
	InjectErrorTypeTagEmpty                           = "Tag empty"
	InjectErrorTypeTaggedConstructorParametersInvalid = "Tagged constructor must have one anonymous struct parameter"
	InjectErrorTypeBindingErrors                      = "Errors with bindings"
)

type InjectError struct {
	errorType string
	tags      map[string]interface{}
	// TODO(pedge): there has to be a better way to do this
	tagOrder []string
}

func (this *InjectError) Error() string {
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

func (this *InjectError) Type() string {
	return this.errorType
}

func (this *InjectError) GetTag(key string) (interface{}, bool) {
	value, ok := this.tags[key]
	return value, ok
}

func (this *InjectError) tagStrings() []string {
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

func (this *injectErrorBuilder) build() *InjectError {
	return &InjectError{this.errorType, this.tags, this.tagOrder}
}
