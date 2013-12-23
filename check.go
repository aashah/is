package main

import (
    "encoding/xml"
    "errors"
    "fmt"
    "os"
    "path/filepath"
)

func init() {
    cmdCheck.Run = runCheck
}

var cmdCheck = &Command {
    UsageLine: "chk [path to modules]",
    Short: "Checks a given list of modules",
    Long: `
Check will attempt to verify the integrity of a module primarily by looking at 
the module's and the Interface SDK's manifest files.

Module Manifest file:
    - Ensure it's valid & complete XML
    - The inputs are provided through the Hardware Manager's Manifest file
    - All the key attributes to the manifest's data exist. Presently, all
    attributes are considered mandatory.
    `,
}

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

        valid, err = checkModuleIntegrity(abs, *flagVerbose)
        if err != nil {
            fmt.Fprintf(os.Stderr, "is: %s\n", err.Error())
        }

        fmt.Println("IS VALID:", valid)
    }
}

func checkModuleIntegrity(moduleRoot string, verbose bool) (valid bool, err error) {

    var hardwareManifest *xmlHardwareManifest
    var moduleManifest *xmlModuleManifest

    if moduleManifest, err = loadModuleManifest(moduleRoot, verbose); err != nil {
        return false, err
    }

    if hardwareManifest, err = loadHardwareManagerManifest(verbose); err != nil {
        return false, err
    }

    /*
     * Checks:
     * - valid & complete XML (*)
     * - Ensure primary directories from module manifest exist
     *     -> Depends greatly on how the code is structured, hard to be sure.
     *     -> Especially since most the important ones are java-reliant
     * - inputs in module manifest can be found in hardware manager manifest (*)
     * - key attributes exist (with verbose display optional ones that aren't provided) (*)
     */
    
    return moduleManifest.isValid(hardwareManifest)
}

func loadModuleManifest(moduleRoot string, verbose bool) (moduleManifest *xmlModuleManifest, err error) {
    var moduleManifestPath string
    var raw []byte

    if moduleManifestPath, err = findFileInsideModulePackage(moduleRoot, "**/manifest.xml", *flagQuick, *flagVerbose); err != nil {
        return nil, err
    }

    if verbose {
        fmt.Printf("[info] is: using module manifest.xml - %s\n", moduleManifestPath)
    }

    if raw, err = readFile(moduleManifestPath); err == nil {
        err = xml.Unmarshal(raw, &moduleManifest)
    }

    if err != nil {
        return nil, err
    }

    if verbose {
        moduleManifest.print()
    }

    return
}

func loadHardwareManagerManifest(verbose bool) (hardwareManifest *xmlHardwareManifest, err error) {
    var sdkPath, hardwareManifestPath string
    var raw []byte

    if sdkPath = os.Getenv("INTERFACESDKROOT"); len(sdkPath) == 0 {
        return nil, errors.New("INTERFACESDKROOT env variable not set")
    }

    sdkPath = filepath.Join(sdkPath, "source", "interfaceSDK")

    // find hw manager manifest
    if hardwareManifestPath, err = findFileInsideModulePackage(sdkPath, "**/hardware_manager_manifest.xml", *flagQuick, *flagVerbose); err != nil {
        return nil, err
    }

    if verbose {
        fmt.Printf("[info] is: using hardware_manager_manifest.xml - %s\n", hardwareManifestPath)
    }

    if raw, err = readFile(hardwareManifestPath); err == nil {
        err = xml.Unmarshal(raw, &hardwareManifest)
    }

    if err != nil {
        return nil, err
    }

    if verbose {
        printHardwareManifest(hardwareManifest)
    }

    return
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