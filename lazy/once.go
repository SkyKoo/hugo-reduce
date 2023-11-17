package lazy

import (
  "sync"
  "sync/atomic"
)

// onceMore is similar to sync.Once.
//
// Additional features are:
// * it can be reset, so the action can be repeated if needed
// * it has methods to check if it's done or in progress
//
type onceMore struct {
  mu sync.Mutex
  lock uint32
  done uint32
}

func (t *onceMore) Done() bool {
  return atomic.LoadUint32(&t.done) == 1
}
