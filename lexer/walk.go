package lexer

import (
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// walk checks if specified file exist.
// If directory is specified, it will return all .ori and .mod files
func walk(config Config) (files *Files, err error) {
	var f Files
	exentions := []string{".ori", ".mod"}
	if config.File != "" && slices.Contains(exentions, filepath.Ext(config.File)) {
		if _, err := os.Stat(config.File); err != nil {
			return nil, err
		}

		f.Files = append(f.Files, config.File)
	}

	if config.Directory != "" {
		if err := filepath.WalkDir(
			config.Directory,
			func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}

				if !d.IsDir() && !strings.Contains(path, ".vendor/") && slices.Contains(exentions, filepath.Ext(d.Name())) {
					f.Files = append(f.Files, path)
				}
				return nil
			},
		); err != nil {
			return nil, err
		}
	}

	if len(f.Files) == 0 {
		return nil, ErrNoFilesFound
	}

	return &f, nil
}
