package main

import (
    "fmt"
)

func init() {
    cmdRun.Run = runRun
}

var cmdRun = &Command {
    UsageLine: "Run -v",
    Short: "Runs the interface SDK",
    Long: `
Description
    `
}

func runRun(cmd *Command, args []string) {
    
}