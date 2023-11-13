package hugolib

import (
	"fmt"
	"sync"

	"github.com/SkyKoo/hugo-reduce/common/para"
	"github.com/SkyKoo/hugo-reduce/deps"
	"github.com/SkyKoo/hugo-reduce/log"
	"github.com/SkyKoo/hugo-reduce/output"
	// "github.com/SkyKoo/hugo-reduce/lazy"
)

type hugoSitesInit struct {
  // Loads the data from all of the /data folders.
  // data *lazy.Init

  // Performs late initialization (before render) of the templates.
  // layouts *lazy.Init
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

  return newHugoSites(cfg, sites...)
}

func createSitesFromConfig(cfg deps.DepsCfg) ([]*Site, error) {
  log.Process("createSitesFromConfig", "start")
  var sites []*Site

  // [en]
  languages := getLanguages(cfg.Cfg)
  for _, lang := range languages {
    var s *Site
    var err error
    cfg.Language = lang
    log.Process("newSite", "create site with DepsCfg with language setup")
    s, err = newSite(cfg)

    if err != nil {
      return nil, err
    }

    sites = append(sites, s)
  }

  log.Process("createSitesFromConfig", "end")
  return sites, nil
}

// NewHugoSites creates a new collection of sites given the input sites, building
// a language configuration based on those.
func newHugoSites(cfg deps.DepsCfg, sites ...*Site) (*HugoSites, error) {
  // Return error at the end. Make the caller decide if it's fatal or not.
  var initErr error

  return nil, initErr
}

