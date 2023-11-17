package hugofs

import (
	"os"

	"github.com/spf13/afero"

	"github.com/SkyKoo/hugo-reduce/hugofs/glob"
	"github.com/SkyKoo/hugo-reduce/hugofs/files"
)

type FileMeta struct {
  Name string
  Filename string
  Path string
  PathWalk string
  OriginalFilename string
  BaseDir string

  SourceRoot string
  MountRoot string
  Module string

  Weight int
  IsOrdered bool
  IsSymlink bool
  IsRootFile bool
  IsProject bool
  Watch bool

  Classifier files.ContentClass

  SkipDir bool

  Lang string
  TranslationBaseName string
  TranslationBaseNameWithExt string
  Translations []string

  Fs afero.Fs
  OpenFunc func() (afero.File, error)
  JoinStateFunc func(name string) (FileMetaInfo, error)

  // Include only files or directories that match.
  InclusionFilter *glob.FilenameFilter
}

type FileMetaInfo interface {
  os.FileInfo
  Meta() *FileMeta
}
