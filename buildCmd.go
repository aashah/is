package main

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

type findTarget func() string

type buildCmd struct {
    name string
    cmd string // name of executable
    build string // command to run on how to build
    target string // final location of jar
}

type buildType struct {
    name string
    files []string // file(s) that must exist to be considered part of this build
    getTarget findTarget
}

var buildList = []*buildCmd{
    {
        name: "Maven",
        cmd: "mvn",
        build: "clean package",
    },
}

var buildTypes = []*buildType{
    {
        name: "Maven",
        files: []string{
            "pom.xml",
        },
        getTarget: func() string {
            // TODO load pom file to find output directory
            // Target dir:
            // - project.build.outputDirectory
            // - Defaults => {project.baseDir}/target
            // Target name:
            // - project.build.finalName
            // - Defaults => {project.artifaceId}-{project.version}
            return "t"
        },
    },
}

func (b *buildCmd) build(dir string) {

}

func getBuildInfo(dir string) *buildType {

}

func getBuildCmd(type string) *buildCmd {

}