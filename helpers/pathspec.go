package helpers

import (
  "github.com/SkyKoo/hugo-reduce/config"
  "github.com/SkyKoo/hugo-reduce/hugofs"
  "github.com/SkyKoo/hugo-reduce/hugolib/filesystems"
  "github.com/SkyKoo/hugo-reduce/hugolib/paths"
)

// PathSpec holds methods that decides how paths in URLs and files in Hugo should look like.
type PathSpec struct {
  *paths.Paths
  *filesystems.BaseFs

  // The file systems to use
  Fs *hugofs.Fs

  // The config provider to use
  Cfg config.Provider
}
