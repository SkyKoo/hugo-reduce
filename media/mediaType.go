package media

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

// DefaultTypes is the default media types supported by Hugo.
var DefaultTypes = Types{
  HTMLType,
  MarkdownType,
  TOMLType,
  TextType,
}

// DecodeTypes takes a list of media type configurations and merges those,
// in the order given, with the Hugo defaults as the last resort.
func DecodeTypes(mms ...map[string]any) (Types, error) {
  var m Types

  // Maps type string to Type. Type string is the full application/svg+xml.
  mmm := make(map[string]Type)
  for _, dt := range DefaultTypes {
  }
}
