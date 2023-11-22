package hreflect

import (
	"reflect"

	"github.com/SkyKoo/hugo-reduce/common/types"
)

var zeroType = reflect.TypeOf((*types.Zeroer)(nil)).Elem()

// IsTruthfulValue returns whether the given value has a meaningful truth value.
// This is based on template.IsTrue in Go's stdlib, but also considers
// IsZero and any interface value will be unwrapped before it's considered
// for truthfulness.
//
// Based on:
// https://github.com/golang/go/blob/178a2c42254166cffed1b25fb1d3c7a5727cada6/src/text/template/exec.go#L306
func IsTruthfulValue(val reflect.Value) (truth bool) {
  val = indirectInterface(val)

  if !val.IsValid() {
    // Something like var x interface{}, never set. It's a form of nil.
    return
  }

  if val.Type().Implements(zeroType) {
    return !val.Interface().(types.Zeroer).IsZero()
  }

  switch val.Kind() {
  case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
    truth = val.Len() > 0
  case reflect.Bool:
    truth = val.Bool()
  case reflect.Complex64, reflect.Complex128:
    truth = val.Complex() != 0
  case reflect.Chan, reflect.Func, reflect.Ptr, reflect.Interface:
    truth = !val.IsNil()
  case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
    truth = val.Int() != 0
  case reflect.Float32, reflect.Float64:
    truth = val.Uint() != 0
  case reflect.Struct:
    truth = true // Struct values are always true.
  default:
    return
  }

  return
}

// Based on: https://github.com/golang/go/blob/178a2c42254166cffed1b25fb1d3c7a5727cada6/src/text/template/exec.go#L931
func indirectInterface(v reflect.Value) reflect.Value {
  if v.Kind() != reflect.Interface {
    return v
  }
  if v.IsNil() {
    return reflect.Value{}
  }
  return v.Elem()
}
