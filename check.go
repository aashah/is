package main

import (
    "encoding/xml"
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
    - Ensure the package dir/class file exists
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
        var valid bool

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

        valid, err = checkModuleIntegrity(abs, manifest, *checkV)
        if err != nil {
            fmt.Fprintf(os.Stderr, "is: %s\n", err.Error())
        }

        fmt.Println("IS VALID:", valid)
    }
}

func checkIntegrityCache(dir string) bool {
    return moduleIntegrityCache[dir]
}

func checkModuleIntegrity(moduleRoot string, manifestPath string, verbose bool) (valid bool, err error) {
    // var hardwareManagerManifest *xmlHWManifest
    var moduleManifest *xmlModuleManifest

    if moduleManifest, err = loadModuleManifest(manifestPath); err != nil {
        return false, err
    }
    printManifest(moduleManifest)

    // verify

    err = errors.New("unimplemented feature - chk")
    if err == nil {
        moduleIntegrityCache[moduleRoot] = true 
    }
    return false, err
}

func loadModuleManifest(manifestPath string) (moduleManifest *xmlModuleManifest, err error) {
    var raw []byte

    if raw, err = readFile(manifestPath); err == nil {
        err = xml.Unmarshal(raw, &moduleManifest)
    }

    return moduleManifest, err
}

func readFile(path string) (data []byte, err error) {
    var fi *os.File
    var fiStat os.FileInfo

    if fi, err = os.Open(path); err != nil {
        return nil, err
    }
    defer fi.Close()

    if fiStat, err = fi.Stat(); err != nil {
        return nil, err
    }

    var raw []byte

    raw = make([]byte, fiStat.Size())

    if _, err = fi.Read(raw); err != nil {
        return nil, err
    }

    return raw, nil
}

func printManifest(manifest *xmlModuleManifest) {
    fmt.Printf("Package: %q\n", manifest.Package)
    fmt.Printf("Class: %q\n", manifest.Class)

    for _, sdk := range manifest.Sdks {
        fmt.Printf("Min: %q\nTarget: %q\n", sdk.Min, sdk.Target)
    }

    for _, module := range manifest.Modules {
        fmt.Printf("Icon: %q\nTitle: %q\nAuthor: %q\nVersion: %q\n",
            module.Icon, module.Title, module.Author, module.Version)
        fmt.Println("Inputs:")
        for _, input := range module.Inputs {
            fmt.Println("-", input.InputType)
        }
    }
}

type xmlHWManifest struct {
    XMLName xml.Name `xml:"manifest"`
    Functionalities []*xmlSupport `xml:"functionalities>supports"`
    Devices []*xmlDevice `xml:"devices>device"`
}

type xmlSupport struct {
    XMLName xml.Name `xml:"supports"`
    Name string `xml:"name,attr"`
    Interface string `xml:"interface,attr"`
}

type xmlDevice struct {
    XMLName xml.Name `xml:"device"`
    Name string `xml:"name,attr"`
    Driver string `xml:"driver,attr"`
}