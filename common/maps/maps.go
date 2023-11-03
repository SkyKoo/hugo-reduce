package maps

import (
  "fmt"
  "github.com/spf13/cast"
  "github.com/SkyKoo/hugo-reduce/types"
)

// ToParamsAndPrepare converts in to Params and prepares it for use.
// If in is nil, an empty map is returned.
// See PrepareParams.
func ToParamsAndPrepare(in any) (Params, bool) {
  if types.IsNil(in) {
    return Params{}, true
  }
  m, err := ToStringMapE(in)
  if err != nil {
    return nil, false
  }
  PrepareParams(m)
  return m, true
}

// ToStringMapE converts in to map[string]interface{}.
func ToStringMapE(in any) (map[string]any, error) {
  switch vv := in.(type) {
  case Params:
    return vv, nil
  case map[string]string:
    var m = map[string]any{}
    for k, v := range vv {
      m[k] = v
    }
    return m, nil

  default:
    return cast.ToStringMapE(in)
  }
}

// MustToParamsAndPrepare calls ToParamsAndPrepare and panic if it fails.
func MustToParamsAndPrepare(in any) Params {
  if p, ok := ToParamsAndPrepare(in); ok {
    return p
  } else {
    panic(fmt.Sprintf("cannot convert %T to maps.Params", in))
  }
}
