package maps

import (
  "github.com/spf13/cast"
  "strings"
)

// Params is a map where all keys are lower case.
type Params map[string]any

// ParamsMergeStrategy tells what strategy to use in Params.Merge.
type ParamsMergeStrategy string

// Set overwrites values in p with values in pp for common or new keys.
// This is done recursively.
func (p Params) Set(pp Params) {
  for k, v := range pp {
    vv, found := p[k]
    if !found {
      p[k] = v
    } else {
      switch vvv := vv.(type) {
      case Params:
        if pv, ok := v.(Params); ok {
          vvv.Set(pv)
        } else {
          p[k] = v
        }
      default:
        p[k] = v
      }
    }
  }
}

const (
  // ParamsMergeStrategyNone Do not merge.
  ParamsMergeStrategyNone ParamsMergeStrategy = "none"
  // ParamsMergeStrategyShallow Only and new keys.
  ParamsMergeStrategyShallow ParamsMergeStrategy = "shallow"
  // ParamsMergeStrategyDeep Add new keys, merge existing..
  ParamsMergeStrategyDeep ParamsMergeStrategy = "deep"

  mergeStrageyKey = "_merge"
)

func toMergeStragegy(v any) ParamsMergeStrategy {
  s := ParamsMergeStrategy(cast.ToString(v))
  switch s {
  case ParamsMergeStrategyDeep, ParamsMergeStrategyNone, ParamsMergeStrategyShallow:
    return s
  default:
    return ParamsMergeStrategyDeep
  }
}

// PrepareParams
// * makes all the keys in the given map lower cased and will do so
// * This will modify the map given.
// * Any nested map[interface{}]interface{}, map[string]interface{}, map[string]string will be converted to Params.
// * Any _merge value will be converted to proper type and value.
func PrepareParams(m Params) {
  for k, v := range m {
    var retyped bool
    lKey := strings.ToLower(k)
    if lKey == mergeStrageyKey {
      v = toMergeStragegy(v)
      retyped = true
    } else {
      switch vv := v.(type) {
      case map[any]any:
        var p Params = cast.ToStringMap(v)
        v = p
        PrepareParams(p) // recurse
        retyped = true
      case map[string]any:
        var p Params = v.(map[string]any)
        v = p
        PrepareParams(p)
        retyped = true
      case map[string]string:
        p := make(Params)
        for k, v := range vv {
          p[k] = v
        }
        v = p
        PrepareParams(p)
        retyped = true
      }
    }

    if retyped || k != lKey {
      delete(m, k)
      m[lKey] = v
    }
  }
}
