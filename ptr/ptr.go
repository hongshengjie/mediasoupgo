package ptr

import (
	"fmt"
	"reflect"
	"time"
)

// AllPtrFieldsNil tests whether all pointer fields in a struct are nil.  This is useful when,
// for example, an API struct is handled by plugins which need to distinguish
// "no plugin accepted this spec" from "this spec is empty".
//
// This function is only valid for structs and pointers to structs.  Any other
// type will cause a panic.  Passing a typed nil pointer will return true.
func AllPtrFieldsNil(obj interface{}) bool {
	v := reflect.ValueOf(obj)
	if !v.IsValid() {
		panic(fmt.Sprintf("reflect.ValueOf() produced a non-valid Value for %#v", obj))
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return true
		}
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() == reflect.Ptr && !v.Field(i).IsNil() {
			return false
		}
	}
	return true
}

// To returns a pointer to the given value.
func To[T any](v T) *T {
	return &v
}

// Deref dereferences ptr and returns the value it points to if no nil, or else
// returns def.
func Deref[T any](ptr *T, def T) T {
	if ptr != nil {
		return *ptr
	}
	return def
}

// Equal returns true if both arguments are nil or both arguments
// dereference to the same value.
func Equal[T comparable](a, b *T) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	if a == nil {
		return true
	}
	return *a == *b
}

// Int returns a pointer to an int.
var Int = To[int]

// IntDeref dereferences the int ptr and returns it if not nil, or else
// returns def.
var IntDeref = Deref[int]

// Int32 returns a pointer to an int32.
var Int32 = To[int32]

// Int32Deref dereferences the int32 ptr and returns it if not nil, or else
// returns def.
var Int32Deref = Deref[int32]

// Int32Equal returns true if both arguments are nil or both arguments
// dereference to the same value.
var Int32Equal = Equal[int32]

// Uint returns a pointer to an uint
var Uint = To[uint]

// UintDeref dereferences the uint ptr and returns it if not nil, or else
// returns def.
var UintDeref = Deref[uint]

// Uint32 returns a pointer to an uint32.
var Uint32 = To[uint32]

// Uint32Deref dereferences the uint32 ptr and returns it if not nil, or else
// returns def.
var Uint32Deref = Deref[uint32]

// Uint32Equal returns true if both arguments are nil or both arguments
// dereference to the same value.
var Uint32Equal = Equal[uint32]

// Int64 returns a pointer to an int64.
var Int64 = To[int64]

// Int64Deref dereferences the int64 ptr and returns it if not nil, or else
// returns def.
var Int64Deref = Deref[int64]

// Int64Equal returns true if both arguments are nil or both arguments
// dereference to the same value.
var Int64Equal = Equal[int64]

// Uint64 returns a pointer to an uint64.
var Uint64 = To[uint64]

// Uint64Deref dereferences the uint64 ptr and returns it if not nil, or else
// returns def.
var Uint64Deref = Deref[uint64]

// Uint64Equal returns true if both arguments are nil or both arguments
// dereference to the same value.
var Uint64Equal = Equal[uint64]

// Bool returns a pointer to a bool.
var Bool = To[bool]

// BoolDeref dereferences the bool ptr and returns it if not nil, or else
// returns def.
var BoolDeref = Deref[bool]

// BoolEqual returns true if both arguments are nil or both arguments
// dereference to the same value.
var BoolEqual = Equal[bool]

// String returns a pointer to a string.
var String = To[string]

// StringDeref dereferences the string ptr and returns it if not nil, or else
// returns def.
var StringDeref = Deref[string]

// StringEqual returns true if both arguments are nil or both arguments
// dereference to the same value.
var StringEqual = Equal[string]

// Float32 returns a pointer to a float32.
var Float32 = To[float32]

// Float32Deref dereferences the float32 ptr and returns it if not nil, or else
// returns def.
var Float32Deref = Deref[float32]

// Float32Equal returns true if both arguments are nil or both arguments
// dereference to the same value.
var Float32Equal = Equal[float32]

// Float64 returns a pointer to a float64.
var Float64 = To[float64]

// Float64Deref dereferences the float64 ptr and returns it if not nil, or else
// returns def.
var Float64Deref = Deref[float64]

// Float64Equal returns true if both arguments are nil or both arguments
// dereference to the same value.
var Float64Equal = Equal[float64]

// Duration returns a pointer to a time.Duration.
var Duration = To[time.Duration]

// DurationDeref dereferences the time.Duration ptr and returns it if not nil, or else
// returns def.
var DurationDeref = Deref[time.Duration]

// DurationEqual returns true if both arguments are nil or both arguments
// dereference to the same value.
var DurationEqual = Equal[time.Duration]
