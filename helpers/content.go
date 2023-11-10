package helpers

import (
  "github.com/SkyKoo/hugo-reduce/config"
  // "github.com/SkyKoo/hugo-reduce/markup"
)

// ContentSpec provides functionality to render markdown content.
type ContentSpec struct {
  // Converters markup.ConverterProvider

  // SummaryLength it the length of the summary that Hugo extracts from a content.
  summaryLength int

  Cfg config.Provider
}
