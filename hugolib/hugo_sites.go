package hugolib

import (
	"fmt"

	"github.com/SkyKoo/hugo-reduce/deps"
)

// NewHugoSites creates HugoSites from the given config.
func NewHugoSites(cfg deps.DepsCfg) (*HugoSites, error) {
  sites, err := createSitesFromConfig(cfg)
  if err != nil {
    return nil, fmt.Errorf("from config: %w", err)
  }

  return newHugoSItes(cfg, sites...)
}
