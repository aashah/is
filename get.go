package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var cmdGet = &Command{
	UsageLine: "get [-v] [-c] [module paths]",
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
	`,
}

var getV = cmdGet.Flag.Bool("v", false, "")
var getC = cmdGet.Flag.Bool("c", false, "")

func init() {
	cmdGet.Run = runGet
}

func runGet(cmd *Command, args []string) {
	// TODO foreach args, attempt to do all the following
	// TODO refactor into function - checkPackage...etc
	vcs := matchVcsPath(args[0])
	if vcs == nil {
		fmt.Println("VCS nil, not okay!")
		return
	}

	// TEST construct path on local file system to where module will go
	sdkpath := os.Getenv("INTERFACESDKROOT")
	targetPath := filepath.Join(sdkpath, "source", "modules")
	for _, folder := range strings.Split(vcs.path, "/") {
		targetPath = filepath.Join(targetPath, folder)
	}
	// TEST check for conflicts (if so, update?)
	// TEST check for downloading or uploading
	fi, err := os.Stat(targetPath)
	if err != nil {
		fmt.Println("Repo doesn't exist yet, lets download not update")
		downloadModule(vcs, targetPath)
	}	else {
		if !fi.IsDir() {
			fmt.Println("Not a directory!!!")
		} else {
			// already exists, let's update rather than download
			fmt.Println(fi.Name(), "already exists, lets update!")
			updateModule(vcs, targetPath)
		}
	}
}

func downloadModule(info *vcsInfo, targetPath string) error {
	// make the directory, call the download command onto that directory
	// if verbose mode, output any data from the command
	fmt.Println(info.vcs, targetPath)
	if err := os.MkdirAll(targetPath, 0777); err != nil {
		return err
	}

	if err := info.vcs.download(targetPath, info.repo, *getV); err != nil {
		return err
	}
	return nil
}

func updateModule(info *vcsInfo, targetPath string) error {
	// cd into the directory, call the update command
	// if verbose mode, output any data
	if err := info.vcs.update(targetPath, *getV); err != nil {
		return err
	}
	return nil
}