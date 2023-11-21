package texttemplate

import (
	"sync"
	"reflect"
  "github.com/SkyKoo/hugo-reduce/tpl/internal/go_templates/texttemplate/parse"
)

// common holds the information shared by related templates.
type common struct {
  tmpl map[string]*Template // Map from name to defined templates.
  muTmpl sync.RWMutex // protects tmpl
  option option
  // We use two maps, one for parsing and one for execution.
  // This separation makes the API cleaner since it doesn't
  // expose reflection to the client.
  muFuncs sync.RWMutex // protects parseFuncs and execFuncs
  parseFuncs FuncMap
  execFuncs map[string]reflect.Value
}

// Template is the representation of a parsed template. The *parse.Tree
// field is exported only for use by html/template and should be treated
// as unexported by all other client.
type Template struct {
  name string
  *parse.Tree
  *common
  leftDelim string
  rightDelim string
}
