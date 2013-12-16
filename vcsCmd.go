package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"regexp"
)

func init() {
	for _, srv := range vcsPaths {
		srv.regexp = regexp.MustCompile(srv.re)
	}
}

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

type vcsPath struct {
	code string
	prefix string
	name string
	re string
	path string
	repo string
	regexp *regexp.Regexp
}

type vcsInfo struct {
	vcs *vcsCmd
	path string
	repo string
}

func (v *vcsCmd) String() string {
	return v.name
}

func (v *vcsCmd) runCmd(dir string, cmdLine string, verbose bool, keyvals map[string]string) ([]byte, error) {
	if keyvals != nil {
		// expand cmd
		cmdLine = expand(keyvals, cmdLine)
	}

	args := strings.Fields(cmdLine)


	// TODO Check if v.cmd exists
	if _, err := exec.LookPath(v.cmd); err != nil {
		fmt.Fprintf(os.Stderr, "is: missing %s command.", v.cmd)
		return nil, err
	}
	fmt.Println("Executing", v.cmd, "with", args)
	// Execute
	cmd := exec.Command(v.cmd, args...)
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	out := buf.Bytes()
	os.Stdout.Write(out)
	return out, err
}

func (v *vcsCmd) download(dir string, repo string, verbose bool) error {
	keyvals := map[string]string {
		"repo": repo,
		"dir": dir,
	}
	_, err := v.runCmd(dir, v.createCmd, verbose, keyvals)
	return err
}

func (v *vcsCmd) update(dir string, verbose bool) error {
	_, err := v.runCmd(dir, v.updateCmd, verbose, nil)
	return err
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

var vcsPaths = []*vcsPath{
	// Github
	{
		code: "git",
		prefix: "github.com",
		path: "{prefix}/{name}/{repo}",
		repo: "git@{prefix}:{name}/{repo}.git",
		re: `^(?P<prefix>github\.com)/(?P<name>[A-Za-z0-9_.\-]+)/(?P<repo>[A-Za-z0-9_.\-]+)(/[A-Za-z0-9_.\-]+)*$`,
	},
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
func matchVcsPath(modulePath string) *vcsInfo {
	// TODO: check for signs of a malformed module path (://, ../, ./)
	for _, vcs := range vcsPaths {
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

		return &vcsInfo{
			vcsMatchCmd,
			expand(segments, vcs.path),
			expand(segments, vcs.repo),
		}
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

