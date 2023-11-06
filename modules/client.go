package modules

import (
  "github.com/spf13/afero"
  "path/filepath"
)

const (
  goModFilename = "go.mod"
)

// ClientConfig configures the module Client.
type ClientConfig struct {
  // Fs to get the source
  Fs afero.Fs

  // If set, it will be run before wo do any duplicate checks for modules
  // etc.
  // It must be set in our case, for default project structure
  HookBeforeFinalize func(m *ModulesConfig) error

  // Absolute path to the project dir.
  WorkingDir string

  // Absolute path to the project's themes dir.
  ThemesDir string

  // Read from config file and transferred
  ModuleConfig Config
}

// Client contains most of the API provided by this package.
type Client struct {
  fs afero.Fs

  ccfg ClientConfig

  // The top level module config
  moduleConfig Config

  // Set when Go modules are initialized in the current repo, that is:
  // a go.mod file exists.
  GoModulesFilename string
}

// NewClient creates a new Client that can be used to manage the Hugo Components
// in a given WorkingDir.
// The Client will resolve the dependencies recursively, but needs the top
// level imports to start out.
func NewClient(cfg ClientConfig) *Client {
  fs := cfg.Fs
  n := filepath.Join(cfg.WorkingDir, goModFilename)
  goModEnable, _ := afero.Exists(fs, n)
  var goModFilename string
  if goModEnable {
    goModFilename = n
  }

  mcfg := cfg.ModuleConfig

  return &Client{
    fs: fs,
    ccfg: cfg,
    moduleConfig: mcfg,
    GoModulesFilename: goModFilename,
  }
}

func (h *Client) Collect() (ModulesConfig, error) {
  mc, coll := h.collect(true)
  if coll.err != nil {
    return mc, coll.err
  }

  if err := (&mc).setActiveMods(); err != nil {
    return mc, err
  }

  if h.ccfg.HookBeforeFinalize != nil {
    if err := h.ccfg.HookBeforeFinalize(&mc); err != nil {
      return mc, err
    }
  }

  if err := (&mc).finalize(); err != nil {
    return mc, err
  }

  return mc, nil
}

func (h *Client) collect(tidy bool) (ModulesConfig, *collector) {
  c := &collector{
    Client: h,
  }

  c.collect()
  if c.err != nil {
    return ModulesConfig{}, c
  }

  return ModulesConfig{
    AllModules: c.modules,
    GoModulesFilename: c.GoModulesFilename,
  }, c
}
