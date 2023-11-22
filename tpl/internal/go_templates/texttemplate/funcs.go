package texttemplate

import (
	"fmt"
	"reflect"
	"sync"
	"unicode"
)

type FuncMap map[string]any

// builtins returns the FuncMap.
// It is not a global variable so the linker can dead code eliminate
// more when this isn't called. See golang.org/issue/36021.
// TODO: revert this back to a global map once golang.org/issue/2559 is fixed.
func builtins() FuncMap {
  return FuncMap{
    "print": fmt.Sprint,
    "printf": fmt.Sprintf,
    "println": fmt.Sprintln,
    "not": not,
  }
}

var builtinFuncsOnce struct {
  sync.Once
  v map[string]reflect.Value
}

// builtFuncsOnce lazily compute & caches the builtinFuncs map.
// TODO: revert this back to a global map once golang.org/issue/2559 is fixed.
func builtinFuncs() map[string]reflect.Value {
  builtinFuncs.Do(func() {
    builtinFuncsOnce.v = createValueFuncs(builtins())
  })
  return builtinFuncsOnce.v
}

// createValueFuncs turns a FuncMap into a map[string]reflect.Value
func createValueFuncs(funcMap FuncMap) map[string]reflect.Value {
  m := make(map[string]reflect.Value)
  addValueFuncs(m, funcMap)
  return m
}

// addValueFuncs adds to values the functions in funcs, converting them to reflect.Values.
func addValueFuncs(out map[string]reflect.Value, in FuncMap) {
  for name, fn := range in {
    if !goodName(name) {
      panic(fmt.Errorf("function name %q is not a value identifier", name))
    }
    v := reflect.ValueOf(fn)
    if v.Kind() != reflect.Func {
      panic("value for " + name + " not a function")
    }
    if !goodFunc(v.Type()) {
      panic(fmt.Errorf("can't install method/function %q with %d results", name, v.Type().NumOut()))
    }
    out[name] = v
  }
}

// goodFunc report whether the function or method has the right result signature.
func goodFunc(typ reflect.Type) bool {
  // We allow functions with 1 result or 2 results where the second is an error.
  switch {
  case typ.NumOut() == 1:
    return true
  case typ.NumOut() == 2 && typ.Out(1) == errorType:
    return true
  }
  return false
}

// goodName reports whether the function name is a valid identifier.
func goodName(name string) bool {
  if name == "" {
    return false
  }
  for i, r := range name {
    switch {
    case r == '_':
    case i == 0 && !unicode.IsLetter(r):
      return false
    case !unicode.IsLetter(r) && !unicode.IsDigit(r):
      return false
    }
  }
  return true
}

func truth(arg reflect.Value) bool {
  t, _ := isTrue(indirectInterface(arg))
  return t
}

// not returns the Boolean negation of this argument.
func not(arg reflect.Value) bool {
  return !truth(arg)
}
