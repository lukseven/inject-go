package inject

import (
	"fmt"
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
	Float32
	Float64
	Complex64
	Complex128
	String
)

var (
	boolReflectType       = reflect.TypeOf(false)
	intReflectType        = reflect.TypeOf(int(0))
	int8ReflectType       = reflect.TypeOf(int8(0))
	int16ReflectType      = reflect.TypeOf(int16(0))
	int32ReflectType      = reflect.TypeOf(int32(0))
	int64ReflectType      = reflect.TypeOf(int64(0))
	uintReflectType       = reflect.TypeOf(uint(0))
	uint8ReflectType      = reflect.TypeOf(uint8(0))
	uint16ReflectType     = reflect.TypeOf(uint16(0))
	uint32ReflectType     = reflect.TypeOf(uint32(0))
	uint64ReflectType     = reflect.TypeOf(uint64(0))
	float32ReflectType    = reflect.TypeOf(float32(0))
	float64ReflectType    = reflect.TypeOf(float64(0))
	complex64ReflectType  = reflect.TypeOf(complex64(0i))
	complex128ReflectType = reflect.TypeOf(complex128(0i))
	stringReflectType     = reflect.TypeOf("")

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
		Float32:    reflect.Float32,
		Float64:    reflect.Float64,
		Complex64:  reflect.Complex64,
		Complex128: reflect.Complex128,
		String:     reflect.String,
	}
	lenConstantKindToReflectKind = len(constantKindToReflectKind)

	constantKindToReflectType = map[ConstantKind]reflect.Type{
		Bool:       boolReflectType,
		Int:        intReflectType,
		Int8:       int8ReflectType,
		Int16:      int16ReflectType,
		Int32:      int32ReflectType,
		Int64:      int64ReflectType,
		Uint:       uintReflectType,
		Uint8:      uint8ReflectType,
		Uint16:     uint16ReflectType,
		Uint32:     uint32ReflectType,
		Uint64:     uint64ReflectType,
		Float32:    float32ReflectType,
		Float64:    float64ReflectType,
		Complex64:  complex64ReflectType,
		Complex128: complex128ReflectType,
		String:     stringReflectType,
	}
	lenConstantKindToReflectType = len(constantKindToReflectType)

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
		reflect.Float32:    Float32,
		reflect.Float64:    Float64,
		reflect.Complex64:  Complex64,
		reflect.Complex128: Complex128,
		reflect.String:     String,
	}

	reflectTypeToConstantKind = map[reflect.Type]ConstantKind{
		boolReflectType:       Bool,
		intReflectType:        Int,
		int8ReflectType:       Int8,
		int16ReflectType:      Int16,
		int32ReflectType:      Int32,
		int64ReflectType:      Int64,
		uintReflectType:       Uint,
		uint8ReflectType:      Uint8,
		uint16ReflectType:     Uint16,
		uint32ReflectType:     Uint32,
		uint64ReflectType:     Uint64,
		float32ReflectType:    Float32,
		float64ReflectType:    Float64,
		complex64ReflectType:  Complex64,
		complex128ReflectType: Complex128,
		stringReflectType:     String,
	}
)

type ConstantKind uint

func (this ConstantKind) ReflectKind() reflect.Kind {
	if int(this) < lenConstantKindToReflectKind {
		return constantKindToReflectKind[this]
	}
	panic(fmt.Sprintf("inject: Unknown ConstantKind: %v", this))
}

func (this ConstantKind) ReflectType() reflect.Type {
	if int(this) < lenConstantKindToReflectType {
		return constantKindToReflectType[this]
	}
	panic(fmt.Sprintf("inject: Unknown ConstantKind: %v", this))
}

func constantKindForReflectKind(reflectKind reflect.Kind) (ConstantKind, bool) {
	constantKind, ok := reflectKindToConstantKind[reflectKind]
	return constantKind, ok
}

func constantKindForReflectType(reflectType reflect.Type) (ConstantKind, bool) {
	constantKind, ok := reflectTypeToConstantKind[reflectType]
	return constantKind, ok
}
