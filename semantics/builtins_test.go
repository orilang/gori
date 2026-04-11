package semantics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSemantics_builtins(t *testing.T) {
	tests := []struct {
		bt       *BuiltinType
		expected string
	}{
		{bt: TBool, expected: "bool"},
		{bt: TInt, expected: "int"},
		{bt: TInt8, expected: "int8"},
		{bt: TInt32, expected: "int32"},
		{bt: TInt64, expected: "int64"},
		{bt: TUInt, expected: "uint"},
		{bt: TUInt8, expected: "uint8"},
		{bt: TUInt32, expected: "uint32"},
		{bt: TUInt64, expected: "uint64"},
		{bt: TFloat32, expected: "float32"},
		{bt: TFloat64, expected: "float64"},
		{bt: nil, expected: "invalidBuiltin"},
	}

	for _, tc := range tests {
		require.Equal(t, tc.expected, tc.bt.String())
	}
}
