package hugolib

import (
  "fmt"
  "sync"

  "github.com/SkyKoo/hugo-reduce/deps"
  "github.com/SkyKoo/hugo-reduce/output"
  "github.com/SkyKoo/hugo-reduce/common/para"
  "github.com/SkyKoo/hugo-reduce/lazy"
)

type hugoSitesInit struct {
  // Loads the data from all of the /data folders.
  data *lazy.Init

  // Performs late initialization (before render) of the templates.
  layouts *lazy.Init
}

// HugoSites represents the sites to build. Each site represents a language.
type HugoSites struct {
  Sites []*Site
  // Render output formats for all sites.
  renderFormats output.Formats

  // The currently rendered Site.
  currentSite *Site

  *deps.Deps

  contentInit sync.Once
  content *pageMaps

  init *hugoSitesInit

  workers *para.Workers
  numWorkers int

  // As loaded form the /data dirs
  data map[string]any
}

// NewHugoSites creates HugoSites from the given config.
func NewHugoSites(cfg deps.DepsCfg) (*HugoSites, error) {
  sites, err := createSitesFromConfig(cfg)
  if err != nil {
    return nil, fmt.Errorf("from config: %w", err)
  }

  return newHugoSItes(cfg, sites...)
}
