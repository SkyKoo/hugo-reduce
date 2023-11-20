package hugolib

import (
	"fmt"
	"strings"
	"sync"

	"github.com/SkyKoo/hugo-reduce/common/para"
	"github.com/SkyKoo/hugo-reduce/config"
	"github.com/SkyKoo/hugo-reduce/deps"
	"github.com/SkyKoo/hugo-reduce/helpers"
	"github.com/SkyKoo/hugo-reduce/hugofs"
	"github.com/SkyKoo/hugo-reduce/lazy"
	"github.com/SkyKoo/hugo-reduce/log"
	"github.com/SkyKoo/hugo-reduce/output"
	"github.com/SkyKoo/hugo-reduce/parser/metadecoders"
	"github.com/SkyKoo/hugo-reduce/source"
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

  return newHugoSites(cfg, sites...)
}

func createSitesFromConfig(cfg deps.DepsCfg) ([]*Site, error) {
  log.Process("createSitesFromConfig", "start")
  var sites []*Site

  // [en]
  // every language has own site
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

  // 3
  log.Process("newHugoSites", "get number of worker")
  numWorkers := config.GetNumWorkerMultiplier()
  if numWorkers > len(sites) {
    numWorkers = len(sites)
  }

  var workers *para.Workers

  log.Process("newHugoSites", "init HugoSites")
  h := &HugoSites{
    Sites: sites,
    workers: workers,
    numWorkers: numWorkers,
    init: &hugoSitesInit{
      data: lazy.New(),
      layouts: lazy.New(),
    },
  }

  log.Process("newHugoSites", "add data to h.init")
  h.init.data.Add(func() (any, error) {
    log.Process("newHugoSites", "h.init run h.loadData")
    err := h.loadData(h.PathSpec.BaseFs.Data.Dirs)
    if err != nil {
      return nil, fmt.Errorf("failed to load data: %w", err)
    }
    return nil, nil
  })

  return nil, initErr
}

func (h *HugoSites) loadData(fis []hugofs.FileMetaInfo) (err error) {
  spec := source.NewSourceSpec(h.PathSpec, nil, nil)

  h.data = make(map[string]any)
  for _, fi := range fis {
    fileSystem := spec.NewFilesystemFromFileMetaInfo(fi)
    files, err := fileSystem.Files()
    if err != nil {
      return err
    }
    for _, r := range files {
      if err := h.handleDataFile(r); err != nil {
        return err
      }
    }
  }

  return
}

func (h *HugoSites) handleDataFile(r source.File) error {
  var current map[string]any

  f, err := r.FileInfo().Meta().Open()
  if err != nil {
    return fmt.Errorf("data: filed to open %q: %w", r.LogicalName(), err)
  }
  defer f.Close()

  // Crawl in data tree to insert data
  current = h.data
  keyParts := strings.Split(r.Dir(), helpers.FilePathSeparator)

  for _, key := range keyParts {
    if key != "" {
      if _, ok := current[key]; !ok {
        current[key] = make(map[string]any)
      }
      current = current[key].(map[string]any)
    }
  }

  data, err := h.readData(r)
  if err != nil {
    return err
  }

  if data == nil {
    return nil
  }

  // filepath.Walk walks the files in lexical order, '/' comes before '.'
  higherPrecedentData := current[r.BaseFileName()]

  switch data.(type) {
  case nil:
  case map[string]any:

    switch higherPrecedentData.(type) {
    case nil:
      current[r.BaseFileName()] = data
    case map[string]any:
      // merge maps: insert entries from data for keys that
      // don't already exist in higherPrecedentData
      higherPrecedentMap := higherPrecedentData.(map[string]any)
      for key, value := range data.(map[string]any) {
        if _, exists := higherPrecedentMap[key]; exists {
          // this warning could happen if
          // 1. A theme used the same key; the main data folder wins
          // 2. A sub folder uses the same key; the sub folder wins
          // TODO(bep) figure out a way to detect 2) above and make that a WARN
          fmt.Printf("Data for key '%s' in path '%s' is overridden by higher precedence data already in the data tree", key, r.Path())
        } else {
          higherPrecedentMap[key] = value
        }
      }
    default:
      // can't merge: higherPrecedentData is not a map
      fmt.Printf("The %T data from '%s' overridden by "+
      "higher precedence %T data already in the data tree", data, r.Path(), higherPrecedentData)
    }

  case []any:
    if higherPrecedentData == nil {
      current[r.BaseFileName()] = data
    } else {
      // we don't merge array data
      fmt.Printf("The %T data from '%s' overridden by "+
      "higher precedence %T data already in the data tree", data, r.Path(), higherPrecedentData)
    }

  default:
    fmt.Printf("unexpected data type %T in file %s", data, r.LogicalName())
  }

  return nil
}

func (h *HugoSites) readData(f source.File) (any, error) {
  file, err := f.FileInfo().Meta().Open()
  if err != nil {
    return nil, fmt.Errorf("readData: failed to open data file: %w", err)
  }
  defer file.Close()
  content := helpers.ReaderToBytes(file)

  format := metadecoders.FormatFromString(f.Ext())
  return metadecoders.Default.Unmarshal(content, format)
}
