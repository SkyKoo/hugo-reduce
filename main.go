package main

import (
  "bytes"
  "fmt"
  "os"
  "path/filepath"

  "github.com/spf13/afero"
  "golang.org/x/tools/txtar"

  "github.com/SkyKoo/hugo-reduce/hugofs"
  "github.com/SkyKoo/hugo-reduce/log"
  "github.com/SkyKoo/hugo-reduce/hugolib"
  "github.com/SkyKoo/hugo-reduce/deps"
)

func main() {
  // 0. local example contents

  log.Process("main", "prepare example project file systems")
  tempDir, clean, err := CreateTempDir(hugofs.Os, "go-hugo-temp-dir")
  if err != nil {
    clean()
    os.Exit(-1)
  }

  var afs afero.Fs
  afs = afero.NewOsFs()
  prepareFs(tempDir, afs)

  // 1. config, this is the most important part
  // 2. themes is a part of modules
  // 3. so, the first load config, and second load modules
  log.Process("main", "load configurations from config.toml and themes")
  cfg, _, err := hugolib.LoadConfig(
    hugolib.ConfigSourceDescriptor{
      WorkingDir: tempDir,
      Fs:         afs,
      Filename:   "config.toml",
    },
  )
  fmt.Printf("%#v\n", cfg)

  // 2. hugo file systems
  log.Process("main", "setup hugo file systems based on machine file system and configurations")
  fs := hugofs.NewFrom(afs, cfg, tempDir)

  // 3. dependencies management
  depsCfg := deps.DepsCfg{Cfg: cfg, Fs: fs}

  // 4. hugo istes
  log.Process("main", "create hugo sites based on deps")
  sites, err := hugolib.NewHugoSites(depsCfg)

  fmt.Println("HugoSites:")
  fmt.Printf("%#v\n", sites)
  fmt.Printf("%#v\n", sites.Sites[0])
  fmt.Println("===temp dir at last > ...")
  fmt.Println(tempDir)
}

func prepareFs(workingDir string, afs afero.Fs) {
  // attention, then content should not have any space at the beginning of line.
  files := `
-- config.toml --
theme = "mytheme"
contentDir = "mycontent"
-- myproject.txt --
Hello project!
-- themes/mytheme/mytheme.txt --
Hello theme!
-- mycontent/blog/post.md --
---
title: "Post Title"
---
### first blog
Hello Blog
-- layouts/index.html --
{{ $entries := (readDir ".") }}
START:|{{ range $entry := $entries }}{{ if not $entry.IsDir }}{{ $entry.Name }}|{{ end }}{{ end }}:END:
-- layouts/_default/single.html --
{{ .content }}
===
Static content
===

  `

  // 1. use txtar parse content, and create right dirs
  // 2. write content to the file
  data := txtar.Parse([]byte(files))
  for _, f := range data.Files { // deal with every file
    filename := filepath.Join(workingDir, f.Name) // whole path
    data := bytes.TrimSuffix(f.Data, []byte("\n"))
    // create dir with filename
    err := afs.MkdirAll(filepath.Dir(filename), 0777)
    if err != nil {
      fmt.Println(err)
    }

    err = afero.WriteFile(afs, filename, data, 0666)
    if err != nil {
      fmt.Println(err)
    }
  }
}

// CreateTempDir creates a temp dir in the given filesystem and
// return the dirname and a func that removes it when done.
func CreateTempDir(fs afero.Fs, prefix string) (string, func(), error) {
  tempDir, err := afero.TempDir(fs, "", prefix)
  if err != nil {
    return "", nil, err
  }

  return tempDir, func() { fs.RemoveAll(tempDir) }, nil
}
