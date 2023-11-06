package modules

import (
  // "github.com/spf13/afero"
  "time"
)

type ModulesConfig struct {
  // All modules, including any disabled.
  AllModules Modules

  // All active Modules.
  ActiveModules Modules

  // Set if this is a Go modules enabled project.
  GoModulesFilename string
}

// Collects and creates a module tree.
type collector struct {
  *Client

  // Store away any non-fatal error and return at the end.
  err error

  // Set to disable any Tidy operation in the end.
  skipTidy bool

  *collected
}

type collected struct {
  // Pick the first and prevent circular loops.
  seen map[string]bool

  // Set if a Go modules enabled project.
  gomods goModules

  // Ordered list if collected modules, including Go Modules and theme
  // components stored below /themes.
  modules Modules
}

type goModule struct {
  Path string // module Path
  Version string // module version
  Versions []string // available module version (with -versions)
  Replace *goModule // replaced by this goModule
  Time *time.Time // time version was created
  Update *goModule // avaiable update, if any (with -u)
  Main bool // is this the main module?
  Indrect bool // is this module  only an indirect dependency of main module?
  Dir string // directory holding files for this module, if any
  GoMod string // path to go.mod file for this module, if any
  Error *goModuleError // error loading module
}

type goModuleError struct {
  Err string // the error itself
}

type goModules []*goModule
