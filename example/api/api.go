package api

import (
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
	foo string
}

type Response struct {
	bar string
}

type Api interface {
	DoStuff(Request) (*Response, error)
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
	return &api{s.awsProvider, s.digitalOceanProvider, s.moreThings}
}

func (this *api) DoStuff(request Request) (*Response, error) {
	return nil, nil
}
