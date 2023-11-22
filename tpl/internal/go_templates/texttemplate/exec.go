package texttemplate

import (
  "fmt"
  "reflect"

	"github.com/SkyKoo/hugo-reduce/common/hreflect"
)

var (
  errorType = reflect.TypeOf((*error)(nil)).Elem()
  fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
  reflectValueType = reflect.TypeOf((*reflect.Value)(nil)).Elem()
)

func isTrue(val reflect.Value) (truth, ok bool) {
  return hreflect.IsTruthfulValue(val), true
}
