package tplimpl

import (
	"reflect"
	"text/template"

	"github.com/SkyKoo/hugo-reduce/deps"
	"github.com/SkyKoo/hugo-reduce/log"
	"github.com/SkyKoo/hugo-reduce/tpl/internal/go_templates/texttemplate"
)

func newTemplateExecuter(d *deps.Deps) (texttemplate.Executer, map[string]reflect.Value) {
  funcs := createFuncMap(d)
  funcsv := make(map[string]reflect.Value)

  for k, v := range funcs {
    vv := reflect.ValueOf(v)
    funcsv[k] = vv
  }

  log.Process("GoFuncs", "map template.GoFuncs to funcMap")
  // Duplicate Go's internal funcs here for faster lookups.
  for k, v := range template.GoFuncs {
    if _, exists := funcsv[k]; !exists {
      vv, ok := v.(reflect.Value)
      if !ok {
        vv = reflect.ValueOf(v)
      }
      funcsv[k] = vv
    }
  }

  log.Process("GoFuncs", "map texttemplate.GoFuncs to funcMap")
  for k, v := range texttemplate.GoFuncs {
    if _, exists := funcsv[k]; !exists {
      funcsv[k] = v
    }
  }

  exeHelper := &templateExecHelper{
    funcs: funcsv,
  }

  return texttemplate.NewExecuter(
    exeHelper,
  ), funcsv
}

func createFuncMap(d *deps.Deps) map[string]any {
  funcMap := template.FuncMap{}

  // Merge the namespace funcs
  for _, nsf := range internal.TemplateFuncsNamespaceRegistry {
    ns := nsf(d)
    if _, exists := funcMap[ns.Name]; exists {
      panic(ns.Name + " is a duplicate template func")
    }
    funcMap[ns.Name] = ns.Context
    for _, mm := range ns.MethodMappings {
      for _, alias := range mm.Aliases {
        if _, exists := funcMap[alias]; exists {
          panic(alias + " is a duplicate template func")
        }
        funcMap[alias] = mm.Method
      }
    }
  }

  return funcMap
}
