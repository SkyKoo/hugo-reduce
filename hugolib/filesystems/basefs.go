package filesystems

import (
	"github.com/SkyKoo/hugo-reduce/hugofs"
	"github.com/spf13/afero"
)

// BaseFs contains the core base filesystems used by Hugo. The name "base" is used
// to underline that even if they can be composites, they all have a base path set to a specific
// resource folder, e.g "/my-project/content". So, no absolute filenames needed.
type BaseFs struct {

  // SourceFilesystems contains the different source file systems.
  *SourceFilesystems

  // The project source.
  SourceFs afero.Fs

  // The filesystem used to publish the rendered site.
  // This usually maps to /my-project/public.
  PublishFs afero.Fs

  // A read-only filesystem starting from the project workDir.
  WorkDir afero.Fs

  // theBigFs *filesystemsCollector
}

// SourceFilesystems contains the different source file systems. These can be
// composite file systems (theme and project etc.), and they have all root
// set to the source type the provides: data, i18n, static, layouts.
type SourceFilesystems struct {
  Content    *SourceFilesystem
  Data       *SourceFilesystem
  I18n       *SourceFilesystem
  Layouts    *SourceFilesystem
  Archetypes *SourceFilesystem
  Assets     *SourceFilesystem

  // Writable filesystem on top the project's resources directory,
  // with any sub module's resource fs layered below.
  ResourcesCache afero.Fs

  // The work folder (may be a composite of project and theme components).
  Work afero.Fs

  // When in multihost we have one static filesystem per language. The sync
  // static files is currently done outside of the Hugo build (where there is
  // a concept of a site per language).
  // When in non-multihost mode there will be one entry in this map with a blank key.
  Static map[string]*SourceFilesystem

  // All the /static dirs (including themes/modules).
  StaticDirs []hugofs.FileMetaInfo
}

// A SourceFilesystem holds the filesystem for a given source type in Hugo (data,
// i18n, layouts, static) and additional metadata to be able to use that filesystem
// in server mode.
type SourceFilesystem struct {
  // Name matches one in files.ComponentFolders
  Name string

  // This is a virtual composite filesystem. It expects path relative to a context.
  Fs afero.Fs

  // This filesystem as separate root directories, starting from project and down
  // to the themes/modules.
  Dirs []hugofs.FileMetaInfo

  // When syncing a source folder to the target (e.g /public), this may
  // be set to publish into a subfolder. This is used for static syncing
  // in multihost mode.
  PublishFolder string
}
