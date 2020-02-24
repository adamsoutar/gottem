// This part syncs up the hard links in the dest dir with those in the source dir
// It uses a pool of goroutines to process directories and create links quickly

// "Why do you need a pool of workers for counting files!?", you cry
// Because the dir I want to sync has over half a million (518,479) files,
// and counting them - let alone linking them - took 13s.
package main

import (
  "sync"
  "io/ioutil"
  "path"
  "os"
  "fmt"
)

const maxConcurrency = 16

func link (cfg config) {
  // Delete the existing destination directory
  var err = os.RemoveAll(cfg.toDir)
  if err != nil {
    panic(err)
  }
  var mErr = os.Mkdir(cfg.toDir, 0775)
  if mErr != nil {
    panic(mErr)
  }

  var wg sync.WaitGroup
  var sem = make(chan int, maxConcurrency)

  wg.Add(1)
  go processDir(cfg.fromDir, "", sem, &wg, cfg)

  wg.Wait()
}

// path, relative path, waitGroup, config
func processDir (pth string, relPth string, sem chan int, wg *sync.WaitGroup, cfg config) {
  // Limit the concurrency
  sem <- 0
  defer wg.Done()
  defer func () { <-sem }()

  var flInfo, err = ioutil.ReadDir(pth)
  if err != nil {
    panic(err)
  }

  // Create the dir in the remote
  var mErr = os.MkdirAll(path.Join(cfg.toDir, relPth), 0775)
  if mErr != nil {
    panic(mErr)
  }

  // Walk the dir contents
  for _, f := range flInfo {
    var fPath = f.Name()
    var isDir = f.IsDir()

    var newAbs = path.Join(pth, fPath)

    if cfg.shouldIgnore(newAbs) {
      continue
    }

    wg.Add(1)

    var newRel = path.Join(relPth, fPath)
    if isDir {
      go processDir(newAbs, newRel, sem, wg, cfg)
    } else {
      go processFile(newAbs, newRel, sem, wg, cfg)
    }
  }
}

func processFile (pth string, relPth string, sem chan int, wg *sync.WaitGroup, cfg config) {
  // Limit the concurrency
  sem <- 0
  defer wg.Done()
  defer func () { <-sem }()

  var newPth = path.Join(cfg.toDir, relPth)
  var err = os.Link(pth, newPth)
  if err != nil {
    fmt.Println("Link error: " + err.Error())
  }
}
