package main

import (
    "github.com/aashah/glob"
    "encoding/xml"
    "errors"
    "fmt"
    "os"
)

type xmlModuleManifest struct {
    XMLName xml.Name `xml:"manifest"`
    Package string `xml:"package,attr"`
    Class string `xml:"class,attr"`

    SDK xmlSDK `xml:"uses-sdk"`
    Module xmlModule `xml:"module"`
}

type xmlSDK struct {
    XMLName xml.Name `xml:"uses-sdk"`

    Min string `xml:"minSdkVersion,attr"`
    Target string `xml:"targetSdkVersion,attr"`
}

type xmlModule struct {
    XMLName xml.Name `xml:"module"`

    Icon string `xml:"icon,attr"`
    Title string `xml:"title,attr"`
    Author string `xml:"author,attr"`
    Version string `xml:"version,attr"`

    Inputs []*xmlInput `xml:"inputs>input"`
    Requires []*xmlRequire `xml:"requires-module"`
}

type xmlInput struct {
    XMLName xml.Name `xml:"input"`

    InputType string `xml:"input-type,attr"`
}

type xmlRequire struct {
    XMLName xml.Name `xml:"requires-module"`
}

func (x *xmlModuleManifest) isValid(hardwareManifest *xmlHardwareManifest) (valid bool, err error) {
    switch {
        case x.Package == "":
            return false, errors.New("manifest>package attribute not set")
        case x.Class == "":
            return false, errors.New("manifest>class attribute not set")
        case x.SDK == (xmlSDK{}):
            return false, errors.New("manifest>uses-sdk not provided")
    }

    return true, nil
}

func (x *xmlModuleManifest) print() {
    fmt.Printf("Package: %q\n", x.Package)
    fmt.Printf("Class: %q\n", x.Class)

    fmt.Printf("Min: %q\nTarget: %q\n", x.SDK.Min, x.SDK.Target)

    module := x.Module
    fmt.Printf("Icon: %q\nTitle: %q\nAuthor: %q\nVersion: %q\n",
        module.Icon, module.Title, module.Author, module.Version)
    fmt.Println("Inputs:")
    for _, input := range module.Inputs {
        fmt.Println("-", input.InputType)
    }
}

func findFileInsideModulePackage(moduleRoot string, pattern string, quick bool, verbose bool) (string, error) {

    var file string
    matches, err := glob.Glob(moduleRoot, pattern);


    if err == nil {
        if verbose {
            fmt.Fprintf(os.Stdout, "[info] is: found %d matches:\n", len(matches))
            for i, match := range matches {
                fmt.Fprintf(os.Stdout, "[info] is: (%d) %s:\n", i, match)
            }
        }
        var whichFile int

        switch {
        case len(matches) <= 0:
            whichFile, err = -1, errors.New(fmt.Sprintf("no matches found - %s\n", pattern))
        case len(matches) == 1:
            whichFile = 0
        case len(matches) > 1:
            if quick {
                whichFile = 0
            } else {
                whichFile, err = promptEntryFromArray("Please pick a file file", matches)            
            }
        }

        if err != nil {
            return file, err
        }
        file = matches[whichFile]
    }

    return file, err
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