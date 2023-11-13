package hugolib

import (
	"time"

	"github.com/SkyKoo/hugo-reduce/deps"
	"github.com/SkyKoo/hugo-reduce/langs"
	"github.com/SkyKoo/hugo-reduce/log"
	"github.com/SkyKoo/hugo-reduce/media"
	"github.com/SkyKoo/hugo-reduce/output"
	"github.com/SkyKoo/hugo-reduce/publisher"
	"github.com/gohugoio/hugo/helpers"
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

  siteCfg siteConfigHolder

  // The func used to title case titles.
  titleFunc func(s string) string

  // newSite with above infos

  // The owning container. When multiple languages, there will be multiple
  // sites.
  h *HugoSites

  *PageCollections

  Sections Taxonomy
  Info *SiteInfo

  // The output formats that we need to render this site in. This slice
  // will be fixed once set.
  // This will be the union of Site.Pages' outputFormats.
  // This slice will be sorted.
  renderFormats output.Formats

  // Logger etc.
  *deps.Deps `json:"-"`

  siteRefLinker

  publisher publisher.Publisher

  // Shortcut to the home page. Note that this may be nil if
  // home page, for some odd reason, is disabled.
  home *pageState
}

type siteRefLinker struct {
  s *Site

  notFoundURL string
}

type SiteInfo struct {
  title string

  relativeURLs bool

  owner *HugoSites
  s *Site
}

// newSite creates a new site with the given configuration.
func newSite(cfg deps.DepsCfg) (*Site, error) {
  var (
    mediaTypesConfig []map[string]any
    outputFormatsConfig []map[string]any

    siteOutputFormatsConfig output.Formats
    siteMediaTypesCofnig media.Types
    err error
  )

  // [{toml}, {html}, {markdown}, {plain}]
  log.Process("media.DecodeTypes", "set default media types")
  siteMediaTypesCofnig, err = media.DecodeTypes(mediaTypesConfig...)
  if err != nil {
    return nil, err
  }

  // {{HTML}, {JSON}, {MARKDOWN}]
  log.Process("output.DecodeFormats", "set default output formats based on media types, and customized output formats configuration")
  siteOutputFormatsConfig, err = output.DecodeFormats(siteMediaTypesCofnig, outputFormatsConfig...)

  if err != nil {
    return nil, err
  }

  // Site output formats source
  log.Process("site output formats", "map siteOutputFormats to every hugo page types(KindPage, KindHome...)")
  outputFormats, err := createSiteOutputFormats(siteOutputFormatsConfig, nil, true)

  if err != nil {
    return nil, err
  }

  // KindTaxonomy, KindTerm like section title
  titleFunc := helpers.GetTitleFunc("")

  siteConfig := siteConfigHolder{
    timeout: 30 * time.Second, // page content output init timeout
  }

  var siteBucket *pagesMapBucket

  s := &Site{
    language: cfg.Language,
    siteBucket: siteBucket,

    outputFormats: outputFormats,
    outputFormatsConfig: siteOutputFormatsConfig,
    mediaTypesConfig: siteMediaTypesCofnig,

    siteCfg: siteConfig,
    titleFunc: titleFunc,
  }

  return s, nil
}

type siteConfigHolder struct {
  timeout time.Duration
}
