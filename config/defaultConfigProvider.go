package config

import (
  "sync"

  "github.com/SkyKoo/hugo-reduce/common/maps"
)

type KeyParams struct {
  Key    string
  Params maps.Params
}

// New creates a Provider backed by an empty maps.Params
func New() Provider {
  return &defaultConfigProvider{
    root: make(maps.Params), // actually maps.Params is a map that all keys are lower case
  }
}

type defaultConfigProvider struct {
  mu    sync.RWMutex
  root  maps.Params

  keyCache sync.Map
}

func (c *defaultConfigProvider) SetDefaultMergeStrategy() {
}


func (c *defaultConfigProvider) Get(k string) any {
  return ""
}

func (c *defaultConfigProvider) GetBool(k string) bool {
  return false
}

func (c *defaultConfigProvider) GetInt(k string) int {
  return 0
}

func (c *defaultConfigProvider) IsSet(k string) bool {
  return false
}

func (c *defaultConfigProvider) GetString(k string) string {
  return ""
}

func (c *defaultConfigProvider) GetParams(k string) maps.Params {
  return nil
}

func (c *defaultConfigProvider) GetStringMap(k string) map[string]any {
  return nil
}

func (c *defaultConfigProvider) GetStringMapString(k string) map[string]string {
  return nil
}

func (c *defaultConfigProvider) GetStringSlice(k string) []string {
  return nil
}

func (c *defaultConfigProvider) Set(k string, v any) {
}

func (c *defaultConfigProvider) SetDefaults(params maps.Params) {
}

func (c *defaultConfigProvider) Merge(k string, v any) {
}

func (c *defaultConfigProvider) WalkParams(walkFn func(params ...KeyParams) bool) {
}

func (c *defaultConfigProvider) getNestedKeyAndMap(key string, create bool) (string, maps.Params) {
  return "", nil
}
