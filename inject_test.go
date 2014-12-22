package inject

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type SimpleInterface interface {
	Foo() string
}

type SimpleStruct struct {
	foo string
}

func (this SimpleStruct) Foo() string {
	return this.foo
}

type SimplePtrStruct struct {
	foo string
}

func (this *SimplePtrStruct) Foo() string {
	return this.foo
}

func TestSimpleStructSingletonDirect(t *testing.T) {
	injector := CreateInjector()
	err := injector.Bind((*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"hello"})
	requireErrNil(t, err)
	container, err := injector.CreateContainer()
	requireErrNil(t, err)
	object, err := container.Get((*SimpleInterface)(nil))
	requireErrNil(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

// TODO(pedge): is there something for this in testify? if not, send a pull request
func requireErrNil(t *testing.T, err error) {
	if err != nil {
		require.Nil(t, err, err.Error())
	}
}
