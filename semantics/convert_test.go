package semantics

import (
	"testing"

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
				b:        &SumType{Name: "Shape", Variants: []SumVariant{{Name: "Circle", Field: []Param{{Name: "radius", Type: TFloat64}}}}},
				expected: true,
			},
		}

		for _, tc := range tests {
			require.Equal(t, tc.expected, IsIdentical(tc.a, tc.b))
		}
	})
}
