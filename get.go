package main

import (
	"fmt"
)

var cmdGet = &Command{
	UsageLine: "get [module",
	Short: "download and package modules",
}

func init() {
	cmdGet.Run = runGet
}

func runGet(cmd *Command, args []string) {
	fmt.Println("Running sdk get")
}
