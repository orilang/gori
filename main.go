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
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err.Error())
	}
}
