package inject

import (
	"bytes"
	"fmt"
	"strings"
)

const (
	injectErrorPrefix = "inject: "
)

type InjectError struct {
	errorType string
	tags      map[string]interface{}
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
	for key, value := range this.tags {
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
}

func newErrorBuilder(errorType string) *injectErrorBuilder {
	return &injectErrorBuilder{errorType, make(map[string]interface{})}
}

func (this *injectErrorBuilder) addTag(key string, value interface{}) *injectErrorBuilder {
	this.tags[key] = value
	return this
}

func (this *injectErrorBuilder) build() *InjectError {
	return &InjectError{this.errorType, this.tags}
}
