package tplimpl

import (
  "github.com/SkyKoo/hugo-reduce/deps"
  "github.com/SkyKoo/hugo-reduce/log"
)

// TemplateProvider manages templates.
type TemplateProvider struct{}

// DefaultTemplateProvider is a globally available TemplateProvider.
var DefaultTemplateProvider *TemplateProvider

// Update updates the Hugo Template System in the provided Deps
// with all the addittional features, templates & functions.
func (*TemplateProvider) Update(d *deps.Deps) error {
  log.Process("templateProvider Update", "new TemplateExec")
  tmpl, err := newTemplateExec(d)
  if err != nil {
    return err
  }
  return tmpl.postTransform()
}

func (*TemplateProvider) Clone(d *deps.Deps) error {
  panic("not implemented")
  return nil
}
