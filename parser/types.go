package parser

import (
	"github.com/orilang/gori/token"
)

// Config holds file or directory to use for tokenization
type Config struct {
	// File to parse
	File string

	// Directory to take as input and list files to parse
	Directory string

	// Output when set to true outputs the AST
	Output bool
}

// Files holds all files to use for tokenization
type Files struct {
	// Files holds the list of files to parse
	Files []string

	// output when set to true outputs the AST
	output bool
}

// Parser holds requirements with the tokens from the Lexer to
// build the Abstract Syntax Tree (AST)
type Parser struct {
	Tokens   []token.Token
	errors   []error
	size     int
	position int
}
