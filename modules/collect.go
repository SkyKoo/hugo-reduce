package modules

type ModulesConfig struct {
  // All modules, including any disabled.
  AllModules Modules

  // All active Modules.
  ActiveModules Modules

  // Set if this is a Go modules enabled project.
  GoModulesFilename string
}
