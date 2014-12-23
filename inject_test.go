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

func TestSimpleStructSingletonDirect(t *testing.T) {
	module := CreateModule()
	err := module.Bind((*SimpleInterface)(nil)).ToSingleton(SimpleStruct{"hello"})
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)
	object, err := injector.Get((*SimpleInterface)(nil))
	require.Nil(t, err, "%v", err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimpleStructSingletonDirectSucceedsWhenPtr(t *testing.T) {
	module := CreateModule()
	err := module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimpleStruct{"hello"})
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)
	object, err := injector.Get((*SimpleInterface)(nil))
	require.Nil(t, err, "%v", err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimplePtrStructSingletonDirect(t *testing.T) {
	module := CreateModule()
	err := module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)
	object, err := injector.Get((*SimpleInterface)(nil))
	require.Nil(t, err, "%v", err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimplePtrStructSingletonDirectFailsWhenNotPtr(t *testing.T) {
	module := CreateModule()
	err := module.Bind((*SimpleInterface)(nil)).ToSingleton(SimplePtrStruct{"hello"})
	require.Contains(t, err.Error(), InjectErrorTypeDoesNotImplement)
}

func TestSimpleStructSingletonIndirect(t *testing.T) {
	module := CreateModule()
	err := module.Bind((*SimpleInterface)(nil)).To(SimpleStruct{})
	require.Nil(t, err, "%v", err)
	err = module.Bind(SimpleStruct{}).ToSingleton(SimpleStruct{"hello"})
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)
	object, err := injector.Get((*SimpleInterface)(nil))
	require.Nil(t, err, "%v", err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimpleStructSingletonIndirectSucceedsWhenPtr(t *testing.T) {
	module := CreateModule()
	err := module.Bind((*SimpleInterface)(nil)).To(&SimpleStruct{})
	require.Nil(t, err, "%v", err)
	err = module.Bind(&SimpleStruct{}).ToSingleton(&SimpleStruct{"hello"})
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)
	object, err := injector.Get((*SimpleInterface)(nil))
	require.Nil(t, err, "%v", err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimplePtrStructSingletonIndirect(t *testing.T) {
	module := CreateModule()
	err := module.Bind((*SimpleInterface)(nil)).To(&SimplePtrStruct{})
	require.Nil(t, err, "%v", err)
	err = module.Bind(&SimplePtrStruct{}).ToSingleton(&SimplePtrStruct{"hello"})
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)
	object, err := injector.Get((*SimpleInterface)(nil))
	require.Nil(t, err, "%v", err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
}

func TestSimplePtrStructSingletonIndirectFailsWhenNotPtr(t *testing.T) {
	module := CreateModule()
	err := module.Bind((*SimpleInterface)(nil)).To(SimplePtrStruct{})
	require.Contains(t, err.Error(), InjectErrorTypeDoesNotImplement)
}

func TestSimplePtrStructIndirectFailsWhenNoFinalBinding(t *testing.T) {
	module := CreateModule()
	err := module.Bind((*SimpleInterface)(nil)).To(&SimplePtrStruct{})
	require.Nil(t, err, "%v", err)
	_, err = CreateInjector(module)()
	require.Contains(t, err.Error(), InjectErrorTypeNoBinding)
}

// ***** simple BindTagged tests *****

func TestTaggedFromNil(t *testing.T) {
	module := CreateModule()
	err := module.BindTagged(nil, "tagOne").ToSingleton(SimpleStruct{"hello"})
	require.Equal(t, ErrNil, err)
}

// TODO(pedge): this test is still useful to make sure if someone does not
// provide a pointer to an interface, it fails, but this is not the correct
// behavior, find a value that is not nil but reflect.TypeOf(...) returns nil
func TestTaggedFromReflectTypeNil(t *testing.T) {
	module := CreateModule()
	err := module.BindTagged((SimpleInterface)(nil), "tagOne").ToSingleton(SimpleStruct{"hello"})
	require.Equal(t, ErrNil, err)
}

func TestTaggedTagEmpty(t *testing.T) {
	module := CreateModule()
	err := module.BindTagged((*SimpleInterface)(nil), "").ToSingleton(SimpleStruct{"hello"})
	require.Contains(t, InjectErrorTypeTagEmpty, err.Error())
}

func TestTaggedSimpleStructSingletonDirect(t *testing.T) {
	module := CreateModule()
	err := module.BindTagged((*SimpleInterface)(nil), "tagOne").ToSingleton(SimpleStruct{"hello"})
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)
	object, err := injector.GetTagged((*SimpleInterface)(nil), "tagOne")
	require.Nil(t, err, "%v", err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	_, err = injector.Get((*SimpleInterface)(nil))
	require.Contains(t, err.Error(), InjectErrorTypeNoBinding)
	_, err = injector.GetTagged((*SimpleInterface)(nil), "tagTwo")
	require.Contains(t, err.Error(), InjectErrorTypeNoBinding)
}

func TestTaggedSimpleStructSingletonDirectSucceedsWhenPtr(t *testing.T) {
	module := CreateModule()
	err := module.BindTagged((*SimpleInterface)(nil), "tagOne").ToSingleton(&SimpleStruct{"hello"})
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)
	object, err := injector.GetTagged((*SimpleInterface)(nil), "tagOne")
	require.Nil(t, err, "%v", err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	_, err = injector.Get((*SimpleInterface)(nil))
	require.Contains(t, err.Error(), InjectErrorTypeNoBinding)
	_, err = injector.GetTagged((*SimpleInterface)(nil), "tagTwo")
	require.Contains(t, err.Error(), InjectErrorTypeNoBinding)
}

func TestTaggedSimplePtrStructSingletonDirect(t *testing.T) {
	module := CreateModule()
	err := module.BindTagged((*SimpleInterface)(nil), "tagOne").ToSingleton(&SimplePtrStruct{"hello"})
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)
	object, err := injector.GetTagged((*SimpleInterface)(nil), "tagOne")
	require.Nil(t, err, "%v", err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	_, err = injector.Get((*SimpleInterface)(nil))
	require.Contains(t, err.Error(), InjectErrorTypeNoBinding)
	_, err = injector.GetTagged((*SimpleInterface)(nil), "tagTwo")
	require.Contains(t, err.Error(), InjectErrorTypeNoBinding)
}

func TestTaggedSimplePtrStructSingletonDirectFailsWhenNotPtr(t *testing.T) {
	module := CreateModule()
	err := module.BindTagged((*SimpleInterface)(nil), "tagOne").ToSingleton(SimplePtrStruct{"hello"})
	require.Contains(t, err.Error(), InjectErrorTypeDoesNotImplement)
}

func TestTaggedSimpleStructSingletonIndirect(t *testing.T) {
	module := CreateModule()
	err := module.BindTagged((*SimpleInterface)(nil), "tagOne").To(SimpleStruct{})
	require.Nil(t, err, "%v", err)
	err = module.Bind(SimpleStruct{}).ToSingleton(SimpleStruct{"hello"})
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)
	object, err := injector.GetTagged((*SimpleInterface)(nil), "tagOne")
	require.Nil(t, err, "%v", err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	_, err = injector.Get((*SimpleInterface)(nil))
	require.Contains(t, err.Error(), InjectErrorTypeNoBinding)
	_, err = injector.GetTagged((*SimpleInterface)(nil), "tagTwo")
	require.Contains(t, err.Error(), InjectErrorTypeNoBinding)
}

func TestTaggedSimpleStructSingletonIndirectSucceedsWhenPtr(t *testing.T) {
	module := CreateModule()
	err := module.BindTagged((*SimpleInterface)(nil), "tagOne").To(&SimpleStruct{})
	require.Nil(t, err, "%v", err)
	err = module.Bind(&SimpleStruct{}).ToSingleton(&SimpleStruct{"hello"})
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)
	object, err := injector.GetTagged((*SimpleInterface)(nil), "tagOne")
	require.Nil(t, err, "%v", err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	_, err = injector.Get((*SimpleInterface)(nil))
	require.Contains(t, err.Error(), InjectErrorTypeNoBinding)
	_, err = injector.GetTagged((*SimpleInterface)(nil), "tagTwo")
	require.Contains(t, err.Error(), InjectErrorTypeNoBinding)
}

func TestTaggedSimplePtrStructSingletonIndirect(t *testing.T) {
	module := CreateModule()
	err := module.BindTagged((*SimpleInterface)(nil), "tagOne").To(&SimplePtrStruct{})
	require.Nil(t, err, "%v", err)
	err = module.Bind(&SimplePtrStruct{}).ToSingleton(&SimplePtrStruct{"hello"})
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)
	object, err := injector.GetTagged((*SimpleInterface)(nil), "tagOne")
	require.Nil(t, err, "%v", err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	_, err = injector.Get((*SimpleInterface)(nil))
	require.Contains(t, err.Error(), InjectErrorTypeNoBinding)
	_, err = injector.GetTagged((*SimpleInterface)(nil), "tagTwo")
	require.Contains(t, err.Error(), InjectErrorTypeNoBinding)
}

func TestTaggedSimplePtrStructSingletonIndirectFailsWhenNotPtr(t *testing.T) {
	module := CreateModule()
	err := module.BindTagged((*SimpleInterface)(nil), "tagOne").To(SimplePtrStruct{})
	require.Contains(t, err.Error(), InjectErrorTypeDoesNotImplement)
}

func TestTaggedSimplePtrStructIndirectFailsWhenNoFinalBinding(t *testing.T) {
	module := CreateModule()
	err := module.BindTagged((*SimpleInterface)(nil), "tagOne").To(&SimplePtrStruct{})
	require.Nil(t, err, "%v", err)
	_, err = CreateInjector(module)
	require.Contains(t, err.Error(), InjectErrorTypeNoBinding)
}

// ***** additional simple tagged tests *****

func TestTaggedSimpleStructSingletonDirectTwoBindings(t *testing.T) {
	module := CreateModule()
	err := module.BindTagged((*SimpleInterface)(nil), "tagOne").ToSingleton(SimpleStruct{"hello"})
	require.Nil(t, err, "%v", err)
	err = module.BindTagged((*SimpleInterface)(nil), SimpleTag{}).ToSingleton(SimpleStruct{"good day"})
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)
	object, err := injector.GetTagged((*SimpleInterface)(nil), "tagOne")
	require.Nil(t, err, "%v", err)
	simpleInterface := object.(SimpleInterface)
	require.Equal(t, "hello", simpleInterface.Foo())
	object, err = injector.GetTagged((*SimpleInterface)(nil), SimpleTag{})
	require.Nil(t, err, "%v", err)
	simpleInterface = object.(SimpleInterface)
	require.Equal(t, "good day", simpleInterface.Foo())
}

// ***** simple provider tests *****

type BarInterface interface {
	Bar() int
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

func createSecondInterface(s SimpleInterface, b BarInterface) (SecondInterface, error) {
	return &SecondPtrStruct{s, b}, nil
}

func createSecondInterfaceErr(s SimpleInterface, b BarInterface) (SecondInterface, error) {
	return nil, errors.New("XYZ")
}

func createSecondInterfaceInjector(injector Injector) (SecondInterface, error) {
	s, err := injector.Get((*SimpleInterface)(nil))
	if err != nil {
		return nil, err
	}
	b, err := injector.Get((*BarInterface)(nil))
	if err != nil {
		return nil, err
	}
	return &SecondPtrStruct{s.(SimpleInterface), b.(BarInterface)}, nil
}

func createSecondInterfaceInjectorErr(injector Injector) (SecondInterface, error) {
	return nil, errors.New("ABC")
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

func TestConstructorDirectInterfaceInjection(t *testing.T) {
	module := CreateModule()
	err := module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	require.Nil(t, err, "%v", err)
	err = module.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	require.Nil(t, err, "%v", err)
	err = module.Bind((*SecondInterface)(nil)).ToConstructor(createSecondInterface)
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)
	object, err := injector.Get((*SecondInterface)(nil))
	require.Nil(t, err, "%v", err)
	secondInterface := object.(SecondInterface)
	require.Equal(t, "hello", secondInterface.Foo().Foo())
	require.Equal(t, 1, secondInterface.Bar().Bar())
}

func TestConstructorInjectorInjection(t *testing.T) {
	module := CreateModule()
	err := module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	require.Nil(t, err, "%v", err)
	err = module.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	require.Nil(t, err, "%v", err)
	err = module.Bind((*SecondInterface)(nil)).ToConstructor(createSecondInterfaceInjector)
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)
	object, err := injector.Get((*SecondInterface)(nil))
	require.Nil(t, err, "%v", err)
	secondInterface := object.(SecondInterface)
	require.Equal(t, "hello", secondInterface.Foo().Foo())
	require.Equal(t, 1, secondInterface.Bar().Bar())
}

func TestConstructorErrReturned(t *testing.T) {
	module := CreateModule()
	err := module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	require.Nil(t, err, "%v", err)
	err = module.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	require.Nil(t, err, "%v", err)
	err = module.Bind((*SecondInterface)(nil)).ToConstructor(createSecondInterfaceErr)
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)
	_, err = injector.Get((*SecondInterface)(nil))
	require.Equal(t, "XYZ", err.Error())
}

func TestConstructorInjectorErrReturned(t *testing.T) {
	module := CreateModule()
	err := module.Bind((*SimpleInterface)(nil)).ToSingleton(&SimplePtrStruct{"hello"})
	require.Nil(t, err, "%v", err)
	err = module.Bind((*BarInterface)(nil)).ToSingleton(&BarPtrStruct{1})
	require.Nil(t, err, "%v", err)
	err = module.Bind((*SecondInterface)(nil)).ToConstructor(createSecondInterfaceInjectorErr)
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)
	_, err = injector.Get((*SecondInterface)(nil))
	require.Equal(t, "ABC", err.Error())
}

func TestConstructorWithEvilCounter(t *testing.T) {
	module := CreateModule()
	err := module.Bind((*BarInterface)(nil)).ToConstructor(createEvilBarInterface)
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)

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
	err := module.Bind((*BarInterface)(nil)).ToSingletonConstructor(createEvilBarInterface)
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)

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
	err := module.Bind((*BarInterface)(nil)).ToSingletonConstructor(createEvilBarInterfaceErr)
	require.Nil(t, err, "%v", err)
	injector, err := CreateInjector(module)
	require.Nil(t, err, "%v", err)

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
