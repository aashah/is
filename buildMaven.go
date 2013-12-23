package main

import (
   "fmt"
   "errors"
   "encoding/xml"
   "path/filepath"
)

var mavenDependencies = map[string]string{
    "pom": "pom.xml",
}
var mavenBuildCmd = &buildCmd{
    name: "Maven",
    cmd: "mvn",
    buildCmd: "clean package",
    params: map[string]string{},
}

var mavenBuildType = &buildType{
    name: "Maven",
    files: mavenDependencies,
    getTarget: findMavenTarget,
}

func findMavenTarget(matches map[string]string) (target string, err error) {
    fmt.Println("matches:", matches)

    var mavenProject *xmlMavenProject
    var raw []byte

    if raw, err = readFile(matches["pom"]); err == nil {
        err = xml.Unmarshal(raw, &mavenProject)
    }

    var OutDirectory string
    var FinalName string

    OutDirectory = mavenProject.Build.OutputDirectory
    FinalName = mavenProject.Build.FinalName

    if mavenProject.Build.OutputDirectory == ""{
        OutDirectory = "target"
    }

    if mavenProject.Build.FinalName == ""{
        switch {
        case mavenProject.ArtifactId == "":
            return "", errors.New("Resorting to {project.artifactId for target name, but element not found or set")
        case mavenProject.Version == "":
            return "", errors.New("Resorting to {project.version for target name, but element not found or set")
        }

        FinalName = fmt.Sprintf("%s-%s", mavenProject.ArtifactId, mavenProject.Version)
    }

    return filepath.Join(OutDirectory, FinalName), nil  
}

type xmlMavenProject struct {
    XMLName xml.Name `xml:"project"`
    ArtifactId string `xml:"artifactId"`
    Version string `xml:"version"`
    Build xmlMavenBuild `xml:"build"`
}

type xmlMavenBuild struct {
    OutputDirectory string `xml:"outputDirectory"`
    FinalName string `xml:"finalName"`
}

func (x *xmlMavenProject) print() {
    fmt.Println("Out dir:", x.Build.OutputDirectory)
    fmt.Println("Final name:", x.Build.FinalName)
}