// This file parses the .gottemconfig for ignore rules and backup paths
package main

import (
  "github.com/mitchellh/go-homedir"
  "io/ioutil"
  "strings"
)

type config struct {
  fromDir string
  toDir string
  ignoreRules []string
}

func (c config) shouldIgnore (path string) bool {
  for _, iR := range c.ignoreRules {
    if strings.Contains(path, iR) {
      return true
    }
  }

  return false
}

func getConfig () config {
  var cfgPath, hdErr = homedir.Expand("~/.gottemconfig")
  if hdErr != nil {
    panic(hdErr)
  }

  var data, flErr = ioutil.ReadFile(cfgPath)
  if flErr != nil {
    panic(flErr)
  }

  return parseConfig(string(data))
}

func parseConfig (cfgStr string) config {
  var rules = strings.Split(cfgStr, "\n")
  var cfg = config {}

  var i = -1
  for _, rule := range rules {
    // Comments and blank lines
    if len(rule) == 0 || rule[0] == '#' {
      continue
    }
    i++

    if i == 0 {
      cfg.fromDir = rule
      continue
    }

    if i == 1 {
      cfg.toDir = rule
      continue
    }

    cfg.ignoreRules = append(cfg.ignoreRules, rule)
  }

  return cfg
}
