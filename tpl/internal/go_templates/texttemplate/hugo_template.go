package texttemplate

import (
	"context"
	"io"
	"reflect"
)

// Preparer prepares the template before execution.
type Preparer interface {
  Prepare() (*Template, error)
}

// Executer executes a given template.
type Executer interface {
  ExecuteWithContext(ctx context.Context, p Preparer, wr io.Writer, data any) error
}

// Export it so we can populate Hugo's func map with it, which makes it faster.
var GoFuncs = builtinFuncs()

func NewExecuter(helper ExecHelper) Executer {
  return &executer{helper: helper}
}

// ExecHelper allows some custom eval hooks.
type ExecHelper interface {
  Init(ctx context.Context, tmpl Preparer)
  GetFunc(ctx context.Context, tmpl Preparer, name string) (reflect.Value, reflect.Value, bool)
  GetMethod(ctx context.Context, tmpl Preparer, receiver reflect.Value, name string) (method reflect.Value, firstArg reflect.Value)
  GetMapValue(ctx context.Context, tmpl Preparer, receiver, key reflect.Value) (reflect.Value, bool)
}

type executer struct {
  helper ExecHelper
}
