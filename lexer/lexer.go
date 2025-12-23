package lexer

import (
	"os"
	"strings"

	"github.com/orilang/gori/token"
)

// NewLexer returns files config to StartLexing
func NewLexer(config Config) (*Files, error) {
	return walk(config)
}

// StartLexing ranges over files to tokenization
func (f *Files) StartLexing() error {
	for _, file := range f.Files {
		f, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		l := New(string(f))
		l.Tokenize()
	}
	return nil
}

// New is used to issue an input to later be used by Tokenize
func New(input string) *Lexer {
	return &Lexer{
		input: input,
	}
}

// Tokenize split every characters to recognize allowed syntaxes
func (l *Lexer) Tokenize() {
	lines := strings.SplitSeq(l.input, "\n")
	var xline uint
	for xdata := range lines {
		xline++
		l.position = 0
		inputSize := len(xdata)
		for index, data := range xdata {
			if index > 0 && index < l.position {
				continue
			}

			dtoken := string(data)
			switch data {
			case ' ', '\n', '\t', '\r':
				l.position++
				continue

			case '{':
				l.Tokens = append(l.Tokens, token.Token{Kind: token.LBrace, Value: dtoken, Line: xline})
				l.position++
				continue

			case '}':
				l.Tokens = append(l.Tokens, token.Token{Kind: token.RBrace, Value: dtoken, Line: xline})
				l.position++
				continue

			// case '[':
			// 	l.Tokens = append(l.Tokens, token.Token{Kind: token.LBracket, Value: dtoken, Line: xline})
			// 	continue

			// case ']':
			// 	l.Tokens = append(l.Tokens, token.Token{Kind: token.RBracket, Value: dtoken, Line: xline})
			// 	continue

			case '(':
				l.Tokens = append(l.Tokens, token.Token{Kind: token.LParen, Value: dtoken, Line: xline})
				l.position++
				continue

			case ')':
				l.Tokens = append(l.Tokens, token.Token{Kind: token.RParen, Value: dtoken, Line: xline})
				l.position++
				continue

			case '+':
				if index+1 < inputSize && xdata[index+1] == '=' {
					l.Tokens = append(l.Tokens, token.Token{
						Kind:  token.PlusEq,
						Value: xdata[index : index+2],
						EOL:   true,
						Line:  xline,
					})
					l.position += 2
					continue
				}

				// if index+1 < inputSize && xdata[index+1] == '+' {
				// 	l.Tokens = append(l.Tokens, token.Token{
				// 		Kind:  token.PPlus,
				// 		Value: xdata[index : index+2],
				// 		EOL:   true,
				// 		Line:  xline,
				// 	})
				// 	l.position += 2
				// 	continue
				// }

				l.Tokens = append(l.Tokens, token.Token{Kind: token.Plus, Value: dtoken, Line: xline})
				l.position++
				continue

			case '-':
				if index+1 < inputSize && xdata[index+1] == '=' {
					l.Tokens = append(l.Tokens, token.Token{
						Kind:  token.MinusEq,
						Value: xdata[index : index+2],
						EOL:   true,
						Line:  xline,
					})
					l.position += 2
					continue
				}

				// if index+1 < inputSize && xdata[index+1] == '-' {
				// 	l.Tokens = append(l.Tokens, token.Token{
				// 		Kind:  token.MMinus,
				// 		Value: xdata[index : index+2],
				// 		EOL:   true,
				// 		Line:  xline,
				// 	})
				// 	l.position += 2
				// 	continue
				// }

				l.Tokens = append(l.Tokens, token.Token{Kind: token.Minus, Value: dtoken, Line: xline})
				l.position++
				continue

			case '*':
				l.Tokens = append(l.Tokens, token.Token{Kind: token.Star, Value: dtoken, Line: xline})
				continue

			case '/':
				if index+1 < inputSize && xdata[index+1] == '/' {
					l.Tokens = append(l.Tokens, token.Token{
						Kind:  token.Comment,
						Value: xdata[index:],
						EOL:   true,
						Line:  xline,
					})
					l.position += len(xdata[index:])
					continue
				}

				l.Tokens = append(l.Tokens, token.Token{Kind: token.Slash, Value: dtoken, Line: xline})
				l.position++
				continue

			case '%':
				l.Tokens = append(l.Tokens, token.Token{Kind: token.Modulo, Value: dtoken, Line: xline})
				l.position++
				continue

			case '!':
				if index+1 < inputSize && xdata[index+1] == '=' {
					l.Tokens = append(l.Tokens, token.Token{
						Kind:  token.Neq,
						Value: xdata[index : index+2],
						EOL:   true,
						Line:  xline,
					})
					l.position += 2
					continue
				}

				// l.Tokens = append(l.Tokens, token.Token{Kind: token.Not, Value: dtoken, Line: xline})
				// l.position++
				// continue

			case '=':
				if index+1 < inputSize && xdata[index+1] == '=' {
					l.Tokens = append(l.Tokens, token.Token{
						Kind:  token.Eq,
						Value: xdata[index : index+2],
						EOL:   true,
						Line:  xline,
					})
					l.position += 2
					continue
				}

				l.Tokens = append(l.Tokens, token.Token{Kind: token.Assign, Value: dtoken, Line: xline})
				l.position++
				continue

			case '>':
				if index+1 < inputSize && xdata[index+1] == '=' {
					l.Tokens = append(l.Tokens, token.Token{
						Kind:  token.Gte,
						Value: xdata[index : index+2],
						EOL:   true,
						Line:  xline,
					})
					l.position += 2
					continue
				}

				l.Tokens = append(l.Tokens, token.Token{Kind: token.Gt, Value: dtoken, Line: xline})
				l.position++
				continue

			case '<':
				if index+1 < inputSize && xdata[index+1] == '=' {
					l.Tokens = append(l.Tokens, token.Token{
						Kind:  token.Lte,
						Value: xdata[index : index+2],
						EOL:   true,
						Line:  xline,
					})
					l.position += 2
					continue
				}

				l.Tokens = append(l.Tokens, token.Token{Kind: token.Lt, Value: dtoken, Line: xline})
				l.position++
				continue

			// case '|':
			// 	if index+1 < inputSize && xdata[index+1] == '|' {
			// 		l.Tokens = append(l.Tokens, token.Token{
			// 			Kind:  token.Or,
			// 			Value: xdata[index : index+2],
			// 			EOL:   true,
			// 			Line:  xline,
			// 		})
			// 		l.position += 2
			// 		continue
			// 	}

			// 	l.Tokens = append(l.Tokens, token.Token{Kind: token.Illegal, Value: dtoken, Line: xline})
			// l.position++
			// 	continue

			// case '&':
			// 	if index+1 < inputSize && xdata[index+1] == '&' {
			// 		l.Tokens = append(l.Tokens, token.Token{
			// 			Kind:  token.And,
			// 			Value: xdata[index : index+2],
			// 			EOL:   true,
			// 			Line:  xline,
			// 		})
			// 		l.position += 2
			// 		continue
			// 	}

			// 	l.Tokens = append(l.Tokens, token.Token{Kind: token.Illegal, Value: dtoken, Line: xline})
			// l.position++
			// 	continue

			// case '.':
			// 	l.Tokens = append(l.Tokens, token.Token{Kind: token.Dot, Value: dtoken, Line: xline})
			// 	l.position++
			// 	continue

			// case ',':
			// 	l.Tokens = append(l.Tokens, token.Token{Kind: token.Comma, Value: dtoken, Line: xline})
			// l.position++
			// 	continue

			// case ';':
			// 	l.Tokens = append(l.Tokens, token.Token{Kind: token.SemiComma, Value: dtoken, Line: xline})
			// l.position++
			// 	continue

			case '"':
				var (
					value []rune
					count int
				)
				datax := xdata[index:]
				sizex := len(datax)
				for xi, v := range datax {
					if v == '"' {
						count++
					}

					value = append(value, v)
					if xi+1 == sizex {
						l.position += len(value)
					}

					if count == 2 {
						break
					}
				}
				result := string(value)
				if count == 2 {
					l.Tokens = append(l.Tokens, token.Token{Kind: token.StringLit, Value: result, Line: xline})
				} else {
					l.Tokens = append(l.Tokens, token.Token{Kind: token.Illegal, Value: result, Line: xline})
				}
				continue

			default:
				// if isLetterDigits(data) {
				// 	var value []rune
				// 	datax := xdata[index:]
				// 	sizex := len(datax)
				// 	for xi, v := range datax {
				// 		if !isLetterDigits(v) {
				// 			l.position += len(value)
				// 			break
				// 		}

				// 		if xi+1 == sizex {
				// 			l.position += sizex
				// 		}
				// 		value = append(value, v)
				// 	}
				// 	result := string(value)
				// 	l.Tokens = append(l.Tokens, token.Token{Kind: token.LookupKeyword(result), Value: result, Line: xline})
				// 	continue
				// }

				if isLetter(data) {
					var value []rune
					datax := xdata[index:]
					sizex := len(datax)
					for xi, v := range datax {
						if !isDigit(v) && !isLetter(v) {
							l.position += len(value)
							break
						}

						value = append(value, v)
						if xi+1 == sizex {
							l.position += len(value)
						}
					}
					result := string(value)
					l.Tokens = append(l.Tokens, token.Token{Kind: token.LookupKeyword(result), Value: result, Line: xline})
					continue
				}

				if isDigit(data) {
					var value []rune
					datax := xdata[index:]
					sizex := len(datax)
					countDot := 0
					for xi, v := range datax {
						if v == '.' || v == '_' {
							if v == '.' {
								countDot++
							}
						} else if !isDigit(v) {
							l.position += len(value)
							break
						}

						value = append(value, v)
						if xi+1 == sizex {
							l.position += len(value)
						}
					}
					result := string(value)

					if strings.Contains(result, ".") {
						if countDot > 1 {
							l.Tokens = append(l.Tokens, token.Token{Kind: token.Illegal, Value: result, Line: xline})
						} else {
							if strings.HasPrefix(result, ".") || strings.HasSuffix(result, ".") || strings.HasSuffix(result, "_") {
								l.Tokens = append(l.Tokens, token.Token{Kind: token.Illegal, Value: result, Line: xline})
							} else {
								l.Tokens = append(l.Tokens, token.Token{Kind: token.FloatLit, Value: result, Line: xline})
							}
						}
					} else {
						l.Tokens = append(l.Tokens, token.Token{Kind: token.IntLit, Value: result, Line: xline})
					}
					continue
				}
			}
		}
	}
}

// isLetter returns wether we found a letter or not.
// Understore is included.
func isLetter(ch rune) bool {
	if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_' {
		return true
	}
	return false
}

// isDigit returns wether we found a digit or not
func isDigit(ch rune) bool {
	if ch >= '0' && ch <= '9' {
		return true
	}
	return false
}
