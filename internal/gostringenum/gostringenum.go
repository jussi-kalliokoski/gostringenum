package gostringenum

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"unicode"
)

// Generate a string enum file.
func Generate(fileSet *token.FileSet, pkgs map[string]*ast.Package, pkgPath string, typeName string) *ast.File {
	pkg, spec, zeroValue, values, valueStrings, _ := findDefinitions(fileSet, pkgs, typeName)
	return generateFile(pkg, spec, zeroValue, values, valueStrings)
}

// GenerateWithTestFile generates a string enum file and an accompanying test
// file.
func GenerateWithTestFile(fileSet *token.FileSet, pkgs map[string]*ast.Package, pkgPath string, typeName string) (*ast.File, *ast.File) {
	pkg, spec, zeroValue, values, valueStrings, valueStringLits := findDefinitions(fileSet, pkgs, typeName)
	return generateFile(pkg, spec, zeroValue, values, valueStrings), generateTestFile(pkg, pkgPath, spec, zeroValue, values, valueStringLits)
}

func findDefinitions(fileSet *token.FileSet, pkgs map[string]*ast.Package, typeName string) (*ast.Package, *ast.TypeSpec, *ast.Ident, []*ast.Ident, map[string]*ast.Ident, map[string]*ast.BasicLit) {
	for _, pkg := range pkgs {
		var spec *ast.TypeSpec
		var zeroValue *ast.Ident
		values := make([]*ast.Ident, 0)
		valueStrings := map[string]*ast.Ident{}
		valueStringLits := map[string]*ast.BasicLit{}

		var visitor ast.Visitor
		var inheritedType ast.Expr
		visitor = visitorFunc(func(node ast.Node) ast.Visitor {
			if typeSpec, isTypeSpec := node.(*ast.TypeSpec); isTypeSpec {
				if typeSpec.Name.Name == typeName {
					if typeIdent, isIdent := typeSpec.Type.(*ast.Ident); !isIdent || typeIdent.Name != "int" {
						pos := fileSet.Position(typeIdent.Pos())
						panic(fmt.Errorf("%s at %s:%d:%d is not an int enum type", typeSpec.Name.Name, pos.Filename, pos.Line, pos.Column))
					}
					spec = typeSpec
				}
			}
			if valueSpec, isValueSpec := node.(*ast.ValueSpec); isValueSpec {
				for i, nameIdent := range valueSpec.Names {
					if strings.HasPrefix(nameIdent.Name, typeName) {
						t := valueSpec.Type
						if t == nil {
							t = inheritedType
						}
						if typeIdent, isIdent := t.(*ast.Ident); !isIdent || typeIdent.Name != typeName {
							pos := fileSet.Position(nameIdent.Pos())
							fmt.Println(t)
							panic(fmt.Errorf("value %s at %s:%d:%d is not of type %s", nameIdent.Name, pos.Filename, pos.Line, pos.Column, typeName))
						}
						if len(values) == 0 && len(valueSpec.Values) > 0 {
							switch v := valueSpec.Values[0].(type) {
							case *ast.Ident:
								if v.Name == "iota" {
									zeroValue = nameIdent
								}
							case *ast.BasicLit:
								if v.Value == "0" {
									zeroValue = nameIdent
								}
							}
						}
						values = append(values, nameIdent)
					} else if strings.HasPrefix(nameIdent.Name, unexported(typeName)+"String") {
						if valueBasicLit, isBasicLit := valueSpec.Values[i].(*ast.BasicLit); !isBasicLit || valueBasicLit.Kind != token.STRING {
							pos := fileSet.Position(nameIdent.Pos())
							panic(fmt.Errorf("%s at %s:%d:%d is not not a string literal", nameIdent.Name, pos.Filename, pos.Line, pos.Column))
						} else {
							valueName := typeName + strings.TrimPrefix(nameIdent.Name, unexported(typeName)+"String")
							valueStrings[valueName] = nameIdent
							valueStringLits[valueName] = valueBasicLit
						}
					}
				}
				if valueSpec.Type != nil {
					inheritedType = valueSpec.Type
				}
			}
			return visitor
		})
		ast.Walk(visitor, pkg)

		if spec != nil {
			for _, v := range values {
				if _, hasStringRepresentation := valueStrings[v.Name]; !hasStringRepresentation {
					missingConstName := unexported(typeName) + "String" + strings.TrimPrefix(v.Name, typeName)
					pos := fileSet.Position(v.Pos())
					panic(fmt.Errorf("%s at %s:%d:%d does not have a string representation - fix this by defining a string constant named %s", v.Name, pos.Filename, pos.Line, pos.Column, missingConstName))
				}
			}
			return pkg, spec, zeroValue, values, valueStrings, valueStringLits
		}
	}
	panic(typeName + "not defined in package")
}

func generateFile(pkg *ast.Package, typeSpec *ast.TypeSpec, zeroValue *ast.Ident, values []*ast.Ident, valueStrings map[string]*ast.Ident) *ast.File {
	if zeroValue != nil {
		return generateFileWithZeroValue(pkg, typeSpec, zeroValue, values, valueStrings)
	}
	receiver := &ast.Ident{Name: unexported(typeSpec.Name.Name)}
	fmtImp := &ast.Ident{Name: "fmt"}
	unsafeImp := &ast.Ident{Name: "unsafe"}
	parseFuncIdent := &ast.Ident{Name: "Parse" + typeSpec.Name.Name}
	stringFuncIdent := &ast.Ident{Name: "String"}
	return &ast.File{
		Name: &ast.Ident{Name: pkg.Name},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{Text: "\n// Code generated. DO NOT EDIT."},
					},
				},
				Specs: []ast.Spec{
					&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"unsafe"`}},
					&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"fmt"`}},
				},
			},
			generateParseFunc(fmtImp, parseFuncIdent, typeSpec, values, valueStrings),
			generateStringFunc(fmtImp, stringFuncIdent, receiver, typeSpec, values, valueStrings),
			generateGoStringFunc(fmtImp, stringFuncIdent, receiver, typeSpec, values, valueStrings),
			generateMarshalTextFunc(unsafeImp, fmtImp, stringFuncIdent, receiver, typeSpec, values, valueStrings),
			generateUnmarshalTextFunc(unsafeImp, parseFuncIdent, receiver, typeSpec),
		},
	}
}

func generateFileWithZeroValue(pkg *ast.Package, typeSpec *ast.TypeSpec, zeroValue *ast.Ident, values []*ast.Ident, valueStrings map[string]*ast.Ident) *ast.File {
	receiver := &ast.Ident{Name: unexported(typeSpec.Name.Name)}
	fmtImp := &ast.Ident{Name: "fmt"}
	unsafeImp := &ast.Ident{Name: "unsafe"}
	fromStringFuncIdent := &ast.Ident{Name: typeSpec.Name.Name + "FromString"}
	stringFuncIdent := &ast.Ident{Name: "String"}
	return &ast.File{
		Name: &ast.Ident{Name: pkg.Name},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{Text: "\n// Code generated. DO NOT EDIT."},
					},
				},
				Specs: []ast.Spec{
					&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"unsafe"`}},
				},
			},
			generateFromStringFunc(fmtImp, fromStringFuncIdent, typeSpec, zeroValue, values, valueStrings),
			generateStringFuncWithZeroValue(fmtImp, stringFuncIdent, receiver, typeSpec, zeroValue, values, valueStrings),
			generateGoStringFuncWithZeroValue(fmtImp, stringFuncIdent, receiver, typeSpec, zeroValue, values, valueStrings),
			generateMarshalTextFuncWithZeroValue(unsafeImp, fmtImp, stringFuncIdent, receiver, typeSpec, zeroValue, values, valueStrings),
			generateUnmarshalTextFuncWithZeroValue(unsafeImp, fromStringFuncIdent, receiver, typeSpec),
		},
	}
}

func generateParseFunc(fmtImp *ast.Ident, parseFuncIdent *ast.Ident, typeSpec *ast.TypeSpec, values []*ast.Ident, valueStrings map[string]*ast.Ident) ast.Decl {
	strIdent := &ast.Ident{Name: "str"}
	caseClauses := make([]ast.Stmt, 0, len(values)+1)
	for _, v := range values {
		caseClauses = append(caseClauses, &ast.CaseClause{
			List: []ast.Expr{valueStrings[v.Name]},
			Body: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						v,
						builtinNil,
					},
				},
			},
		})
	}
	caseClauses = append(caseClauses, &ast.CaseClause{
		Body: []ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.BasicLit{Kind: token.INT, Value: "0"},
					errorf(fmtImp, fmt.Sprintf("not a %s: %s", typeSpec.Name.Name, "%q"), strIdent),
				},
			},
		},
	})
	return &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "\n// " + parseFuncIdent.Name + " parses a " + typeSpec.Name.Name + " from its string representation."},
			},
		},
		Name: parseFuncIdent,
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{strIdent},
						Type:  builtinString,
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: typeSpec.Name},
					{Type: builtinError},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.SwitchStmt{
					Tag: strIdent,
					Body: &ast.BlockStmt{
						List: caseClauses,
					},
				},
			},
		},
	}
}

func generateFromStringFunc(fmtImp *ast.Ident, fromStringFuncIdent *ast.Ident, typeSpec *ast.TypeSpec, zeroValue *ast.Ident, values []*ast.Ident, valueStrings map[string]*ast.Ident) ast.Decl {
	strIdent := &ast.Ident{Name: "str"}
	caseClauses := make([]ast.Stmt, 0, len(values)+1)
	for _, v := range values {
		if v == zeroValue {
			continue
		}
		caseClauses = append(caseClauses, &ast.CaseClause{
			List: []ast.Expr{valueStrings[v.Name]},
			Body: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						v,
					},
				},
			},
		})
	}
	caseClauses = append(caseClauses, &ast.CaseClause{
		Body: []ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					zeroValue,
				},
			},
		},
	})
	return &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "\n// " + fromStringFuncIdent.Name + " parses a " + typeSpec.Name.Name + " from its string representation."},
			},
		},
		Name: fromStringFuncIdent,
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{strIdent},
						Type:  builtinString,
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{Type: typeSpec.Name},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.SwitchStmt{
					Tag: strIdent,
					Body: &ast.BlockStmt{
						List: caseClauses,
					},
				},
			},
		},
	}
}

func generateStringFunc(fmtImp *ast.Ident, stringFuncIdent *ast.Ident, receiver *ast.Ident, typeSpec *ast.TypeSpec, values []*ast.Ident, valueStrings map[string]*ast.Ident) ast.Decl {
	caseClauses := make([]ast.Stmt, 0, len(values)+1)
	for _, v := range values {
		caseClauses = append(caseClauses, &ast.CaseClause{
			List: []ast.Expr{v},
			Body: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						valueStrings[v.Name],
					},
				},
			},
		})
	}
	caseClauses = append(caseClauses, &ast.CaseClause{
		Body: []ast.Stmt{
			&ast.ExprStmt{
				X: &ast.CallExpr{
					Fun: builtinPanic,
					Args: []ast.Expr{
						errorf(fmtImp, fmt.Sprintf("not a %s: %s", typeSpec.Name.Name, "%d"), receiver),
					},
				},
			},
		},
	})
	return &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "\n// String returns the string representation of a " + typeSpec.Name.Name + "."},
				{Text: "//"},
				{Text: "// Will panic if the value is not a valid " + typeSpec.Name.Name + "."},
			},
		},
		Name: stringFuncIdent,
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{receiver},
					Type:  typeSpec.Name,
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: builtinString,
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.SwitchStmt{
					Tag: receiver,
					Body: &ast.BlockStmt{
						List: caseClauses,
					},
				},
			},
		},
	}
}

func generateStringFuncWithZeroValue(fmtImp *ast.Ident, stringFuncIdent *ast.Ident, receiver *ast.Ident, typeSpec *ast.TypeSpec, zeroValue *ast.Ident, values []*ast.Ident, valueStrings map[string]*ast.Ident) ast.Decl {
	caseClauses := make([]ast.Stmt, 0, len(values)+1)
	for _, v := range values {
		if v == zeroValue {
			continue
		}
		caseClauses = append(caseClauses, &ast.CaseClause{
			List: []ast.Expr{v},
			Body: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						valueStrings[v.Name],
					},
				},
			},
		})
	}
	caseClauses = append(caseClauses, &ast.CaseClause{
		Body: []ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					valueStrings[zeroValue.Name],
				},
			},
		},
	})
	return &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "\n// String returns the string representation of a " + typeSpec.Name.Name + "."},
			},
		},
		Name: stringFuncIdent,
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{receiver},
					Type:  typeSpec.Name,
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: builtinString,
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.SwitchStmt{
					Tag: receiver,
					Body: &ast.BlockStmt{
						List: caseClauses,
					},
				},
			},
		},
	}
}

func generateGoStringFunc(fmtImp *ast.Ident, stringFuncIdent *ast.Ident, receiver *ast.Ident, typeSpec *ast.TypeSpec, values []*ast.Ident, valueStrings map[string]*ast.Ident) ast.Decl {
	caseClauses := make([]ast.Stmt, 0, len(values)+1)
	for _, v := range values {
		caseClauses = append(caseClauses, &ast.CaseClause{
			List: []ast.Expr{v},
			Body: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						valueStrings[v.Name],
					},
				},
			},
		})
	}
	caseClauses = append(caseClauses, &ast.CaseClause{
		Body: []ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{X: fmtImp, Sel: &ast.Ident{Name: "Sprintf"}},
						Args: []ast.Expr{
							&ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf(`"%s{invalid %s}"`, typeSpec.Name.Name, "%d")},
							receiver,
						},
					},
				},
			},
		},
	})
	return &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "\n// GoString implements fmt.GoStringer."},
			},
		},
		Name: &ast.Ident{Name: "GoString"},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{receiver},
					Type:  typeSpec.Name,
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: builtinString,
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.SwitchStmt{
					Tag: receiver,
					Body: &ast.BlockStmt{
						List: caseClauses,
					},
				},
			},
		},
	}
}

func generateGoStringFuncWithZeroValue(fmtImp *ast.Ident, stringFuncIdent *ast.Ident, receiver *ast.Ident, typeSpec *ast.TypeSpec, zeroValue *ast.Ident, values []*ast.Ident, valueStrings map[string]*ast.Ident) ast.Decl {
	return &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "\n// GoString implements fmt.GoStringer."},
			},
		},
		Name: &ast.Ident{Name: "GoString"},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{receiver},
					Type:  typeSpec.Name,
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: builtinString,
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{X: receiver, Sel: stringFuncIdent},
						},
					},
				},
			},
		},
	}
}

func generateMarshalTextFunc(unsafeImp *ast.Ident, fmtImp *ast.Ident, stringFuncIdent *ast.Ident, receiver *ast.Ident, typeSpec *ast.TypeSpec, values []*ast.Ident, valueStrings map[string]*ast.Ident) ast.Decl {
	strIdent := &ast.Ident{Name: "str"}
	caseClauses := make([]ast.Stmt, 0, len(values)+1)
	for _, v := range values {
		caseClauses = append(caseClauses, &ast.CaseClause{
			List: []ast.Expr{v},
			Body: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{strIdent},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{valueStrings[v.Name]},
				},
			},
		})
	}
	caseClauses = append(caseClauses, &ast.CaseClause{
		Body: []ast.Stmt{
			&ast.ReturnStmt{
				Results: []ast.Expr{
					builtinNil,
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{X: fmtImp, Sel: &ast.Ident{Name: "Errorf"}},
						Args: []ast.Expr{
							&ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf(`"%s{invalid %s}"`, typeSpec.Name.Name, "%d")},
							receiver,
						},
					},
				},
			},
		},
	})
	return &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "\n// MarshalText implements encoding.TextMarshaler."},
				{Text: "//"},
				{Text: "// The returned byte slice is read-only and writing to it will panic."},
			},
		},
		Name: &ast.Ident{Name: "MarshalText"},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{receiver},
					Type:  typeSpec.Name,
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.ArrayType{Elt: builtinByte},
					},
					{
						Type: builtinError,
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{strIdent},
								Type:  builtinString,
							},
						},
					},
				},
				&ast.SwitchStmt{
					Tag:  receiver,
					Body: &ast.BlockStmt{List: caseClauses},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.StarExpr{
							X: &ast.CallExpr{
								Fun: &ast.StarExpr{
									X: &ast.ArrayType{Elt: builtinByte},
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{X: unsafeImp, Sel: &ast.Ident{Name: "Pointer"}},
										Args: []ast.Expr{
											&ast.UnaryExpr{
												Op: token.AND,
												X:  strIdent,
											},
										},
									},
								},
							},
						},
						builtinNil,
					},
				},
			},
		},
	}
}

func generateMarshalTextFuncWithZeroValue(unsafeImp *ast.Ident, fmtImp *ast.Ident, stringFuncIdent *ast.Ident, receiver *ast.Ident, typeSpec *ast.TypeSpec, zeroValue *ast.Ident, values []*ast.Ident, valueStrings map[string]*ast.Ident) ast.Decl {
	strIdent := &ast.Ident{Name: "str"}
	return &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "\n// MarshalText implements encoding.TextMarshaler."},
				{Text: "//"},
				{Text: "// The returned byte slice is read-only and writing to it will panic."},
			},
		},
		Name: &ast.Ident{Name: "MarshalText"},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{receiver},
					Type:  typeSpec.Name,
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: &ast.ArrayType{Elt: builtinByte},
					},
					{
						Type: builtinError,
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						strIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{X: receiver, Sel: stringFuncIdent},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						&ast.StarExpr{
							X: &ast.CallExpr{
								Fun: &ast.StarExpr{
									X: &ast.ArrayType{Elt: builtinByte},
								},
								Args: []ast.Expr{
									&ast.CallExpr{
										Fun: &ast.SelectorExpr{X: unsafeImp, Sel: &ast.Ident{Name: "Pointer"}},
										Args: []ast.Expr{
											&ast.UnaryExpr{
												Op: token.AND,
												X:  strIdent,
											},
										},
									},
								},
							},
						},
						builtinNil,
					},
				},
			},
		},
	}
}

func generateUnmarshalTextFunc(unsafeImp *ast.Ident, parseFuncIdent *ast.Ident, receiver *ast.Ident, typeSpec *ast.TypeSpec) ast.Decl {
	textIdent := &ast.Ident{Name: "text"}
	errIdent := &ast.Ident{Name: "err"}
	return &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "\n// UnmarshalText implements encoding.TextUnmarshaler."},
			},
		},
		Name: &ast.Ident{Name: "UnmarshalText"},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{receiver},
					Type:  &ast.StarExpr{X: typeSpec.Name},
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{textIdent},
						Type:  &ast.ArrayType{Elt: builtinByte},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: builtinError,
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{errIdent},
								Type:  builtinError,
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.StarExpr{X: receiver},
						errIdent,
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: parseFuncIdent,
							Args: []ast.Expr{
								&ast.StarExpr{
									X: &ast.CallExpr{
										Fun: &ast.StarExpr{X: builtinString},
										Args: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{X: unsafeImp, Sel: &ast.Ident{Name: "Pointer"}},
												Args: []ast.Expr{
													&ast.UnaryExpr{
														Op: token.AND,
														X:  textIdent,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						errIdent,
					},
				},
			},
		},
	}
}

func generateUnmarshalTextFuncWithZeroValue(unsafeImp *ast.Ident, fromStringFuncIdent *ast.Ident, receiver *ast.Ident, typeSpec *ast.TypeSpec) ast.Decl {
	textIdent := &ast.Ident{Name: "text"}
	return &ast.FuncDecl{
		Doc: &ast.CommentGroup{
			List: []*ast.Comment{
				{Text: "\n// UnmarshalText implements encoding.TextUnmarshaler."},
			},
		},
		Name: &ast.Ident{Name: "UnmarshalText"},
		Recv: &ast.FieldList{
			List: []*ast.Field{
				{
					Names: []*ast.Ident{receiver},
					Type:  &ast.StarExpr{X: typeSpec.Name},
				},
			},
		},
		Type: &ast.FuncType{
			Params: &ast.FieldList{
				List: []*ast.Field{
					{
						Names: []*ast.Ident{textIdent},
						Type:  &ast.ArrayType{Elt: builtinByte},
					},
				},
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					{
						Type: builtinError,
					},
				},
			},
		},
		Body: &ast.BlockStmt{
			List: []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						&ast.StarExpr{X: receiver},
					},
					Tok: token.ASSIGN,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: fromStringFuncIdent,
							Args: []ast.Expr{
								&ast.StarExpr{
									X: &ast.CallExpr{
										Fun: &ast.StarExpr{X: builtinString},
										Args: []ast.Expr{
											&ast.CallExpr{
												Fun: &ast.SelectorExpr{X: unsafeImp, Sel: &ast.Ident{Name: "Pointer"}},
												Args: []ast.Expr{
													&ast.UnaryExpr{
														Op: token.AND,
														X:  textIdent,
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
				&ast.ReturnStmt{
					Results: []ast.Expr{
						builtinNil,
					},
				},
			},
		},
	}
}

func generateTestFile(pkg *ast.Package, pkgPath string, typeSpec *ast.TypeSpec, zeroValue *ast.Ident, values []*ast.Ident, valueStringLits map[string]*ast.BasicLit) *ast.File {
	if zeroValue != nil {
		return generateTestFileWithZeroValue(pkg, pkgPath, typeSpec, zeroValue, values, valueStringLits)
	}
	pkgIdent := &ast.Ident{Name: pkg.Name}
	receiver := &ast.Ident{Name: unexported(typeSpec.Name.Name)}
	testingImp := &ast.Ident{Name: "testing"}
	tType := &ast.StarExpr{X: &ast.SelectorExpr{X: testingImp, Sel: &ast.Ident{Name: "T"}}}
	t := &ast.Ident{Name: "t"}
	return &ast.File{
		Name: &ast.Ident{Name: pkg.Name + "_test"},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{Text: "\n// Code generated. DO NOT EDIT."},
					},
				},
				Specs: []ast.Spec{
					&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"testing"`}},
					&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", pkgPath)}},
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{Name: "Test" + typeSpec.Name.Name},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{Names: []*ast.Ident{t}, Type: tType},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						generateParseTest(pkgIdent, receiver, typeSpec, t, tType, values, valueStringLits),
						generateStringTest(pkgIdent, receiver, typeSpec, t, tType, values, valueStringLits),
						generateGoStringTest(pkgIdent, receiver, typeSpec, t, tType, values, valueStringLits),
						generateMarshalTextTest(pkgIdent, receiver, typeSpec, t, tType, values, valueStringLits),
						generateUnmarshalTextTest(pkgIdent, receiver, typeSpec, t, tType, values, valueStringLits),
					},
				},
			},
		},
	}
}

func generateTestFileWithZeroValue(pkg *ast.Package, pkgPath string, typeSpec *ast.TypeSpec, zeroValue *ast.Ident, values []*ast.Ident, valueStringLits map[string]*ast.BasicLit) *ast.File {
	pkgIdent := &ast.Ident{Name: pkg.Name}
	receiver := &ast.Ident{Name: unexported(typeSpec.Name.Name)}
	testingImp := &ast.Ident{Name: "testing"}
	tType := &ast.StarExpr{X: &ast.SelectorExpr{X: testingImp, Sel: &ast.Ident{Name: "T"}}}
	t := &ast.Ident{Name: "t"}
	return &ast.File{
		Name: &ast.Ident{Name: pkg.Name + "_test"},
		Decls: []ast.Decl{
			&ast.GenDecl{
				Tok: token.IMPORT,
				Doc: &ast.CommentGroup{
					List: []*ast.Comment{
						{Text: "\n// Code generated. DO NOT EDIT."},
					},
				},
				Specs: []ast.Spec{
					&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: `"testing"`}},
					&ast.ImportSpec{Path: &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", pkgPath)}},
				},
			},
			&ast.FuncDecl{
				Name: &ast.Ident{Name: "Test" + typeSpec.Name.Name},
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{Names: []*ast.Ident{t}, Type: tType},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						generateFromStringTest(pkgIdent, receiver, typeSpec, t, tType, values, valueStringLits),
						generateStringTestWithZeroValue(pkgIdent, receiver, typeSpec, t, tType, values, valueStringLits),
						generateGoStringTestWithZeroValue(pkgIdent, receiver, typeSpec, t, tType, values, valueStringLits),
						generateMarshalTextTestWithZeroValue(pkgIdent, receiver, typeSpec, t, tType, values, valueStringLits),
						generateUnmarshalTextTestWithZeroValue(pkgIdent, receiver, typeSpec, t, tType, values, valueStringLits),
					},
				},
			},
		},
	}
}

func generateParseTest(pkgIdent *ast.Ident, receiver *ast.Ident, typeSpec *ast.TypeSpec, t *ast.Ident, tType ast.Expr, values []*ast.Ident, valueStringLits map[string]*ast.BasicLit) ast.Stmt {
	parseIdent := &ast.Ident{Name: "Parse" + typeSpec.Name.Name}
	expectedIdent := &ast.Ident{Name: "expected"}
	receivedIdent := &ast.Ident{Name: "received"}
	errIdent := &ast.Ident{Name: "err"}
	testCases := make([]ast.Stmt, 0)
	for _, v := range values {
		testCases = append(testCases, &ast.ExprStmt{
			X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", v.Name)}, []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						expectedIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.SelectorExpr{X: pkgIdent, Sel: v},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						receivedIdent,
						errIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{X: pkgIdent, Sel: parseIdent},
							Args: []ast.Expr{
								valueStringLits[v.Name],
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{X: errIdent, Op: token.NEQ, Y: builtinNil},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Fatalf"}},
									Args: []ast.Expr{
										&ast.BasicLit{Kind: token.STRING, Value: `"expected no error, got %#v"`},
										errIdent,
									},
								},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{X: expectedIdent, Op: token.NEQ, Y: receivedIdent},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Fatalf"}},
									Args: []ast.Expr{
										&ast.BasicLit{Kind: token.STRING, Value: `"expected %#v, got %#v"`},
										expectedIdent,
										receivedIdent,
									},
								},
							},
						},
					},
				},
			}),
		})
	}
	testCases = append(testCases, &ast.ExprStmt{
		X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: `"Invalid value"`}, []ast.Stmt{
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					builtinNothing,
					errIdent,
				},
				Tok: token.DEFINE,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{X: pkgIdent, Sel: parseIdent},
						Args: []ast.Expr{
							&ast.CallExpr{
								Fun: builtinString,
								Args: []ast.Expr{
									&ast.CompositeLit{
										Type: &ast.ArrayType{Elt: builtinByte},
										Elts: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "0"}},
									},
								},
							},
						},
					},
				},
			},
			&ast.IfStmt{
				Cond: &ast.BinaryExpr{
					X:  errIdent,
					Op: token.EQL,
					Y:  builtinNil,
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Fatal"}},
								Args: []ast.Expr{
									&ast.BasicLit{Kind: token.STRING, Value: `"expected an error, got nil"`},
								},
							},
						},
					},
				},
			},
		}),
	})
	return &ast.ExprStmt{
		X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", parseIdent.Name)}, testCases),
	}
}

func generateFromStringTest(pkgIdent *ast.Ident, receiver *ast.Ident, typeSpec *ast.TypeSpec, t *ast.Ident, tType ast.Expr, values []*ast.Ident, valueStringLits map[string]*ast.BasicLit) ast.Stmt {
	fromStringIdent := &ast.Ident{Name: typeSpec.Name.Name + "FromString"}
	expectedIdent := &ast.Ident{Name: "expected"}
	receivedIdent := &ast.Ident{Name: "received"}
	testCases := make([]ast.Stmt, 0)
	for _, v := range values {
		testCases = append(testCases, &ast.ExprStmt{
			X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", v.Name)}, []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						expectedIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.SelectorExpr{X: pkgIdent, Sel: v},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						receivedIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{X: pkgIdent, Sel: fromStringIdent},
							Args: []ast.Expr{
								valueStringLits[v.Name],
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{X: expectedIdent, Op: token.NEQ, Y: receivedIdent},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Fatalf"}},
									Args: []ast.Expr{
										&ast.BasicLit{Kind: token.STRING, Value: `"expected %#v, got %#v"`},
										expectedIdent,
										receivedIdent,
									},
								},
							},
						},
					},
				},
			}),
		})
	}
	return &ast.ExprStmt{
		X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", fromStringIdent.Name)}, testCases),
	}
}

func generateStringTest(pkgIdent *ast.Ident, receiver *ast.Ident, typeSpec *ast.TypeSpec, t *ast.Ident, tType ast.Expr, values []*ast.Ident, valueStringLits map[string]*ast.BasicLit) ast.Stmt {
	expectedIdent := &ast.Ident{Name: "expected"}
	receivedIdent := &ast.Ident{Name: "received"}
	testCases := make([]ast.Stmt, 0)
	for _, v := range values {
		testCases = append(testCases, &ast.ExprStmt{
			X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", v.Name)}, []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						expectedIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						valueStringLits[v.Name],
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						receivedIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{X: &ast.SelectorExpr{X: pkgIdent, Sel: v}, Sel: &ast.Ident{Name: "String"}},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{X: expectedIdent, Op: token.NEQ, Y: receivedIdent},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Fatalf"}},
									Args: []ast.Expr{
										&ast.BasicLit{Kind: token.STRING, Value: `"expected %q, got %q"`},
										expectedIdent,
										receivedIdent,
									},
								},
							},
						},
					},
				},
			}),
		})
	}
	testCases = append(testCases, &ast.ExprStmt{
		X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: `"Invalid value"`}, []ast.Stmt{
			&ast.DeferStmt{
				Call: &ast.CallExpr{
					Fun: &ast.FuncLit{
						Type: &ast.FuncType{
							Params: &ast.FieldList{},
						},
						Body: &ast.BlockStmt{
							List: []ast.Stmt{
								&ast.IfStmt{
									Cond: &ast.BinaryExpr{
										X:  &ast.CallExpr{Fun: builtinRecover},
										Op: token.EQL,
										Y:  builtinNil,
									},
									Body: &ast.BlockStmt{
										List: []ast.Stmt{
											&ast.ExprStmt{
												X: &ast.CallExpr{
													Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Fatal"}},
													Args: []ast.Expr{
														&ast.BasicLit{Kind: token.STRING, Value: `"expected a panic"`},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names:  []*ast.Ident{receiver},
							Type:   &ast.SelectorExpr{X: pkgIdent, Sel: typeSpec.Name},
							Values: []ast.Expr{&ast.UnaryExpr{Op: token.XOR, X: &ast.BasicLit{Kind: token.INT, Value: "0"}}},
						},
					},
				},
			},
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					builtinNothing,
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{X: receiver, Sel: &ast.Ident{Name: "String"}},
					},
				},
			},
		}),
	})
	return &ast.ExprStmt{
		X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: `"String"`}, testCases),
	}
}

func generateStringTestWithZeroValue(pkgIdent *ast.Ident, receiver *ast.Ident, typeSpec *ast.TypeSpec, t *ast.Ident, tType ast.Expr, values []*ast.Ident, valueStringLits map[string]*ast.BasicLit) ast.Stmt {
	expectedIdent := &ast.Ident{Name: "expected"}
	receivedIdent := &ast.Ident{Name: "received"}
	testCases := make([]ast.Stmt, 0)
	for _, v := range values {
		testCases = append(testCases, &ast.ExprStmt{
			X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", v.Name)}, []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						expectedIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						valueStringLits[v.Name],
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						receivedIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{X: &ast.SelectorExpr{X: pkgIdent, Sel: v}, Sel: &ast.Ident{Name: "String"}},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{X: expectedIdent, Op: token.NEQ, Y: receivedIdent},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Fatalf"}},
									Args: []ast.Expr{
										&ast.BasicLit{Kind: token.STRING, Value: `"expected %q, got %q"`},
										expectedIdent,
										receivedIdent,
									},
								},
							},
						},
					},
				},
			}),
		})
	}
	return &ast.ExprStmt{
		X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: `"String"`}, testCases),
	}
}

func generateGoStringTest(pkgIdent *ast.Ident, receiver *ast.Ident, typeSpec *ast.TypeSpec, t *ast.Ident, tType ast.Expr, values []*ast.Ident, valueStringLits map[string]*ast.BasicLit) ast.Stmt {
	expectedIdent := &ast.Ident{Name: "expected"}
	receivedIdent := &ast.Ident{Name: "received"}
	testCases := make([]ast.Stmt, 0)
	for _, v := range values {
		testCases = append(testCases, &ast.ExprStmt{
			X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", v.Name)}, []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						expectedIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						valueStringLits[v.Name],
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						receivedIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{X: &ast.SelectorExpr{X: pkgIdent, Sel: v}, Sel: &ast.Ident{Name: "GoString"}},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{X: expectedIdent, Op: token.NEQ, Y: receivedIdent},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Fatalf"}},
									Args: []ast.Expr{
										&ast.BasicLit{Kind: token.STRING, Value: `"expected %q, got %q"`},
										expectedIdent,
										receivedIdent,
									},
								},
							},
						},
					},
				},
			}),
		})
	}
	testCases = append(testCases, &ast.ExprStmt{
		X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: `"Invalid value"`}, []ast.Stmt{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names:  []*ast.Ident{receiver},
							Type:   &ast.SelectorExpr{X: pkgIdent, Sel: typeSpec.Name},
							Values: []ast.Expr{&ast.UnaryExpr{Op: token.XOR, X: &ast.BasicLit{Kind: token.INT, Value: "0"}}},
						},
					},
				},
			},
			&ast.AssignStmt{
				Lhs: []ast.Expr{
					builtinNothing,
				},
				Tok: token.ASSIGN,
				Rhs: []ast.Expr{
					&ast.CallExpr{
						Fun: &ast.SelectorExpr{X: receiver, Sel: &ast.Ident{Name: "GoString"}},
					},
				},
			},
		}),
	})
	return &ast.ExprStmt{
		X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: `"GoString"`}, testCases),
	}
}

func generateGoStringTestWithZeroValue(pkgIdent *ast.Ident, receiver *ast.Ident, typeSpec *ast.TypeSpec, t *ast.Ident, tType ast.Expr, values []*ast.Ident, valueStringLits map[string]*ast.BasicLit) ast.Stmt {
	expectedIdent := &ast.Ident{Name: "expected"}
	receivedIdent := &ast.Ident{Name: "received"}
	testCases := make([]ast.Stmt, 0)
	for _, v := range values {
		testCases = append(testCases, &ast.ExprStmt{
			X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", v.Name)}, []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						expectedIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						valueStringLits[v.Name],
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						receivedIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{X: &ast.SelectorExpr{X: pkgIdent, Sel: v}, Sel: &ast.Ident{Name: "GoString"}},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{X: expectedIdent, Op: token.NEQ, Y: receivedIdent},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Fatalf"}},
									Args: []ast.Expr{
										&ast.BasicLit{Kind: token.STRING, Value: `"expected %q, got %q"`},
										expectedIdent,
										receivedIdent,
									},
								},
							},
						},
					},
				},
			}),
		})
	}
	return &ast.ExprStmt{
		X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: `"GoString"`}, testCases),
	}
}

func generateMarshalTextTest(pkgIdent *ast.Ident, receiver *ast.Ident, typeSpec *ast.TypeSpec, t *ast.Ident, tType ast.Expr, values []*ast.Ident, valueStringLits map[string]*ast.BasicLit) ast.Stmt {
	marshalTextIdent := &ast.Ident{Name: "MarshalText"}
	expectedIdent := &ast.Ident{Name: "expected"}
	receivedIdent := &ast.Ident{Name: "received"}
	errIdent := &ast.Ident{Name: "err"}
	testCases := make([]ast.Stmt, 0)
	for _, v := range values {
		testCases = append(testCases, &ast.ExprStmt{
			X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", v.Name)}, []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						expectedIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.ArrayType{Elt: builtinByte},
							Args: []ast.Expr{
								valueStringLits[v.Name],
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						receivedIdent,
						errIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{X: &ast.SelectorExpr{X: pkgIdent, Sel: v}, Sel: marshalTextIdent},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{X: errIdent, Op: token.NEQ, Y: builtinNil},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Fatalf"}},
									Args: []ast.Expr{
										&ast.BasicLit{Kind: token.STRING, Value: `"expected no error, got %#v"`},
										errIdent,
									},
								},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  &ast.CallExpr{Fun: builtinString, Args: []ast.Expr{expectedIdent}},
						Op: token.NEQ,
						Y:  &ast.CallExpr{Fun: builtinString, Args: []ast.Expr{receivedIdent}},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Fatalf"}},
									Args: []ast.Expr{
										&ast.BasicLit{Kind: token.STRING, Value: `"expected %#v, got %#v"`},
										expectedIdent,
										receivedIdent,
									},
								},
							},
						},
					},
				},
			}),
		})
	}
	testCases = append(testCases, &ast.ExprStmt{
		X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: `"Invalid value"`}, []ast.Stmt{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{receiver},
							Type:  &ast.SelectorExpr{X: pkgIdent, Sel: typeSpec.Name},
						},
					},
				},
			},
			&ast.IfStmt{
				Cond: &ast.BinaryExpr{
					X:  errIdent,
					Op: token.EQL,
					Y:  builtinNil,
				},
				Init: &ast.AssignStmt{
					Lhs: []ast.Expr{
						builtinNothing,
						errIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{X: receiver, Sel: marshalTextIdent},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Fatal"}},
								Args: []ast.Expr{
									&ast.BasicLit{Kind: token.STRING, Value: `"expected an error, got nil"`},
								},
							},
						},
					},
				},
			},
		}),
	})
	return &ast.ExprStmt{
		X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", marshalTextIdent.Name)}, testCases),
	}
}

func generateMarshalTextTestWithZeroValue(pkgIdent *ast.Ident, receiver *ast.Ident, typeSpec *ast.TypeSpec, t *ast.Ident, tType ast.Expr, values []*ast.Ident, valueStringLits map[string]*ast.BasicLit) ast.Stmt {
	marshalTextIdent := &ast.Ident{Name: "MarshalText"}
	expectedIdent := &ast.Ident{Name: "expected"}
	receivedIdent := &ast.Ident{Name: "received"}
	errIdent := &ast.Ident{Name: "err"}
	testCases := make([]ast.Stmt, 0)
	for _, v := range values {
		testCases = append(testCases, &ast.ExprStmt{
			X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", v.Name)}, []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						expectedIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.ArrayType{Elt: builtinByte},
							Args: []ast.Expr{
								valueStringLits[v.Name],
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						receivedIdent,
						errIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{X: &ast.SelectorExpr{X: pkgIdent, Sel: v}, Sel: marshalTextIdent},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{X: errIdent, Op: token.NEQ, Y: builtinNil},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Fatalf"}},
									Args: []ast.Expr{
										&ast.BasicLit{Kind: token.STRING, Value: `"expected no error, got %#v"`},
										errIdent,
									},
								},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{
						X:  &ast.CallExpr{Fun: builtinString, Args: []ast.Expr{expectedIdent}},
						Op: token.NEQ,
						Y:  &ast.CallExpr{Fun: builtinString, Args: []ast.Expr{receivedIdent}},
					},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Fatalf"}},
									Args: []ast.Expr{
										&ast.BasicLit{Kind: token.STRING, Value: `"expected %#v, got %#v"`},
										expectedIdent,
										receivedIdent,
									},
								},
							},
						},
					},
				},
			}),
		})
	}
	return &ast.ExprStmt{
		X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", marshalTextIdent.Name)}, testCases),
	}
}

func generateUnmarshalTextTest(pkgIdent *ast.Ident, receiver *ast.Ident, typeSpec *ast.TypeSpec, t *ast.Ident, tType ast.Expr, values []*ast.Ident, valueStringLits map[string]*ast.BasicLit) ast.Stmt {
	unmarshalTextIdent := &ast.Ident{Name: "UnmarshalText"}
	expectedIdent := &ast.Ident{Name: "expected"}
	receivedIdent := &ast.Ident{Name: "received"}
	errIdent := &ast.Ident{Name: "err"}
	testCases := make([]ast.Stmt, 0)
	for _, v := range values {
		testCases = append(testCases, &ast.ExprStmt{
			X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", v.Name)}, []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						expectedIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.SelectorExpr{X: pkgIdent, Sel: v},
					},
				},
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{receivedIdent},
								Type:  &ast.SelectorExpr{X: pkgIdent, Sel: typeSpec.Name},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						errIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{X: receivedIdent, Sel: unmarshalTextIdent},
							Args: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.ArrayType{Elt: builtinByte},
									Args: []ast.Expr{
										valueStringLits[v.Name],
									},
								},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{X: errIdent, Op: token.NEQ, Y: builtinNil},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Fatalf"}},
									Args: []ast.Expr{
										&ast.BasicLit{Kind: token.STRING, Value: `"expected no error, got %#v"`},
										errIdent,
									},
								},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{X: expectedIdent, Op: token.NEQ, Y: receivedIdent},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Fatalf"}},
									Args: []ast.Expr{
										&ast.BasicLit{Kind: token.STRING, Value: `"expected %#v, got %#v"`},
										expectedIdent,
										receivedIdent,
									},
								},
							},
						},
					},
				},
			}),
		})
	}
	testCases = append(testCases, &ast.ExprStmt{
		X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: `"Invalid value"`}, []ast.Stmt{
			&ast.DeclStmt{
				Decl: &ast.GenDecl{
					Tok: token.VAR,
					Specs: []ast.Spec{
						&ast.ValueSpec{
							Names: []*ast.Ident{receiver},
							Type:  &ast.SelectorExpr{X: pkgIdent, Sel: typeSpec.Name},
						},
					},
				},
			},
			&ast.IfStmt{
				Cond: &ast.BinaryExpr{
					X:  errIdent,
					Op: token.EQL,
					Y:  builtinNil,
				},
				Init: &ast.AssignStmt{
					Lhs: []ast.Expr{
						errIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{X: receiver, Sel: unmarshalTextIdent},
							Args: []ast.Expr{
								&ast.CompositeLit{
									Type: &ast.ArrayType{Elt: builtinByte},
									Elts: []ast.Expr{&ast.BasicLit{Kind: token.INT, Value: "0"}},
								},
							},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Fatal"}},
								Args: []ast.Expr{
									&ast.BasicLit{Kind: token.STRING, Value: `"expected an error, got nil"`},
								},
							},
						},
					},
				},
			},
		}),
	})
	return &ast.ExprStmt{
		X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", unmarshalTextIdent.Name)}, testCases),
	}
}

func generateUnmarshalTextTestWithZeroValue(pkgIdent *ast.Ident, receiver *ast.Ident, typeSpec *ast.TypeSpec, t *ast.Ident, tType ast.Expr, values []*ast.Ident, valueStringLits map[string]*ast.BasicLit) ast.Stmt {
	unmarshalTextIdent := &ast.Ident{Name: "UnmarshalText"}
	expectedIdent := &ast.Ident{Name: "expected"}
	receivedIdent := &ast.Ident{Name: "received"}
	errIdent := &ast.Ident{Name: "err"}
	testCases := make([]ast.Stmt, 0)
	for _, v := range values {
		testCases = append(testCases, &ast.ExprStmt{
			X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", v.Name)}, []ast.Stmt{
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						expectedIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.SelectorExpr{X: pkgIdent, Sel: v},
					},
				},
				&ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{
							&ast.ValueSpec{
								Names: []*ast.Ident{receivedIdent},
								Type:  &ast.SelectorExpr{X: pkgIdent, Sel: typeSpec.Name},
							},
						},
					},
				},
				&ast.AssignStmt{
					Lhs: []ast.Expr{
						errIdent,
					},
					Tok: token.DEFINE,
					Rhs: []ast.Expr{
						&ast.CallExpr{
							Fun: &ast.SelectorExpr{X: receivedIdent, Sel: unmarshalTextIdent},
							Args: []ast.Expr{
								&ast.CallExpr{
									Fun: &ast.ArrayType{Elt: builtinByte},
									Args: []ast.Expr{
										valueStringLits[v.Name],
									},
								},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{X: errIdent, Op: token.NEQ, Y: builtinNil},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Fatalf"}},
									Args: []ast.Expr{
										&ast.BasicLit{Kind: token.STRING, Value: `"expected no error, got %#v"`},
										errIdent,
									},
								},
							},
						},
					},
				},
				&ast.IfStmt{
					Cond: &ast.BinaryExpr{X: expectedIdent, Op: token.NEQ, Y: receivedIdent},
					Body: &ast.BlockStmt{
						List: []ast.Stmt{
							&ast.ExprStmt{
								X: &ast.CallExpr{
									Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Fatalf"}},
									Args: []ast.Expr{
										&ast.BasicLit{Kind: token.STRING, Value: `"expected %#v, got %#v"`},
										expectedIdent,
										receivedIdent,
									},
								},
							},
						},
					},
				},
			}),
		})
	}
	return &ast.ExprStmt{
		X: tRun(t, tType, &ast.BasicLit{Kind: token.STRING, Value: fmt.Sprintf("%q", unmarshalTextIdent.Name)}, testCases),
	}
}

func tRun(t *ast.Ident, tType ast.Expr, name ast.Expr, bodyStmts []ast.Stmt) *ast.CallExpr {
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{X: t, Sel: &ast.Ident{Name: "Run"}},
		Args: []ast.Expr{
			name,
			&ast.FuncLit{
				Type: &ast.FuncType{
					Params: &ast.FieldList{
						List: []*ast.Field{
							{Names: []*ast.Ident{t}, Type: tType},
						},
					},
				},
				Body: &ast.BlockStmt{
					List: bodyStmts,
				},
			},
		},
	}
}

func errorf(fmtImp *ast.Ident, format string, args ...ast.Expr) ast.Expr {
	return &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   fmtImp,
			Sel: &ast.Ident{Name: "Errorf"},
		},
		Args: append(
			[]ast.Expr{
				&ast.BasicLit{
					Kind:  token.STRING,
					Value: fmt.Sprintf("%q", format),
				},
			}, args...),
	}
}

func unexported(name string) string {
	prev := 0
	for pos, char := range name {
		if unicode.IsLower(char) {
			if prev == 0 {
				return strings.ToLower(name[:pos]) + name[pos:]
			}
			return strings.ToLower(name[:prev]) + name[prev:]
		}
		prev = pos
	}
	return strings.ToLower(name)
}

type visitorFunc func(node ast.Node) (w ast.Visitor)

func (fn visitorFunc) Visit(node ast.Node) (w ast.Visitor) {
	return fn(node)
}

var builtinString = &ast.Ident{Name: "string"}
var builtinError = &ast.Ident{Name: "error"}
var builtinByte = &ast.Ident{Name: "byte"}
var builtinNil = &ast.Ident{Name: "nil"}
var builtinNothing = &ast.Ident{Name: "_"}
var builtinPanic = &ast.Ident{Name: "panic"}
var builtinRecover = &ast.Ident{Name: "recover"}
