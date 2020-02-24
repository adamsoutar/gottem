package main

import (
  "fmt"
)

func main () {
  var cfg = getConfig()

  fmt.Println("Gottem is linking " + cfg.fromDir + " to " + cfg.toDir + "...")

  link(cfg)

  fmt.Println("Done, enjoy your backup!")
}
