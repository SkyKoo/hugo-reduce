package hreflect

import (
  "reflect"
)

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

  return
}
