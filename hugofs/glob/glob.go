package glob

import (
  "github.com/gobwas/glob"
)

type FilenameFilter struct {
  shouldInclude func(filename string) bool
  inclusions []glob.Glob
  dirInclusions []glob.Glob
  exinclusions []glob.Glob
  isWindows bool
}
