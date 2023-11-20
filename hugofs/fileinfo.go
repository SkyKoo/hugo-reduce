package hugofs

import (
	"os"
  "errors"
  "runtime"

  "golang.org/x/text/unicode/norm"
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

func (f *FileMeta) Open() (afero.File, error) {
  if f.OpenFunc == nil {
    return nil, errors.New("OpenFunc not set")
  }
  return f.OpenFunc()
}

type FileMetaInfo interface {
  os.FileInfo
  Meta() *FileMeta
}

type fileInfoMeta struct {
  os.FileInfo

  m *FileMeta
}

func (fi *fileInfoMeta) Meta() *FileMeta {
  return fi.m
}

func fileInfosToFileMetaInfos(fis []os.FileInfo) []FileMetaInfo {
  fims := make([]FileMetaInfo, len(fis))
  for i, v := range fis {
    fims[i] = v.(FileMetaInfo)
  }
  return fims
}

func normalizeFilename(filename string) string {
  if filename == "" {
    return ""
  }
  if runtime.GOOS == "darwin" {
    // When a file system is HFS+, its filepath is in NFD form.
    return norm.NFC.String(filename)
  }
  return filename
}
