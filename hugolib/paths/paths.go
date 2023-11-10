package paths

import (
  "github.com/SkyKoo/hugo-reduce/config"
  "github.com/SkyKoo/hugo-reduce/hugofs"
  "github.com/SkyKoo/hugo-reduce/langs"
  "github.com/SkyKoo/hugo-reduce/modules"
)

type Paths struct {
  Fs *hugofs.Fs
  Cfg config.Provider

  BaseURL
  BaseURLString string

  // If the baseURL contains a base path, e.g. https://example.com/docs, then "/docs" will be the BasePath.
  BasePath string

  // Directories
  // TODO(bep) when we have trimmed down most of the dirs usage outside of this package, make
  // these into an interface.
  ThemesDir string
  WorkingDir string

  // Directories to store Resource related artifacts.
  AbsResourcesDir string

  AbsPublishDir string

  // pagination path handling
  PaginatePath string

  Language *langs.Language
  Languages langs.Languages
  LanguagesDefaultFirst langs.Languages

  AllModules modules.Modules
}
