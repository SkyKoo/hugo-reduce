package deps

import (
	"github.com/SkyKoo/hugo-reduce/config"
	"github.com/SkyKoo/hugo-reduce/hugofs"
	"github.com/SkyKoo/hugo-reduce/langs"
	"github.com/SkyKoo/hugo-reduce/hugo/media"
	"github.com/SkyKoo/hugo-reduce/hugo/output"
)

// Deps holds dependencies used by many.
// There will be normally only one instance of deps in play
// at a given time, i.e. one per Site built.
type Deps struct {
  // The PathSpec to use
  // *helpers.PathSpec `json:"-"`
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
