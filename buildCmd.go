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
    
}

var buildTypes = []*buildType{
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
    
}

func (b *buildCmd) build(dir string) {

}

func getBuildInfo(dir string) *buildType {

}

func getBuildCmd(type string) *buildCmd {

}