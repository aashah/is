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

        valid, err = checkModuleIntegrity(abs, *checkV)
        if err != nil {
            fmt.Fprintf(os.Stderr, "is: %s\n", err.Error())
        }

        fmt.Println("IS VALID:", valid)
    }
}

func checkIntegrityCache(dir string) bool {
    return moduleIntegrityCache[dir]
}

func checkModuleIntegrity(moduleRoot string, verbose bool) (valid bool, err error) {    
    var hardwareManifest *xmlHardwareManifest
    var moduleManifest *xmlModuleManifest

    if moduleManifest, err = loadModuleManifest(moduleRoot); err != nil {
        return false, err
    }

    printModuleManifest(moduleManifest)

    if hardwareManifest, err = loadHardwareManagerManifest(); err != nil {
        return false, err
    }

    printHardwareManifest(hardwareManifest)

    err = errors.New("unimplemented feature - chk")
    if err == nil {
        moduleIntegrityCache[moduleRoot] = true 
    }
    return false, err
}

func loadModuleManifest(moduleRoot string) (moduleManifest *xmlModuleManifest, err error) {
    var moduleManifestPath string
    var raw []byte

    if moduleManifestPath, err = findFileInsideModulePackage(moduleRoot, "**/manifest.xml", *checkQ, *checkV); err != nil {
        return nil, err
    }

    if raw, err = readFile(moduleManifestPath); err == nil {
        err = xml.Unmarshal(raw, &moduleManifest)
    }

    return
}

func loadHardwareManagerManifest() (hardwareManifest *xmlHardwareManifest, err error) {
    var sdkPath, hardwareManifestPath string
    var raw []byte

    if sdkPath = os.Getenv("INTERFACESDKROOT"); len(sdkPath) == 0 {
        return nil, errors.New("INTERFACESDKROOT env variable not set")
    }

    sdkPath = filepath.Join(sdkPath, "source", "interfaceSDK")

    // find hw manager manifest
    if hardwareManifestPath, err = findFileInsideModulePackage(sdkPath, "**/hardware_manager_manifest.xml", *checkQ, *checkV); err != nil {
        return nil, err
    }

    if raw, err = readFile(hardwareManifestPath); err == nil {
        err = xml.Unmarshal(raw, &hardwareManifest)
    }

    if err != nil {
        return nil, err
    }

    return
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

func printModuleManifest(manifest *xmlModuleManifest) {
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

func printHardwareManifest(manifest *xmlHardwareManifest) {
    fmt.Println("Printing hardware manager manifest")

    fmt.Println("-----------------")
    for _, support := range manifest.Functionalities {
        fmt.Printf("Name: %q\nInterface%q\n", support.Name, support.Interface)
    }

    fmt.Println("-----------------")

    for _, device := range manifest.Devices {
        fmt.Printf("Name: %q\nDriver: %q\n", device.Name, device.Driver)
        fmt.Println("SUPPORTS:")
        for _, support := range device.Provides {
            fmt.Println("-", support.Name)
        }
    }
}

type xmlHardwareManifest struct {
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
    Provides []*xmlDeviceSupport `xml:"provides"`
}

type xmlDeviceSupport struct {
    XMLName xml.Name `xml:"provides"`
    Name string `xml:"name,attr"`
}