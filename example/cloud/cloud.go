package cloud

import (
	"errors"
	"gopkg.in/peter-edge/inject.v1"
	"gopkg.in/peter-edge/inject.v1/example/stuff"
)

func CreateModule() inject.Module {
	module := inject.CreateModule()
	module.BindTagged((*Provider)(nil), "aws").ToSingletonConstructor(createAwsProvider)
	module.BindTagged((*Provider)(nil), "digitalOcean").ToSingletonConstructor(createDigitalOceanProvider)
	return module
}

type Command struct {
	path string
}

type Result struct {
	exitCode int
}

type Instance interface {
	RunCommand(Command) (*Result, error)
}

type Provider interface {
	NewInstance() (Instance, error)
}

type awsProvider struct {
	stuffService stuff.StuffService
}

func createAwsProvider(stuffService stuff.StuffService) (Provider, error) {
	return &awsProvider{stuffService}
}

func (this *awsProvider) NewInstance() (Instance, error) {
	return nil, errors.New("Not implemented")
}

type digitalOceanProvider struct {
	stuffService stuff.StuffService
}

func createDigitalOceanProvider(stuffService stuff.StuffService) (Provider, error) {
	return &digitalPceanProvider{stuffService}
}

func (this *digitalOceanProvider) NewInstance() (Instance, error) {
	return nil, errors.New("Not implemented")
}
