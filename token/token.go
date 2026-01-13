package token

// LookupKeyword lookup for the current keyword list
// and returns its kind if found otherwise identifier
func LookupKeyword(s string) Kind {
	if token, ok := keywords[s]; ok {
		return token
	}
	return Ident
}

// IsBuiltinType returns true when builtin types is found otherwise false
func IsBuiltinType(k Kind) bool {
	return builtinTypes[k]
}

// IsPrefix returns true when the provided kind is found is the list
func IsPrefix(k Kind) bool {
	return prefix[k]
}

// IsInfix returns true when the provided kind is found is the list
func IsInfix(k Kind) bool {
	return infix[k]
}

// IsOperator returns true when the provided kind is found is the list
func IsOperator(k Kind) bool {
	return operator[k]
}
