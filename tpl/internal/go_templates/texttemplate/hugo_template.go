package texttemplate

import (
  "io"
  "context"
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
