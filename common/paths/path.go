package paths

import (
  "path/filepath"
)

// AbsPathify creates an absolute path if given a working dir and a relative path.
// If already absolute, the path is just cleaned.
func AbsPathify(workingDir, inPath string) string {
  if filepath.IsAbs(inPath) {
    return filepath.Clean(inPath)
  }
  return filepath.Join(workingDir, inPath)
}
