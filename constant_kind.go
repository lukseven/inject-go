package inject

import (
	"reflect"
)

const (
	Bool ConstantKind = iota
	Int
	Int8
	Int16
	Int32
	Int64
	Uint
	Uint8
	Uint16
	Uint32
	Uint64
	Uintptr
	Float32
	Float64
	Complex64
	Complex128
	String
)

var (
	reflectKindToConstantKind = map[reflect.Kind]ConstantKind{
		reflect.Bool:       Bool,
		reflect.Int:        Int,
		reflect.Int8:       Int8,
		reflect.Int16:      Int16,
		reflect.Int32:      Int32,
		reflect.Int64:      Int64,
		reflect.Uint:       Uint,
		reflect.Uint8:      Uint8,
		reflect.Uint16:     Uint16,
		reflect.Uint32:     Uint32,
		reflect.Uint64:     Uint64,
		reflect.Uintptr:    Uintptr,
		reflect.Float32:    Float32,
		reflect.Float64:    Float64,
		reflect.Complex64:  Complex64,
		reflect.Complex128: Complex128,
		reflect.String:     String,
	}
	constantKindToReflectKind = map[ConstantKind]reflect.Kind{
		Bool:       reflect.Bool,
		Int:        reflect.Int,
		Int8:       reflect.Int8,
		Int16:      reflect.Int16,
		Int32:      reflect.Int32,
		Int64:      reflect.Int64,
		Uint:       reflect.Uint,
		Uint8:      reflect.Uint8,
		Uint16:     reflect.Uint16,
		Uint32:     reflect.Uint32,
		Uint64:     reflect.Uint64,
		Uintptr:    reflect.Uintptr,
		Float32:    reflect.Float32,
		Float64:    reflect.Float64,
		Complex64:  reflect.Complex64,
		Complex128: reflect.Complex128,
		String:     reflect.String,
	}
	lenConstantKindToReflectKind = len(constantKindToReflectKind)
)

type ConstantKind uint

func ConstantKindOf(reflectKind reflect.Kind) ConstantKind {
	constantKind, ok := reflectKindToConstantKind[reflectKind]
	if !ok {
		panic("inject: No ConstantKind for reflect.Kind")
	}
	return constantKind
}

func (this ConstantKind) ReflectKind() reflect.Kind {
	if int(this) < lenConstantKindToReflectKind {
		return constantKindToReflectKind[this]
	}
	panic("inject: Unknown ConstantKind")
}
