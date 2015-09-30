package stuff // import "go.pedge.io/inject/example/stuff"

import (
	"go.pedge.io/inject"
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

func (s *stuffService) DoStuff(ss string) (int, error) {
	if ss == "pwd" {
		return 0, nil
	} else {
		return -1, nil
	}
}
