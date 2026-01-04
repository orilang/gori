package walk

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
