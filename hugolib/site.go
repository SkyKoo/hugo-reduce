package hugolib

import (
	"github.com/SkyKoo/hugo-reduce/langs"
	"github.com/SkyKoo/hugo-reduce/media"
	"github.com/SkyKoo/hugo-reduce/output"
)

// Site contains all the information relevant for constucting a static
// site. The basic flow of information is as follows:
//
// 1. A list of Files is parsed and then converted into Pages.
//
// 2. Pages contain sections (based on the file they were generated from),
//    aliases and slugs (included in a pages frontmatter) which are they
//    various targets that will get generated. There will be canoical
//    listing. The canonical path can be overruled based on a pattern.
//
// 3. Taxonomies are created via configuration and will present some aspect of
//    the final page and typically a perm url.
//
// 4. All Pages are passed through a template based on their desired
//    layout based on numerous different elements.
//
// 5. The entire collection of files is written to disk.
type Site struct {
  language *langs.Language
  siteBucket *pagesMapBucket

  // Output formats defined in site config per Page Kind, or some defaults
  // if not set.
  // Output formats defined in Page front matter will override these.
  outputFormats map[string]output.Formats

  // All the output formats and media types available for this site.
  // These values will be merged from the Hugu defaults, the site config and,
  // finally, the language settings.
  outputFormatsConfig output.Formats
  mediaTypesConfig media.Types
}
