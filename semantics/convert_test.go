package semantics

import (
	"testing"

	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/require"
)

func TestSemantics_convert(t *testing.T) {
	t.Run("is_identical", func(t *testing.T) {
		tests := []struct {
			a, b     Type
			expected bool
		}{
			{a: TBool, b: nil, expected: false},
			{a: TBool, b: TInt, expected: false},
			{a: TInt, b: TInt, expected: true},
			{a: &NamedType{Name: "UserID", UnderlyingType: TInt}, b: TInt, expected: false},
			{a: &NamedType{Name: "UserID", UnderlyingType: TInt}, b: &NamedType{Name: "OrderID", UnderlyingType: TInt}, expected: false},
			{a: &NamedType{Name: "UserID", UnderlyingType: TInt}, b: &NamedType{Name: "UserID", UnderlyingType: TInt}, expected: true},
			{a: &ArrayType{Len: 1, Elem: TInt}, b: TInt, expected: false},
			{a: &ArrayType{Len: 1, Elem: TInt}, b: &ArrayType{Len: 1, Elem: TString}, expected: false},
			{a: &ArrayType{Len: 1, Elem: TInt}, b: &ArrayType{Len: 2, Elem: TString}, expected: false},
			{a: &ArrayType{Len: 1, Elem: TInt}, b: &ArrayType{Len: 1, Elem: TInt}, expected: true},
			{a: &SliceType{Elem: TInt}, b: TInt, expected: false},
			{a: &SliceType{Elem: TInt}, b: &SliceType{Elem: TString}, expected: false},
			{a: &SliceType{Elem: TInt}, b: &SliceType{Elem: TInt}, expected: true},
			{a: &MapType{Kind: MapHash, Key: TInt, Value: TString}, b: TInt, expected: false},
			{a: &MapType{Kind: MapHash, Key: TInt, Value: TString}, b: &MapType{Kind: MapHash, Key: TInt, Value: TInt}, expected: false},
			{a: &MapType{Kind: MapHash, Key: TInt, Value: TString}, b: &MapType{Kind: MapRegular, Key: TInt, Value: TString}, expected: false},
			{a: &MapType{Kind: MapHash, Key: TInt, Value: TString}, b: &MapType{Kind: MapHash, Key: TInt, Value: TString}, expected: true},
			{a: &StructType{Fields: []StructField{{Name: "Age", Type: TInt}}}, b: TBool, expected: false},
			{a: &StructType{Fields: []StructField{{Name: "Age", Type: TInt}}}, b: &StructType{Fields: []StructField{{Name: "Name", Type: TString}, {Name: "Age", Type: TInt}}}, expected: false},
			{a: &StructType{Fields: []StructField{{Name: "Age", Type: TInt}}}, b: &StructType{Fields: []StructField{{Name: "Age", Type: TString}}}, expected: false},
			{a: &StructType{Fields: []StructField{{Name: "Age", Type: TInt}}}, b: &StructType{Fields: []StructField{{Name: "Age", Type: TInt}}}, expected: true},
			{a: &FuncMethod{Name: "test", FuncType: &FuncType{Params: []Param{{Name: "a", Type: TInt}}}}, b: TInt, expected: false},
			{
				a: &FuncMethod{
					Name: "test", FuncType: &FuncType{
						Params: []Param{{Name: "a", Type: TInt}},
					},
				},
				b: &FuncMethod{
					Name: "test", FuncType: &FuncType{
						Params: []Param{{Name: "a", Type: TString}},
					},
				},
				expected: false,
			},
			{
				a: &FuncMethod{
					Name: "test", FuncType: &FuncType{
						Params:  []Param{{Name: "a", Type: TInt}},
						Results: []Param{{Name: "b", Type: TInt}},
					},
				},
				b: &FuncMethod{
					Name: "test", FuncType: &FuncType{
						Params:  []Param{{Name: "a", Type: TInt}},
						Results: []Param{{Name: "b", Type: TString}},
					},
				},
				expected: false,
			},
			{
				a: &FuncMethod{
					Name: "test", FuncType: &FuncType{
						Params:  []Param{{Name: "a", Type: TInt}},
						Results: []Param{{Name: "b", Type: TInt}},
					},
				},
				b: &FuncMethod{
					Name: "test", FuncType: &FuncType{
						Params:  []Param{{Name: "a", Type: TInt}},
						Results: []Param{{Name: "b", Type: TInt}},
					},
				},
				expected: true,
			},
			{
				a:        &InterfaceType{Methods: []FuncMethod{{Name: "test", FuncType: &FuncType{Params: []Param{{Name: "a", Type: TInt}}}}}},
				b:        &InterfaceType{Methods: []FuncMethod{{Name: "test", FuncType: &FuncType{Params: []Param{{Name: "a", Type: TString}}}}}},
				expected: false,
			},
			{
				a:        &InterfaceType{Methods: []FuncMethod{{Name: "test", FuncType: &FuncType{Params: []Param{{Name: "a", Type: TInt}}}}}},
				b:        &InterfaceType{Methods: []FuncMethod{{Name: "test", FuncType: &FuncType{Params: []Param{{Name: "a", Type: TInt}}}}}},
				expected: true,
			},
			{a: &Enum{Name: "Color", Variants: []string{"Red", "Blue", "Green"}}, b: TBool, expected: false},
			{a: &Enum{Name: "Color", Variants: []string{"Red", "Blue", "Green"}}, b: &Enum{Name: "Color", Variants: []string{"red", "Blue", "Green"}}, expected: false},
			{a: &Enum{Name: "Color", Variants: []string{"Red", "Blue", "Green"}}, b: &Enum{Name: "Color", Variants: []string{"Red", "Blue", "Green"}}, expected: true},
			{
				a:        &SumType{Name: "Shape", Variants: []SumVariant{{Name: "Circle", Field: []Param{{Name: "radius", Type: TFloat64}}}}},
				b:        TBool,
				expected: false,
			},
			{
				a:        &SumType{Name: "Shape", Variants: []SumVariant{{Name: "Circle", Field: []Param{{Name: "radius", Type: TFloat64}}}}},
				b:        &SumType{Name: "Shape", Variants: []SumVariant{{Name: "Circle", Field: []Param{{Name: "radius", Type: TFloat32}}}}},
				expected: false,
			},
			{
				a:        &SumType{Name: "Shape", Variants: []SumVariant{{Name: "Circle", Field: []Param{{Name: "radius", Type: TFloat64}}}}},
				b:        &SumType{Name: "Shape", Variants: []SumVariant{{Name: "Circle", Field: []Param{{Name: "radius", Type: TFloat64}, {Name: "perimeter", Type: TFloat64}}}}},
				expected: false,
			},
			{
				a:        &SumType{Name: "Shape", Variants: []SumVariant{{Name: "Circle", Field: []Param{{Name: "radius", Type: TFloat64}}}}},
				b:        &SumType{Name: "Shape", Variants: []SumVariant{{Name: "Circle", Field: []Param{{Name: "radius", Type: TFloat64}}}}},
				expected: true,
			},
		}

		for _, tc := range tests {
			require.Equal(t, tc.expected, IsIdentical(tc.a, tc.b))
		}
	})

	t.Run("is_assignable_to", func(t *testing.T) {
		tests := []struct {
			targetType, valueType Type
			expected              bool
		}{
			{targetType: TBool, valueType: &UntypeNilType{}, expected: false},
			{targetType: TBool, valueType: TBool, expected: true},
			{targetType: &ArrayType{Len: 1, Elem: TInt}, valueType: TBool, expected: false},
			{targetType: &ArrayType{Len: 1, Elem: TInt}, valueType: nil, expected: false},
			{targetType: &SliceType{Elem: TInt}, valueType: TBool, expected: false},
			{targetType: &SliceType{Elem: TInt}, valueType: &UntypeNilType{}, expected: true},
			{targetType: &MapType{Kind: MapHash, Key: TInt, Value: TString}, valueType: TBool, expected: false},
			{targetType: &MapType{Kind: MapHash, Key: TInt, Value: TString}, valueType: &UntypeNilType{}, expected: true},
			{targetType: &FuncMethod{Name: "test", FuncType: &FuncType{Params: []Param{{Name: "a", Type: TInt}}}}, valueType: TBool, expected: false},
			{targetType: &FuncMethod{Name: "test", FuncType: &FuncType{Params: []Param{{Name: "a", Type: TInt}}}}, valueType: &UntypeNilType{}, expected: false},
			{targetType: &InterfaceType{Methods: []FuncMethod{{Name: "test", FuncType: &FuncType{Params: []Param{{Name: "a", Type: TInt}}}}}}, valueType: TBool, expected: false},
			{targetType: &InterfaceType{Methods: []FuncMethod{{Name: "test", FuncType: &FuncType{Params: []Param{{Name: "a", Type: TInt}}}}}}, valueType: &UntypeNilType{}, expected: false},
			{targetType: &StructType{Fields: []StructField{{Name: "Age", Type: TInt}}}, valueType: &UntypeNilType{}, expected: false},
			{targetType: &StructType{Fields: []StructField{{Name: "Age", Type: TInt}}}, valueType: &StructType{Fields: []StructField{{Name: "Age", Type: TInt}}}, expected: true},
			{targetType: &Enum{Name: "Color", Variants: []string{"Red", "Blue", "Green"}}, valueType: nil, expected: false},
			{targetType: &Enum{Name: "Color", Variants: []string{"Red", "Blue", "Green"}}, valueType: &Enum{Name: "Color", Variants: []string{"Red", "Blue", "Green"}}, expected: true},
			{targetType: &InvalidType{}, valueType: &InvalidType{}, expected: true},
		}

		for _, tc := range tests {
			require.Equal(t, tc.expected, IsAssignableTo(tc.targetType, tc.valueType))
		}
	})

	t.Run("is_numeric", func(t *testing.T) {
		tests := []struct {
			src      Type
			expected bool
		}{
			{src: TBool, expected: false},
			{src: TInt, expected: true},
			{src: &ArrayType{Len: 1, Elem: TInt}, expected: false},
			{src: &NamedType{Name: "UserID", UnderlyingType: TInt}, expected: true},
		}

		for _, tc := range tests {
			require.Equal(t, tc.expected, IsNumeric(tc.src))
		}
	})

	t.Run("is_integer", func(t *testing.T) {
		tests := []struct {
			src, dst Type
			expected bool
		}{
			{src: TBool, expected: false},
			{src: &ArrayType{Len: 1, Elem: TInt}, expected: false},
			{src: TInt, expected: true},
			{src: &NamedType{Name: "UserID", UnderlyingType: TInt}, expected: true},
		}

		for _, tc := range tests {
			require.Equal(t, tc.expected, IsInteger(tc.src))
		}
	})

	t.Run("is_bool", func(t *testing.T) {
		tests := []struct {
			src, dst Type
			expected bool
		}{
			{src: TBool, expected: true},
			{src: TInt, expected: false},
			{src: &ArrayType{Len: 1, Elem: TInt}, expected: false},
			{src: &NamedType{Name: "old", UnderlyingType: TBool}, expected: true},
		}

		for _, tc := range tests {
			require.Equal(t, tc.expected, IsBool(tc.src))
		}
	})

	t.Run("is_string", func(t *testing.T) {
		tests := []struct {
			src, dst Type
			expected bool
		}{
			{src: TString, expected: true},
			{src: TBool, expected: false},
			{src: TInt, expected: false},
			{src: &ArrayType{Len: 1, Elem: TInt}, expected: false},
			{src: &NamedType{Name: "Name", UnderlyingType: TString}, expected: true},
		}

		for _, tc := range tests {
			require.Equal(t, tc.expected, IsString(tc.src))
		}
	})

	t.Run("is_invalid", func(t *testing.T) {
		tests := []struct {
			src, dst Type
			expected bool
		}{
			{src: &InvalidType{}, expected: true},
			{src: TBool, expected: false},
			{src: TInt, expected: false},
			{src: &ArrayType{Len: 1, Elem: TInt}, expected: false},
		}

		for _, tc := range tests {
			require.Equal(t, tc.expected, IsInvalid(tc.src))
		}
	})

	t.Run("is_convertible_to", func(t *testing.T) {
		tests := []struct {
			src, dst Type
			expected bool
		}{
			{src: &InvalidType{}, dst: nil, expected: true},
			{src: TBool, dst: TInt, expected: false},
			{src: TInt, expected: false},
			{src: TInt, dst: TFloat, expected: true},
			{src: TInt, dst: TInt, expected: true},
			{src: TString, dst: TString, expected: true},
			{src: &NamedType{Name: "Age", UnderlyingType: TInt}, expected: false},
		}

		for _, tc := range tests {
			require.Equal(t, tc.expected, IsConvertibleTo(tc.src, tc.dst))
		}
	})

	t.Run("supports_binary_op", func(t *testing.T) {
		tests := []struct {
			src      Type
			op       token.Kind
			expected bool
		}{
			{src: TBool, op: token.And, expected: true},
			{src: &NamedType{Name: "IsValid", UnderlyingType: TBool}, op: token.And, expected: true},
			{src: TInt, expected: false},
			{src: TInt, op: token.Plus, expected: true},
			{src: TInt, op: token.Modulo, expected: true},
			{src: TInt, op: token.Eq, expected: true},
			{src: TInt, op: token.Lt, expected: true},
			{src: TInt, op: token.Modulo, expected: true},
			{src: TFloat, op: token.Modulo, expected: false},
		}

		for _, tc := range tests {
			require.Equal(t, tc.expected, SupportsBinaryOp(tc.src, tc.op))
		}
	})

	t.Run("supports_unary_op", func(t *testing.T) {
		tests := []struct {
			src      Type
			op       token.Kind
			expected bool
		}{
			{src: TInt, op: token.Plus, expected: true},
			{src: TBool, op: token.Not, expected: true},
			{src: &NamedType{Name: "UserID", UnderlyingType: TInt}, op: token.Plus, expected: true},
			{src: &NamedType{Name: "IsValid", UnderlyingType: TBool}, op: token.Not, expected: true},
			{src: TString, expected: false},
			{src: TString, op: token.PPlus, expected: false},
		}

		for _, tc := range tests {
			require.Equal(t, tc.expected, SupportsUnaryOp(tc.src, tc.op))
		}
	})

	t.Run("is_comparable", func(t *testing.T) {
		tests := []struct {
			src      Type
			expected bool
		}{
			{src: TInt, expected: true},
			{src: TBool, expected: true},
			{src: &NamedType{Name: "UserID", UnderlyingType: TInt}, expected: true},
			{src: &NamedType{Name: "IsValid", UnderlyingType: TBool}, expected: true},
			{src: TString, expected: true},
		}

		for _, tc := range tests {
			require.Equal(t, tc.expected, IsComparable(tc.src))
		}
	})

	t.Run("is_ordered", func(t *testing.T) {
		tests := []struct {
			src      Type
			expected bool
		}{
			{src: TInt, expected: true},
			{src: TBool, expected: false},
			{src: &NamedType{Name: "UserID", UnderlyingType: TInt}, expected: true},
			{src: &NamedType{Name: "IsValid", UnderlyingType: TBool}, expected: false},
			{src: TString, expected: true},
		}

		for _, tc := range tests {
			require.Equal(t, tc.expected, IsOrdered(tc.src))
		}
	})

	t.Run("is_untyped_nil_type", func(t *testing.T) {
		tests := []struct {
			src      Type
			expected bool
		}{
			{src: &UntypeNilType{}, expected: true},
			{src: TBool, expected: false},
		}

		for _, tc := range tests {
			require.Equal(t, tc.expected, IsUntypedNilType(tc.src))
		}
	})

	t.Run("is_nil_assignable", func(t *testing.T) {
		tests := []struct {
			src      Type
			expected bool
		}{
			{src: &SliceType{}, expected: true},
			{src: &MapType{}, expected: true},
			{src: TBool, expected: false},
		}

		for _, tc := range tests {
			require.Equal(t, tc.expected, IsNilAssignable(tc.src))
		}
	})
}
