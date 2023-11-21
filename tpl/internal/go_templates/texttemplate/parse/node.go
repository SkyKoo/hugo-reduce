package parse

import (
  "strings"
)

// A Node is an element in the parse tree. The interface is trivial.
// The interface contains an unexported method so that only
// types local to this package can satisfy it.
type Node interface {
  Type() NodeType
  String() string
  // Copy does a deep copy of the Node and all its components.
  // To avoid type assertions, some XxxNodes also have specialized
  // CopyXxx methods that return *XxxNode
  Copy() Node
  Position() Pos // byte position of start of node in full original input string
  // tree returns the containing *Tree.
  // It is unexported so all implementations of Node are in this package.
  tree() *Tree
  // writeTo writes the String output to the builder.
  writeTo(*strings.Builder)
}

// NodeType identifies the type of a parse tree node.
type NodeType int

// Pos represents a byte position in the original input text from which
// this template was parsed.
type Pos int

func (p Pos) Position() Pos {
  return p
}

// ListNode holds a sequence of nodes.
type ListNode struct {
  NodeType
  Pos
  tr *Tree
  Nodes []Node // The element nodes in lexical order
}
