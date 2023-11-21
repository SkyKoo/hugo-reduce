package parse

// item represents a token or text string returned from the scanner.
type item struct {
  typ itemType // The type of this item.
  pos Pos // The starting position, in bytes, of this item in the input string.
  val string // The value of this item.
  line int // The line number at the start of this item.
}

// itemType identifies the type of lex items.
type itemType int

// lexer holds the state of the scanner.
type lexer struct {
  name string // the name of the input; used only for error reports
  input string // the string being scanned
  leftDelim string // start of action
  rightDelim string // end of action
  emitComment bool // emit itemComment tokens.
  pos Pos // current position in the input
  start Pos // start position of this item
  width Pos // width of last rune read from input
  items chan item // channel of scanned items
  parenDepth int // nesting depth of ( ) exprs
  line int // 1+number of newlines seen
  startLine int // start line of this item
  breakOK bool // break keyword allowed
  continueOK bool // continue keyword allowed
}
