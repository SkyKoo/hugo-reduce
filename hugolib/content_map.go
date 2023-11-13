package hugolib

import (
  "github.com/armon/go-radix"
)

type contentTree struct {
  Name string
  *radix.Tree
}

type contentTrees []*contentTree

type contentMap struct {
  // View of regular pages, sections, and taxonomies.
  pageTrees contentTrees

  // View of pages, sections, taxonomies, and resources.
  bundleTrees contentTrees

  // Stores page bundles keyed by its path's directory or the base filename,
  // e.g. "blog/post.md" => /blog/post", "blog/post/index.md" => "/blog/post"
  // These are the "regular pages" and all of them are bundles.
  pages *contentTree

  // Section nodes.
  sections *contentTree

  // Resources stored per bundle below a common prefix, e.g. "/blog/post__hb_".
  resources *contentTree
}
