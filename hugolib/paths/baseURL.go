package paths

import (
  "net/url"
)

// A BaseURL in Hugo is normally on the form scheme://path, but the
// form scheme: is also valid (mailto:hugo@rules.com).
type BaseURL struct {
  url *url.URL
  urlStr string
}
