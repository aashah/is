package main

import (
	"fmt"
	"flag"
	"os"
	"strings"
)

// Approach on handling subcommands has largely been influenced by
// golang's approach towards the `go` cli command.
// This can be found at src/cmd/go/main.go in the golang source at
// http://code.google.com/p/go/source/browse/src/cmd/go/

type Command struct {
	// Runs the command
	Run func(cmd *Command, args []string)

	// Usage and short information on the command
	UsageLine string
	Short string
	Long string

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
}

func generalUsage() {
	fmt.Println(`
Is is a tool for managing the Colorado School of Mines Interface SDK exhibit.
Use "is help [command]" for more information about a command.
Use "is [command] [command args] to execute the command.

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
			cmd.Usage();
			os.Exit(0)
		}
	}
	generalUsage()
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("usage...")
		generalUsage()
	}

	if args[0] == "help" {
		processHelp(args[1])
		fmt.Println("Help...")
	}

	// check if we have the correct environmetn variables
	sdkpath := os.Getenv("INTERFACESDKROOT")
	if len(sdkpath) == 0 {
		fmt.Fprintf(os.Stderr, `
is: INTERFACESDKROOT environment variable not set\n
		`)
		os.Exit(0)
	} else {
		// check for invalid path for the environment variable
		if strings.HasPrefix(sdkpath, "~") {
			fmt.Fprintf(os.Stderr, `is: INTERFACESDKROOT can not start with shell
			metacharacter '~': %q\n`, sdkpath)
		}
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			cmd.Flag.Parse(args[1:])
			cmd.Run(cmd, cmd.Flag.Args())
			os.Exit(0)
			return
		}
	}
}


