package main

import (
    "fmt"
)

func init() {
    cmdMv.Run = runMv
}

var cmdMv = &Command {
    UsageLine: "mv -v [path to modules]",
    Short: "mvs a given list of modules",
    Long: `
mv will attempt to verify the integrity of a module by scanning for errors in
the module's manifest file.

TODO: add more mvs
    `
}

func runMv(cmd *Command, args []string) {
    
}