package commons

type Config struct {
	// File to parse
	File string

	// Directory to take as input and list files to parse
	Directory string

	// Output when set to true outputs the result
	Output bool
}

// Files holds all files to use for tokenization
type Files struct {
	// Files holds the list of files to parse
	Files []string

	// output when set to true outputs the AST
	output bool
}
