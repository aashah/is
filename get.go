package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	cmdGet.Run = runGet
}

var cmdGet = &Command{
	UsageLine: "get [-v] [-b] [module paths]",
	Short:     "download and package modules",
	Long: `
Get downloads/updates modules as well as building the module for the interface
sdk. Get will also check for errors in an individual module's structure before
building.

By default, get will download the module. If the module already exists, get will
update the package if possible.

Example usages: 
	"is get github.com/aashah/Pong": Downloads and checks the integrity of the
	module located at github.com/aashah/Pong

	"is get -b github.com/aashah/Pong": Downloads, checks the module, and builds
	the module to be immediately used under the sdk

	"is get -b github.com/aashah/Pong github.com/aashah/fractalModule": Downloads
	several modules

Flags:
	-v [Verbose]: Prints detailed information on the status of get as it retrieves
	and builds each module.
	
	-b [Build]: Attempts to build the module and place it into the appropriate
	directory.
	`,
}

var getV = cmdGet.Flag.Bool("v", false, "")
var getB = cmdGet.Flag.Bool("b", false, "")

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
		err = downloadModule(vcs, targetPath)
	} else {
		if !fi.IsDir() {
			err = errors.New(fmt.Sprintf("is: found file, not directory: %s\n", fi.Name()))
		} else {
			// already exists, let's update rather than download
			fmt.Println(fi.Name(), "already exists, lets update!")
			err = updateModule(vcs, targetPath)
		}
	}

	if err == nil {
		// err = integrity of module
		err = errors.New("is: incomplete validation of module. do not build yet")
		if err != nil && *getB {
			// attempt to build module
		}
	} else {
		// error downloading/updating
		fmt.Print(err)
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
