package lexer

import (
	"github.com/orilang/gori/token"
)

// Config holds file or directory to use for tokenization
type Config struct {
	// File to parse
	File string

	// Directory to take as input and list files to parse
	Directory string
}

// LexerFiles holds all files to use for tokenization
type Files struct {
	// Files holds the list of files to parse
	Files []string
}

// Lexer holds requirements to parse tokens
type Lexer struct {
	Tokens   []token.Token
	input    []byte
	position int
	line     int
	column   int
	size     int
}
