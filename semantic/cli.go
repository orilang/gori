package semantic

import (
	"os"

	"github.com/orilang/gori/lexer"
	"github.com/orilang/gori/parser"
	"github.com/orilang/gori/walk"
)

// NewTypeChecker returns files config to StartParsing
func NewTypeChecker(config Config) (*Files, error) {
	w, err := walk.Walk(walk.Config{File: config.File, Directory: config.Directory})
	if err != nil {
		return nil, err
	}

	return &Files{
		Files:  w.Files,
		output: config.Output,
	}, nil
}

// StartChecking ranges over files to type check the AST
func (f *Files) StartTypeChecking() error {
	for _, file := range f.Files {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		l := lexer.New(data)
		l.Tokenize()
		p := parser.New(l.Tokens)
		pr := p.ParseFile()

		if len(p.Errors) > 0 {
			return p.Errors[0]
		}

		check := NewChecker()
		result := check.Check(pr)
		for _, v := range result {
			return v.Err
		}
	}
	return nil
}
