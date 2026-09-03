package commands

import (
	"context"

	"github.com/orilang/gori/commons"
	"github.com/orilang/gori/lower"
	"github.com/orilang/gori/walk"
	"github.com/urfave/cli/v3"
)

func Lower() *cli.Command {
	var app commons.Config

	return &cli.Command{
		Name:  "lower",
		Usage: "option to type lower file or directory",
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
				Usage:       "print the LIR",
				Destination: &app.Output,
				Value:       true,
			},
		},
		Action: func(ctx context.Context, _ *cli.Command) error {
			if app.File == "" && app.Directory == "" {
				return walk.ErrNoFileOrDirectoryPassed
			}

			c, err := lower.NewLowerCLI(app)
			if err != nil {
				return err
			}

			return c.StartLowering()
		},
	}
}
