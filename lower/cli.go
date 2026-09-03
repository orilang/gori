package lower

import (
	"os"

	"github.com/orilang/gori/commons"
	"github.com/orilang/gori/lexer"
	"github.com/orilang/gori/parser"
	"github.com/orilang/gori/semantic"
	"github.com/orilang/gori/walk"
)

// NewLower initialize Lower requirements
func NewLowerCLI(config commons.Config) (*Files, error) {
	w, err := walk.Walk(walk.Config{File: config.File, Directory: config.Directory})
	if err != nil {
		return nil, err
	}

	return &Files{
		Files:  w.Files,
		output: config.Output,
	}, nil
}

// StartLowering ranges over files to lower type checker output
func (f *Files) StartLowering() error {
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

		check := semantic.NewChecker()
		program, diagnostics := check.Check(pr)
		if diagnostics.HasErrors() {
			for _, v := range diagnostics {
				return v.Err
			}
		}

		lower := NewLower(f.output)
		_, ld := lower.Lower(program)
		if ld.HasErrors() {
			for _, v := range ld {
				return v.Err
			}
		}
	}
	return nil
}
