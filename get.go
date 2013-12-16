package main

import (
	"fmt"
)

var cmdGet = &Command{
	UsageLine: "get [-v verbose] [-c check/validate module] [modules]",
	Short: "download and package modules",
}

var getV = cmdGet.Flag.Bool("v", false, "")
var getC = cmdGet.Flag.Bool("c", false, "")

func init() {
	cmdGet.Run = runGet
}

func runGet(cmd *Command, args []string) {
	fmt.Println("Running sdk get")
}
