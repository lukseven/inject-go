package api

import (
	"bytes"
	"fmt"
	"gopkg.in/peter-edge/inject.v1"
	"gopkg.in/peter-edge/inject.v1/example/cloud"
	"gopkg.in/peter-edge/inject.v1/example/more"
)

func CreateModule() inject.Module {
	module := inject.CreateModule()
	module.Bind((*Api)(nil)).ToTaggedSingletonConstructor(createApi)
	return module
}

type Request struct {
	Provider string
	Foo      string
}

type Response struct {
	Bar string
	Baz int
}

type Api interface {
	Do(Request) (*Response, error)
}

type api struct {
	awsProvider          cloud.Provider
	digitalOceanProvider cloud.Provider
	moreThings           more.MoreThings
}

func createApi(s struct {
	awsProvider          cloud.Provider `inject:"aws"`
	digitalOceanProvider cloud.Provider `inject:"digitalOcean"`
	moreThings           more.MoreThings
}) (Api, error) {
	return &api{s.awsProvider, s.digitalOceanProvider, s.moreThings}, nil
}

func (this *api) Do(request Request) (*Response, error) {
	provider, err := this.getProvider(request.Provider)
	if err != nil {
		return nil, err
	}
	instance, err := provider.NewInstance()
	if err != nil {
		return nil, err
	}
	result, err := instance.RunCommand(cloud.Command{"ls"})
	if err != nil {
		return nil, err
	}
	s, err := this.moreThings.MoreStuffToDo(1)
	if err != nil {
		return nil, err
	}
	var buffer bytes.Buffer
	buffer.WriteString(request.Foo)
	buffer.WriteString(" ")
	buffer.WriteString(result.Message)
	buffer.WriteString(" ")
	buffer.WriteString(s)
	return &Response{buffer.String(), result.ExitCode}, nil
}

func (this *api) getProvider(provider string) (cloud.Provider, error) {
	switch provider {
	case "aws":
		return this.awsProvider, nil
	case "digitalOcean":
		return this.digitalOceanProvider, nil
	default:
		return nil, fmt.Errorf("api: Unknown provider %v", provider)
	}
}
