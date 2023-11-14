package page

// import "strings"

const (
  KindPage = "page"

  // The rest are node types; home page, sections etc.

  KindHome = "home"
  KindSection = "section"

  // Note tha before Hugo 0.73 these were confusingly named
  // taxonomy (now: term)
  // taxonomyTerm (now: taxonomy)
  KindTaxonomy = "taxonomy"
  KindTerm = "term"
)
