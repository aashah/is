package main

import (
    "errors"
    "fmt"
    "os"
    "path/filepath"
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
    for _, dir := range args {
        // check validity of argument
        fi, err := os.Stat(dir)
        if err != nil {
            fmt.Fprintf(os.Stderr, "is: could not find directory - %s\n", dir)
            continue
        }
        if !fi.IsDir() {
            fmt.Fprintf(os.Stderr, "is: found file, not directory - %s\n", dir)
            continue
        }

        abs, err := filepath.Abs(dir)
        if err != nil {

            continue
        }
        err = buildModule(abs, *checkB)
        if err != nil {
            fmt.Fprintf(os.Stderr, "is: %s\n", err.Error())
        }
    }
}

func buildModule(dir string, verbose bool) error {
    err := errors.New("unimplemented feature - build")
    return err
}