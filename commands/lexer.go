package commands

import (
	"context"

	"github.com/orilang/gori/lexer"
	"github.com/orilang/gori/walk"
	"github.com/urfave/cli/v3"
)

func Lexer() *cli.Command {
	var app lexer.Config

	return &cli.Command{
		Name:  "lex",
		Usage: "option to parse file or directory",
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
		},
		Action: func(ctx context.Context, _ *cli.Command) error {
			if app.File == "" && app.Directory == "" {
				return walk.ErrNoFileOrDirectoryPassed
			}

			lex, err := lexer.NewLexer(app)
			if err != nil {
				return err
			}

			return lex.StartLexing()
		},
	}
}
