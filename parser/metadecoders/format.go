package metadecoders

import (
  "path/filepath"
  "strings"
)


type Format string

const (
  // These are the supported metdata formats in Hugo. Most of these are also
  // supported as /data formats.
  ORG  Format = "org"
  JSON Format = "json"
  TOML Format = "toml"
  YAML Format = "yaml"
  CSV  Format = "csv"
  XML  Format = "xml"
)

// FormatFromString turns formatStr, typically a file extension without any ".",
// into a Format. It returns an empty string for unknown formats.
func FormatFromString(formatStr string) Format {
  formatStr = strings.ToLower(formatStr)
  if strings.Contains(formatStr, ".") {
    // recuse the filepath contains more than one "."
    formatStr = strings.TrimPrefix(filepath.Ext(formatStr), ".")
  }
  switch formatStr {
  case "toml":
    return TOML
  }

  return ""
}
