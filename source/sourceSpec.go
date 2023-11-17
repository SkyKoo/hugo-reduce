package source

import (
	"github.com/spf13/afero"
	"github.com/SkyKoo/hugo-reduce/helpers"
	"github.com/SkyKoo/hugo-reduce/hugofs/glob"
)

// SourceSpec abstracts language-specific file creation.
// TODO(bep) rename to Spec
type SourceSpec struct {
  *helpers.PathSpec

  SourceFs afero.Fs

  Languages map[string]any
  DefaultContentLanguage string
}

// NewSourceSpec initializes SourceSpec using languages the given filesystem and PathSpec.
func NewSourceSpec(ps *helpers.PathSpec, inclusionFilter *glob.FilenameFilter, fs afero.Fs) *SourceSpec {
  cfg := ps.Cfg
  defaultLang := cfg.GetString("defaultContentLanguage")
  languages := cfg.GetStringMap("languages")

  return &SourceSpec{
    PathSpec: ps,
    SourceFs: fs,
    Languages: languages,
    DefaultContentLanguage: defaultLang,
  }
}
