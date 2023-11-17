package helpers

import (
  "strings"

  "github.com/jdkato/prose/transform"
)

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
