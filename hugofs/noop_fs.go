package hugofs

import (
  "errors"
  "os"
  "time"

  "github.com/spf13/afero"
)

var (
  errNoOp = errors.New("this is a filesystem that does nothing and this operation is not supported")
  _ afero.Fs = (*noOpFs)(nil)

  // NoOpFs provides a no-op filesystem that implements the afero.Fs
  // interface.
  NoOpFs = &noOpFs{}
)

type noOpFs struct {
}

func (fs noOpFs) Create(name string) (afero.File, error) {
  return nil, errNoOp
}

func (fs noOpFs) Mkdir(name string, perm os.FileMode) error {
  return nil
}

func (fs noOpFs) MkdirAll(name string, perm os.FileMode) error {
  return nil
}

func (fs noOpFs) Open(name string) (afero.File, error) {
  return nil, os.ErrNotExist
}

func (fs noOpFs) OpenFile(name string, flag int, perm os.FileMode) (afero.File, error) {
  return nil, os.ErrNotExist
}

func (fs noOpFs) Remove(name string) error {
  return nil
}

func (fs noOpFs) RemoveAll(path string) error {
  return nil
}

func (fs noOpFs) Rename(oldname string, newname string) error {
  return errNoOp
}

func (fs noOpFs) Stat(name string) (os.FileInfo, error) {
  return nil, os.ErrNotExist
}

func (fs noOpFs) Name() string {
  return "noOpFs"
}

func (fs noOpFs) Chmod(name string, mode os.FileMode) error {
  return errNoOp
}

func (fs noOpFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
  return errNoOp
}

func (fs noOpFs) Chown(name string, uid int, gid int) error {
  return errNoOp
}
