package hugofs

import (
  "github.com/spf13/afero"
)

// Os points to the (read) Os filesystem
var Os = &afero.OsFs{}
