package main

import (
	"context"
	"log"
	"os"

	"github.com/orilang/gori/commands"
	"github.com/urfave/cli/v3"
)

func main() {
	usage := "A new cli for Ori purposes"
	description := "Gori is Ori lexer and parser"

	cmd := cli.Command{
		Name:                  "gori",
		Usage:                 usage,
		Description:           description,
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			commands.Lexer(),
			commands.Parse(),
			commands.Check(),
			commands.Lower(),
		},
	}

	// remove date/timespamp from log output
	// https://stackoverflow.com/questions/48629988/remove-timestamp-prefix-from-go-logger
	log.SetFlags(0)

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err.Error())
		return
	}
	log.Println("No errors found")
}
