package main

import "fmt"

const (
    major int = 0
    minor int = 0
    patch int = 0
)

func init() {
    cmdVersion.Run = runVersion
}

var cmdVersion = &Command{
    UsageLine: "version",
    Short: "prints the program version",
    Long: "prints the program version",
}

func runVersion(cmd *Command, args []string) {
    fmt.Printf("is version: %d.%d.%d\n", 
                major,
                minor,
                patch)
}