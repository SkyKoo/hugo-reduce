> 说明, 本文中分析的源码是经过精简之后的，来自于 github.com/sunwei/hugo-playground，所以其中的数据结构内容是与源码不同的，但是结构上是差不多的，先了解这个精简的源码之后再看 hugo 真正的源码会简单一些。

# I. 配置加载
hugo 的配置加载分为几个步骤：
1. 创建对应的目录，对应的文件
2. 解析配置文件的内容
3. 加载配置文件 config.toml 和 theme, theme 作为 modules 对待

# II. 依赖 package
依赖的 package 主要有:
1. 虚拟文件系统，github.comf/spf13/afero， 封装了各种路径、文件的管理操作接口，管理和记录所有根文件相关的操作和数据
2. text 解析，golang.org/x/tools/txtar，用来将特定格式的文本内容解析，根据内容用 afero 创建对应的文件

# III. 主要的数据结构
想要了解 hugo 源码，必须要了解其中重要的数据结构，其体现了作者的设计思想，非常重要！

**configLoader、ConfigSourceDescriptor 和 Provider**
```go
type configLoader struct {
  cfg config.Provider
  ConfigSourceDescriptor
}

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
```

这三个数据结构呈现这样一种关系：
1. configLoader 中嵌套了 Provider，并组合了 Configourceescriptor，意味着 configLoader 可以直接使用 ConfigSourceDescriptor 的方法。
2. 其中 Provider 是一个 Interface，它是一系列方法的集合，以及配置内容，方便之后存入不同的对象实现；
3. ConfigSourceDescriptor 描述的是一个具体的配置文件，它有文件名，以及操作它的虚拟文件系统对象。
4. 所以， configLoader 表达的就是一个配置文件和操作这个配置文件的工具集合，和配置文件读取到内存中的内容。

# IV. 配置加载调用过程
**1. main 函数中调用 hugolib.LoadConfig 函数, 作为入口**
```go
// main.go
func main() {
  ...
  cfg, _, err := hugolib.LoadConfig(
    hugolib.ConfigSourceDescriptor{
      WorkingDir: tempDir,
      Fs:         afs,
      Filename:   "config.toml",
    },
  )
  ...
}
```
**2. hugolib.LoadConfig 中又调用 loadConfig 和 loadModulesConfig 分别处理 config 文件 和 mudolues 文件**
```go
// config.go
// 1. first, load config file
// 2. second, load modules
// 3. Important things, how to save config and modules that have been loaded.
func LoadConfig(d ConfigSourceDescriptor) (config.Provider, []string, error) {
  ...
  filename, err := l.loadConfig(d.Filename)
  ...
  if err := l.applyConfigDefaults(); err != nil {
    return l.cfg, configFiles, err
  }
  ...
  modulesConfig, err := l.loadModulesConfig()
  ...
  return l.cfg, configFiles, err
}
```
**3. loadConfig 读取 config 配置文件（hugo.toml) 并将其所有配置转换为 map[string]interface{} 格式的内容保存**
```go
func (l configLoader) loadConfig(configName string) (string, error) {
  ...
  m, err := config.FromFileToMap(l.Fs, filename)
  ...
  l.cfg.Set("", m)
  return filename, nil
}
```
**4. 调用 applyConfigDefaults 函数将默认的配置填充到输入的配置中，规则是如果输入配置中缺少了某一个配置项，则从默认中配置中将对应的默认配置填充进去，从函数中可以看到目前 hugo 的默认配置项**
```go
func (l configLoader) applyConfigDefaults() error {
  defaultSettings := maps.Params{
    "cleanDestinationDir": false,
    "watch": false,
    "resourceDir": "resources",
    "publishDir": "public",
    "themesDir": "themes",
    "buildDrafts": false,
    "buildFuture": false,
    "buildExpired": false,
    "environment": "production",
    "uglyURLs": false,
    "verbose": false,
    "ignoreCache": false,
    "canonifyURLs": false,
    "relativeURLs": false,
    "removePathAccents": false,
    "titleCaseStyle": "AP",
    "taxonomies": maps.Params{"tag": "tags", "category": "categories"},
    "premalinks": maps.Params{"a": "b"},
    "sitemap": maps.Params{"priority": -1, "filename": "sitemap.xml"},
    "disableLiveReload": false,
    "pluralizeListTitles": true,
    "footnoteAnchorPrefix": "",
    "footnoteReturnLinkContents": "",
    "newContentEditor": "",
    "paginate": 10,
    "paginatePath": "page",
    "summaryLength": 70,
    "rssLimit": -1,
    "sectionPagesMenu": "",
    "disablePathToLower": false,
    "hasCJKLanguage": false,
    "enableEmoji": false,
    "defaultContentLanguage": "en",
    "enableMissingTranslationPlaceholders": false,
    "enableGitInfo": false,
    "ignoreFiles": make([]string, 0),
    "disableAliases": false,
    "debug": false,
    "disableFastRender": false,
    "timeout": "30s",
    "enableInlineShortcodes": false,
  }

  l.cfg.SetDefaults(defaultSettings)

  return nil
}
```
**5. 调用 loadModulesConfig 函数加载 modules，主题 theme 是作为一个 module 处理的，所以这里也可以将处理主题的加载**

# V. 其他细节
### 递归处理 map[string]interface{} 类型的配置函数
利用 golang 的类型断言进行递归处理
```go
type Params map[string]any

// PrepareParams
// * makes all the keys in the given map lower cased and will do so
// * This will modify the map given.
// * Any nested map[interface{}]interface{}, map[string]interface{}, map[string]string
// * will be converted to Params.
// * Any _merge value will be converted to proper type and value.
func PrepareParams(m Params) {
  for k, v := range m {
    var retyped bool
    lKey := strings.ToLower(k)
    if lKey == mergeStrageyKey {
      v = toMergeStragegy(v)
      retyped = true
    } else {
      switch vv := v.(type) {
      case map[any]any:
        var p Params = cast.ToStringMap(v)
        v = p
        PrepareParams(p) // recurse
        retyped = true
      case map[string]any:
        var p Params = v.(map[string]any)
        v = p
        PrepareParams(p)
        retyped = true
      case map[string]string:
        p := make(Params)
        for k, v := range vv {
          p[k] = v
        }
        v = p
        PrepareParams(p)
        retyped = true
      }
    }

    if retyped || k != lKey {
      delete(m, k)
      m[lKey] = v
    }
  }
}
```
```go
// stringifyMapKeys recurse into in and changes all instances of
// map[interface{}]interface{} to map[string]interface{}. This is useful to
// work around the impedance mismatch between JSON and YAML unmarshaling that's
// described here: https://github.com/go-yaml/yaml/issues/139
//
// Inspired by https://github.com/stripe/stripe-mock, MIT licensed
func stringifyMapKeys(in any) (any, bool) {
  switch in := in.(type) {
  case []any:
    for i, v := range in {
      if vv, replaced := stringifyMapKeys(v); replaced {
        in[i] = vv
      }
    }
  case map[string]any:
    for k, v := range in {
      if vv, changed := stringifyMapKeys(v); changed {
        in[k] = vv
      }
    }
  case map[any]any:
    res := make(map[string]any)
    var (
      ok  bool
      err error
    )
    for k, v := range in {
      var ks string
      if ks, ok = k.(string); !ok {
        ks, err = cast.ToStringE(k)
        if err != nil {
          ks = fmt.Sprintf("%s", k)
        }
      }
      if vv, replaced := stringifyMapKeys(v); replaced {
        res[ks] = vv
      } else {
        res[ks] = v
      }
    }
    return res, true
  }

  return nil, false
}
```
### toml 和 yaml 的依赖 package
- 处理 toml 的 package 是 github.com/pelletier/go-toml/v2
- 处理 yaml 的 package 是 gopkg.in/yaml.v2
