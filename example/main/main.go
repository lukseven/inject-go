package main

import (
	"gopkg.in/peter-edge/inject.v1"
	"gopkg.in/peter-edge/inject.v1/example/api"
	"gopkg.in/peter-edge/inject.v1/example/cloud"
	"gopkg.in/peter-edge/inject.v1/example/more"
	"gopkg.in/peter-edge/inject.v1/example/stuff"
	"os"
)

func main() {
	injector, err := inject.CreateInjector(
		api.CreateModule(),
		cloud.CreateModule(),
		more.CreateModule(),
		stuff.CreateModule(),
	)
	if err != nil {
		panic(err)
	}
	obj, err := injector.Get((*api.Api)(nil))
	if err != nil {
		panic(err)
	}
	api := obj.(api.Api)
	response, err := api.Do(api.Request{os.Args[1], "this is fun"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v %v\n", response.Bar, response.Baz)
}
