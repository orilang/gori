package token

type Token struct {
	Kind     Kind
	Value    string
	Line     uint
	EOL, EOF bool
}
