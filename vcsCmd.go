package main

import (
	"fmt"
	"strings"
	"regexp"
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

var vcsGit = &vcsCmd{
	name: "Git",
	cmd: "git",

	createCmd: "clone {repo} {dir}",
	updateCmd: "pull --ff-only",
}

var vcsList = []*vcsCmd{
	vcsGit,
}

type vcsPath struct {
	code string
	prefix string
	name string
	re string
	regexp *regexp.Regexp
}

var vcsPaths = []*vcsPath{
	// Github
	{
		code: "git",
		prefix: "github.com/",
		re: `^(?P<prefix>github\.com/)(?P<name>[A-Za-z0-9_.\-]+)/(?P<repo>[A-Za-z0-9_.\-]+)(/[A-Za-z0-9_.\-]+)*$`,
	},
}

func init() {
	for _, srv := range vcsPaths {
		srv.regexp = regexp.MustCompile(srv.re)
	}
}

func vcsByCmd(cmd string) *vcsCmd {
	for _, vcs := range vcsList {
		if vcs.cmd == cmd {
			return vcs
		}
	}
	return nil
}

// Given a module path such as github.com/user/repo, return the appropriate
// handle for that version control system and where it should ideally be stored
// on the file system
func matchVcsPath(modulePath string) (*vcsCmd, map[string]string) {
	// TODO: check for signs of a malformed module path (://, ../, ./)
	for _, vcs := range vcsPaths {
		fmt.Println("Checking", modulePath, "against", vcs.prefix)
		if !strings.HasPrefix(modulePath, vcs.prefix) {
			continue
		}

		matches := vcs.regexp.FindStringSubmatch(modulePath)
		if matches == nil {
			fmt.Errorf("no match found against %s\n", modulePath)
			continue
		}

		segments := make(map[string]string)
		for i, name := range vcs.regexp.SubexpNames() {
			if name != "" && segments[name] == "" {
				segments[name] = matches[i]
			}
		}

		if segments["name"] == "" || segments["repo"] == "" {
			fmt.Errorf("unknown version control system %s\n", modulePath)
			return nil, nil
		}

		vcsMatchCmd := vcsByCmd(vcs.code)
		if vcsMatchCmd == nil {
			fmt.Errorf("unknown version control system %s\n", modulePath)
			return nil, nil
		}
		return vcsMatchCmd, segments
	}
	return nil, nil
}

// TODO: match package name with vcsCmd
// TODO: Run the vcsCmd on a given directory
// - Check for conflicts (permissions, directory exists (do I need to only
// update?))

