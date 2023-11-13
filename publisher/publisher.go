package publisher

import (
  "io"
  "net/url"

  "github.com/SkyKoo/hugo-reduce/output"
)

// Publisher publishes a result file.
type Publisher interface {
  Publish(d Descriptor) error
}

// Descriptor describes the needed publishing chain for an item.
type Descriptor struct {
  // The content to publish.
  Src io.Reader

  // The OutputFormat of the this content.
  OutputFormat output.Format

  // Where to publish this content. This is a filesystem-relative path.
  TargetPath string

  // Configuration that trigger pre-processing.
  // LiveReload script will be injected if this is != nil
  LiveReloadBaseURL *url.URL

  // Enable to inject the Hugo generated tag in the header. Is currently only
  // injected on the home page for HTML type of output formats.
  AddHugoGeneratorTag bool

  // If set, will replace all relative URLs with this one.
  AbsURLPath string

  // Enable to minify the output using the OutputFormat defined above to
  // pick the correct minifier configuration.
  Minify bool
}
