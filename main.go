package main

import (
  "fmt"
  "strconv"
  "os"
)

func main () {
  var maxConcurrency = 16
  var args = os.Args[1:]

  if len(args) > 0 {
    if i, err := strconv.Atoi(args[0]); err == nil {
      maxConcurrency = i
    }
  }

  var cfg = getConfig()

  fmt.Println("Gottem is linking " + cfg.fromDir + " to " + cfg.toDir + "...")

  link(cfg, maxConcurrency)

  fmt.Println("Done, enjoy your backup!")
}
