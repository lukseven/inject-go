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

	boolConstant       = false
	intConstant        = int(0)
	int8Constant       = int8(0)
	int16Constant      = int16(0)
	int32Constant      = int32(0)
	int64Constant      = int64(0)
	uintConstant       = uint(0)
	uint8Constant      = uint8(0)
	uint16Constant     = uint16(0)
	uint32Constant     = uint32(0)
	uint64Constant     = uint64(0)
	float32Constant    = float32(0)
	float64Constant    = float64(0)
	complex64Constant  = complex64(0i)
	complex128Constant = complex128(0i)
	stringConstant     = ""
)

var (
	boolReflectType       = reflect.TypeOf(boolConstant)
	intReflectType        = reflect.TypeOf(intConstant)
	int8ReflectType       = reflect.TypeOf(int8Constant)
	int16ReflectType      = reflect.TypeOf(int16Constant)
	int32ReflectType      = reflect.TypeOf(int32Constant)
	int64ReflectType      = reflect.TypeOf(int64Constant)
	uintReflectType       = reflect.TypeOf(uintConstant)
	uint8ReflectType      = reflect.TypeOf(uint8Constant)
	uint16ReflectType     = reflect.TypeOf(uint16Constant)
	uint32ReflectType     = reflect.TypeOf(uint32Constant)
	uint64ReflectType     = reflect.TypeOf(uint64Constant)
	float32ReflectType    = reflect.TypeOf(float32Constant)
	float64ReflectType    = reflect.TypeOf(float64Constant)
	complex64ReflectType  = reflect.TypeOf(complex64Constant)
	complex128ReflectType = reflect.TypeOf(complex128Constant)
	stringReflectType     = reflect.TypeOf(stringConstant)

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

	constantKindToConstant = map[ConstantKind]interface{}{
		Bool:       boolConstant,
		Int:        intConstant,
		Int8:       int8Constant,
		Int16:      int16Constant,
		Int32:      int32Constant,
		Int64:      int64Constant,
		Uint:       uintConstant,
		Uint8:      uint8Constant,
		Uint16:     uint16Constant,
		Uint32:     uint32Constant,
		Uint64:     uint64Constant,
		Float32:    float32Constant,
		Float64:    float64Constant,
		Complex64:  complex64Constant,
		Complex128: complex128Constant,
		String:     stringConstant,
	}
	lenConstantKindToConstant = len(constantKindToConstant)
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

func (this ConstantKind) constant() interface{} {
	if int(this) < lenConstantKindToConstant {
		return constantKindToConstant[this]
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
