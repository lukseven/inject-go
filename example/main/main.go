package main

import (
	"fmt"
	"os"

	"github.com/peter-edge/go-inject"
	"github.com/peter-edge/go-inject/example/api"
	"github.com/peter-edge/go-inject/example/cloud"
	"github.com/peter-edge/go-inject/example/more"
	"github.com/peter-edge/go-inject/example/stuff"
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
