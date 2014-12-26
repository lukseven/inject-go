package inject

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"sync/atomic"
	"testing"
)

const (
	goRoutineIterations = 100
)

// ***** error tests *****

func TestInjectErrorWithoutTags(t *testing.T) {
	errorBuilder := newErrorBuilder("foo")
	injectError := errorBuilder.build()
	require.Equal(t, "inject: foo", injectError.Error())
	require.Equal(t, "foo", injectError.Type())
	_, ok := injectError.GetTag("tagOne")
	require.False(t, ok)
}

func TestInjectErrorWithTags(t *testing.T) {
	errorBuilder := newErrorBuilder("foo")
	errorBuilder = errorBuilder.addTag("tagOne", 1)
	errorBuilder = errorBuilder.addTag("tagTwo", "two")
	injectError := errorBuilder.build()
	// making sure that the order of the tags is the same and that this does not rely
	// only on a map - this is a bad and non-deterministic way to do this but fix later
	for i := 0; i < 100; i++ {
		require.Equal(t, "inject: foo tags{tagOne:1 tagTwo:two}", injectError.Error())
	}
	require.Equal(t, "foo", injectError.Type())
	tagOne, ok := injectError.GetTag("tagOne")
	require.True(t, ok)
	require.Equal(t, 1, tagOne)
	tagTwo, ok := injectError.GetTag("tagTwo")
	require.True(t, ok)
	require.Equal(t, "two", tagTwo)
}

// ***** simple bind tests *****

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
	module := CreateModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"hello"})
	injector, err := CreateInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimpleStructSingletonDirectSucceedsWhenPtr(t *testing.T) {
	module := CreateModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimpleStruct{"hello"})
	injector, err := CreateInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimplePtrStructSingletonDirect(t *testing.T) {
	module := CreateModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	injector, err := CreateInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimplePtrStructSingletonDirectFailsWhenNotPtr(t *testing.T) {
	module := CreateModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(SimplePtrStruct{"hello"})
	_, err := CreateInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeDoesNotImplement)
}

func TestSimpleStructSingletonIndirect(t *testing.T) {
	module := CreateModule()
	module.Bind((*SimpleInterface)(nil)).To(SimpleStruct{})
	module.Bind(SimpleStruct{}).ToSingleton(SimpleStruct{"hello"})
	injector, err := CreateInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimpleStructSingletonIndirectSucceedsWhenPtr(t *testing.T) {
	module := CreateModule()
	module.Bind((*SimpleInterface)(nil)).To(&SimpleStruct{})
	module.Bind(&SimpleStruct{}).ToSingleton(&SimpleStruct{"hello"})
	injector, err := CreateInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimplePtrStructSingletonIndirect(t *testing.T) {
	module := CreateModule()
	module.Bind((*SimpleInterface)(nil)).To(&SimplePtrStruct{})
	module.Bind(&SimplePtrStruct{}).ToSingleton(&SimplePtrStruct{"hello"})
	injector, err := CreateInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimplePtrStructSingletonIndirectFailsWhenNotPtr(t *testing.T) {
	module := CreateModule()
	module.Bind((*SimpleInterface)(nil)).To(SimplePtrStruct{})
	_, err := CreateInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeDoesNotImplement)
}

func TestSimplePtrStructIndirectFailsWhenNoFinalBinding(t *testing.T) {
	module := CreateModule()
	module.Bind((*SimpleInterface)(nil)).To(&SimplePtrStruct{})
	_, err := CreateInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoFinalBinding)
}

// ***** simple BindTagged tests *****

func TestTaggedTagEmpty(t *testing.T) {
	module := CreateModule()
	module.BindTagged("", (*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"hello"})
	_, err := CreateInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeTagEmpty)
}

func TestTaggedSimpleStructSingletonDirect(t *testing.T) {
	module := CreateModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"hello"})
	injector, err := CreateInjector(module)
	require.NoError(t, err)
	object, err := injector.GetTagged("tagOne", (*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	_, err = injector.Get((*SimpleInterface)(nil))
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoBinding)
	_, err = injector.GetTagged("tagTwo", (*SimpleInterface)(nil))
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoBinding)
}

func TestTaggedSimpleStructSingletonDirectSucceedsWhenPtr(t *testing.T) {
	module := CreateModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(&SimpleStruct{"hello"})
	injector, err := CreateInjector(module)
	require.NoError(t, err)
	object, err := injector.GetTagged("tagOne", (*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	_, err = injector.Get((*SimpleInterface)(nil))
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoBinding)
	_, err = injector.GetTagged("tagTwo", (*SimpleInterface)(nil))
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoBinding)
}

func TestTaggedSimplePtrStructSingletonDirect(t *testing.T) {
	module := CreateModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	injector, err := CreateInjector(module)
	require.NoError(t, err)
	object, err := injector.GetTagged("tagOne", (*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	_, err = injector.Get((*SimpleInterface)(nil))
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoBinding)
	_, err = injector.GetTagged("tagTwo", (*SimpleInterface)(nil))
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoBinding)
}

func TestTaggedSimplePtrStructSingletonDirectFailsWhenNotPtr(t *testing.T) {
	module := CreateModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(SimplePtrStruct{"hello"})
	_, err := CreateInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeDoesNotImplement)
}

func TestTaggedSimpleStructSingletonIndirect(t *testing.T) {
	module := CreateModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).To(SimpleStruct{})
	module.Bind(SimpleStruct{}).ToSingleton(SimpleStruct{"hello"})
	injector, err := CreateInjector(module)
	require.NoError(t, err)
	object, err := injector.GetTagged("tagOne", (*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	_, err = injector.Get((*SimpleInterface)(nil))
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoBinding)
	_, err = injector.GetTagged("tagTwo", (*SimpleInterface)(nil))
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoBinding)
}

func TestTaggedSimpleStructSingletonIndirectSucceedsWhenPtr(t *testing.T) {
	module := CreateModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).To(&SimpleStruct{})
	module.Bind(&SimpleStruct{}).ToSingleton(&SimpleStruct{"hello"})
	injector, err := CreateInjector(module)
	require.NoError(t, err)
	object, err := injector.GetTagged("tagOne", (*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	_, err = injector.Get((*SimpleInterface)(nil))
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoBinding)
	_, err = injector.GetTagged("tagTwo", (*SimpleInterface)(nil))
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoBinding)
}

func TestTaggedSimplePtrStructSingletonIndirect(t *testing.T) {
	module := CreateModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).To(&SimplePtrStruct{})
	module.Bind(&SimplePtrStruct{}).ToSingleton(&SimplePtrStruct{"hello"})
	injector, err := CreateInjector(module)
	require.NoError(t, err)
	object, err := injector.GetTagged("tagOne", (*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	_, err = injector.Get((*SimpleInterface)(nil))
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoBinding)
	_, err = injector.GetTagged("tagTwo", (*SimpleInterface)(nil))
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoBinding)
}

func TestTaggedSimplePtrStructSingletonIndirectFailsWhenNotPtr(t *testing.T) {
	module := CreateModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).To(SimplePtrStruct{})
	_, err := CreateInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeDoesNotImplement)
}

func TestTaggedSimplePtrStructIndirectFailsWhenNoFinalBinding(t *testing.T) {
	module := CreateModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).To(&SimplePtrStruct{})
	_, err := CreateInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoFinalBinding)
}

// ***** additional simple tagged tests *****

func TestTaggedSimpleStructSingletonDirectTwoBindings(t *testing.T) {
	module := CreateModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"hello"})
	module.BindTagged("tagTwo", (*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"good day"})
	injector, err := CreateInjector(module)
	require.NoError(t, err)
	object, err := injector.GetTagged("tagOne", (*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	object, err = injector.GetTagged("tagTwo", (*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface = object.(SimpleInterface)
	require.Equal(t, "good day", simpleInterface.Foo())
}

// ***** simple constructor tests *****

type BarInterface interface {
	Bar() int
}

type BarStruct struct {
	bar int
}

func (this BarStruct) Bar() int {
	return this.bar
}

type BarPtrStruct struct {
	bar int
}

func (this *BarPtrStruct) Bar() int {
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

type UnboundInterface interface {
	Baz() string
}

func createSecondInterface(s SimpleInterface, b BarInterface) (SecondInterface, error) {
	return &SecondPtrStruct{s, b}, nil
}

func createSecondInterfaceErr(s SimpleInterface, b BarInterface) (SecondInterface, error) {
	return nil, errors.New("XYZ")
}

func createSecondInterfaceErrNoBinding(s SimpleInterface, b BarInterface, u UnboundInterface) (SecondInterface, error) {
	return &SecondPtrStruct{s, b}, nil
}

type BarInterfaceError struct {
	BarInterface
	err error
}

var evilCounter int32

func createEvilBarInterface() (BarInterface, error) {
	value := atomic.AddInt32(&evilCounter, 1)
	return &BarPtrStruct{int(value)}, nil
}

func createEvilBarInterfaceErr() (BarInterface, error) {
	value := atomic.AddInt32(&evilCounter, 1)
	return nil, fmt.Errorf("XYZ %v", value)
}

func TestMultipleBindingErrors(t *testing.T) {
	module := CreateModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(SimplePtrStruct{"hello"})
	module.Bind((*BarInterface)(nil)).ToSingleton(BarPtrStruct{1})
	_, err := CreateInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), "1:inject:")
	require.Contains(t, err.Error(), "2:inject:")
	require.Contains(t, err.Error(), injectErrorTypeDoesNotImplement)
}

func TestConstructorSimple(t *testing.T) {
	module := CreateModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	module.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	module.Bind((*SecondInterface)(nil)).ToConstructor(createSecondInterface)
	injector, err := CreateInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SecondInterface)(nil))
	require.NoError(t, err)
	secondInterface := object.(SecondInterface)
	require.Equal(t, "hello", secondInterface.Foo().Foo())
	require.Equal(t, 1, secondInterface.Bar().Bar())
}

func TestConstructorErrReturned(t *testing.T) {
	module := CreateModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	module.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	module.Bind((*SecondInterface)(nil)).ToConstructor(createSecondInterfaceErr)
	injector, err := CreateInjector(module)
	require.NoError(t, err)
	_, err = injector.Get((*SecondInterface)(nil))
	require.Equal(t, "XYZ", err.Error())
}

func TestConstructorErrNoBindingReturned(t *testing.T) {
	module := CreateModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	module.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	module.Bind((*SecondInterface)(nil)).ToConstructor(createSecondInterfaceErrNoBinding)
	_, err := CreateInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoBinding)
	require.Contains(t, err.Error(), "inject.UnboundInterface")
}

func TestConstructorWithEvilCounter(t *testing.T) {
	module := CreateModule()
	module.Bind((*BarInterface)(nil)).ToConstructor(createEvilBarInterface)
	injector, err := CreateInjector(module)
	require.NoError(t, err)

	evilCounter = int32(0)
	evilChan := make(chan BarInterfaceError)

	// TODO(pedge): i know this is a terrible way to do concurrency testing
	for i := 0; i < goRoutineIterations; i++ {
		go func() {
			barInterface, err := injector.Get((*BarInterface)(nil))
			if err != nil {
				evilChan <- BarInterfaceError{nil, err}
			} else {
				evilChan <- BarInterfaceError{barInterface.(BarInterface), nil}
			}
		}()
	}
	count := 0
	for i := 0; i < goRoutineIterations; i++ {
		barInterfaceErr := <-evilChan
		require.Nil(t, barInterfaceErr.err, "%v", barInterfaceErr.err)
		bar := barInterfaceErr.Bar()
		count += bar
	}
	// cute
	require.Equal(t, (goRoutineIterations*(goRoutineIterations+1))/2, count)

	close(evilChan)
}

func TestSingletonConstructorWithEvilCounter(t *testing.T) {
	module := CreateModule()
	module.Bind((*BarInterface)(nil)).ToSingletonConstructor(createEvilBarInterface)
	injector, err := CreateInjector(module)
	require.NoError(t, err)

	evilCounter = int32(0)
	evilChan := make(chan BarInterfaceError)

	// TODO(pedge): i know this is a terrible way to do concurrency testing
	for i := 0; i < goRoutineIterations; i++ {
		go func() {
			barInterface, err := injector.Get((*BarInterface)(nil))
			if err != nil {
				evilChan <- BarInterfaceError{nil, err}
			} else {
				evilChan <- BarInterfaceError{barInterface.(BarInterface), nil}
			}
		}()
	}
	for i := 0; i < goRoutineIterations; i++ {
		barInterfaceErr := <-evilChan
		require.Nil(t, barInterfaceErr.err, "%v", barInterfaceErr.err)
		require.Equal(t, 1, barInterfaceErr.Bar())
	}

	close(evilChan)
}

func TestSingletonConstructorWithEvilCounterErr(t *testing.T) {
	module := CreateModule()
	module.Bind((*BarInterface)(nil)).ToSingletonConstructor(createEvilBarInterfaceErr)
	injector, err := CreateInjector(module)
	require.NoError(t, err)

	evilCounter = int32(0)
	evilChan := make(chan BarInterfaceError)

	// TODO(pedge): i know this is a terrible way to do concurrency testing
	for i := 0; i < goRoutineIterations; i++ {
		go func() {
			barInterface, err := injector.Get((*BarInterface)(nil))
			if err != nil {
				evilChan <- BarInterfaceError{nil, err}
			} else {
				evilChan <- BarInterfaceError{barInterface.(BarInterface), nil}
			}
		}()
	}
	for i := 0; i < goRoutineIterations; i++ {
		barInterfaceErr := <-evilChan
		require.Equal(t, "XYZ 1", barInterfaceErr.err.Error())
	}

	close(evilChan)
}

func TestSingletonConstructorWithEvilCounterMultipleInjectors(t *testing.T) {
	module := CreateModule()
	module.Bind((*BarInterface)(nil)).ToSingletonConstructor(createEvilBarInterface)
	injector1, err := CreateInjector(module)
	require.NoError(t, err)
	injector2, err := CreateInjector(module)
	require.NoError(t, err)
	injector3, err := CreateInjector(module)
	require.NoError(t, err)

	evilCounter = int32(0)

	barInterface1, err := injector1.Get((*BarInterface)(nil))
	require.NoError(t, err)
	barInterface2, err := injector2.Get((*BarInterface)(nil))
	require.NoError(t, err)
	barInterface3, err := injector3.Get((*BarInterface)(nil))
	require.NoError(t, err)
	require.Equal(t, 1, barInterface1.(BarInterface).Bar())
	require.Equal(t, 2, barInterface2.(BarInterface).Bar())
	require.Equal(t, 3, barInterface3.(BarInterface).Bar())
	barInterface1, err = injector1.Get((*BarInterface)(nil))
	require.NoError(t, err)
	barInterface2, err = injector2.Get((*BarInterface)(nil))
	require.NoError(t, err)
	barInterface3, err = injector3.Get((*BarInterface)(nil))
	require.NoError(t, err)
	require.Equal(t, 1, barInterface1.(BarInterface).Bar())
	require.Equal(t, 2, barInterface2.(BarInterface).Bar())
	require.Equal(t, 3, barInterface3.(BarInterface).Bar())
}

// ***** tagged constructors

func createSecondInterfaceTaggedOne(str struct {
	S SimpleInterface `injectTag:"tagOne"`
	B BarInterface    `injectTag:"tagTwo"`
}) (SecondInterface, error) {
	return &SecondPtrStruct{str.S, str.B}, nil
}

func createSecondInterfaceTaggedOneHasNoTag(str struct {
	S SimpleInterface `injectTag:"tagOne"`
	B BarInterface
}) (SecondInterface, error) {
	return &SecondPtrStruct{str.S, str.B}, nil
}

func createSecondInterfaceTaggedNoTags(str struct {
	S SimpleInterface
	B BarInterface
}) (SecondInterface, error) {
	return &SecondPtrStruct{str.S, str.B}, nil
}

func TestTaggedConstructorSimple(t *testing.T) {
	module := CreateModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	module.BindTagged("tagTwo", (*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	module.Bind((*SecondInterface)(nil)).ToTaggedConstructor(createSecondInterfaceTaggedOne)
	injector, err := CreateInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SecondInterface)(nil))
	require.NoError(t, err)
	secondInterface := object.(SecondInterface)
	require.Equal(t, "hello", secondInterface.Foo().Foo())
	require.Equal(t, 1, secondInterface.Bar().Bar())
}

func TestTaggedConstructorOneHasNoTag(t *testing.T) {
	module := CreateModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	module.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	module.Bind((*SecondInterface)(nil)).ToTaggedConstructor(createSecondInterfaceTaggedOneHasNoTag)
	injector, err := CreateInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SecondInterface)(nil))
	require.NoError(t, err)
	secondInterface := object.(SecondInterface)
	require.Equal(t, "hello", secondInterface.Foo().Foo())
	require.Equal(t, 1, secondInterface.Bar().Bar())
}

func TestTaggedConstructorNoTags(t *testing.T) {
	module := CreateModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	module.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	module.Bind((*SecondInterface)(nil)).ToTaggedConstructor(createSecondInterfaceTaggedNoTags)
	injector, err := CreateInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SecondInterface)(nil))
	require.NoError(t, err)
	secondInterface := object.(SecondInterface)
	require.Equal(t, "hello", secondInterface.Foo().Foo())
	require.Equal(t, 1, secondInterface.Bar().Bar())
}

func TestTaggedConstructorOneHasNoTagMultipleBindings(t *testing.T) {
	module := CreateModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	module.BindTagged("tagTwo", (*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"goodbye"})
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"another"})
	module.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	module.Bind((*SecondInterface)(nil)).ToTaggedConstructor(createSecondInterfaceTaggedOneHasNoTag)
	injector, err := CreateInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SecondInterface)(nil))
	require.NoError(t, err)
	secondInterface := object.(SecondInterface)
	require.Equal(t, "hello", secondInterface.Foo().Foo())
	require.Equal(t, 1, secondInterface.Bar().Bar())
	object, err = injector.GetTagged("tagTwo", (*SimpleInterface)(nil))
	require.NoError(t, err)
	require.Equal(t, "goodbye", object.(SimpleInterface).Foo())
	object, err = injector.Get((*SimpleInterface)(nil))
	require.NoError(t, err)
	require.Equal(t, "another", object.(SimpleInterface).Foo())
}

func createModuleForTheTestBelow(t *testing.T) Module {
	module := CreateModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	module.BindTagged("tagTwo", (*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"goodbye"})
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"another"})
	module.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	module.Bind((*SecondInterface)(nil)).ToTaggedConstructor(createSecondInterfaceTaggedOneHasNoTag)
	return module
}

// had a situation where I thought I might have a pointer issue, keeping this test anyways for now
func TestTaggedConstructorOneHasNoTagMultipleBindingsModuleFromFunction(t *testing.T) {
	injector, err := CreateInjector(createModuleForTheTestBelow(t))
	require.NoError(t, err)
	object, err := injector.Get((*SecondInterface)(nil))
	require.NoError(t, err)
	secondInterface := object.(SecondInterface)
	require.Equal(t, "hello", secondInterface.Foo().Foo())
	require.Equal(t, 1, secondInterface.Bar().Bar())
	object, err = injector.GetTagged("tagTwo", (*SimpleInterface)(nil))
	require.NoError(t, err)
	require.Equal(t, "goodbye", object.(SimpleInterface).Foo())
	object, err = injector.Get((*SimpleInterface)(nil))
	require.NoError(t, err)
	require.Equal(t, "another", object.(SimpleInterface).Foo())
}

// ***** Call, CallTagged tests *****

func getSecondInterface(s SimpleInterface, b BarInterface) (SecondInterface, string, error) {
	return &SecondPtrStruct{s, b}, "hello", nil
}

func getSecondInterfaceErrNoBinding(s SimpleInterface, b BarInterface, u UnboundInterface) (SecondInterface, string, error) {
	return &SecondPtrStruct{s, b}, "hello", nil
}

func getSecondInterfaceTaggedOne(str struct {
	S SimpleInterface `injectTag:"tagOne"`
	B BarInterface    `injectTag:"tagTwo"`
}) (SecondInterface, string, error) {
	return &SecondPtrStruct{str.S, str.B}, "hello", nil
}

func getSecondInterfaceTaggedOneHasNoTag(str struct {
	S SimpleInterface `injectTag:"tagOne"`
	B BarInterface
}) (SecondInterface, string, error) {
	return &SecondPtrStruct{str.S, str.B}, "hello", nil
}

func getSecondInterfaceTaggedNoTags(str struct {
	S SimpleInterface
	B BarInterface
}) (SecondInterface, string, error) {
	return &SecondPtrStruct{str.S, str.B}, "hello", nil
}

func TestCallAndCallTaggedSimple(t *testing.T) {
	module := CreateModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"hello"})
	module.Bind((*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"another"})
	module.BindTagged("tagTwo", (*BarInterface)(nil)).ToSingleton(BarStruct{1})
	module.Bind((*BarInterface)(nil)).ToSingleton(BarStruct{2})
	injector, err := CreateInjector(module)
	require.NoError(t, err)

	values, err := injector.Call(getSecondInterface)
	require.NoError(t, err)
	secondPtrStruct := values[0].(*SecondPtrStruct)
	str := values[1].(string)
	require.Equal(t, SecondPtrStruct{SimpleStruct{"another"}, BarStruct{2}}, *secondPtrStruct)
	require.Equal(t, "hello", str)
	require.Nil(t, values[2])

	values, err = injector.Call(getSecondInterfaceErrNoBinding)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoBinding)

	values, err = injector.CallTagged(getSecondInterfaceTaggedOne)
	require.NoError(t, err)
	secondPtrStruct = values[0].(*SecondPtrStruct)
	str = values[1].(string)
	require.Equal(t, SecondPtrStruct{SimpleStruct{"hello"}, BarStruct{1}}, *secondPtrStruct)
	require.Equal(t, "hello", str)
	require.Nil(t, values[2])

	values, err = injector.CallTagged(getSecondInterfaceTaggedOneHasNoTag)
	require.NoError(t, err)
	secondPtrStruct = values[0].(*SecondPtrStruct)
	str = values[1].(string)
	require.Equal(t, SecondPtrStruct{SimpleStruct{"hello"}, BarStruct{2}}, *secondPtrStruct)
	require.Equal(t, "hello", str)
	require.Nil(t, values[2])

	values, err = injector.CallTagged(getSecondInterfaceTaggedNoTags)
	require.NoError(t, err)
	secondPtrStruct = values[0].(*SecondPtrStruct)
	str = values[1].(string)
	require.Equal(t, SecondPtrStruct{SimpleStruct{"another"}, BarStruct{2}}, *secondPtrStruct)
	require.Equal(t, "hello", str)
	require.Nil(t, values[2])
}
