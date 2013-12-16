package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	// TODO foreach args, attempt to do all the following
	// TODO refactor into function - checkPackage...etc
	vcs := matchVcsPath(args[0])
	if vcs != nil {
		fmt.Println("VCS okay!")
	} else {
		fmt.Println("vcs nil")
		return
	}

	// TEST construct path on local file system to where module will go
	sdkpath := os.Getenv("INTERFACESDKROOT")
	targetPath := filepath.Join(sdkpath, "source", "modules")
	for _, folder := range strings.Split(vcs.repo, "/") {
		targetPath = filepath.Join(targetPath, folder)
	}
	// TEST check for conflicts (if so, update?)
	// TEST check for downloading or uploading
	fi, err := os.Stat(targetPath)
	if err != nil {
		fmt.Println("Repo doesn't exist yet, lets download not update")
		if *getU {
			fmt.Println("They wanted to update...but that's not right")
		}
	}
	if err == nil {
		// already exists, let's update rather than download
		fmt.Println(fi.Name(), "already exists, lets update!")
	}
}

func init() {
	cmdGet.Run = runGet
}


