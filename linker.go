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
  "strconv"
  "github.com/mitchellh/go-homedir"
)

type errors struct {
  sync.Mutex
  errs []error
}

func (e *errors) append (err error) {
  e.Lock()
  e.errs = append(e.errs, err)
  e.Unlock()
}

func link (cfg config, maxConcurrency int) {
  // Create the dest folder if not exists
  var err = os.MkdirAll(cfg.toDir, 0775)
  if err != nil {
    panic(err)
  }

  // Delete the contents of the destination directory
  fmt.Println("Emptying the destination directory...")
  emptyDir(cfg.toDir)

  var wg sync.WaitGroup
  var sem = make(chan int, maxConcurrency)
  // nonFatalErrors
  var nFE errors

  fmt.Println("Walking the source dir and linking files...")
  wg.Add(1)
  go processDir(cfg.fromDir, "", sem, &nFE, &wg, cfg)

  wg.Wait()

  logErrors(&nFE)
}

// path, relative path, waitGroup, config
func processDir (pth string, relPth string, sem chan int, nFE *errors, wg *sync.WaitGroup, cfg config) {
  // Limit the concurrency
  sem <- 0
  defer wg.Done()
  defer func () { <-sem }()

  var flInfo, err = ioutil.ReadDir(pth)
  if err != nil {
    nFE.append(err)
    return
  }

  // Create the dir in the remote
  var mErr = os.MkdirAll(path.Join(cfg.toDir, relPth), 0775)
  if mErr != nil {
    nFE.append(mErr)
    return
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
      go processDir(newAbs, newRel, sem, nFE, wg, cfg)
    } else {
      go processFile(newAbs, newRel, sem, nFE, wg, cfg)
    }
  }
}

func processFile (pth string, relPth string, sem chan int, nFE *errors, wg *sync.WaitGroup, cfg config) {
  // Limit the concurrency
  sem <- 0
  defer wg.Done()
  defer func () { <-sem }()

  var newPth = path.Join(cfg.toDir, relPth)
  var err = os.Link(pth, newPth)
  if err != nil {
    nFE.append(err)
  }
}

func logErrors (nFE *errors) {
  var logStr = ""

  for _, e := range nFE.errs {
    logStr += e.Error() + "\n"
  }

  fmt.Println(strconv.Itoa(len(nFE.errs)) + " non-fatal errors occured during your backup. Check ~/gottem.log")

  // I find this pattern tiresome - I'm probably doing it stupidly.
  // Can someone suggest a better way?
  var pth, err = homedir.Expand("~/gottem.log")
  if err != nil {
    panic(err)
  }
  var log, lErr = os.Create(pth)
  if lErr != nil {
    panic(lErr)
  }
  var _, wErr = log.WriteString(logStr)
  if wErr != nil {
    panic(wErr)
  }
  var cErr = log.Close()
  if cErr != nil {
    panic(cErr)
  }
}

func emptyDir (pth string) {
  // Empty a directory without deleting it.
  // G B&S freaks out if you actually delete it.
  var dI, err = ioutil.ReadDir(pth)
  if err != nil {
    panic(err)
  }

  for _, f := range dI {
    var thisPth = path.Join(pth, f.Name())
    var fErr error

    if f.IsDir() {
      fErr = os.RemoveAll(thisPth)
    } else {
      fErr = os.Remove(thisPth)
    }

    if fErr != nil {
      panic(fErr)
    }
  }
}
