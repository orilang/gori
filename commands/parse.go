package commands

import (
	"context"

	"github.com/orilang/gori/parser"
	"github.com/orilang/gori/walk"
	"github.com/urfave/cli/v3"
)

func Parse() *cli.Command {
	var app parser.Config

	return &cli.Command{
		Name:  "parse",
		Usage: "option to parse file or directory and dump the AST",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				Usage:       "file to use",
				Destination: &app.File,
			},
			&cli.StringFlag{
				Name:        "directory",
				Aliases:     []string{"d"},
				Usage:       "directory to use",
				Destination: &app.Directory,
			},
			&cli.BoolFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "print the AST",
				Destination: &app.Output,
				Value:       true,
			},
		},
		Action: func(ctx context.Context, _ *cli.Command) error {
			if app.File == "" && app.Directory == "" {
				return walk.ErrNoFileOrDirectoryPassed
			}

			p, err := parser.NewParser(app)
			if err != nil {
				return err
			}

			return p.StartParsing()
		},
	}
}
