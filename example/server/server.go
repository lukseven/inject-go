package server

import (
	"gopkg.in/peter-edge/inject.v1"
	"gopkg.in/peter-edge/inject.v1/example/api"
	"gopkg.in/peter-edge/inject.v1/example/cloud"
	"gopkg.in/peter-edge/inject.v1/example/more"
	"gopkg.in/peter-edge/inject.v1/example/stuff"
)

func Run() error {
	injector, err := inject.CreateInjector(
		api.CreateModule(),
		cloud.CreateModule(),
		more.CreateModule(),
		stuff.CreateModule(),
	)
	if err != nil {
		return err
	}
	obj, err := injector.Get((*api.Api)(nil))
	if err != nil {
		return nil
	}
	api := obj.(api.Api)
	fmt.Printf("I got my api! %v", api)
}
