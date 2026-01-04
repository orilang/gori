package token

// Token holds requirements for lexer
type Token struct {
	Kind   Kind
	Value  string
	Line   int
	Column int
}
