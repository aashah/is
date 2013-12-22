package main

var mavenDependencies = map[string]string{
    "pom": "pom.xml",
}
var mavenBuildCmd = &buildCmd{
    name: "Maven",
    cmd: "mvn",
    buildCmd: "clean package",
}

var mavenBuildType = &buildType{
    name: "Maven",
    files: []string{
        "pom.xml",
    },
    getTarget: findMavenTarget,
}

func findMavenTarget(matches[] string) string {
    // TODO load pom file to find output directory
    // Target dir:
    // - project.build.outputDirectory
    // - Defaults => {project.baseDir}/target
    // Target name:
    // - project.build.finalName
    // - Defaults => {project.artifaceId}-{project.version}
    return "MAVEN TARGET"
}