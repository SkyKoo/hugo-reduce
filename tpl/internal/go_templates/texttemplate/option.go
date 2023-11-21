package texttemplate

// missingKeyAction defines how to respond to indexing a map with a key that is not present.
type missingKeyAction int

type option struct {
  missingKey missingKeyAction
}
