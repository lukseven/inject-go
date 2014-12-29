package inject

import (
	"fmt"
	"reflect"
)

const (
	boolConstantKind constantKind = iota
	intConstantKind
	int8ConstantKind
	int16ConstantKind
	int32ConstantKind
	int64ConstantKind
	uintConstantKind
	uint8ConstantKind
	uint16ConstantKind
	uint32ConstantKind
	uint64ConstantKind
	float32ConstantKind
	float64ConstantKind
	complex64ConstantKind
	complex128ConstantKind
	stringConstantKind

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

	constantKindToReflectKind = map[constantKind]reflect.Kind{
		boolConstantKind:       reflect.Bool,
		intConstantKind:        reflect.Int,
		int8ConstantKind:       reflect.Int8,
		int16ConstantKind:      reflect.Int16,
		int32ConstantKind:      reflect.Int32,
		int64ConstantKind:      reflect.Int64,
		uintConstantKind:       reflect.Uint,
		uint8ConstantKind:      reflect.Uint8,
		uint16ConstantKind:     reflect.Uint16,
		uint32ConstantKind:     reflect.Uint32,
		uint64ConstantKind:     reflect.Uint64,
		float32ConstantKind:    reflect.Float32,
		float64ConstantKind:    reflect.Float64,
		complex64ConstantKind:  reflect.Complex64,
		complex128ConstantKind: reflect.Complex128,
		stringConstantKind:     reflect.String,
	}
	lenConstantKindToReflectKind = len(constantKindToReflectKind)

	constantKindToReflectType = map[constantKind]reflect.Type{
		boolConstantKind:       boolReflectType,
		intConstantKind:        intReflectType,
		int8ConstantKind:       int8ReflectType,
		int16ConstantKind:      int16ReflectType,
		int32ConstantKind:      int32ReflectType,
		int64ConstantKind:      int64ReflectType,
		uintConstantKind:       uintReflectType,
		uint8ConstantKind:      uint8ReflectType,
		uint16ConstantKind:     uint16ReflectType,
		uint32ConstantKind:     uint32ReflectType,
		uint64ConstantKind:     uint64ReflectType,
		float32ConstantKind:    float32ReflectType,
		float64ConstantKind:    float64ReflectType,
		complex64ConstantKind:  complex64ReflectType,
		complex128ConstantKind: complex128ReflectType,
		stringConstantKind:     stringReflectType,
	}
	lenConstantKindToReflectType = len(constantKindToReflectType)

	reflectKindToConstantKind = map[reflect.Kind]constantKind{
		reflect.Bool:       boolConstantKind,
		reflect.Int:        intConstantKind,
		reflect.Int8:       int8ConstantKind,
		reflect.Int16:      int16ConstantKind,
		reflect.Int32:      int32ConstantKind,
		reflect.Int64:      int64ConstantKind,
		reflect.Uint:       uintConstantKind,
		reflect.Uint8:      uint8ConstantKind,
		reflect.Uint16:     uint16ConstantKind,
		reflect.Uint32:     uint32ConstantKind,
		reflect.Uint64:     uint64ConstantKind,
		reflect.Float32:    float32ConstantKind,
		reflect.Float64:    float64ConstantKind,
		reflect.Complex64:  complex64ConstantKind,
		reflect.Complex128: complex128ConstantKind,
		reflect.String:     stringConstantKind,
	}

	reflectTypeToConstantKind = map[reflect.Type]constantKind{
		boolReflectType:       boolConstantKind,
		intReflectType:        intConstantKind,
		int8ReflectType:       int8ConstantKind,
		int16ReflectType:      int16ConstantKind,
		int32ReflectType:      int32ConstantKind,
		int64ReflectType:      int64ConstantKind,
		uintReflectType:       uintConstantKind,
		uint8ReflectType:      uint8ConstantKind,
		uint16ReflectType:     uint16ConstantKind,
		uint32ReflectType:     uint32ConstantKind,
		uint64ReflectType:     uint64ConstantKind,
		float32ReflectType:    float32ConstantKind,
		float64ReflectType:    float64ConstantKind,
		complex64ReflectType:  complex64ConstantKind,
		complex128ReflectType: complex128ConstantKind,
		stringReflectType:     stringConstantKind,
	}

	constantKindToConstant = map[constantKind]interface{}{
		boolConstantKind:       boolConstant,
		intConstantKind:        intConstant,
		int8ConstantKind:       int8Constant,
		int16ConstantKind:      int16Constant,
		int32ConstantKind:      int32Constant,
		int64ConstantKind:      int64Constant,
		uintConstantKind:       uintConstant,
		uint8ConstantKind:      uint8Constant,
		uint16ConstantKind:     uint16Constant,
		uint32ConstantKind:     uint32Constant,
		uint64ConstantKind:     uint64Constant,
		float32ConstantKind:    float32Constant,
		float64ConstantKind:    float64Constant,
		complex64ConstantKind:  complex64Constant,
		complex128ConstantKind: complex128Constant,
		stringConstantKind:     stringConstant,
	}
	lenConstantKindToConstant = len(constantKindToConstant)
)

type constantKind uint

func (this constantKind) reflectKind() reflect.Kind {
	if int(this) < lenConstantKindToReflectKind {
		return constantKindToReflectKind[this]
	}
	panic(fmt.Sprintf("inject: Unknown constantKind: %v", this))
}

func (this constantKind) reflectType() reflect.Type {
	if int(this) < lenConstantKindToReflectType {
		return constantKindToReflectType[this]
	}
	panic(fmt.Sprintf("inject: Unknown constantKind: %v", this))
}

func (this constantKind) constant() interface{} {
	if int(this) < lenConstantKindToConstant {
		return constantKindToConstant[this]
	}
	panic(fmt.Sprintf("inject: Unknown constantKind: %v", this))
}

func constantKindForReflectKind(reflectKind reflect.Kind) (constantKind, bool) {
	value, ok := reflectKindToConstantKind[reflectKind]
	return value, ok
}

func constantKindForReflectType(reflectType reflect.Type) (constantKind, bool) {
	value, ok := reflectTypeToConstantKind[reflectType]
	return value, ok
}
