package langs

import (
	"fmt"
	"strings"
	"time"

	"github.com/SkyKoo/hugo-reduce/config"
)

// These are the settings that should only be looked up in the global Viper
// config and not per language.
// This list may not be complete, but contains only settings that we know
// will be looked up in both.
// This isn't perfect, but it is ultimately the user who shoots him/herself in
// the foot.
// See the pathSpec.
var globalOnlySettings = map[string]bool{
  strings.ToLower("defaultContentLanguage"): true,
  strings.ToLower("multiligual"):            true,
  strings.ToLower("assetDir"):               true,
  strings.ToLower("resourceDir"):            true,
  strings.ToLower("build"):                  true,
}

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

// IsMultihost returns whether there are more than one language and at least one of
// the languages has baseURL specificed on the language level.
func (l Languages) IsMultihost() bool {
  if len(l) <= 1 {
    return false
  }

  for _, lang := range l {
    if lang.GetLocal("baseURL") != nil {
      return true
    }
  }
  return false
}

// GetLocal gets a configuration value set on language level. It Will
// not fall back to any global value.
// It will return nil if a value with the given key cannot be found.
// For internal use.
func (l *Language) GetLocal(key string) any {
  if l == nil {
    panic("language not set")
  }
  key = strings.ToLower(key)
  if !globalOnlySettings[key] {
    return l.LocalCfg.Get(key)
  }
  return nil
}
