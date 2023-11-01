package hugolib

import (
  // "fmt"

  "github.com/spf13/afero"

  "github.com/SkyKoo/hugo-reduce/config"
  "github.com/SkyKoo/hugo-reduce/log"
  // "github.com/SkyKoo/hugo-reduce/modules"
)

// ConfigSourceDescriptor describes where to find the config (e.g. config.toml etc.).
type ConfigSourceDescriptor struct {
  Fs afero.Fs

  // Path to the config file to use, e.g. /my/procect/config.toml
  // Multiple config files supported, e.g. 'config.toml, abc.toml'
  // In our case, one config file is just right Filename string
  Filename string

  // The project's working dir. Is used to look for addtitional theme config.
  WorkingDir string
}

func (d ConfigSourceDescriptor) configFileDir() string {
  return d.WorkingDir
}

type configLoader struct {
  cfg config.Provider
  ConfigSourceDescriptor
}

func (l configLoader) loadConfig(configName string) (string, error) {
  return "", nil
}

func (l configLoader) applyConfigDefaults() error {
  return nil
}

/*
func (l configLoader) loadModulesConfig() (modules.Config, error) {
  return nil, nil
}
*/

func LoadConfig(d ConfigSourceDescriptor) (config.Provider, []string, error) {
  var configFiles []string
  l := configLoader{ConfigSourceDescriptor: d, cfg: config.New()}
  log.Process("LoadConfig", "start init configLoader")
  filename, err := l.loadConfig(d.Filename)
  if err == nil {
    configFiles = append(configFiles, filename)
  } else {
    return nil, nil, err
  }

  log.Process("LoadConfig", "apply config defaults")
  if err := l.applyConfigDefaults(); err != nil {
    return l.cfg, configFiles, err
  }

  /*
  log.Process("LoadConfig", "load modules config")
  modulesConfig, err := l.loadModulesConfig()
  fmt.Printf("%#v\n", modulesConfig)
  if err != nil {
    fmt.Println("load modules config err...")
    return l.cfg, configFiles, err
  }
  */

  return l.cfg, configFiles, err
}
