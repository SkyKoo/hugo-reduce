package source

import (
  "github.com/SkyKoo/hugo-reduce/helpers"
  "github.com/spf13/afero"
)

// SourceSpec abstracts language-specific file creation.
// TODO(bep) rename to Spec
type SourceSpec struct {
  *helpers.PathSpec

  SourceFs afero.Fs

  Language map[string]any
  DefaultContentLanguage string
}
