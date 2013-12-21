package main

import (
    "github.com/aashah/glob"
    "encoding/xml"
    "errors"
    "fmt"
    "os"
)

type xmlManifest struct {
    XMLName xml.Name `xml:"manifest"`
    Package string `xml:"package,attr"`
    Class string `xml:"class,attr"`

    Sdks []*xmlSDK `xml:"uses-sdk"`
    Modules []*xmlModule `xml:"module"`
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