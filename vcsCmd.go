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
	repo string
	regexp *regexp.Regexp
}

var vcsPaths = []*vcsPath{
	// Github
	{
		code: "git",
		prefix: "github.com",
		repo: "{prefix}/{name}/{repo}",
		re: `^(?P<prefix>github\.com/)(?P<name>[A-Za-z0-9_.\-]+)/(?P<repo>[A-Za-z0-9_.\-]+)(/[A-Za-z0-9_.\-]+)*$`,
	},
}

type vcsInstance struct {
	vcs *vcsCmd
	repo string
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
func matchVcsPath(modulePath string) *vcsInstance {
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

		vcsMatchCmd := vcsByCmd(vcs.code)
		if vcsMatchCmd == nil {
			fmt.Errorf("unknown version control system %s\n", modulePath)
			return nil
		}

		vcsRepo := expand(segments, vcs.repo)
		return &vcsInstance{vcsMatchCmd, vcsRepo}
	}
	return nil
}

func expand(list map[string]string, str string) string {
	ret := str
	for key, val := range list {
		ret = strings.Replace(ret, "{" + key + "}", val, -1)
	}
	return ret
}

// TODO: match package name with vcsCmd
// TODO: Run the vcsCmd on a given directory
// - Check for conflicts (permissions, directory exists (do I need to only
// update?))

