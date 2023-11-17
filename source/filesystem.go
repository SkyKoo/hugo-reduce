package source

import (
	"fmt"
	"sync"

	"github.com/SkyKoo/hugo-reduce/hugofs"
)

// Filesystem represents a source filesystem.
type Filesystem struct {
  files []File
  filesInit sync.Once
  filesInitErr error

  Base string

  fi hugofs.FileMetaInfo

  SourceSpec
}

func (sp SourceSpec) NewFilesystemFromFileMetaInfo(fi hugofs.FileMetaInfo) *Filesystem {
  return &Filesystem{SourceSpec: sp, fi: fi}
}

// Files returns a slice of readable files.
func (f *Filesystem) Files() ([]File, error) {
  f.filesInit.Do(func() {
    err := f.captureFiles()
    if err != nil {
      f.filesInitErr = fmt.Errorf("capture files: %w", err)
    }
  })
  return f.files, f.filesInitErr
}

func (f *Filesystem) captureFiles() error {
  walker := func(path string, fi hugofs.FileMetaInfo, err error) error {
    if err != nil {
      return err
    }

    if fi.IsDir() {
      return nil
    }

    meta := fi.Meta()
    filename := meta.Filename

    b, err := f.shouldRead(filename, fi)
    if err != nil {
      return err
    }

    if b {
      err = f.add(filename, fi)
    }

    return err
  }

  w := hugofs.NewWalkway(hugofs.WalkwayConfig{
    Fs: f.SourceFs,
    Info: f.fi,
    Root: f.Base,
    WalkFn: walker,
  })

  return w.Walk()
}
