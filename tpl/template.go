package tpl

import (
  "context"
  "io"
  "reflect"

  "github.com/SkyKoo/hugo-reduce/output"
  "github.com/SkyKoo/hugo-reduce/tpl/internal/go_templates/texttemplate"
)

// TemplateManager manages the collection of templates.
type TemplateManager interface {
  TemplateHandler
  TemplateFuncGetter
  AddTemplate(name, tpl string) error
  MarkReady() error
}

// TemplateHandler finds and executes templates.
type TemplateHandler interface {
  TemplateFinder
  Execute(t Template, wr io.Writer, data any) error
  ExecuteWithContext(ctx context.Context, t Template, wr io.Writer, data any) error
  LookupLayout(d output.LayoutDescriptor, f output.Format) (Template, bool, error)
  HasTemplate(name string) bool
}

// Template is the common interface between text/template and html/template.
type Template interface {
  Name() string
  Prepare() (*texttemplate.Template, error)
}

// TemplateFinder finds templates.
type TemplateFinder interface {
  TemplateLookup
}

// TemplateFuncGetter allows to find a template func by name.
type TemplateFuncGetter interface {
  GetFunc(name string) (reflect.Value, bool)
}

// TemplateParseFinder provides both parsing and finding.
type TemplateParseFinder interface {
  TemplateParser
  TemplateFinder
}

// TemplateParser is used to parse ad-hoc templates, e.g. in the Resource chain.
type TemplateParser interface {
  Parse(name, tpl string) (Template, error)
}

// TemplateVariants describes the possible variants of a template.
// All of these may be empty.
type TemplateVariants struct {
  Language string
  OutputFormat output.Format
}

type TemplateLookup interface {
  Lookup(name string) (Template, bool)
}
