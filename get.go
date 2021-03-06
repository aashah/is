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
	UsageLine: "get [-b] [module paths]",
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
	-b [Build]: Attempts to build the module and place it into the appropriate
	directory.

See "go help chk" to understand more on what is meant by the integrity of a
module.
	`,
}

var getB = cmdGet.Flag.Bool("b", false, "")

func runGet(cmd *Command, args []string) {
	for _, modulePath := range args {
		info := matchVcsPath(modulePath)
		if info == nil {
			fmt.Println("VCS nil, not okay!")
			continue
		}
		err := getModule(info)
		if err != nil {
			fmt.Fprintf(os.Stderr, "is: %s\n", err.Error())
		}
	}
}

func getModule(info *vcsInfo) error {
	// TEST construct path on local file system to where module will go
	sdkpath := os.Getenv("INTERFACESDKROOT")
	targetPath := filepath.Join(sdkpath, "source", "modules")
	for _, folder := range strings.Split(info.path, "/") {
		targetPath = filepath.Join(targetPath, folder)
	}

	// TEST check for conflicts (if so, update?)
	// TEST check for downloading or uploading
	fi, err := os.Stat(targetPath)
	if err != nil {
		err = downloadModule(info, targetPath)
	} else {
		if !fi.IsDir() {
			err = errors.New(fmt.Sprintf("is: found file, not directory - %s\n", fi.Name()))
		} else {
			// already exists, let's update rather than download
			err = updateModule(info, targetPath)
		}
	}

	if err != nil {
		return err
	}

	// FUTURE When the `chk` routine becomes more stable, it may become more
	// appropriate to re-introduce the cache & make it mandatory before
	// executing the `build` routine. However, in it's present state, it remains
	// incomplete; therefor, it's best to ignore it for now.

	if *getB {
		err = buildModule(targetPath, *flagVerbose)
		if err != nil {
			return err
		}
	}
	return nil
}

func downloadModule(info *vcsInfo, targetPath string) error {
	// make the directory, call the download command onto that directory
	// if verbose mode, output any data from the command
	if err := os.MkdirAll(targetPath, 0777); err != nil {
		return err
	}

	if err := info.vcs.download(targetPath, info.repo, *flagVerbose); err != nil {
		return err
	}
	return nil
}

func updateModule(info *vcsInfo, targetPath string) error {
	// cd into the directory, call the update command
	// if verbose mode, output any data
	if err := info.vcs.update(targetPath, *flagVerbose); err != nil {
		return err
	}
	return nil
}
