package parse

// import "github.com/gobwas/glob/syntax/lexer"

// Tree is the representation of a single parsed template.
type Tree struct {
  Name string  // name of the template represented by the tree.
  ParseName string // name of the top-level template during parsing, for error messages.
  Root *ListNode // top-level root of the tree.
  Mode Mode // parsing mode.
  text string // text parsed to create the template (or its parent)
  // Parsing only; cleared after parse.
  funcs []map[string]any
  lex *lexer
  token [3]item // three-token lookahead for parser.
  peekCount int
  vars []string // variables defined at the monent.
  treeSet map[string]*Tree
  actionLine int // line of left delim starting actionLine
  rangeDepth int
}

// A mode value is a set of flags (or 0). MOdes control parser behavior.
type Mode uint
