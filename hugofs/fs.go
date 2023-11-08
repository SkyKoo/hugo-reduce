package hugofs

import (
	"os"

	"github.com/SkyKoo/hugo-reduce/common/paths"
	"github.com/SkyKoo/hugo-reduce/config"
	"github.com/SkyKoo/hugo-reduce/log"
	"github.com/spf13/afero"
)

// Os points to the (read) Os filesystem
var Os = &afero.OsFs{}

// Fs holds the core filesystems used by Hugo.
type Fs struct {
  // Source is Hugo's source file system.
  // Note that this will always be a "plain" Afero filesystem:
  // * afero.OsFs when running in production
  // * afero.MemMapFs for many of the tests.
  Source afero.Fs

  // PublishDir is where Hugo publishes its rendered content.
  // It's mounted inside publishDir (default /public).
  PublishDir afero.Fs

  // WorkingDirReadOnly is a read-only file system
  // restricted to the project working dir.
  WorkingDirReadOnly afero.Fs
}

// NewFrom creates a new Fs based  on the provided Afero Fs
// as source and destination file systems.
// Useful for testing.
func NewFrom(fs afero.Fs, cfg config.Provider, wd string) *Fs {
  return newFs(fs, fs, cfg, wd)
}

func newFs(source, destination afero.Fs, cfg config.Provider, wd string) *Fs {
  cfg.Set("workingDir", wd)
  workingDir := cfg.GetString("workingDir")
  publishDir := cfg.GetString("publishDir")

  absPublishDir := paths.AbsPathify(workingDir, publishDir)

  // Make sure we always hava the /public folder ready to use.
  if err := source.MkdirAll(absPublishDir, 0777); err != nil && !os.IsExist(err) {
    panic(err)
  }
  log.Process("newFs", "create /public folder")

  log.Process("newFs", "new base path fs &BasePathFs{}")
  pubFs := afero.NewBasePathFs(destination, absPublishDir)

  return &Fs{
    Source:             source,
    PublishDir:         pubFs,
    WorkingDirReadOnly: getWorkingDirFsReadOnly(source, workingDir),
  }
}

func getWorkingDirFsReadOnly(base afero.Fs, workingDir string) afero.Fs {
  if workingDir == "" {
    return afero.NewReadOnlyFs(base)
  }
  return afero.NewBasePathFs(afero.NewReadOnlyFs(base), workingDir)
}
