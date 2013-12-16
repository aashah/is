package main

import (
    "fmt"
)

func init() {
    cmdBuild.Run = runBuild
}

var cmdBuild = &Command {
    UsageLine: "build [-v] [path to modules]",
    Short: "builds a given list of modules",
    Long: `
Build will attempt to a compile a list of modules given that they are structured
correctly. If a module has not yet passed the integrity check, then build will
call the check routine to ensure the module is structured correctly.

See "go help check" to understand more on what is meant by the integrity of a
module.
    `,
}

func runBuild(cmd *Command, args []string) {
    fmt.Println("Running build")
}

func buildModule(dir string, verbose bool) error {
    return nil
}