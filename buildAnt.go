package main

// NOT DONE NOR TESTED
import (
    // "fmt"
    // "errors"
    // "encoding/xml"
    // "path/filepath"
)

var antDependencies = map[string]string{
    "ant": "build.xml",
}

var antBuildCmd = &buildCmd{
    name: "Ant",
    cmd: "ant",
    buildCmd: "compile jar",
    params: nil,
}

var antbuildType = &buildType{
    name: "Ant",
    files: antDependencies,
    getTarget: findAntTarget,
}

func findAntTarget(matches map[string]string) (target string, err error) {
    return nil, nil
}