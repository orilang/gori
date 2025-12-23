package lexer

import "errors"

var (
	ErrNoFileOrDirectoryPassed = errors.New("no file or directory passed")
	ErrNoFilesFound            = errors.New("no files found")
)
