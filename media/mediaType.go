package media

import (
  "strings"
  "sort"
)

const (
  defaultDelimiter = "."
)

// Type (also known as MIME type and content type) is a two-part identifier for
// file formats and format contents transmitted on the Internet.
// For Hugo's use case, we use the top-level type name / subtype name + suffix.
// One example would be application/svg+xml
// If suffix is not provided, the sub type will be used.
// See // https//en.wikipedia.org/wiki/Media_type
type Type struct {
  MainType  string `json:"mainType"`  // i.e. text
  SubType   string `json:"subType"`   // i.e. html
  Delimiter string `json:"delimiter"` // i.e. "."

  // FirstSuffix holds the first suffix defined for this Type.
  FirstSuffix SuffixInfo `json:"firstSuffix"`

  // This is the optional suffix after the "+" in the MIME type,
  // e.g. "xml" in "application/rss+xml".
  mimeSuffix string

  // E.g. "jpg,jpeg"
  // Stored as a string to make Type comparable.
  suffixesCSV string
}

// SuffixInfo holds information about a Type's suffix.
type SuffixInfo struct {
  Suffix     string `json:"suffix"`
  FullSuffix string `json:"fullSuffix"`
}

// Types is a slice of media types.
type Types []Type

func (t Types) Len() int { return len(t) }
func (t Types) Swap(i, j int) { t[i], t[j] = t[j], t[i] }
func (t Types) Less(i, j int) bool { return t[i].Type() < t[j].Type() }

func (m *Type) init() {
  m.FirstSuffix.FullSuffix = ""
  m.FirstSuffix.Suffix = ""
  if suffixes := m.Suffixes(); suffixes != nil {
    m.FirstSuffix.Suffix = suffixes[0]
    m.FirstSuffix.FullSuffix = m.Delimiter + m.FirstSuffix.Suffix
  }
}

// Suffixes returns all valid file suffixes for this type.
func (m Type) Suffixes() []string {
  if m.suffixesCSV == "" {
    return nil
  }

  return strings.Split(m.suffixesCSV, ",")
}

// Type returns a string representing the main- and sub-type of a media type, e.g. "text/css".
// A suffix identifier will be appended after a "+" if set, e.g. "image/svg+xml".
// Hugo will register a set of default media types.
// These can be overridden by the user in the configurations,
// by defining a media type with the same Type.
func (m Type) Type() string {
  // Examples are
  // image/svg+xml
  // text/css
  if m.mimeSuffix != "" {
    return m.MainType + "/" + m.SubType + "+" + m.mimeSuffix
  }
  return m.MainType + "/" + m.SubType
}

// Definitions from https://developer.mozilla.org/en-US/docs/Web/HTTP/Basics_of_HTTP/MIME_types etc.
// Note that from Hugo 0.44 we only set Suffix if it is part of the MIME type.
var (
  HTMLType = newMediaType("text", "html", []string{"html"})

  JSONType = newMediaType("application", "json", []string{"json"})
  TOMLType = newMediaType("application", "toml", []string{"toml"})

  MarkdownType = newMediaType("text", "markdown", []string{"md", "markdown"})

  OctetType = newMediaType("application", "octet-stream", nil)

  TextType = newMediaType("text", "plain", []string{"txt"})
)

// DefaultTypes is the default media types supported by Hugo.
var DefaultTypes = Types{
  HTMLType,
  MarkdownType,
  TOMLType,
  TextType,
}

func newMediaType(main, sub string, suffixes []string) Type {
  t := Type{MainType: main, SubType: sub, suffixesCSV: strings.Join(suffixes, ","), Delimiter: defaultDelimiter}
  t.init()
  return t
}

// DecodeTypes takes a list of media type configurations and merges those,
// in the order given, with the Hugo defaults as the last resort.
func DecodeTypes(mms ...map[string]any) (Types, error) {
  var m Types

  // Maps type string to Type. Type string is the full application/svg+xml.
  mmm := make(map[string]Type)
  for _, dt := range DefaultTypes {
    mmm[dt.Type()] = dt
  }

  // no customized media type

  for _, v := range mmm {
    m = append(m, v)
  }
  sort.Sort(m)

  return m, nil
}
