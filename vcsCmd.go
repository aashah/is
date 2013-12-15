package main

import (
	// "fmt"
	// "strings"
)

// We look towards `go get` again as to how best to approach the issue of
// pulling from a remote repository.
// This is influenced by the golang source code src/cmd/go/vcs.go found at:
// http://code.google.com/p/go/source/browser/src/cmd/go/vcs.go

type vcsCmd struct {
	name string
	cmd string

	createCmd string
	updateCmd string
}

func (v *vcsCmd) String() string {
	return v.name
}

var vcsList = []*vcsCmd{
	vcsGit,
}

// TODO: match package name with vcsCmd
// TODO: Run the vcsCmd on a given directory
// - Check for conflicts (permissions, directory exists (do I need to only
// update?))

