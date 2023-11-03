package langs

import (
	"fmt"
	"time"

	"github.com/SkyKoo/hugo-reduce/config"
)

// Language manages specific-language configuration.
type Language struct {
  Lang string
  Weight int // for sort

  // If set per language, this tells Hugo that all content files without any
  // language indicator (e.g. my-page.en.md) is in this language.
  // This is usually a path relative to the workding dir, but it can be an
  // absolute directory reference. It is what we get.
  // For internal use.
  ContentDir string

  // Global config.
  // For internal use.
  Cfg config.Provider

  // Language specific config.
  // For internal use.
  LocalCfg config.Provider

  // Composite config.
  // For internal use.
  config.Provider

  location *time.Location

  // Error during initialization. Will fail the buld.
  initErr error
}

// For internal use.
func (l *Language) String() string {
  return l.Lang
}

// NewLanguage creates a new language.
func NewLanguage(lang string, cfg config.Provider) *Language {
  localCfg := config.New()
  compositeConfig := config.NewCompositeConfig(cfg, localCfg)

  l := &Language{
    Lang: lang,
    ContentDir: cfg.GetString("contentDir"),
    Cfg: cfg,
    LocalCfg: localCfg,
    Provider: compositeConfig,
  }

  if err := l.loadLocation(cfg.GetString("timeZone")); err != nil {
    l.initErr = err
  }

  return l
}

// NewDefaultLanguage creates the default language for config.Provider.
// If not otherwise specified the default is "en".
func NewDefaultLanguage(cfg config.Provider) *Language {
  defaultLang := cfg.GetString("defaultContentLanguage")

  if defaultLang == "" {
    defaultLang = "en"
  }

  return NewLanguage(defaultLang, cfg)
}

type Languages []*Language

func (l *Language) loadLocation(tzStr string) error {
  location, err := time.LoadLocation(tzStr)
  if err != nil {
    return fmt.Errorf("invalid timeZone for language %q: %w", l.Lang, err)
  }
  l.location = location

  return nil
}
