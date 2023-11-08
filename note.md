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

### 配置相关数据结构
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

### 模块相关数据结构
**Config、Mount 和 Inport**
```go
type Import struct {
  Path                string // Module path
  pathProjectReplaced bool // Set when Path is replaced in project config.
  IgnoreConfig        bool // Ignore any config in config.toml (will still follow imports).
  IgnoreImport        bool  // Do not follow any configured imports.
  NoMounts            bool  // Do not mount any folder in this import.
  NoVendor            bool  // Never vendor this import (only allowed in main project).
  Disable             bool  // Turn off this module.
  Mounts              []Mount
}

type Mount struct {
  Source string // relative path in source repo, e.g. "scss"
  Target string // relative target path, e.g. "assets/boostrap/scss"

  Lang string // any language code associated with this mount.
}

// Config holds a module config.
type Config struct {
  Mounts  []Mount
  Imports []Import

  // Meta info about this module (license information etc.).
  Params map[string]any
}
```
三者的关系：
1. Config 包含了 Mount 和 Import 的列表
2. Import 中包含了 Mount 列表，应该是用于 Module 中包含了其他 Moudle 的情况

**ClientConfig 和 Client**
```go
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
```
Client 是 hugo 用来加载 module 的类，之所以叫作 Client 它其实是一个 go module 的客户端

**collector、collected 和 goModules**
```go
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
  Path     string         // module Path
  Version  string         // module version
  Versions []string       // available module version (with -versions)
  Replace  *goModule      // replaced by this goModule
  Time     *time.Time     // time version was created
  Update   *goModule      // avaiable update, if any (with -u)
  Main     bool           // is this the main module?
  Indrect  bool           // is this module  only an indirect dependency of main module?
  Dir      string         // directory holding files for this module, if any
  GoMod    string         // path to go.mod file for this module, if any
  Error    *goModuleError // error loading module
}
```
收集器

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
  l.applyConfigDefaults()
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
**5. 调用 loadModulesConfig 函数读取 modules 的配置，主题 theme 是作为一个 module 处理的，所以这里也可以将主题的配置进行加载**
```go
  ...
  modulesConfig, err := l.loadModulesConfig()
  ...
}
```

**6. 调用 collectModules 函数收集 Modules, 同时定义了一个 hook 函数，作为回调**
```go
  // Need to run these after the modules are loaded, but before
  // they are finalized.
  collectHook := func(m *modules.ModulesConfig) error {
    mods := m.ActiveModules
    l.loadLanguageSettings(nil)

    // Apply default project mounts.
    // Default folder structure for hugo project
    modules.ApplyProjectConfigDefaults(l.cfg, mods[0])
  }

  _, modulesConfigFiles, modulesCollectErr := l.collectModules(modulesConfig, l.cfg, collectHook)
```
可以看到 hook 函数中做了两件事，针对每一个收集到的 Module 都需要加载 语言配置 和 Module 的默认配置

同时还要区分活跃和不活跃的模块，并不是所有加载的模块都有效，只有活跃的才有效

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
