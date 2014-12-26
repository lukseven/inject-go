package cloud

import (
	"gopkg.in/peter-edge/inject.v1"
	"gopkg.in/peter-edge/inject.v1/example/stuff"
)

func CreateModule() inject.Module {
	module := inject.CreateModule()
	module.BindTagged("aws", (*Provider)(nil)).ToSingletonConstructor(createAwsProvider)
	module.BindTagged("digitalOcean", (*Provider)(nil)).ToSingletonConstructor(createDigitalOceanProvider)
	return module
}

type Command struct {
	Path string
}

type Result struct {
	Message  string
	ExitCode int
}

type Instance interface {
	RunCommand(Command) (*Result, error)
}

type Provider interface {
	NewInstance() (Instance, error)
}

type instance struct {
	data         string
	stuffService stuff.StuffService
}

func (this *instance) RunCommand(command Command) (*Result, error) {
	i, err := this.stuffService.DoStuff("pwd")
	if err != nil {
		return nil, err
	}
	if command.Path == "ls" {
		return &Result{this.data, i}, nil
	} else {
		return &Result{this.data, 1}, nil
	}
}

type awsProvider struct {
	stuffService stuff.StuffService
}

func createAwsProvider(stuffService stuff.StuffService) (Provider, error) {
	return &awsProvider{stuffService}, nil
}

func (this *awsProvider) NewInstance() (Instance, error) {
	return &instance{"aws can do stuff", this.stuffService}, nil
}

type digitalOceanProvider struct {
	stuffService stuff.StuffService
}

func createDigitalOceanProvider(stuffService stuff.StuffService) (Provider, error) {
	return &digitalOceanProvider{stuffService}, nil
}

func (this *digitalOceanProvider) NewInstance() (Instance, error) {
	return &instance{"digitalOcean can also do stuff", this.stuffService}, nil
}
