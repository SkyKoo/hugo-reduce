package helpers

import (
  "fmt"
  "io"
  "strings"
  "crypto/md5"
  "encoding/hex"
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

// MD5String takes a string and returns its MD5 hash.
func MD5String(f string) string {
  h := md5.New()
  h.Write([]byte(f))
  return hex.EncodeToString(h.Sum([]byte{}))
}

// Deprecated informs about a deprecation, but only once for a given set of arguments' values.
// If the err flag is enabled, it logs as an ERROR (will exit with -1) and the text will
// point at the next Hugo release.
// The idea is two remove an item in two Hugo releases to give users and theme authors
// plenty of time to fix their templates.
func Deprecated(item, alternative string, err bool) {
  if err {
    fmt.Printf("%s is deprecated and will be removed in Hugo %s. %s", item, "0.1", alternative)
  } else {
    fmt.Printf("%s is deprecated and will be removed in a future release. %s%s", item, alternative, "warnPanicMessage")
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
