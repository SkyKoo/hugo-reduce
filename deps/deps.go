package deps

import (
  "github.com/SkyKoo/hugo-reduce/config"
  "github.com/SkyKoo/hugo-reduce/hugofs"
  "github.com/SkyKoo/hugo-reduce/langs"
  "github.com/SkyKoo/hugo-reduce/media"
  "github.com/SkyKoo/hugo-reduce/output"
  "github.com/SkyKoo/hugo-reduce/helpers"
  "github.com/SkyKoo/hugo-reduce/tpl"
  "github.com/SkyKoo/hugo-reduce/resources"
  "github.com/SkyKoo/hugo-reduce/resources/page"
  "github.com/SkyKoo/hugo-reduce/source"
)

// Deps holds dependencies used by many.
// There will be normally only one instance of deps in play
// at a given time, i.e. one per Site built.
type Deps struct {
  // The PathSpec to use
  *helpers.PathSpec `json:"-"`

  // The templates to use. This will usually implement the full tpl.TemplateManager.
  tmpl tpl.TemplateHandler

  // We use this to parse and execute ad-hoc text templates.
  textTmpl tpl.TemplateParseFinder

  // All the output formats available for the current site.
  OutputFormatsConfig output.Formats

  // The Resource Spec to use
  ResourceSpec *resources.Spec

  // The SourceSpec to use
  SourceSpec *source.SourceSpec `json:"-"`

  // The ContentSpec to use
  *helpers.ContentSpec `json:"-"`

  // The site building.
  Site page.Site

  // The file systems to use.
  Fs *hugofs.Fs `json:"-"`

  templateProvider ResourceProvider

  // The configuration to use
  Cfg config.Provider `json:"-"`

  // The language in use. TODO(bep) consolidate with site
  Language *langs.Language

  // The translation func to use
  Translate func(translationID string, templateData any) string `json:"-"`
}

// DepsCfg contains configuration options tha can be used to configure Hugo
// on a global level, i.e. logging etc.
// Nil values will be given default values.
type DepsCfg struct {
  // The language to use.
  Language *langs.Language

  // The file systems to use
  Fs *hugofs.Fs

  // The configuration to use.
  Cfg config.Provider

  // The media types configured.
  MediaTypes media.Types

  // The out formats configured.
  OutputFormats output.Formats

  // Template handling.
  TemplateProvider ResourceProvider
}

// ResourceProvider is used to create and refresh, and clone resources needed.
type ResourceProvider interface {
  Update(deps *Deps) error
  Clone(deps *Deps) error
}
