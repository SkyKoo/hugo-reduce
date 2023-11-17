package source

import (
  "github.com/SkyKoo/hugo-reduce/hugofs"
)

// File represents a source file.
// This is a temporary construct until we resolve page.Page conflicts.
// TODO(bep) remove this construct once we have resolved page deprecations
type File interface {
  fileOverlap
  FileWithoutOverlap
}

// Temporary to solve duplicate/deprecated names in page.Page
type fileOverlap interface {
  // Path gets the relative path including file name and extension.
  // The directory is relative to the content root.
  Path() string

  // Section is first directory below the content root.
  Section() string

  IsZero() bool
}

type FileWithoutOverlap interface {

  // Filename get the full path and filename to the file.
  Filename() string

  // Dir gets the name of the directory that contains this file.
  // The directory is relative to the content root.
  Dir() string

  // Extension is an alias to Ext().
  // Deprecated: Use Ext instead.
  Extension() string

  // Ext gets the file extension, i.e "myblogpost.md" will return "md".
  Ext() string

  // LogicalName is filename and extension of the file.
  LogicalName() string

  // BaseFileName is a filename without extension.
  BaseFileName() string

  // TranslationBaseName is a filename with no extension,
  // not even the optional language extension part.
  TranslationBaseName() string

  // ContentBaseName is a either TranslationBaseName or name of containing folder
  // if file is a leaf bundle.
  ContentBaseName() string

  // UniqueID is the MD5 hash of the file's path and is for most practical applications,
  // Hugo content files being one of them, considered to be unique.
  UniqueID() string

  FileInfo() hugofs.FileMetaInfo
}
