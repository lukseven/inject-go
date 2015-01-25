package main

import (
	"fmt"
	"os"

	"github.com/peter-edge/inject"
	"github.com/peter-edge/inject/example/api"
	"github.com/peter-edge/inject/example/cloud"
	"github.com/peter-edge/inject/example/more"
	"github.com/peter-edge/inject/example/stuff"
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
	apiObj := obj.(api.Api)
	response, err := apiObj.Do(api.Request{os.Args[1], "this is fun"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v %v\n", response.Bar, response.Baz)
}
