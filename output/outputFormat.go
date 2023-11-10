package output

import (
	"github.com/SkyKoo/hugo-reduce/media"
)

// Format represents an output representation, usually to a file on disk.
type Format struct {
  // The Name is used as an identifier. Internal output formats (i.e. HTML and RSS)
  // can be overridden by providing a new definition for those types.
  Name string `json:"name"`

  MediaType media.Type `json:"-"`

  // Must be set to a value when there are two or more conflicting mediatype for the same resource.
  Path string `json:"path"`

  // The base output file name used when not using "ugly URLs", defaults to "index".
  BaseName string `json:"baseName"`

  // The value to use for rel links
  //
  // See https://www.w3schools.com/tags/att_link_rel.asp
  //
  // Amp has a special requirement in the department, see:
  // https://www.ampproject.org/docs/guides/deploy/discovery
  // I.e.:
  // <link rel="amphtml" href="https://www.example.com/url/to/amp/document.html">
  Rel string `json:"rel"`

  // The protocol to use, i.e. "webcal://", Defaults to the protocol of the baseURL.
  Protocol string `json:"protocol"`

  // IsPlainText decides whether to use text/template or html/template
  // as template parser.
  IsPlainText bool `json:"isPlainText"`

  // IsHTML returns whether this format is int the HTML family. This includes
  // HTML, AMP etc. This is used to decide when to create alias redirects etc.
  IsHTML bool `json:"isHTML"`

  // Enable if it doesn't make sence to include this format in an alternative
  // format listing, CSS being one good expample.
  // Note that we use the term "alternative" and not "alternate" here, as it
  // does not necessarily replace the other format, it is an alternative representation.
  NotAlternative bool `json:"notAlternative"`

  // Setting this will make this output format control the value of
  // .Premalink and .RelPremalink for a rendered Page.
  // If not set, these values will point to the main (first) output format
  // configured. That is probably the behaviour you want in most situations,
  // as you probably don't want to link back to the RSS version of a page, as an
  // example. AMP would, however, be a good example of an output format where this
  // behaviour is wanted.
  Premalinkable bool `json:"premalinkable"`

  // Setting this to a non-zero value will be used as the first sort criteria.
  Weight int `json:"weight"`
}

// Formats is a slice of Format.
type Formats []Format
