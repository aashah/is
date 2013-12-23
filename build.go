package main

import (
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
        err = buildModule(abs, *flagVerbose)
        if err != nil {
            fmt.Fprintf(os.Stderr, "is: %s\n", err.Error())
        }
    }
}

func buildModule(dir string, verbose bool) (err error) {
    var info *buildInfo
    var target string

    if info, err = getBuildInfo(dir, true, verbose); err != nil {
        fmt.Println(err)
        return
    }

    // call build
    if err = info.build(dir, verbose); err != nil {
        fmt.Println(err)
        return
    }

    // find jar
    if target, err = findFileInsideModulePackage(dir, info.target, true, verbose); err != nil {
        return err
    }

    if verbose {
        fmt.Println("Found target:", target)        
    }
    
   
    // move file
    // prepare destination
    var dstDirectory string
    dstDirectory = os.Getenv("INTERFACESDKROOT")
    dstDirectory = filepath.Join(dstDirectory, "development")
    dstDirectory = filepath.Join(dstDirectory, "modules")

    if err = copyFile(target, dstDirectory); err != nil {
        return err
    }

    return
}