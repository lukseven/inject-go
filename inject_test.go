package inject

import (
	"errors"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	goRoutineIterations = 100
)

// ***** simple bind tests *****

type SimpleInterface interface {
	Foo() string
}

type SimpleStruct struct {
	foo string
}

func (s SimpleStruct) Foo() string {
	return s.foo
}

type SimplePtrStruct struct {
	foo string
}

func (s *SimplePtrStruct) Foo() string {
	return s.foo
}

func TestSimpleStructSingletonDirect(t *testing.T) {
	module := NewModule()
	module.Bind((*SimpleInterface)(nil), SimpleStruct{}).ToSingleton(SimpleStruct{"hello"})
	injector, err := NewInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	object, err = injector.Get(SimpleStruct{})
	require.NoError(t, err)
	simpleStruct := object.(SimpleStruct)
	require.Equal(t, "hello", simpleStruct.Foo())
}

func TestSimpleStructSingletonDirectSucceedsWhenPtr(t *testing.T) {
	module := NewModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimpleStruct{"hello"})
	injector, err := NewInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimplePtrStructSingletonDirect(t *testing.T) {
	module := NewModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	injector, err := NewInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimplePtrStructSingletonDirectFailsWhenNotPtr(t *testing.T) {
	module := NewModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(SimplePtrStruct{"hello"})
	_, err := NewInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNotAssignable)
}

func TestSimpleStructSingletonIndirect(t *testing.T) {
	module := NewModule()
	module.BindInterface((*SimpleInterface)(nil)).To(SimpleStruct{})
	module.Bind(SimpleStruct{}).ToSingleton(SimpleStruct{"hello"})
	injector, err := NewInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimpleStructSingletonIndirectSucceedsWhenPtr(t *testing.T) {
	module := NewModule()
	module.BindInterface((*SimpleInterface)(nil)).To(&SimpleStruct{})
	module.Bind(&SimpleStruct{}).ToSingleton(&SimpleStruct{"hello"})
	injector, err := NewInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimplePtrStructSingletonIndirect(t *testing.T) {
	module := NewModule()
	module.BindInterface((*SimpleInterface)(nil)).To(&SimplePtrStruct{})
	module.Bind(&SimplePtrStruct{}).ToSingleton(&SimplePtrStruct{"hello"})
	injector, err := NewInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SimpleInterface)(nil))
	require.NoError(t, err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimplePtrStructSingletonIndirectFailsWhenNotPtr(t *testing.T) {
	module := NewModule()
	module.BindInterface((*SimpleInterface)(nil)).To(SimplePtrStruct{})
	_, err := NewInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNotAssignable)
}

func TestSimplePtrStructIndirectFailsWhenNoFinalBinding(t *testing.T) {
	module := NewModule()
	module.BindInterface((*SimpleInterface)(nil)).To(&SimplePtrStruct{})
	_, err := NewInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoFinalBinding)
}

// ***** simple BindTagged tests *****

func TestTaggedTagEmpty(t *testing.T) {
	module := NewModule()
	module.BindTagged("", (*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"hello"})
	_, err := NewInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeTagEmpty)
}

func TestTaggedSimpleStructSingletonDirect(t *testing.T) {
	module := NewModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"hello"})
	injector, err := NewInjector(module)
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
	module := NewModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(&SimpleStruct{"hello"})
	injector, err := NewInjector(module)
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
	module := NewModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	injector, err := NewInjector(module)
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
	module := NewModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(SimplePtrStruct{"hello"})
	_, err := NewInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNotAssignable)
}

func TestTaggedSimpleStructSingletonIndirect(t *testing.T) {
	module := NewModule()
	module.BindTaggedInterface("tagOne", (*SimpleInterface)(nil)).To(SimpleStruct{})
	module.Bind(SimpleStruct{}).ToSingleton(SimpleStruct{"hello"})
	injector, err := NewInjector(module)
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
	module := NewModule()
	module.BindTaggedInterface("tagOne", (*SimpleInterface)(nil)).To(&SimpleStruct{})
	module.Bind(&SimpleStruct{}).ToSingleton(&SimpleStruct{"hello"})
	injector, err := NewInjector(module)
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
	module := NewModule()
	module.BindTaggedInterface("tagOne", (*SimpleInterface)(nil)).To(&SimplePtrStruct{})
	module.Bind(&SimplePtrStruct{}).ToSingleton(&SimplePtrStruct{"hello"})
	injector, err := NewInjector(module)
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
	module := NewModule()
	module.BindTaggedInterface("tagOne", (*SimpleInterface)(nil)).To(SimplePtrStruct{})
	_, err := NewInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNotAssignable)
}

func TestTaggedSimplePtrStructIndirectFailsWhenNoFinalBinding(t *testing.T) {
	module := NewModule()
	module.BindTaggedInterface("tagOne", (*SimpleInterface)(nil)).To(&SimplePtrStruct{})
	_, err := NewInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoFinalBinding)
}

// ***** additional simple tagged tests *****

func TestTaggedSimpleStructSingletonDirectTwoBindings(t *testing.T) {
	module := NewModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"hello"})
	module.BindTagged("tagTwo", (*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"good day"})
	injector, err := NewInjector(module)
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

func (b BarStruct) Bar() int {
	return b.bar
}

type BarPtrStruct struct {
	bar int
}

func (b *BarPtrStruct) Bar() int {
	return b.bar
}

type SecondInterface interface {
	Foo() SimpleInterface
	Bar() BarInterface
}

type SecondPtrStruct struct {
	foo SimpleInterface
	bar BarInterface
}

func (s *SecondPtrStruct) Foo() SimpleInterface {
	return s.foo
}

func (s *SecondPtrStruct) Bar() BarInterface {
	return s.bar
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
	module := NewModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(SimplePtrStruct{"hello"})
	module.Bind((*BarInterface)(nil)).ToSingleton(BarPtrStruct{1})
	_, err := NewInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), "1:inject:")
	require.Contains(t, err.Error(), "2:inject:")
	require.Contains(t, err.Error(), injectErrorTypeNotAssignable)
}

func TestConstructorSimple(t *testing.T) {
	module := NewModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	module.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	module.Bind((*SecondInterface)(nil)).ToConstructor(createSecondInterface)
	injector, err := NewInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SecondInterface)(nil))
	require.NoError(t, err)
	secondInterface := object.(SecondInterface)
	require.Equal(t, "hello", secondInterface.Foo().Foo())
	require.Equal(t, 1, secondInterface.Bar().Bar())
}

func TestConstructorErrReturned(t *testing.T) {
	module := NewModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	module.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	module.Bind((*SecondInterface)(nil)).ToConstructor(createSecondInterfaceErr)
	injector, err := NewInjector(module)
	require.NoError(t, err)
	_, err = injector.Get((*SecondInterface)(nil))
	require.Equal(t, "XYZ", err.Error())
}

func TestConstructorErrNoBindingReturned(t *testing.T) {
	module := NewModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	module.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	module.Bind((*SecondInterface)(nil)).ToConstructor(createSecondInterfaceErrNoBinding)
	_, err := NewInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoBinding)
	require.Contains(t, err.Error(), "inject.UnboundInterface")
}

func TestConstructorWithEvilCounter(t *testing.T) {
	module := NewModule()
	module.Bind((*BarInterface)(nil)).ToConstructor(createEvilBarInterface)
	injector, err := NewInjector(module)
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
	module := NewModule()
	module.Bind((*BarInterface)(nil)).ToSingletonConstructor(createEvilBarInterface)
	injector, err := NewInjector(module)
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
	module := NewModule()
	module.Bind((*BarInterface)(nil)).ToSingletonConstructor(createEvilBarInterfaceErr)
	injector, err := NewInjector(module)
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
	module := NewModule()
	module.Bind((*BarInterface)(nil)).ToSingletonConstructor(createEvilBarInterface)
	injector1, err := NewInjector(module)
	require.NoError(t, err)
	injector2, err := NewInjector(module)
	require.NoError(t, err)
	injector3, err := NewInjector(module)
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
	S SimpleInterface `inject:"tagOne"`
	B BarInterface    `inject:"tagTwo"`
}) (SecondInterface, error) {
	return &SecondPtrStruct{str.S, str.B}, nil
}

func createSecondInterfaceTaggedOneHasNoTag(str struct {
	S SimpleInterface `inject:"tagOne"`
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
	module := NewModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	module.BindTagged("tagTwo", (*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	module.Bind((*SecondInterface)(nil)).ToTaggedConstructor(createSecondInterfaceTaggedOne)
	injector, err := NewInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SecondInterface)(nil))
	require.NoError(t, err)
	secondInterface := object.(SecondInterface)
	require.Equal(t, "hello", secondInterface.Foo().Foo())
	require.Equal(t, 1, secondInterface.Bar().Bar())
}

func TestTaggedConstructorOneHasNoTag(t *testing.T) {
	module := NewModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	module.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	module.Bind((*SecondInterface)(nil)).ToTaggedConstructor(createSecondInterfaceTaggedOneHasNoTag)
	injector, err := NewInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SecondInterface)(nil))
	require.NoError(t, err)
	secondInterface := object.(SecondInterface)
	require.Equal(t, "hello", secondInterface.Foo().Foo())
	require.Equal(t, 1, secondInterface.Bar().Bar())
}

func TestTaggedConstructorNoTags(t *testing.T) {
	module := NewModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	module.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	module.Bind((*SecondInterface)(nil)).ToTaggedConstructor(createSecondInterfaceTaggedNoTags)
	injector, err := NewInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SecondInterface)(nil))
	require.NoError(t, err)
	secondInterface := object.(SecondInterface)
	require.Equal(t, "hello", secondInterface.Foo().Foo())
	require.Equal(t, 1, secondInterface.Bar().Bar())
}

func TestTaggedConstructorOneHasNoTagMultipleBindings(t *testing.T) {
	module := NewModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	module.BindTagged("tagTwo", (*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"goodbye"})
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"another"})
	module.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	module.Bind((*SecondInterface)(nil)).ToTaggedConstructor(createSecondInterfaceTaggedOneHasNoTag)
	injector, err := NewInjector(module)
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

func newModuleForTheTestBelow(t *testing.T) Module {
	module := NewModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	module.BindTagged("tagTwo", (*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"goodbye"})
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"another"})
	module.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	module.Bind((*SecondInterface)(nil)).ToTaggedConstructor(createSecondInterfaceTaggedOneHasNoTag)
	return module
}

// had a situation where I thought I might have a pointer issue, keeping this test anyways for now
func TestTaggedConstructorOneHasNoTagMultipleBindingsModuleFromFunction(t *testing.T) {
	injector, err := NewInjector(newModuleForTheTestBelow(t))
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
	S SimpleInterface `inject:"tagOne"`
	B BarInterface    `inject:"tagTwo"`
}) (SecondInterface, string, error) {
	return &SecondPtrStruct{str.S, str.B}, "hello", nil
}

func getSecondInterfaceTaggedOneHasNoTag(str struct {
	S SimpleInterface `inject:"tagOne"`
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

func getSecondInterfaceTaggedNoTagsAndError(str struct {
	S SimpleInterface
	B BarInterface
}) (SecondInterface, string, error) {
	return &SecondPtrStruct{str.S, str.B}, "hello", errors.New("an error")
}

func TestCallAndCallTaggedSimple(t *testing.T) {
	module := NewModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"hello"})
	module.Bind((*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"another"})
	module.BindTagged("tagTwo", (*BarInterface)(nil)).ToSingleton(BarStruct{1})
	module.Bind((*BarInterface)(nil)).ToSingleton(BarStruct{2})
	injector, err := NewInjector(module)
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

	values, err = injector.CallTagged(getSecondInterfaceTaggedNoTagsAndError)
	require.NoError(t, err)
	secondPtrStruct = values[0].(*SecondPtrStruct)
	str = values[1].(string)
	require.Equal(t, SecondPtrStruct{SimpleStruct{"another"}, BarStruct{2}}, *secondPtrStruct)
	require.Equal(t, "hello", str)
	err = values[2].(error)
	require.Equal(t, errors.New("an error"), err)
}

// ***** Populate tests *****

type PopulateStructTwoTags struct {
	S SimpleInterface `inject:"tagOne"`
	B BarInterface    `inject:"tagTwo"`
}

type PopulateStructNoBinding struct {
	S SimpleInterface `inject:"tagOne"`
	B BarInterface    `inject:"tagTwo"`
	U UnboundInterface
}

type PopulateStructOneTag struct {
	S SimpleInterface `inject:"tagOne"`
	B BarInterface
}

type PopulateStructNoTags struct {
	S SimpleInterface
	B BarInterface
}

func TestPopulateSimple(t *testing.T) {
	module := NewModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"hello"})
	module.Bind((*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"another"})
	module.BindTagged("tagTwo", (*BarInterface)(nil)).ToSingleton(BarStruct{1})
	module.Bind((*BarInterface)(nil)).ToSingleton(BarStruct{2})
	injector, err := NewInjector(module)
	require.NoError(t, err)

	populateStructTwoTags := PopulateStructTwoTags{}
	err = injector.Populate(&populateStructTwoTags)
	require.NoError(t, err)
	require.Equal(t, PopulateStructTwoTags{SimpleStruct{"hello"}, BarStruct{1}}, populateStructTwoTags)

	populateStructNoBinding := PopulateStructNoBinding{}
	err = injector.Populate(&populateStructNoBinding)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNoBinding)

	populateStructOneTag := PopulateStructOneTag{}
	err = injector.Populate(&populateStructOneTag)
	require.NoError(t, err)
	require.Equal(t, PopulateStructOneTag{SimpleStruct{"hello"}, BarStruct{2}}, populateStructOneTag)

	populateStructNoTags := PopulateStructNoTags{}
	err = injector.Populate(&populateStructNoTags)
	require.NoError(t, err)
	require.Equal(t, PopulateStructNoTags{SimpleStruct{"another"}, BarStruct{2}}, populateStructNoTags)
}

// ***** BindTaggedConstant tests *****

type PopulateStructOneTagWithInt struct {
	S SimpleInterface `inject:"tagOne"`
	B BarInterface
	I int `inject:"intTag"`
}

type PopulateStructOneTagWithIntErr struct {
	S SimpleInterface `inject:"tagOne"`
	B BarInterface
	I int
}

func TestBindTaggedConstantSimple(t *testing.T) {
	module := NewModule()
	module.BindTaggedBool("boolTrue").ToSingleton(true)
	module.BindTaggedBool("boolFalse").ToSingleton(false)
	module.BindTaggedInt("int10").ToConstructor(func() (int, error) { return 10, nil })
	injector, err := NewInjector(module)
	require.NoError(t, err)

	boolTrue, err := injector.GetTaggedBool("boolTrue")
	require.NoError(t, err)
	require.Equal(t, true, boolTrue)
	boolFalse, err := injector.GetTaggedBool("boolFalse")
	require.NoError(t, err)
	require.Equal(t, false, boolFalse)
}

func TestBindTaggedConstantWrongType(t *testing.T) {
	module := NewModule()
	module.BindTaggedInt8("tagOne").ToSingleton(1)
	_, err := NewInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNotAssignable)
}

func TestBindTaggedConstantPopulate(t *testing.T) {
	module := NewModule()
	module.BindTagged("tagOne", (*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"hello"})
	module.Bind((*BarInterface)(nil)).ToSingleton(BarStruct{2})
	module.BindTaggedInt("intTag").ToSingleton(10)
	injector, err := NewInjector(module)
	require.NoError(t, err)

	populateStructOneTagWithInt := PopulateStructOneTagWithInt{}
	err = injector.Populate(&populateStructOneTagWithInt)
	require.NoError(t, err)
	require.Equal(t, PopulateStructOneTagWithInt{SimpleStruct{"hello"}, BarStruct{2}, 10}, populateStructOneTagWithInt)

	populateStructOneTagWithIntErr := PopulateStructOneTagWithIntErr{}
	err = injector.Populate(&populateStructOneTagWithIntErr)
	require.Error(t, err)
	require.Contains(t, err.Error(), injectErrorTypeNotSupportedYet)
}

func TestOverride(t *testing.T) {
	module := NewModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	module.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	module.Bind((*SecondInterface)(nil)).ToConstructor(createSecondInterface)

	override1 := NewModule()
	override1.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"yo, yo!"})

	override2 := NewModule()
	override2.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{3})

	injector, err := NewInjector(Override(module).With(override1, override2))
	require.NoError(t, err)
	object, err := injector.Get((*SecondInterface)(nil))
	require.NoError(t, err)
	secondInterface := object.(SecondInterface)
	require.Equal(t, "yo, yo!", secondInterface.Foo().Foo())
	require.Equal(t, 3, secondInterface.Bar().Bar())
}

func TestOverrideErrorsInSource(t *testing.T) {
	module := NewModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(SimplePtrStruct{"hello"})
	module.Bind((*BarInterface)(nil)).ToSingleton(BarPtrStruct{1})

	override := NewModule()
	override.Bind((*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"hello"})
	override.Bind((*BarInterface)(nil)).ToSingleton(BarStruct{1})

	_, err := NewInjector(Override(module).With(override))
	require.Error(t, err)
	require.Contains(t, err.Error(), "1:inject:")
	require.Contains(t, err.Error(), "2:inject:")
	require.Contains(t, err.Error(), injectErrorTypeNotAssignable)
}

func TestOverrideErrorsInOverride(t *testing.T) {
	module := NewModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"hello"})
	module.Bind((*BarInterface)(nil)).ToSingleton(BarStruct{1})

	override := NewModule()
	override.Bind((*SimpleInterface)(nil)).ToSingleton(SimplePtrStruct{"hello"})
	override.Bind((*BarInterface)(nil)).ToSingleton(BarPtrStruct{1})

	_, err := NewInjector(Override(module).With(override))
	require.Error(t, err)
	require.Contains(t, err.Error(), "1:inject:")
	require.Contains(t, err.Error(), "2:inject:")
	require.Contains(t, err.Error(), injectErrorTypeNotAssignable)
}

func TestOverrideErrorsInBoth(t *testing.T) {
	module := NewModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"hello"})
	module.Bind((*BarInterface)(nil)).ToSingleton(BarPtrStruct{1})

	override := NewModule()
	override.Bind((*BarInterface)(nil)).ToSingleton(BarStruct{1})
	override.Bind((*SimpleInterface)(nil)).ToSingleton(SimplePtrStruct{"hello"})

	_, err := NewInjector(Override(module).With(override))
	require.Error(t, err)
	require.Contains(t, err.Error(), "1:inject:")
	require.Contains(t, err.Error(), "2:inject:")
	require.Contains(t, err.Error(), injectErrorTypeNotAssignable)
}

func TestInstall(t *testing.T) {
	sub1 := NewModule()
	sub1.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})

	sub2 := NewModule()
	sub2.Bind((*SecondInterface)(nil)).ToConstructor(createSecondInterface)

	module := NewModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	module.Install(sub1)
	module.Install(sub2)

	injector, err := NewInjector(module)
	require.NoError(t, err)
	object, err := injector.Get((*SecondInterface)(nil))
	require.NoError(t, err)
	secondInterface := object.(SecondInterface)
	require.Equal(t, "hello", secondInterface.Foo().Foo())
	require.Equal(t, 1, secondInterface.Bar().Bar())
}

func TestInstallDuplicateBinding(t *testing.T) {
	sub1 := NewModule()
	sub1.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"}) // duplicate binding

	module := NewModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})

	module.Install(sub1)

	_, err := NewInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), "1:inject:")
	require.Contains(t, err.Error(), injectErrorTypeAlreadyBound)
}

func TestInstallErrorsInSource(t *testing.T) {
	sub1 := NewModule()
	sub1.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	sub1.Bind((*SecondInterface)(nil)).ToConstructor(createSecondInterface)

	module := NewModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(SimplePtrStruct{"hello"}) // invalid binding

	module.Install(sub1)

	_, err := NewInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), "1:inject:")
	require.Contains(t, err.Error(), injectErrorTypeNotAssignable)
}

func TestInstallErrorsInTarget(t *testing.T) {
	sub := NewModule()
	sub.Bind((*BarInterface)(nil)).ToSingleton(BarPtrStruct{1}) // invalid binding
	sub.Bind((*SecondInterface)(nil)).ToConstructor(createSecondInterface)

	module := NewModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})

	module.Install(sub)

	_, err := NewInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), "1:inject:")
	require.Contains(t, err.Error(), injectErrorTypeNotAssignable)
}

func TestInstallErrorsInBoth(t *testing.T) {
	sub := NewModule()
	sub.Bind((*BarInterface)(nil)).ToSingleton(BarPtrStruct{1}) // invalid binding
	sub.Bind((*SecondInterface)(nil)).ToConstructor(createSecondInterface)

	module := NewModule()
	module.Bind((*SimpleInterface)(nil)).ToSingleton(SimplePtrStruct{"hello"}) // invalid binding

	module.Install(sub)

	_, err := NewInjector(module)
	require.Error(t, err)
	require.Contains(t, err.Error(), "1:inject:")
	require.Contains(t, err.Error(), "2:inject:")
	require.Contains(t, err.Error(), injectErrorTypeNotAssignable)
}
