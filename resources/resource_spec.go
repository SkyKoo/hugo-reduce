package resources

import (
  "github.com/SkyKoo/hugo-reduce/helpers"
  "github.com/SkyKoo/hugo-reduce/media"
  "github.com/SkyKoo/hugo-reduce/output"
)

type Spec struct {
  *helpers.PathSpec

  MediaTypes media.Types
  OutputFormats output.Formats
}
