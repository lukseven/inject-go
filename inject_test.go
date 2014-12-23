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

type SimpleTag struct{}

// ***** simple Bind tests *****

func TestFromNil(t *testing.T) {
	injector := CreateInjector()
	err := injector.Bind(nil).ToSingleton(SimpleStruct{"hello"})
	require.Equal(t, ErrNil, err)
}

// TODO(pedge): this test is still useful to make sure if someone does not
// provide a pointer to an interface, it fails, but this is not the correct
// behavior, find a value that is not nil but reflect.TypeOf(...) returns nil
func TestFromReflectTypeNil(t *testing.T) {
	injector := CreateInjector()
	err := injector.Bind((SimpleInterface)(nil)).ToSingleton(SimpleStruct{"hello"})
	require.Equal(t, ErrNil, err)
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
	requireErrMsgContains(t, err, noBindingToSingletonOrProviderMsg)
}

// ***** simple BindTagged tests *****

func TestTaggedFromNil(t *testing.T) {
	injector := CreateInjector()
	err := injector.BindTagged(nil, "tagOne").ToSingleton(SimpleStruct{"hello"})
	require.Equal(t, ErrNil, err)
}

// TODO(pedge): this test is still useful to make sure if someone does not
// provide a pointer to an interface, it fails, but this is not the correct
// behavior, find a value that is not nil but reflect.TypeOf(...) returns nil
func TestTaggedFromReflectTypeNil(t *testing.T) {
	injector := CreateInjector()
	err := injector.BindTagged((SimpleInterface)(nil), "tagOne").ToSingleton(SimpleStruct{"hello"})
	require.Equal(t, ErrNil, err)
}

func TestTaggedTagNil(t *testing.T) {
	injector := CreateInjector()
	err := injector.BindTagged((*SimpleInterface)(nil), nil).ToSingleton(SimpleStruct{"hello"})
	require.Equal(t, ErrNil, err)
}

func TestTaggedSimpleStructSingletonDirect(t *testing.T) {
	injector := CreateInjector()
	err := injector.BindTagged((*SimpleInterface)(nil), "tagOne").ToSingleton(SimpleStruct{"hello"})
	requireErrNil(t, err)
	container, err := injector.CreateContainer()
	requireErrNil(t, err)
	object, err := container.GetTagged((*SimpleInterface)(nil), "tagOne")
	requireErrNil(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	_, err = container.Get((*SimpleInterface)(nil))
	requireErrMsgContains(t, err, noBindingMsg)
	_, err = container.GetTagged((*SimpleInterface)(nil), "tagTwo")
	requireErrMsgContains(t, err, noBindingMsg)
}

func TestTaggedSimpleStructSingletonDirectSucceedsWhenPtr(t *testing.T) {
	injector := CreateInjector()
	err := injector.BindTagged((*SimpleInterface)(nil), "tagOne").ToSingleton(&SimpleStruct{"hello"})
	requireErrNil(t, err)
	container, err := injector.CreateContainer()
	requireErrNil(t, err)
	object, err := container.GetTagged((*SimpleInterface)(nil), "tagOne")
	requireErrNil(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	_, err = container.Get((*SimpleInterface)(nil))
	requireErrMsgContains(t, err, noBindingMsg)
	_, err = container.GetTagged((*SimpleInterface)(nil), "tagTwo")
	requireErrMsgContains(t, err, noBindingMsg)
}

func TestTaggedSimplePtrStructSingletonDirect(t *testing.T) {
	injector := CreateInjector()
	err := injector.BindTagged((*SimpleInterface)(nil), "tagOne").ToSingleton(&SimplePtrStruct{"hello"})
	requireErrNil(t, err)
	container, err := injector.CreateContainer()
	requireErrNil(t, err)
	object, err := container.GetTagged((*SimpleInterface)(nil), "tagOne")
	requireErrNil(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	_, err = container.Get((*SimpleInterface)(nil))
	requireErrMsgContains(t, err, noBindingMsg)
	_, err = container.GetTagged((*SimpleInterface)(nil), "tagTwo")
	requireErrMsgContains(t, err, noBindingMsg)
}

func TestTaggedSimplePtrStructSingletonDirectFailsWhenNotPtr(t *testing.T) {
	injector := CreateInjector()
	err := injector.BindTagged((*SimpleInterface)(nil), "tagOne").ToSingleton(SimplePtrStruct{"hello"})
	require.Equal(t, ErrDoesNotImplement, err)
}

func TestTaggedSimpleStructSingletonIndirect(t *testing.T) {
	injector := CreateInjector()
	err := injector.BindTagged((*SimpleInterface)(nil), "tagOne").To(SimpleStruct{})
	requireErrNil(t, err)
	err = injector.Bind(SimpleStruct{}).ToSingleton(SimpleStruct{"hello"})
	requireErrNil(t, err)
	container, err := injector.CreateContainer()
	requireErrNil(t, err)
	object, err := container.GetTagged((*SimpleInterface)(nil), "tagOne")
	requireErrNil(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	_, err = container.Get((*SimpleInterface)(nil))
	requireErrMsgContains(t, err, noBindingMsg)
	_, err = container.GetTagged((*SimpleInterface)(nil), "tagTwo")
	requireErrMsgContains(t, err, noBindingMsg)
}

func TestTaggedSimpleStructSingletonIndirectSucceedsWhenPtr(t *testing.T) {
	injector := CreateInjector()
	err := injector.BindTagged((*SimpleInterface)(nil), "tagOne").To(&SimpleStruct{})
	requireErrNil(t, err)
	err = injector.Bind(&SimpleStruct{}).ToSingleton(&SimpleStruct{"hello"})
	requireErrNil(t, err)
	container, err := injector.CreateContainer()
	requireErrNil(t, err)
	object, err := container.GetTagged((*SimpleInterface)(nil), "tagOne")
	requireErrNil(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	_, err = container.Get((*SimpleInterface)(nil))
	requireErrMsgContains(t, err, noBindingMsg)
	_, err = container.GetTagged((*SimpleInterface)(nil), "tagTwo")
	requireErrMsgContains(t, err, noBindingMsg)
}

func TestTaggedSimplePtrStructSingletonIndirect(t *testing.T) {
	injector := CreateInjector()
	err := injector.BindTagged((*SimpleInterface)(nil), "tagOne").To(&SimplePtrStruct{})
	requireErrNil(t, err)
	err = injector.Bind(&SimplePtrStruct{}).ToSingleton(&SimplePtrStruct{"hello"})
	requireErrNil(t, err)
	container, err := injector.CreateContainer()
	requireErrNil(t, err)
	object, err := container.GetTagged((*SimpleInterface)(nil), "tagOne")
	requireErrNil(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	_, err = container.Get((*SimpleInterface)(nil))
	requireErrMsgContains(t, err, noBindingMsg)
	_, err = container.GetTagged((*SimpleInterface)(nil), "tagTwo")
	requireErrMsgContains(t, err, noBindingMsg)
}

func TestTaggedSimplePtrStructSingletonIndirectFailsWhenNotPtr(t *testing.T) {
	injector := CreateInjector()
	err := injector.BindTagged((*SimpleInterface)(nil), "tagOne").To(SimplePtrStruct{})
	require.Equal(t, ErrDoesNotImplement, err)
}

func TestTaggedSimplePtrStructIndirectFailsWhenNoFinalBinding(t *testing.T) {
	injector := CreateInjector()
	err := injector.BindTagged((*SimpleInterface)(nil), "tagOne").To(&SimplePtrStruct{})
	requireErrNil(t, err)
	_, err = injector.CreateContainer()
	requireErrMsgContains(t, err, noBindingToSingletonOrProviderMsg)
}

// ***** additional simple tagged tests *****

func TestTaggedSimpleStructSingletonDirectTwoBindings(t *testing.T) {
	injector := CreateInjector()
	err := injector.BindTagged((*SimpleInterface)(nil), "tagOne").ToSingleton(SimpleStruct{"hello"})
	requireErrNil(t, err)
	err = injector.BindTagged((*SimpleInterface)(nil), SimpleTag{}).ToSingleton(SimpleStruct{"good day"})
	requireErrNil(t, err)
	container, err := injector.CreateContainer()
	requireErrNil(t, err)
	object, err := container.GetTagged((*SimpleInterface)(nil), "tagOne")
	requireErrNil(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	object, err = container.GetTagged((*SimpleInterface)(nil), SimpleTag{})
	requireErrNil(t, err)
	simpleInterface = object.(SimpleInterface)
	require.Equal(t, "good day", simpleInterface.Foo())
}

// ***** helpers *****

// TODO(pedge): is there something for this in testify? if not, send a pull request
func requireErrNil(t *testing.T, err error) {
	if err != nil {
		require.Nil(t, err, err.Error())
	}
}

func requireErrMsgContains(t *testing.T, err error, s string) {
	require.NotNil(t, err)
	require.Contains(t, err.Error(), s)
}
