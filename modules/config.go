package modules

import (
	"strings"

	"github.com/SkyKoo/hugo-reduce/config"
	"github.com/SkyKoo/hugo-reduce/log"
)

type Import struct {
  Path string  // Module path.
  pathProjectReplaced bool  // Set when Path is replaced in project config.
  IgnoreConfig bool  // Ignore any config in config.toml (will still follow imports).
  IgnoreImport bool  // Do not follow any configured imports.
  NoMounts bool  // Do not mount any folder in this import.
  NoVendor bool  // Never vendor this import (only allowed in main project).
  Disable bool  // Turn off this module.
  Mounts []Mount
}

type Mount struct {
  Source string // relative path in source repo. e.g. "scss"
  Target string // relative Target path, e.g. "assets/boostrap/scss"

  Lang string // any language code associated with this mount.
}

// Used as key to remove duplicates.
func (m Mount) key() string {
  return strings.Join([]string{"", m.Source, m.Target}, "/")
}

// Config holds a module config.
type Config struct {
  Mounts []Mount
  Imports []Import

  // Meta info about this module (license infomation etc.).
  Params map[string]any
}

// DecodeConfig creates a modules Config from a given Hugo configuration.
func DecodeConfig(cfg config.Provider) (Config, error) {
  return decodeConfig(cfg)
}

func decodeConfig(cfg config.Provider) (Config, error) {
  c := DefaultModuleConfig

  if cfg == nil {
    // applyThemeConfig
    return c, nil
  }

  themeSet := cfg.IsSet("theme")

  if themeSet {
    // [mytheme]
    log.Process("decodeConfig", "set mytheme as Imports in DefaultModuleConfig, config{}")
    imports := config.GetStringSlicePreserveString(cfg, "theme")
    for _, imp := range imports {
      c.Imports = append(c.Imports, Import{
        Path: imp,
      })
    }
  }
  return c, nil
}

var DefaultModuleConfig = Config{}
