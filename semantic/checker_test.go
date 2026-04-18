package semantic

import (
	"testing"

	"github.com/orilang/gori/ast"
	"github.com/orilang/gori/lexer"
	"github.com/orilang/gori/parser"
	"github.com/orilang/gori/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSemantic_checker(t *testing.T) {
	t.Run("x1", func(t *testing.T) {
		data :=
			`package main
type UserID int

type User struct {
	id UserID
}
`
		scope := &Scope{
			Parent: nil,
			Symbols: map[string]*Symbol{
				"UserID": {
					Name: "UserID",
					Kind: SymType,
					Type: &NamedType{
						UnderlyingType: TInt,
					},
				},
				"User": {
					Name: "User",
					Kind: SymType,
					Type: &StructType{
						Fields: []StructField{
							{
								Name: "id",
								Type: &NamedType{
									Name:           "UserID",
									UnderlyingType: TInt,
								},
							},
						},
					},
				},
			},
		}

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		// fmt.Printf("%s\n", ast.Dump(pr))
		check := NewChecker()

		assert.Equal(t, 0, len(check.Check(pr)))
		assert.Equal(t, scope.Parent, check.pkgScope.Parent)
		assert.Equal(t, scope.Symbols["UserID"].Name, check.pkgScope.Symbols["UserID"].Name)
		assert.Equal(t, scope.Symbols["UserID"].Kind, check.pkgScope.Symbols["UserID"].Kind)
		xx := check.pkgScope.Symbols["UserID"].Type.(*NamedType)
		assert.Equal(t, TInt, xx.UnderlyingType)

		assert.Equal(t, scope.Symbols["User"].Name, check.pkgScope.Symbols["User"].Name)
		assert.Equal(t, scope.Symbols["User"].Kind, check.pkgScope.Symbols["User"].Kind)
		st := check.pkgScope.Symbols["User"].Type.(*StructType)
		xy := st.Fields[0].Type.(*NamedType)
		assert.Equal(t, TInt, TInt, xy.UnderlyingType)
	})

	t.Run("x1_duplicate", func(t *testing.T) {
		data :=
			`package main
type UserID int
type UserID int

type User struct {
	id UserID
	id UserID
}
`
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		// fmt.Printf("%s\n", ast.Dump(pr))
		check := NewChecker()
		assert.Greater(t, len(check.Check(pr)), 0)
	})

	t.Run("x2", func(t *testing.T) {
		data :=
			`package main
type UserID int
type test interface{
	x(a string, b bool) (c float64, d float32, e uint, f uint32, g uint64 )
	y(a string, b bool) (int8, int32, int64)
	z(a string, b bool) (c uint8, d float, e UserID)
}
`
		scope := &Scope{
			Parent: nil,
			Symbols: map[string]*Symbol{
				"UserID": {
					Name: "UserID",
					Kind: SymType,
					Type: &NamedType{
						UnderlyingType: TInt,
					},
				},
				"test": {
					Name: "test",
					Kind: SymType,
					Type: &InterfaceType{
						Methods: []FuncMethod{
							{
								Name: "x",
								FuncType: &FuncType{
									Params: []Param{
										{Name: "a", Type: TString},
										{Name: "b", Type: TBool},
									},
									Results: []Param{
										{Name: "c", Type: TFloat64},
										{Name: "d", Type: TFloat32},
										{Name: "e", Type: TUInt},
										{Name: "f", Type: TUInt32},
										{Name: "g", Type: TUInt64},
									},
								},
							},
							{
								Name: "y",
								FuncType: &FuncType{
									Params: []Param{
										{Name: "a", Type: TString},
										{Name: "b", Type: TBool},
									},
									Results: []Param{
										{Type: TInt8},
										{Type: TInt32},
										{Type: TInt64},
									},
								},
							},
							{
								Name: "z",
								FuncType: &FuncType{
									Params: []Param{
										{Name: "a", Type: TString},
										{Name: "b", Type: TBool},
									},
									Results: []Param{
										{Name: "c", Type: TUInt8},
										{Name: "d", Type: TFloat},
										{Name: "e", Type: &NamedType{
											UnderlyingType: TInt,
										}},
									},
								},
							},
						},
					},
				},
			},
		}

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		check := NewChecker()

		assert.Equal(t, 0, len(check.Check(pr)))
		assert.Equal(t, scope.Parent, check.pkgScope.Parent)
		assert.Equal(t, scope.Symbols["UserID"].Name, check.pkgScope.Symbols["UserID"].Name)
		assert.Equal(t, scope.Symbols["UserID"].Kind, check.pkgScope.Symbols["UserID"].Kind)
		xx := check.pkgScope.Symbols["UserID"].Type.(*NamedType)
		assert.Equal(t, TInt, xx.UnderlyingType)

		assert.Equal(t, scope.Symbols["test"].Name, check.pkgScope.Symbols["test"].Name)
		source := scope.Symbols["test"].Type.(*InterfaceType)
		destination := check.pkgScope.Symbols["test"].Type.(*InterfaceType)
		assert.Equal(t, len(source.Methods), len(destination.Methods))

		for k := range len(source.Methods) {
			assert.Equal(t, source.Methods[k].Name, destination.Methods[k].Name)
			assert.Equal(t, source.Methods[k].FuncType.Params, destination.Methods[k].FuncType.Params)
			assert.Equal(t, len(source.Methods[k].FuncType.Results), len(destination.Methods[k].FuncType.Results))

			for r := range len(source.Methods[k].FuncType.Results) {
				if k == len(source.Methods)-1 && r == len(source.Methods[k].FuncType.Results)-1 {
					src := source.Methods[k].FuncType.Results[r].Type.(*NamedType)
					dst := destination.Methods[k].FuncType.Results[r].Type.(*NamedType)
					assert.Equal(t, src.Name, dst.Name)
					assert.Equal(t, src.UnderlyingType, dst.UnderlyingType)
				} else {
					assert.Equal(t, source.Methods[k].FuncType.Results[r].Name, destination.Methods[k].FuncType.Results[r].Name)
					assert.Equal(t, source.Methods[k].FuncType.Results[r].Type, destination.Methods[k].FuncType.Results[r].Type)
				}
			}
		}
	})

	t.Run("x2_duplicate", func(t *testing.T) {
		data :=
			`package main
type test interface{
	x(a string, b bool) (c float64, d float32, e uint, f uint32, g uint64 )
	x(a string, b bool) (c float64, d float32, e uint, f uint32, g uint64, g uint64 )
	y(a string, a string, b bool) (int8, int32, zzz)
`
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		// fmt.Printf("%s\n", ast.Dump(pr))
		check := NewChecker()
		assert.Greater(t, len(check.Check(pr)), 0)
	})

	t.Run("x3", func(t *testing.T) {
		data :=
			`package main
type Color enum {
  Red;Blue;Green;Yellow
}
`
		scope := &Scope{
			Parent: nil,
			Symbols: map[string]*Symbol{
				"Color": {
					Name: "Color",
					Kind: SymType,
					Type: &EnumType{
						Variants: []string{"Red", "Blue", "Green", "Yellow"},
					},
				},
			},
		}

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		// fmt.Printf("%s\n", ast.Dump(pr))
		check := NewChecker()

		assert.Equal(t, 0, len(check.Check(pr)))
		assert.Equal(t, scope.Parent, check.pkgScope.Parent)
		src := scope.Symbols["Color"].Type.(*EnumType)
		dst := check.pkgScope.Symbols["Color"].Type.(*EnumType)
		assert.Equal(t, src.Variants, dst.Variants)
	})

	t.Run("x3_duplicate", func(t *testing.T) {
		data :=
			`package main
type Color enum {
  Red;Blue;Green;Yellow;Yellow
}
`
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		// fmt.Printf("%s\n", ast.Dump(pr))
		check := NewChecker()
		assert.Greater(t, len(check.Check(pr)), 0)
	})

	t.Run("x4", func(t *testing.T) {
		data :=
			`package main
type test sum {
  Circle(radius float);Rect(w float, h float);None
}
`
		scope := &Scope{
			Parent: nil,
			Symbols: map[string]*Symbol{
				"test": {
					Name: "test",
					Kind: SymType,
					Type: &SumType{
						Name: "test",
						Variants: []SumVariant{
							{
								Name: "Circle",
								Field: []Param{
									{Name: "radius", Type: TFloat},
								},
							},
							{
								Name: "Rect",
								Field: []Param{
									{Name: "w", Type: TFloat},
									{Name: "h", Type: TFloat},
								},
							},
							{Name: "None"},
						},
					},
				},
			},
		}

		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		// fmt.Printf("%s\n", ast.Dump(pr))
		check := NewChecker()

		assert.Equal(t, 0, len(check.Check(pr)))
		assert.Equal(t, scope.Parent, check.pkgScope.Parent)
		assert.Equal(t, scope.Symbols["test"].Name, check.pkgScope.Symbols["test"].Name)
		assert.Equal(t, scope.Symbols["test"].Kind, check.pkgScope.Symbols["test"].Kind)

		src := scope.Symbols["test"].Type.(*SumType)
		dst := check.pkgScope.Symbols["test"].Type.(*SumType)
		assert.Equal(t, len(src.Variants), len(dst.Variants))
		for k := range len(src.Variants) {
			assert.Equal(t, src.Variants[k], dst.Variants[k])
		}
	})

	t.Run("x4_duplicate", func(t *testing.T) {
		data :=
			`package main
type test sum {
  Circle(radius float);Rect(w float, h float);None;None
}
`
		lex, err := lexer.NewLexer(lexer.Config{StringOnly: true})
		require.NoError(t, err)
		parser := parser.New(lex.FetchTokensFromString(data))
		pr := parser.ParseFile()
		// fmt.Printf("%s\n", ast.Dump(pr))
		check := NewChecker()
		assert.Greater(t, len(check.Check(pr)), 0)
	})

	t.Run("type_decl_name_lookup", func(t *testing.T) {
		assert.Equal(t, "", typeDeclName(nil))
	})

	t.Run("declare_type", func(t *testing.T) {
		check := NewChecker()
		check.declareTypeSymbol(nil)
	})

	t.Run("lookup", func(t *testing.T) {
		check := NewChecker()
		assert.Nil(t, check.pkgScope.Lookup(""))
	})

	t.Run("resolve_type", func(t *testing.T) {
		check := NewChecker()
		assert.Nil(t, check.resolveType(nil))
	})

	t.Run("resolve_named_type", func(t *testing.T) {
		check := NewChecker()
		assert.Nil(t, check.resolveNamedType(&ast.NamedType{Parts: []token.Token{{Kind: token.KWIf}}}))
		p := []token.Token{
			{
				Kind:  token.Ident,
				Value: "x",
			},
			{
				Kind:  token.Ident,
				Value: "y",
			},
		}
		assert.Equal(t, TInvalid, check.resolveNamedType(&ast.NamedType{Parts: p}))
	})
}
