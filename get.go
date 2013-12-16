package main

import (
	"fmt"
)

var cmdGet = &Command{
	UsageLine: "get [-v] [-c] [-u] [modules]",
	Short: "download and package modules",
	Long: `
Get downloads/updates modules as well as building the module for the interface
sdk. Get also can check for errors in an individual modules' structure before
building.

By default, get will download the module. If the module already exists, get will
update the package if possible.

Flags:
	-v [Verbose]: Prints detailed information on the status of get as it retrieves
	and builds each module.
	
	-c [Check/Validate]: Checks the integrity of the module. See {wiki_link} for
	more information regarding what constitutes a valid module structure.

	-u [Update]: Attempts to update a given package. This will not download the
	module in the event of it not existing.

	`,
}

var getV = cmdGet.Flag.Bool("v", false, "")
var getC = cmdGet.Flag.Bool("c", false, "")
var getU = cmdGet.Flag.Bool("u", false, "")

func runGet(cmd *Command, args []string) {
	fmt.Println("Running sdk get", args)

	vcs, segments := matchVcsPath(args[1])
	if vcs != nil {
		fmt.Println(vcs.name, segments)
	} else {
		fmt.Println("vcs nil", segments)
	}

	// TODO construct path on local file system to where module will go
	// TODO check for conflicts (if so, update?)
	// TODO check for downloading or uploading
}

func init() {
	cmdGet.Run = runGet
}


