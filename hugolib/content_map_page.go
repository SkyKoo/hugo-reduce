package hugolib

import (
  "sync"

  "github.com/SkyKoo/hugo-reduce/common/para"
  "github.com/SkyKoo/hugo-reduce/resources/page"
)

type pageMap struct {
  s *Site
  *contentMap
}

type pageMaps struct {
  workers *para.Workers
  pmaps []*pageMap
}

type pagesMapBucket struct {
  owner *pageState // The branch node

  *pagesMapBucketPages
}

type pagesMapBucketPages struct {
  pagesInit sync.Once
  pages page.Pages

  pagesAndSectionsInit sync.Once
  pagesAndSections page.Pages

  sectionsInit sync.Once
  sections page.Pages
}
