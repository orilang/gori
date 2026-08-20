package commands

import (
	"context"

	"github.com/orilang/gori/semantic"
	"github.com/orilang/gori/walk"
	"github.com/urfave/cli/v3"
)

func Check() *cli.Command {
	var app semantic.Config

	return &cli.Command{
		Name:  "check",
		Usage: "option to type check file or directory",
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

			tp, err := semantic.NewTypeChecker(app)
			if err != nil {
				return err
			}

			return tp.StartTypeChecking()
		},
	}
}
