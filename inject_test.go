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
	require.Nil(t, err, "%v", err)
	container, err := injector.CreateContainer()
	require.Nil(t, err, "%v", err)
	object, err := container.Get((*SimpleInterface)(nil))
	require.Nil(t, err, "%v", err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimpleStructSingletonDirectSucceedsWhenPtr(t *testing.T) {
	injector := CreateInjector()
	err := injector.Bind((*SimpleInterface)(nil)).ToSingleton(&SimpleStruct{"hello"})
	require.Nil(t, err, "%v", err)
	container, err := injector.CreateContainer()
	require.Nil(t, err, "%v", err)
	object, err := container.Get((*SimpleInterface)(nil))
	require.Nil(t, err, "%v", err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimplePtrStructSingletonDirect(t *testing.T) {
	injector := CreateInjector()
	err := injector.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	require.Nil(t, err, "%v", err)
	container, err := injector.CreateContainer()
	require.Nil(t, err, "%v", err)
	object, err := container.Get((*SimpleInterface)(nil))
	require.Nil(t, err, "%v", err)
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
	require.Nil(t, err, "%v", err)
	err = injector.Bind(SimpleStruct{}).ToSingleton(SimpleStruct{"hello"})
	require.Nil(t, err, "%v", err)
	container, err := injector.CreateContainer()
	require.Nil(t, err, "%v", err)
	object, err := container.Get((*SimpleInterface)(nil))
	require.Nil(t, err, "%v", err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimpleStructSingletonIndirectSucceedsWhenPtr(t *testing.T) {
	injector := CreateInjector()
	err := injector.Bind((*SimpleInterface)(nil)).To(&SimpleStruct{})
	require.Nil(t, err, "%v", err)
	err = injector.Bind(&SimpleStruct{}).ToSingleton(&SimpleStruct{"hello"})
	require.Nil(t, err, "%v", err)
	container, err := injector.CreateContainer()
	require.Nil(t, err, "%v", err)
	object, err := container.Get((*SimpleInterface)(nil))
	require.Nil(t, err, "%v", err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimplePtrStructSingletonIndirect(t *testing.T) {
	injector := CreateInjector()
	err := injector.Bind((*SimpleInterface)(nil)).To(&SimplePtrStruct{})
	require.Nil(t, err, "%v", err)
	err = injector.Bind(&SimplePtrStruct{}).ToSingleton(&SimplePtrStruct{"hello"})
	require.Nil(t, err, "%v", err)
	container, err := injector.CreateContainer()
	require.Nil(t, err, "%v", err)
	object, err := container.Get((*SimpleInterface)(nil))
	require.Nil(t, err, "%v", err)
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
	require.Nil(t, err, "%v", err)
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
	require.Nil(t, err, "%v", err)
	container, err := injector.CreateContainer()
	require.Nil(t, err, "%v", err)
	object, err := container.GetTagged((*SimpleInterface)(nil), "tagOne")
	require.Nil(t, err, "%v", err)
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
	require.Nil(t, err, "%v", err)
	container, err := injector.CreateContainer()
	require.Nil(t, err, "%v", err)
	object, err := container.GetTagged((*SimpleInterface)(nil), "tagOne")
	require.Nil(t, err, "%v", err)
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
	require.Nil(t, err, "%v", err)
	container, err := injector.CreateContainer()
	require.Nil(t, err, "%v", err)
	object, err := container.GetTagged((*SimpleInterface)(nil), "tagOne")
	require.Nil(t, err, "%v", err)
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
	require.Nil(t, err, "%v", err)
	err = injector.Bind(SimpleStruct{}).ToSingleton(SimpleStruct{"hello"})
	require.Nil(t, err, "%v", err)
	container, err := injector.CreateContainer()
	require.Nil(t, err, "%v", err)
	object, err := container.GetTagged((*SimpleInterface)(nil), "tagOne")
	require.Nil(t, err, "%v", err)
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
	require.Nil(t, err, "%v", err)
	err = injector.Bind(&SimpleStruct{}).ToSingleton(&SimpleStruct{"hello"})
	require.Nil(t, err, "%v", err)
	container, err := injector.CreateContainer()
	require.Nil(t, err, "%v", err)
	object, err := container.GetTagged((*SimpleInterface)(nil), "tagOne")
	require.Nil(t, err, "%v", err)
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
	require.Nil(t, err, "%v", err)
	err = injector.Bind(&SimplePtrStruct{}).ToSingleton(&SimplePtrStruct{"hello"})
	require.Nil(t, err, "%v", err)
	container, err := injector.CreateContainer()
	require.Nil(t, err, "%v", err)
	object, err := container.GetTagged((*SimpleInterface)(nil), "tagOne")
	require.Nil(t, err, "%v", err)
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
	require.Nil(t, err, "%v", err)
	_, err = injector.CreateContainer()
	requireErrMsgContains(t, err, noBindingToSingletonOrProviderMsg)
}

// ***** additional simple tagged tests *****

func TestTaggedSimpleStructSingletonDirectTwoBindings(t *testing.T) {
	injector := CreateInjector()
	err := injector.BindTagged((*SimpleInterface)(nil), "tagOne").ToSingleton(SimpleStruct{"hello"})
	require.Nil(t, err, "%v", err)
	err = injector.BindTagged((*SimpleInterface)(nil), SimpleTag{}).ToSingleton(SimpleStruct{"good day"})
	require.Nil(t, err, "%v", err)
	container, err := injector.CreateContainer()
	require.Nil(t, err, "%v", err)
	object, err := container.GetTagged((*SimpleInterface)(nil), "tagOne")
	require.Nil(t, err, "%v", err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	object, err = container.GetTagged((*SimpleInterface)(nil), SimpleTag{})
	require.Nil(t, err, "%v", err)
	simpleInterface = object.(SimpleInterface)
	require.Equal(t, "good day", simpleInterface.Foo())
}

// ***** simple provider tests *****

type BarInterface interface {
	Bar() string
}

type BarPtrStruct struct {
	bar string
}

func (this *BarPtrStruct) Bar() string {
	return this.bar
}

type SecondInterface interface {
	Foo() SimpleInterface
	Bar() BarInterface
}

type SecondPtrStruct struct {
	foo SimpleInterface
	bar BarInterface
}

func (this *SecondPtrStruct) Foo() SimpleInterface {
	return this.foo
}

func (this *SecondPtrStruct) Bar() BarInterface {
	return this.bar
}

func createSecondInterface(s SimpleInterface, b BarInterface) (SecondInterface, error) {
	return &SecondPtrStruct{s, b}, nil
}

func TestProviderOne(t *testing.T) {
	injector := CreateInjector()
	err := injector.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	require.Nil(t, err, "%v", err)
	err = injector.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{"good day"})
	require.Nil(t, err, "%v", err)
	err = injector.Bind((*SecondInterface)(nil)).ToProvider(createSecondInterface)
	require.Nil(t, err, "%v", err)
	container, err := injector.CreateContainer()
	require.Nil(t, err, "%v", err)
	object, err := container.Get((*SecondInterface)(nil))
	require.Nil(t, err, "%v", err)
	secondInterface := object.(SecondInterface)
	require.Equal(t, "hello", secondInterface.Foo().Foo())
	require.Equal(t, "good day", secondInterface.Bar().Bar())
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
