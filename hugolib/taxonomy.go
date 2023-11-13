package hugolib

import "github.com/SkyKoo/hugo-reduce/resources/page"

// A Taxonomy is a map of keywords to a list of pages.
// For example
//    TagTaxonomy['technology'] = page.WeightedPages
//    TagTaxonomy['go'] = page.WeightedPages
type Taxonomy map[string]page.WeightedPages
