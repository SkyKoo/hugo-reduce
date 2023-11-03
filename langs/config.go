package langs

import (
  "fmt"
  "sort"
  "strings"

  "github.com/SkyKoo/hugo-reduce/config"
)

type LanguagesConfig struct {
  Languages Languages // []*Language
}

func LoadLanguageSettings(cfg config.Provider, oldLangs Languages) (c LanguagesConfig, err error) {
  defaultLang := strings.ToLower(cfg.GetString("defaultContentLanguage"))
  if defaultLang == "" {
    defaultLang = "en"
    cfg.Set("defaultContentLanguage", defaultLang)
  }

  var languages map[string]any

  languagesFromConfig := cfg.GetParams("languages")
  disableLanguages := cfg.GetStringSlice("disableLanguages")

  if len(disableLanguages) == 0 {
    languages = languagesFromConfig
  } else {
    panic("there is no disabled language")
  }

  var languages2 Languages
  if len(languages) == 0 {
    languages2 = append(languages2, NewDefaultLanguage(cfg))
  } else {
    panic("no languages config params supported")
  }

  // oldLangs is nil

  // The defaultContentLanguage is something the user has to decide, but it needs
  // to match a language in the language definition list.
  langExists := false
  for _, lang := range languages2 {
    if lang.Lang == defaultLang {
      langExists = true
      break
    }
  }

  if !langExists {
    return c, fmt.Errorf("site config value %q for defaultContentLanguage does not match any language definition", defaultLang)
  }

  c.Languages = languages2

  sortedDefaultFirst := make(Languages, len(c.Languages))
  for i, v := range c.Languages {
    sortedDefaultFirst[i] = v
  }
  sort.Slice(sortedDefaultFirst, func(i, j int) bool {
    li, lj := sortedDefaultFirst[i], sortedDefaultFirst[j]
    if li.Lang == defaultLang {
      return true
    }

    if lj.Lang == defaultLang {
      return false
    }

    return i < j
  })

  cfg.Set("languagesSorted", c.Languages)  // ["en"]
  cfg.Set("languagesSortedDefaultFirst", sortedDefaultFirst)  // ["en"]
  cfg.Set("multilingual", len(languages2) > 1)  // false

  for _, language := range c.Languages {
    if language.initErr != nil {
      return c, language.initErr
    }
  }

  return c, nil
}

var languages map[string]any
