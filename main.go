package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

var flagVerbose *bool
var flagQuick *bool

// Approach on handling subcommands has largely been influenced by
// golang's approach towards the `go` cli command.
// This can be found at src/cmd/go/main.go in the golang source at
// http://code.google.com/p/go/source/browse/src/cmd/go/

type Command struct {
	// Runs the command
	Run func(cmd *Command, args []string)

	// Usage and short information on the command
	UsageLine string
	Short     string
	Long      string

	// Flags specific for this command
	Flag flag.FlagSet
}

func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
}

var commands = []*Command{
	cmdGet,
	cmdMv,
	cmdBuild,
	cmdCheck,
	cmdRun,
	cmdVersion,
}

func generalUsage() {
	fmt.Println(`
Is is a tool for managing the Colorado School of Mines Interface SDK exhibit.
Use "is help [flags] [command]" for more information about a command.
Use "is [command] [command args] to execute the command.

Available flags:
	-v [verbose]: Prints detailed information on the status of get as it retrieves
		and builds each module.
	-q [quick]: Uses the first option rather than prompting the user on how to
	    proceed.

Available commands:
	`)
	for _, cmd := range commands {
		fmt.Println(cmd.Name(), ": ", cmd.Short)
	}
	os.Exit(0)
}

func processHelp(helpCmd string) {
	for _, cmd := range commands {
		if cmd.Name() == helpCmd {
			cmd.Usage()
			os.Exit(0)
		}
	}
	generalUsage()
}

func main() {
	flagVerbose = flag.Bool("v", false, "t")
	flagQuick = flag.Bool("q", false, "")
	flag.Parse()
	args := flag.Args()

	fmt.Println(*flagVerbose, *flagQuick, args)

	if len(args) < 1 {
		fmt.Println("usage...")
		generalUsage()
	}

	if args[0] == "help" {
		if len(args) < 2 {
			generalUsage()
		} else {
			processHelp(args[1])
		}
	}

	// check if we have the correct environment variables
	sdkpath := os.Getenv("INTERFACESDKROOT")
	if len(sdkpath) == 0 {
		fmt.Fprintln(os.Stderr, "is: INTERFACESDKROOT environment variable not set")
		os.Exit(1)
	} else if strings.HasPrefix(sdkpath, "~") {
		// check for invalid path for the environment variable
		fmt.Fprintf(os.Stderr, `is: INTERFACESDKROOT can not start with shell metacharacter '~': %q\n`, sdkpath)
		os.Exit(1)
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			cmd.Flag.Parse(args[1:])
			cmd.Run(cmd, cmd.Flag.Args())
			os.Exit(0)
		}
	}

	generalUsage()
}
