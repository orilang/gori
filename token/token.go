package token

func LookupKeyword(s string) Kind {
	if token, ok := keywords[s]; ok {
		return token
	}
	return Ident
}

func IsBuiltinType(k Kind) bool {
	if builtinTypes[k] {
		return true
	}
	return false
}
