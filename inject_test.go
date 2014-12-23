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

func TestSimpleStructSingletonDirectSucceedsWhenPtr(t *testing.T) {
	injector := CreateInjector()
	err := injector.Bind((*SimpleInterface)(nil)).ToSingleton(&SimpleStruct{"hello"})
	requireErrNil(t, err)
	container, err := injector.CreateContainer()
	requireErrNil(t, err)
	object, err := container.Get((*SimpleInterface)(nil))
	requireErrNil(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimplePtrStructSingletonDirect(t *testing.T) {
	injector := CreateInjector()
	err := injector.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	requireErrNil(t, err)
	container, err := injector.CreateContainer()
	requireErrNil(t, err)
	object, err := container.Get((*SimpleInterface)(nil))
	requireErrNil(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimplePtrStructSingletonDirectFailsWhenNotPtr(t *testing.T) {
	injector := CreateInjector()
	err := injector.Bind((*SimpleInterface)(nil)).ToSingleton(SimplePtrStruct{"hello"})
	require.Equal(t, ErrDoesNotImplement, err)
}

func TestSimpleStructSingletonIndirect(t *testing.T) {
	injector := CreateInjector()
	err := injector.Bind((*SimpleInterface)(nil)).To(SimpleStruct{})
	requireErrNil(t, err)
	err = injector.Bind(SimpleStruct{}).ToSingleton(SimpleStruct{"hello"})
	requireErrNil(t, err)
	container, err := injector.CreateContainer()
	requireErrNil(t, err)
	object, err := container.Get((*SimpleInterface)(nil))
	requireErrNil(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimpleStructSingletonIndirectSucceedsWhenPtr(t *testing.T) {
	injector := CreateInjector()
	err := injector.Bind((*SimpleInterface)(nil)).To(&SimpleStruct{})
	requireErrNil(t, err)
	err = injector.Bind(&SimpleStruct{}).ToSingleton(&SimpleStruct{"hello"})
	requireErrNil(t, err)
	container, err := injector.CreateContainer()
	requireErrNil(t, err)
	object, err := container.Get((*SimpleInterface)(nil))
	requireErrNil(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimplePtrStructSingletonIndirect(t *testing.T) {
	injector := CreateInjector()
	err := injector.Bind((*SimpleInterface)(nil)).To(&SimplePtrStruct{})
	requireErrNil(t, err)
	err = injector.Bind(&SimplePtrStruct{}).ToSingleton(&SimplePtrStruct{"hello"})
	requireErrNil(t, err)
	container, err := injector.CreateContainer()
	requireErrNil(t, err)
	object, err := container.Get((*SimpleInterface)(nil))
	requireErrNil(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimplePtrStructSingletonIndirectFailsWhenNotPtr(t *testing.T) {
	injector := CreateInjector()
	err := injector.Bind((*SimpleInterface)(nil)).To(SimplePtrStruct{})
	require.Equal(t, ErrDoesNotImplement, err)
}

func TestSimplePtrStructIndirectFailsWhenNoFinalBinding(t *testing.T) {
	injector := CreateInjector()
	err := injector.Bind((*SimpleInterface)(nil)).To(&SimplePtrStruct{})
	requireErrNil(t, err)
	_, err = injector.CreateContainer()
	require.NotNil(t, err)
}

// ***** HELPERS *****

// TODO(pedge): is there something for this in testify? if not, send a pull request
func requireErrNil(t *testing.T, err error) {
	if err != nil {
		require.Nil(t, err, err.Error())
	}
}
