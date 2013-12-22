package builds

import (
    "fmt"
    // "os"
    // "os/exec"
    // "regexp"
    // "strings"
)

func init() {
    fmt.Println("Loading build helper file")
}

type findTarget func(matches[] string) string

type buildCmd struct {
    name string
    cmd string // name of executable
    buildCmd string // command to run on how to build
    target string // final location of jar
}

type buildType struct {
    name string
    files []string // file(s) that must exist to be considered part of this build
    getTarget findTarget
}

var buildList = []*buildCmd{
    mavenBuildCmd,
}

var buildTypes = []*buildType{
    mavenBuildType,
}

func (b *buildCmd) build(dir string) {

}

func getBuildInfo(dir string, quick bool, verbose bool) (*buildType, error) {
    for _, bType := range buildTypes {
        // ensure all files have a match
        var matches []string
        for _, filePattern := range bType.files {
            if path, err := findFileInsideModulePackage(dir, filePattern, quick, verbose); err == nil {
                matches = append(matches, path)
            }
        }

        if len(matches) == len(bType.files) {
            // Found a matching build system
            fmt.Println("Using", bType.name, "->", bType.getTarget(matches))
        }
    }
    return nil, nil
}

func getBuildCmd(buildTypeName string) *buildCmd {
    for _, cmd := range buildList {
        if cmd.name == buildTypeName {
            return cmd
        }
    }
    return nil
}