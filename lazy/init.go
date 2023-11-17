package lazy

import "sync"

// Init holds a graph of lazily initialized dependencies.
type Init struct {
  // Used in tests
  initCount uint64

  mu sync.Mutex

  prev *Init
  children []*Init

  init onceMore
  out any
  err error
  f func() (any, error)
}

// New creates a new empty Init.
func New() *Init {
  return &Init{}
}

func (ini *Init) add(branch bool, initFn func() (any, error)) *Init {
  ini.mu.Lock()
  defer ini.mu.Unlock()

  if branch {
    return &Init{
      f: initFn,
      prev: ini,
    }
  }

  ini.checkDone()
  ini.children = append(ini.children, &Init{
    f: initFn,
  })

  return ini
}

func (ini *Init) checkDone() {
  if ini.init.Done() {
    panic("init cannot be added to after it has run")
  }
}

// Add adds a func as a new child dependency.
func (ini *Init) Add(initFn func() (any, error)) *Init {
  if ini == nil {
    ini = New()
  }
  return ini.add(false, initFn)
}
