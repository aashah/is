package main

import (
    "fmt"
)

func init() {
    cmdMv.Run = runMv
}

var cmdMv = &Command {
    UsageLine: "mv [-b] [module root directory] [destination directory]",
    Short: "mv moves pre-existing modules into the appropriate location",
    Long: `
Moves pre-existing modules into the appropriate location. mv will also check
the integrity of the module and attempt to build if the -b flag is set

TODO: add more mvs
    `,
}

func runMv(cmd *Command, args []string) {
    fmt.Println("Running mv")
}