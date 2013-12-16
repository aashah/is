package main

import (
    "fmt"
)

var moduleIntegrityCache map[string]bool

func init() {
    cmdCheck.Run = runCheck
    moduleIntegrityCache = make(map[string]bool)
}

var cmdCheck = &Command {
    UsageLine: "Check -v [path to modules]",
    Short: "Checks a given list of modules",
    Long: `
Check will attempt to verify the integrity of a module by scanning for errors in
the module's manifest file.

TODO: add more checks
    `,
}

func runCheck(cmd *Command, args []string) {
    fmt.Println("Running check")
}

func checkModuleIntegrity(dir string) error {
    return nil
}