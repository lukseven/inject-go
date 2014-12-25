package more

import (
	"fmt"
	"gopkg.in/peter-edge/inject.v1"
)

func CreateModule() inject.Module {
	module := inject.CreateModule()
	module.Bind((*MoreThings)(nil)).ToSingleton(&moreThings{})
	return module
}

type MoreThings interface {
	MoreStuffToDo(int) (string, error)
}

type moreThings struct{}

func (this *moreThings) MoreStuffToDo(i int) (string, error) {
	return fmt.Sprintf("but there's not much to do here %v", i), nil
}
