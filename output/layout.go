package output

// LayoutDescriptor descirbes how a layout should be chosen. This is
// typically built from a Page.
type LayoutDescriptor struct {
  Type string
  Section string

  // E.g. "page", but also used for the _markup render kinds, e.g. "render-image".
  Kind string

  // Comma-separated list of kind variants, e.g. "go, json" as variants which would find "render-codeblock-go.html"
  KindVariants string

  Lang string
  Layout string
  // LayoutOverride indicates what we should only look for the above layout.
  LayoutOverride bool

  RenderingHook bool
  Baseof bool
}
