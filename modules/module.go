package modules

type Module interface {

  // Config The decoded module config and mounts.
  Config() Config

  // Dir Directory holding files for this module.
  Dir() string

  // IsGoMod Returns whether this is a Go Module.
  IsGoMod() bool

  // Path Returns the path to this module.
  // This will either be the module path, e.g. "github.com/gohugoio/myshortcodes",
  // or the path below your /theme folder, e.g. "mytheme".
  Path() string

  // Owner In the dependency three, this is the first module that defines this module
  // as a dependency
  Owner() Module

  // Mounts Any directory remappings.
  Mounts() []Mount
}

type Modules []Module
