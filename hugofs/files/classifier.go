package files

import "path/filepath"

const (
  ComponentFolderArchetypes = "archetypes"
  ComponentFolderStatic = "static"
  ComponentFolderLayouts = "layouts"
  ComponentFolderContent = "content"
  ComponentFolderData = "data"
  ComponentFolderAssets = "assets"
  ComponentFolderI18n = "i18n"

  FolderResources = "resources"
  FolderJSConfig = "_jsconfig" // Mounted below /assets with postcss.config.js etc.
)

var (
  JsConfigFolderMountPrefix = filepath.Join(ComponentFolderAssets, FolderJSConfig)

  ComponentFolders = []string{
    ComponentFolderArchetypes,
    ComponentFolderStatic,
    ComponentFolderLayouts,
    ComponentFolderContent,
    ComponentFolderData,
    ComponentFolderAssets,
    ComponentFolderI18n,
  }

  ComponentFoldersSet = make(map[string]bool)
)

type ContentClass string

const (
  ContentClassLeaf ContentClass = "leaf"
  ContentClassBranch ContentClass = "branch"
  ContentClassFile ContentClass = "zfile" // Sort below
  ContentClassContent ContentClass = "zcontent"
)

func (c ContentClass) IsBundle() bool {
  return c == ContentClassLeaf || c == ContentClassBranch
}
