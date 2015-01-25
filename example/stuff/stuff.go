package stuff

import (
	"github.com/peter-edge/inject"
)

func CreateModule() inject.Module {
	module := inject.CreateModule()
	module.Bind((*StuffService)(nil)).ToConstructor(createStuffService)
	return module
}

type StuffService interface {
	DoStuff(string) (int, error)
}

type stuffService struct{}

func createStuffService() (StuffService, error) {
	return &stuffService{}, nil
}

func (this *stuffService) DoStuff(s string) (int, error) {
	if s == "pwd" {
		return 0, nil
	} else {
		return -1, nil
	}
}
