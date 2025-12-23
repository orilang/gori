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
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		// log.Fatalf("Error occured while executing the program with error: %s", err.Error())
		log.Fatal(err.Error())
	}
}
