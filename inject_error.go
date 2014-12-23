package inject

import (
	"bytes"
	"fmt"
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
		buffer.WriteString(" tags:{ ")
		for key, value := range this.tags {
			buffer.WriteString(key)
			buffer.WriteString(":")
			if stringer, ok := value.(fmt.Stringer); ok {
				buffer.WriteString(fmt.Sprintf("%v", stringer.String()))
			} else {
				buffer.WriteString(fmt.Sprintf("%v", value))
			}
			buffer.WriteString(" ")
		}
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
