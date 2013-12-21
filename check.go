package main

import (
    // "encoding/xml"
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
Check will attempt to verify the integrity of a module primarily by looking at 
the module's and the Interface SDK's manifest files.

Module Manifest file:
    - Ensure it's valid XML
    - Ensure the class file exists
    - The inputs are provided through the Hardware Manager's Manifest file
    - All the key attributes to the manifest's data exist (package, class, uses-sdk)

Flags:
    -v [Verbose]:
    -q [Quick]: Uses the first option rather than prompting the user on how to
    proceed.
    `,
}

var checkV = cmdCheck.Flag.Bool("v", false, "")
var checkQ = cmdCheck.Flag.Bool("q", false, "")

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

        // Try to find the manifest file
        manifest, err := findFileInsideModulePackage(abs, "**/manifest.xml", *checkQ, *checkV)

        // Error could be from globbing, number of found matches, or from
        // the prompting of which manifest to pick
        if err != nil {
            continue
        }

        if *checkV {
            fmt.Fprintf(os.Stdout, "[info] is: using manifest.xml - %s\n", manifest)
        }

        err = checkModuleIntegrity(abs, manifest, *checkV)
        if err != nil {
            fmt.Fprintf(os.Stderr, "is: %s\n", err.Error())
        }
    }
}

func checkIntegrityCache(dir string) bool {
    return moduleIntegrityCache[dir]
}

func checkModuleIntegrity(moduleRoot string, manifest string, verbose bool) error {
    err := errors.New("unimplemented feature - chk")

    if err == nil {
        moduleIntegrityCache[moduleRoot] = true 
    }
    return err
}

type Manifest struct {
    pkg string `xml:"package, attr"`
}