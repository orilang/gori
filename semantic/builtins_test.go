package semantic

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSemantics_builtins(t *testing.T) {
	t.Run("builtins", func(t *testing.T) {
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
			{bt: TFloat, expected: "float"},
			{bt: TFloat32, expected: "float32"},
			{bt: TFloat64, expected: "float64"},
			{bt: nil, expected: "invalidBuiltin"},
		}

		for _, tc := range tests {
			require.Equal(t, tc.expected, tc.bt.String())
		}
	})

	t.Run("string", func(t *testing.T) {
		tests := []struct {
			tp       Type
			expected string
		}{
			{tp: &NamedType{Name: "UserID", UnderlyingType: TInt}, expected: "UserID"},
			{tp: &ArrayType{Len: 1, Elem: TInt}, expected: "[1]int"},
			{tp: &SliceType{Elem: TInt}, expected: "[]int"},
			{tp: &MapType{Kind: MapHash, Key: TInt, Value: TString}, expected: "hashmap[int]string"},
			{tp: &MapType{Kind: MapRegular, Key: TInt, Value: TString}, expected: "map[int]string"},
			{tp: &StructType{Fields: []StructField{{Name: "Age", Type: TInt}}}, expected: "struct"},
			{tp: &Param{Name: "Age", Type: TInt}, expected: "param"},
			{tp: &FuncMethod{Name: "test", FuncType: &FuncType{Params: []Param{{Name: "a", Type: TInt}}}}, expected: "funcMethod"},
			{tp: &InterfaceType{Methods: []FuncMethod{{Name: "test", FuncType: &FuncType{Params: []Param{{Name: "a", Type: TInt}}}}}}, expected: "interface"},
			{tp: &EnumType{Name: "Color", Variants: []string{"Red", "Blue", "Green"}}, expected: "enum"},
			{tp: &SumType{Name: "Shape", Variants: []SumVariant{{Name: "Circle", Field: []Param{{Name: "radius", Type: TFloat64}}}}}, expected: "sum"},
			{tp: &InvalidType{}, expected: "<invalid>"},
			{tp: &UntypedNilType{}, expected: "nil"},
			{tp: &FuncType{}, expected: "funcType"},
		}

		for _, tc := range tests {
			require.Equal(t, tc.expected, tc.tp.String())
		}
	})
}
