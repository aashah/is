package main

import (
    "github.com/aashah/glob"
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
        // Try to find the manifest file
        manifestGlob := "/**/manifest.xml"
        
        var manifest string
        if matches, err := glob.Glob(abs, manifestGlob); err != nil {
            var whichManifest int

            switch {
            case len(matches) <= 0:
                err = errors.New("no manifest.xml found")
            case len(matches) == 1:
                whichManifest = 0
            case len(matches) > 1:
                whichManifest, err = promptEntryFromArray("Please pick a manifest file", matches)
            }

            manifest = matches[whichManifest]
        }
        // Error could be from globbing, number of found matches, or from
        // the prompting of which manifest to pick
        if err != nil {
            continue
        }

        return
        err = checkModuleIntegrity(abs, manifest, *checkB)
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

func promptEntryFromArray(prompt string, entries []string) (idx int, e error) {
    idx = -1
    for idx < 0 || idx >= len(entries) {
        // prompt each entry
        fmt.Println(prompt)
        for i, entry := range entries {
            fmt.Printf("(%d) %s\n", i, entry)
        }

        _, err := fmt.Scanf("%d", &idx)
        if err != nil {
            return -1, err
        }
    }
    return
}