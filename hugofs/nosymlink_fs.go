package hugofs

import (
  "errors"
)

var ErrPermissionSymlink = errors.New("symlinks not allowed in this filesystem")
