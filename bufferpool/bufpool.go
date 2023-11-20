package bufferpool

import (
  "bytes"
  "sync"
)

var bufferPool = &sync.Pool{
  New: func() any {
    return &bytes.Buffer{}
  },
}

// GetBuffer returns a buffer from the pool.
func GetBuffer() (buf *bytes.Buffer) {
  return bufferPool.Get().(*bytes.Buffer)
}

// PutBuffer returns a buffer to the pool.
// The buffer is reset before it it put back into circulation.
func PutBuffer(buf *bytes.Buffer) {
  buf.Reset()
  bufferPool.Put(buf)
}
