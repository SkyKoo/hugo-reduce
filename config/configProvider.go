package config

import (
  "github.com/SkyKoo/hugo-reduce/common/maps"
  "github.com/SkyKoo/hugo-reduce/types"
)

// Provider provides the configuration settings for Hugo.
type Provider interface {
  GetString(key string) string
  GetInt(key string) int
  GetBool(key string) bool
  GetParams(key string) maps.Params
  GetStringMap(key string) map[string]any
  GetStringMapString(key string) map[string]string
  GetStringSlice(key string) []string
  Get(key string) any
  Set(key string, value any)
  Merge(key string, value any)
  SetDefaults(params maps.Params)
  SetDefaultMergeStrategy()
  WalkParams(walkFn func(params ...KeyParams) bool)
  IsSet(key string) bool
}

func GetStringSlicePreserveString(cfg Provider, key string) []string {
  sd := cfg.Get(key)
  return types.ToStringSlicePreserveString(sd)
}
