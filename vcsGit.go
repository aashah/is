package main

var vcsGit = &vcsCmd{
	name: "Git",
	cmd: "git",

	createCmd: "clone {repo} {dir}",
	updateCmd: "pull --ff-only",
}
