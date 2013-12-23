package main

import (
    "errors"
    // "fmt"
    // "os"
    // "os/exec"
    // "regexp"
    // "strings"
)

func init() {
}

type findTarget func(matches map[string]string) (string, error)

type buildCmd struct {
    name string
    cmd string // name of executable
    buildCmd string // command to run on how to build
    params map[string]string
}

type buildType struct {
    name string
    files map[string]string // file(s) that must exist to be considered part of this build
    getTarget findTarget
}

type buildInfo struct {
    cmd *buildCmd
    target string // final location of jar
}

var buildList = []*buildCmd{
    mavenBuildCmd,
    antBuildCmd,
}

var buildTypes = []*buildType{
    mavenBuildType,
    antBuildType,
}

func (b *buildInfo) build(dir string, verbose bool) (err error) {
    return runCmd(dir, b.cmd.cmd, b.cmd.buildCmd, verbose, b.cmd.params)
}

func getBuildInfo(dir string, quick bool, verbose bool) (info *buildInfo, err error) {
    for _, bType := range buildTypes {
        // ensure all files have a match
        matches := make(map[string]string)
        for patternName, pattern := range bType.files {
            var path string
            if path, err = findFileInsideModulePackage(dir, pattern, quick, verbose); err != nil {
                break                
            }
            matches[patternName] = path
        }
        if len(matches) == len(bType.files) {
            var target string
            if target, err = bType.getTarget(matches); err != nil {
                return nil, err
            }

            info = &buildInfo{}
            if cmd := getBuildCmd(bType.name); cmd != nil {
                info.cmd = cmd
                info.target = target
                return
            }


            return nil, errors.New("Unknown error retrieving build information")
        }

    }
    return nil, errors.New("Could not find matching build command")
}

func getBuildCmd(buildTypeName string) *buildCmd {
    for _, cmd := range buildList {
        if cmd.name == buildTypeName {
            return cmd
        }
    }
    return nil
}