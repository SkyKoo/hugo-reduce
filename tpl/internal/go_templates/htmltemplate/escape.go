package template

import (
  template "github.com/SkyKoo/hugo-reduce/tpl/internal/go_templates/texttemplate"
)

// funcMap maps command names to functions that render their inputs safe.
var funcMap = template.FuncMap{
  "_html_template_htmlescaper": htmlEscaper,
}
