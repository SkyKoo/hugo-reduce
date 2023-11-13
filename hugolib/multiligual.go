package hugolib

import (
	"github.com/SkyKoo/hugo-reduce/config"
	"github.com/SkyKoo/hugo-reduce/langs"
	"github.com/SkyKoo/hugo-reduce/log"
)

func getLanguages(cfg config.Provider) langs.Languages {
  log.Process("NewLanguages", "create multiple languages, only 'en' in our case")
  return langs.Languages{langs.NewDefaultLanguage(cfg)}
}
