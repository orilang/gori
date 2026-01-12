package token

// LookupKeyword lookup for the current keyword list
// and returns its kind if found otherwise identifier
func LookupKeyword(s string) Kind {
	if token, ok := keywords[s]; ok {
		return token
	}
	return Ident
}

// IsBuiltinType returns true builtin types is found otherwise false
func IsBuiltinType(k Kind) bool {
	return builtinTypes[k]
}

func IsPrefix(k Kind) bool {
	return prefix[k]
}

func IsInfix(k Kind) bool {
	return infix[k]
}
