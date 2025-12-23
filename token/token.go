package token

func LookupKeyword(s string) Kind {
	if token, ok := keywords[s]; ok {
		return token
	}
	return Ident
}
