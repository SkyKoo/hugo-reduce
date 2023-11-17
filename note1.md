> 说明, 本文中分析的源码是经过精简之后的，来自于 github.com/sunwei/hugo-playground，所以其中的数据结构内容是与源码不同的，但是结构上是差不多的，先了解这个精简的源码之后再看 hugo 真正的源码会简单一些。

# I. 构建 Site
hugo 中创建 Site 对象有以下几个步骤:
1. 构建依赖
2. 创建 Site，如果是多语言，则每个语言对应一个 Site

# II. 主要的数据结构

### 1. 描述站点的数据结构
**HugoSites**
```go
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

// Site contains all the information relevant for constucting a static
// site. The basic flow of information is as follows:
//
// 1. A list of Files is parsed and then converted into Pages.
//
// 2. Pages contain sections (based on the file they were generated from),
//    aliases and slugs (included in a pages frontmatter) which are they
//    various targets that will get generated. There will be canoical
//    listing. The canonical path can be overruled based on a pattern.
//
// 3. Taxonomies are created via configuration and will present some aspect of
//    the final page and typically a perm url.
//
// 4. All Pages are passed through a template based on their desired
//    layout based on numerous different elements.
//
// 5. The entire collection of files is written to disk.
type Site struct {
  language *langs.Language
  siteBucket *pagesMapBucket

  // Output formats defined in site config per Page Kind, or some defaults
  // if not set.
  // Output formats defined in Page front matter will override these.
  outputFormats map[string]output.Formats

  // All the output formats and media types available for this site.
  // These values will be merged from the Hugu defaults, the site config and,
  // finally, the language settings.
  outputFormatsConfig output.Formats
  mediaTypesConfig media.Types

  siteCfg siteConfigHolder

  // The func used to title case titles.
  titleFunc func(s string) string

  // newSite with above infos

  // The owning container. When multiple languages, there will be multiple
  // sites.
  h *HugoSites

  *PageCollections

  Sections Taxonomy
  Info *SiteInfo

  // The output formats that we need to render this site in. This slice
  // will be fixed once set.
  // This will be the union of Site.Pages' outputFormats.
  // This slice will be sorted.
  renderFormats output.Formats

  // Logger etc.
  *deps.Deps `json:"-"`

  siteRefLinker

  publisher publisher.Publisher

  // Shortcut to the home page. Note that this may be nil if
  // home page, for some odd reason, is disabled.
  home *pageState
}
```
- Site 结构是 Hugo 中最重要的结构，它描述了一个 Site 站点从构建到渲染出来需要用到的所有资源数据。
- HugoSites 结构可以看作是 Site 的容器，由于 Hugo 支持多语言，所以 Hugo 会对每一个 Language 创建一个 Site，这些 Site 就存在于 Hugoites 结构中。

### 2. 管理依赖的数据结构
hugo 中会将构建 Site 的所有资源事先准备好放在 依赖（Deps）对象中

**DepsCfg, Deps**
```go
// DepsCfg contains configuration options tha can be used to configure Hugo
// on a global level, i.e. logging etc.
// Nil values will be given default values.
type DepsCfg struct {
  // The language to use.
  Language *langs.Language

  // The file systems to use
  Fs *hugofs.Fs

  // The configuration to use.
  Cfg config.Provider

  // The media types configured.
  MediaTypes media.Types

  // The out formats configured.
  OutputFormats output.Formats

  // Template handling.
  TemplateProvider ResourceProvider
}

// Deps holds dependencies used by many.
// There will be normally only one instance of deps in play
// at a given time, i.e. one per Site built.
type Deps struct {
  // The PathSpec to use
  *helpers.PathSpec `json:"-"`

  // The templates to use. This will usually implement the full tpl.TemplateManager.
  tmpl tpl.TemplateHandler

  // We use this to parse and execute ad-hoc text templates.
  textTmpl tpl.TemplateParseFinder

  // All the output formats available for the current site.
  OutputFormatsConfig output.Formats

  // The Resource Spec to use
  ResourceSpec *resources.Spec

  // The SourceSpec to use
  SourceSpec *source.SourceSpec `json:"-"`

  // The ContentSpec to use
  *helpers.ContentSpec `json:"-"`

  // The site building.
  Site page.Site

  // The file systems to use.
  Fs *hugofs.Fs `json:"-"`

  templateProvider ResourceProvider

  // The configuration to use
  Cfg config.Provider `json:"-"`

  // The language in use. TODO(bep) consolidate with site
  Language *langs.Language

  // The translation func to use
  Translate func(translationID string, templateData any) string `json:"-"`
}
```
- Deps 管理的是构建 Site 需要的所有依赖资源，DepsCfg 则是这些资源的配置
