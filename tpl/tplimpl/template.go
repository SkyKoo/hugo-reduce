package tplimpl

import (
  "fmt"

	"github.com/SkyKoo/hugo-reduce/deps"
	"github.com/SkyKoo/hugo-reduce/log"
)

func newTemplateExec(d *deps.Deps) (*templateExec, error) {
  exec, funcs := newTmmplateExecuter(d)
  funcMap := make(map[string]any)
  for k, v := range funcs {
    funcMap[k] = v.Interface()
  }

  log.Process("newTemplateNamespace", "with funcMap")
  log.Process("NewLayoutHandler", "to process layout request")
  h := &templateHandler{
    main: newTemplateNamespace(funcMap),

    Deps: d,
    layoutHandler: output.NewLayoutHandler(),
  }

  if err := h.loadEmbedded(); err != nil {
    fmt.Println("error load embedded")
    return nil, err
  }

  if err != h.loadTemplates(); err != nil {
    fmt.Println("error load templates")
    fmt.Println("%#v", err)
    return nil, err
  }

  e := &templateExec{
    d: d,
    executor: exec,
    funcs: funcs,
    templateHandler: h,
  }

  d.SetTmpl(e)
  d.SetTextTmpl(newStandaloneTextTemplate(funcMap))

  return e, nil
}
