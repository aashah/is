package main

import (
    "errors"
    "fmt"
    "os"
    "path/filepath"
)

var moduleIntegrityCache map[string]bool

func init() {
    cmdCheck.Run = runCheck
    moduleIntegrityCache = make(map[string]bool)
}

var cmdCheck = &Command {
    UsageLine: "chk [-v] [path to modules]",
    Short: "Checks a given list of modules",
    Long: `
Check will attempt to verify the integrity of a module by scanning for errors in
the module's manifest file.

TODO: add more checks
    `,
}

var checkB = cmdCheck.Flag.Bool("b", false, "")

func runCheck(cmd *Command, args []string) {
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
        err = checkModuleIntegrity(abs, *checkB)
        if err != nil {
            fmt.Fprintf(os.Stderr, "is: %s\n", err.Error())
        }
    }
}

func checkIntegrityCache(dir string) bool {
    return moduleIntegrityCache[dir]
}

func checkModuleIntegrity(dir string, verbose bool) error {
    err := errors.New("unimplemented feature - chk")

    if err == nil {
        moduleIntegrityCache[dir] = true
    }
    return err
}