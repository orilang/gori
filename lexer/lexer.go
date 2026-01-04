package lexer

import (
	"fmt"
	"os"

	"github.com/orilang/gori/token"
	"github.com/orilang/gori/walk"
)

// NewLexer returns files config to StartLexing
func NewLexer(config Config) (*Files, error) {
	w, err := walk.Walk(walk.Config{File: config.File, Directory: config.Directory})
	if err != nil {
		return nil, err
	}

	return &Files{
		Files: w.Files,
	}, nil
}

// StartLexing ranges over files to tokenization
func (f *Files) StartLexing() error {
	for _, file := range f.Files {
		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		l := New(data)
		l.Tokenize()
	}
	return nil
}

// New is used to issue an input to later be used by Tokenize
func New(input []byte) *Lexer {
	return &Lexer{
		input:  input,
		line:   1,
		column: 1,
		size:   len(input),
	}
}

// next appends the new data to the current token list
func (l *Lexer) newToken(kind token.Kind, data []byte, line, column int) {
	l.Tokens = append(l.Tokens, token.Token{
		Kind:   kind,
		Value:  string(data),
		Line:   line,
		Column: column,
	})
}

func (l *Lexer) advance(pos int, newLine bool) {
	l.position += pos
	if newLine {
		l.column = 1
		l.line++
		return
	}
	l.column += pos
}

func (l *Lexer) Tokenize() {
	for l.position < l.size {
		var tok []byte
		v := l.input[l.position]
		// fmt.Printf("INPUT '%s' line %d position %d column %d\n", string(v), l.line, l.position, l.column)

		switch {
		case isWhitespace(v):
			l.skipWhitespace()

		case v == '=':
			line, column := l.line, l.column
			if l.compareNextToken('=') {
				tok = append(tok, v, v)
				l.newToken(token.Eq, tok, line, column)
				l.advance(2, false)
			} else {
				tok = append(tok, v)
				l.newToken(token.Assign, tok, line, column)
				l.advance(1, false)
			}

		case v == ':':
			line, column := l.line, l.column
			if l.compareNextToken('=') {
				tok = append(tok, v, '=')
				l.newToken(token.Define, tok, line, column)
				l.advance(2, false)
			} else {
				tok = append(tok, v)
				l.newToken(token.Colon, tok, line, column)
				l.advance(1, false)
			}

		case v == '/':
			line, column := l.line, l.column
			if l.compareNextToken('/') {
				l.singleLineComment()
			} else if l.compareNextToken('*') {
				l.multiLineComment()
			} else if l.compareNextToken('=') {
				tok = append(tok, v, '=')
				l.newToken(token.SlashEq, tok, line, column)
				l.advance(2, false)
			} else {
				tok = append(tok, v)
				l.newToken(token.Slash, tok, line, column)
				l.advance(1, false)
			}

		case v == '+':
			line, column := l.line, l.column
			if l.compareNextToken('=') {
				tok = append(tok, v, '=')
				l.newToken(token.PlusEq, tok, line, column)
				l.advance(2, false)
			} else if l.compareNextToken('+') {
				tok = append(tok, v, '+')
				l.newToken(token.PPlus, tok, line, column)
				l.advance(2, false)
			} else {
				tok = append(tok, v)
				l.newToken(token.Plus, tok, line, column)
				l.advance(1, false)
			}

		case v == '-':
			line, column := l.line, l.column
			if l.compareNextToken('=') {
				tok = append(tok, v, '=')
				l.newToken(token.MinusEq, tok, line, column)
				l.advance(2, false)
			} else if l.compareNextToken('-') {
				tok = append(tok, v, '-')
				l.newToken(token.MMinus, tok, line, column)
				l.advance(2, false)
			} else {
				tok = append(tok, v)
				l.newToken(token.Minus, tok, line, column)
				l.advance(1, false)
			}

		case v == '*':
			line, column := l.line, l.column
			if l.compareNextToken('=') {
				tok = append(tok, v, '=')
				l.newToken(token.StarEq, tok, line, column)
				l.advance(2, false)
			} else {
				tok = append(tok, v)
				l.newToken(token.Star, tok, line, column)
				l.advance(1, false)
			}

		case v == '%':
			line, column := l.line, l.column
			tok = append(tok, v)
			l.newToken(token.Modulo, tok, line, column)
			l.advance(1, false)

		case v == '!':
			line, column := l.line, l.column
			if l.compareNextToken('=') {
				tok = append(tok, v, '=')
				l.newToken(token.Neq, tok, line, column)
				l.advance(2, false)
			} else {
				tok = append(tok, v)
				l.newToken(token.Not, tok, line, column)
				l.advance(1, false)
			}

		case v == '|':
			line, column := l.line, l.column
			if l.compareNextToken('|') {
				tok = append(tok, v, '|')
				l.newToken(token.Or, tok, line, column)
				l.advance(2, false)
			} else {
				tok = append(tok, v)
				l.newToken(token.Illegal, tok, line, column)
				l.advance(1, false)
			}

		case v == '<':
			line, column := l.line, l.column
			if l.compareNextToken('=') {
				tok = append(tok, v, '=')
				l.newToken(token.Lte, tok, line, column)
				l.advance(2, false)
			} else {
				tok = append(tok, v)
				l.newToken(token.Lt, tok, line, column)
				l.advance(1, false)
			}

		case v == '>':
			line, column := l.line, l.column
			if l.compareNextToken('=') {
				tok = append(tok, v, '=')
				l.newToken(token.Gte, tok, line, column)
				l.advance(2, false)
			} else {
				tok = append(tok, v)
				l.newToken(token.Gt, tok, line, column)
				l.advance(1, false)
			}

		case v == '&':
			line, column := l.line, l.column
			if l.compareNextToken('&') {
				tok = append(tok, v, '&')
				l.newToken(token.And, tok, line, column)
				l.advance(2, false)
			} else {
				tok = append(tok, v)
				l.newToken(token.Illegal, tok, line, column)
				l.advance(1, false)
			}

		// must be check, handled twice
		case v == '.':
			if ch, ok := l.fetchNextToken(); ok && isDigit(ch) {
				l.number()
			} else {
				line, column := l.line, l.column
				tok = append(tok, v)
				l.newToken(token.Dot, tok, line, column)
				l.advance(1, false)
			}

		case v == '_':
			if ch, ok := l.fetchNextToken(); ok && (ch == '.' || isDigit(ch)) {
				l.number()
			} else if ch, ok := l.fetchNextToken(); ok && isLetter(ch) {
				l.identOrKeyword()
			} else {
				line, column := l.line, l.column
				tok = append(tok, v)
				l.newToken(token.Illegal, tok, line, column)
				l.advance(1, false)
			}

		case v == ',':
			line, column := l.line, l.column
			tok = append(tok, v)
			l.newToken(token.Comma, tok, line, column)
			l.advance(1, false)

		case v == ';':
			line, column := l.line, l.column
			tok = append(tok, v)
			l.newToken(token.SemiComma, tok, line, column)
			l.advance(1, false)

		case v == '(':
			line, column := l.line, l.column
			tok = append(tok, v)
			l.newToken(token.LParen, tok, line, column)
			l.advance(1, false)

		case v == ')':
			line, column := l.line, l.column
			tok = append(tok, v)
			l.newToken(token.RParen, tok, line, column)
			l.advance(1, false)

		case v == '[':
			line, column := l.line, l.column
			tok = append(tok, v)
			l.newToken(token.LBracket, tok, line, column)
			l.advance(1, false)

		case v == ']':
			line, column := l.line, l.column
			tok = append(tok, v)
			l.newToken(token.RBracket, tok, line, column)
			l.advance(1, false)

		case v == '{':
			line, column := l.line, l.column
			tok = append(tok, v)
			l.newToken(token.LBrace, tok, line, column)
			l.advance(1, false)

		case v == '}':
			line, column := l.line, l.column
			tok = append(tok, v)
			l.newToken(token.RBrace, tok, line, column)
			l.advance(1, false)

		case v == '"':
			l.stringLit()

		case isLetter(v):
			l.identOrKeyword()

		case isDigit(v):
			l.number()

		default:
			line, column := l.line, l.column
			tok = append(tok, v)
			l.newToken(token.Illegal, tok, line, column)
			l.advance(1, false)
		}
	}
	l.newToken(token.EOF, nil, l.line, l.column)

	for _, v := range l.Tokens {
		fmt.Printf("Kind %d value %s line %d column %d\n", v.Kind, v.Value, v.Line, v.Column)
	}
}

// skipWhitespace skips any white space characters
func (l *Lexer) skipWhitespace() {
	for _, v := range l.input[l.position:] {
		if !isWhitespace(v) {
			break
		}
		if v == '\n' {
			l.advance(1, true)
		} else {
			l.advance(1, false)
		}
	}
}

// compareNextToken compares if next token match the provided one
func (l *Lexer) compareNextToken(ch byte) bool {
	if l.position+1 < l.size && l.input[l.position+1] == ch {
		return true
	}
	return false
}

// fetchNextToken returns the next token.
// bool is the to true when a token is find
func (l *Lexer) fetchNextToken() (byte, bool) {
	if l.position+1 < l.size {
		return l.input[l.position+1], true
	}
	return 0, false
}

// identOrKeyword parses the token and appends token list
func (l *Lexer) identOrKeyword() {
	var tok []byte
	line, column := l.line, l.column
	for _, v := range l.input[l.position:] {
		if !isLetter(v) && !isDigit(v) && v != '_' {
			break
		}
		tok = append(tok, v)
	}
	l.advance(len(tok), false)
	if tok[0] == '_' {
		l.newToken(token.Illegal, tok, line, column)
		return
	}
	l.newToken(token.LookupKeyword(string(tok)), tok, line, column)
}

// number parses the token and appends token list
func (l *Lexer) number() {
	var (
		tok            []byte
		dot, undescore int
		prev           byte
		illegal        bool
	)
	line, column := l.line, l.column
	for i, v := range l.input[l.position:] {
		if v == '.' || v == '_' {
			if i > 0 && ((prev == '.' && v == '_') || (prev == '_' && v == '.') || (prev == '_' && v == '_')) {
				illegal = true
			}
			if v == '.' {
				dot++
			}
			if v == '_' {
				undescore++
			}
		} else if !isDigit(v) {
			break
		}
		prev = v
		tok = append(tok, v)
	}

	l.advance(len(tok), false)
	if len(tok) > 1 {
		if illegal {
			l.newToken(token.Illegal, tok, line, column)
			return
		}

		last := tok[len(tok)-1]
		switch {
		case dot == 1:
			if last == '.' || last == '_' {
				l.newToken(token.Illegal, tok, line, column)
				return
			}
			l.newToken(token.FloatLit, tok, line, column)

		case dot > 1:
			l.newToken(token.Illegal, tok, line, column)

		case undescore > 0:
			if tok[0] == '_' || last == '_' {
				l.newToken(token.Illegal, tok, line, column)
				return
			}
			l.newToken(token.IntLit, tok, line, column)

		default:
			l.newToken(token.IntLit, tok, line, column)
		}
		return
	}
	l.newToken(token.IntLit, tok, line, column)
}

// stringLit parses the token and appends token list
func (l *Lexer) stringLit() {
	var (
		tok   []byte
		count int
		prev  byte
	)
	line, column := l.line, l.column
	for _, v := range l.input[l.position:] {
		tok = append(tok, v)
		if prev != '\\' && v == '"' {
			count++
			if count == 2 {
				break
			}
		}
		prev = v
	}

	l.advance(len(tok), false)
	if count == 2 {
		l.newToken(token.StringLit, tok, line, column)
		return
	}
	l.newToken(token.Illegal, tok, line, column)
}

// singleLineComment parses single line comment and appends token list
func (l *Lexer) singleLineComment() {
	var tok []byte
	line, column := l.line, l.column
	for _, v := range l.input[l.position:] {
		if v == '\n' {
			break
		}
		tok = append(tok, v)
	}

	l.advance(len(tok), false)
	l.newToken(token.Comment, tok, line, column)
}

// multiLineComment parses multi line comments like /* */ and appends token list
func (l *Lexer) multiLineComment() {
	var (
		tok []byte
		ok  bool
	)
	line, column := l.line, l.column
	compareNextToken := func(ch byte, pos int) bool {
		if pos+1 < l.size && l.input[pos+1] == ch {
			return true
		}
		return false
	}

	pos := l.position
	for _, v := range l.input[l.position:] {
		if v == '*' && compareNextToken('/', pos) {
			tok = append(tok, v, '/')
			ok = true
			l.advance(2, false)
			break
		}
		if v == '\n' {
			l.advance(1, true)
		} else {
			l.advance(1, false)
		}
		tok = append(tok, v)
		pos++
	}
	if ok {
		l.newToken(token.Comment, tok, line, column)
		return
	}
	l.newToken(token.Illegal, tok, line, column)
}

// isLetter returns wether we found a letter or not.
// Understore is included.
func isLetter(ch byte) bool {
	if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' {
		return true
	}
	return false
}

// isDigit returns wether we found a digit or not
func isDigit(ch byte) bool {
	if ch >= '0' && ch <= '9' {
		return true
	}
	return false
}

// isWhitespace returns wether we found a space character or not
func isWhitespace(ch byte) bool {
	if ch == ' ' || ch == '\n' || ch == '\r' || ch == '\t' {
		return true
	}
	return false
}
