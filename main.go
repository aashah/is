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
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Short))
	os.Exit(2)
}

var commands = []*Command{
	cmdGet,
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Println("usage...")
		os.Exit(2)
	}

	if args[0] == "help" {
		fmt.Println("Help...")
		os.Exit(0)
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			cmd.Flag.Parse(args[1:])
			cmd.Run(cmd, args)
			os.Exit(0)
			return
		}
	}
}


