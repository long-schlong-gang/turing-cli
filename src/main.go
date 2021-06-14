package main

import (
	"os"

	"github.com/teris-io/cli"

	"github.com/long-schlong-gang/turing-cli/commands"
)

func main() {

	app := cli.New("Cypher CLI").
		WithCommand(commands.Decypher)

	os.Exit(app.Run(os.Args, os.Stdout))
}

// TODO: Replace all constant vars with functions so they don't get diddled with by shitty devs
