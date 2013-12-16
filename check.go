package main

import (
    "fmt"
)

func init() {
    cmdCheck.Run = runCheck
}

var cmdCheck = &Command {
    UsageLine: "Check -v [path to modules]",
    Short: "Checks a given list of modules",
    Long: `
Check will attempt to verify the integrity of a module by scanning for errors in
the module's manifest file.

TODO: add more checks
    `
}

func runCheck(cmd *Command, args []string) {
    
}