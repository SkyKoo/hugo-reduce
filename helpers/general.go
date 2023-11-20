package helpers

import (
  "io"
  "strings"
  "path/filepath"

  "github.com/jdkato/prose/transform"

  bp "github.com/SkyKoo/hugo-reduce/bufferpool"
)

// FilePathSeparator as defined by os.Separator.
const FilePathSeparator = string(filepath.Separator)

// GetTitleFunc returns a func that can be used to transform a string to
// title case.
//
// The supported styles are
//
// - "Go" (strings.Title)
// - "AP" (see http://www.apstylebook.com/)
// - "Chicago" (see http://www.chicagomanualofstyle.org.home.html)
//
// If an unknown or empty style is provided, AP style is what you get.
func GetTitleFunc(style string) func(s string) string {
  switch strings.ToLower(style) {
  case "go":
    return strings.Title
  default:
    tc := transform.NewTitleConverter(transform.APStyle)
    return tc.Title
  }
}

// ReaderToBytes takes an io.Reader argument, reads from it
// and returns bytes.
func ReaderToBytes(lines io.Reader) []byte {
  if lines == nil {
    return []byte{}
  }
  b := bp.GetBuffer()
  defer bp.PutBuffer(b)

  b.ReadFrom(lines)

  bc := make([]byte, b.Len())
  copy(bc, b.Bytes())
  return bc
}
