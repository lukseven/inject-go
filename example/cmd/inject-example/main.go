package main

import (
	"fmt"
	"os"

	"go.pedge.io/inject"
	"go.pedge.io/inject/example/api"
	"go.pedge.io/inject/example/cloud"
	"go.pedge.io/inject/example/more"
	"go.pedge.io/inject/example/stuff"
)

func main() {
	if err := do(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func do() error {
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
		return err
	}
	apiObj := obj.(api.Api)
	response, err := apiObj.Do(api.Request{Provider: os.Args[1], Foo: "this is fun"})
	if err != nil {
		return err
	}
	fmt.Printf("%v %v\n", response.Bar, response.Baz)
	return nil
}
