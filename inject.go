package inject

const (
	InjectErrorTypeNil                            = "Parameter is nil"
	InjectErrorTypeReflectTypeNil                 = "reflect.TypeOf() returns nil"
	InjectErrorTypeUnknownBinderType              = "Unknown binder type"
	InjectErrorTypeUnknownValidBinderType         = "Unknown valid binder type"
	InjectErrorTypeNotInterfacePtr                = "Binding with Binder.ToType() and from is not an interface pointer"
	InjectErrorTypeDoesNotImplement               = "to binding does not implement from binding"
	InjectErrorTypeNotSupportedYet                = "Binding type not supported yet, feel free to help!"
	InjectErrorTypeNotAssignable                  = "Binding not assignable"
	InjectErrorTypeConstructorNotFunction         = "Constructor is not a function"
	InjectErrorTypeConstructorReturnValuesInvalid = "Constructor can only have two return values, the first providing the value, the second being an error"
	InjectErrorTypeInvalidReturnFromConstructor   = "Invalid return values from constructor"
	InjectErrorTypeBindingTypeIntermediate        = "Binding type is intermediate"
	InjectErrorTypeBindingTypeUnknown             = "Binding type is unknown"
	InjectErrorTypeFinalBindingTypeUnknown        = "Final binding type is unknown"
	InjectErrorTypeIntermediateBinding            = "Trying to get for an intermediate binding"
	InjectErrorTypeFinalBinding                   = "Trying to get bindingKey for a final binding"
	InjectErrorTypeCannotCastModule               = "Cannot cast Module to internal module type"
	InjectErrorTypeNoBinding                      = "No binding for binding key"
	InjectErrorTypeNoFinalBinding                 = "No final binding for binding key"
	InjectErrorTypeAlreadyBound                   = "Already found a binding for this binding key"
	InjectErrorTypeTagEmpty                       = "Tag empty"
)

type Module interface {
	Bind(from interface{}) Builder
	BindTagged(from interface{}, tag string) Builder
}

func CreateModule() Module { return createModule() }

type Builder interface {
	To(to interface{}) error
	ToSingleton(singleton interface{}) error
	ToConstructor(constructor interface{}) error
	ToSingletonConstructor(constructor interface{}) error
	ToTaggedConstructor(constructor interface{}) error
	ToTaggedSingletonConstructor(constructor interface{}) error
}

type Injector interface {
	Get(from interface{}) (interface{}, error)
	GetTagged(from interface{}, tag string) (interface{}, error)
}

func CreateInjector(modules ...Module) (Injector, error) { return createInjector(modules) }
