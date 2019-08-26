package gostringenum_test

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/jussi-kalliokoski/gostringenum/internal/codegen"
	"github.com/jussi-kalliokoski/gostringenum/internal/gostringenum"
)

func TestGenerate(t *testing.T) {
	t.Run("fixtures", func(t *testing.T) {
		const testDataDir = "testdata"
		testCases, err := ioutil.ReadDir(testDataDir)
		mustNotError(t, err)
		for _, testCase := range testCases {
			name := testCase.Name()
			t.Run(name, func(t *testing.T) {
				dirName := path.Join(testDataDir, name)
				fileSet := token.NewFileSet()
				pkgs, err := parser.ParseDir(fileSet, dirName, nil, 0)
				mustNotError(t, err)
				pkgPath, err := codegen.GetImportPath(dirName)
				mustNotError(t, err)
				typeNames := make([]string, 0)
				for _, pkg := range pkgs {
					for _, file := range pkg.Files {
						for _, obj := range file.Scope.Objects {
							typeSpec, isTypeSpec := obj.Decl.(*ast.TypeSpec)
							if !isTypeSpec {
								continue
							}
							ident, isIdentType := typeSpec.Type.(*ast.Ident)
							if !isIdentType || ident.Name != "int" {
								continue
							}
							typeNames = append(typeNames, typeSpec.Name.Name)
						}
					}
				}
				for _, typeName := range typeNames {
					t.Run(typeName, func(t *testing.T) {
						implFilename := path.Join(dirName, strings.ToLower(typeName)+"_encoding.go")
						testFilename := path.Join(dirName, strings.ToLower(typeName)+"_encoding_test.go")
						implData, err := ioutil.ReadFile(implFilename)
						fresh := false
						if os.IsNotExist(err) {
							fresh = true
							err = nil
						}
						mustNotError(t, err)
						testData, err := ioutil.ReadFile(testFilename)
						withTest := false
						if os.IsNotExist(err) {
							withTest = true
							err = nil
						}
						mustNotError(t, err)
						recordSnapshots := fresh || os.Getenv("RECORD_SNAPSHOTS") != ""
						if withTest && !fresh {
							implFile := gostringenum.Generate(fileSet, pkgs, pkgPath, typeName)
							var implBuf bytes.Buffer
							mustNotError(t, format.Node(&implBuf, token.NewFileSet(), implFile))
							if recordSnapshots {
								mustNotError(t, ioutil.WriteFile(implFilename, implBuf.Bytes(), 0644))
								t.Errorf("recorded new snapshot %s", implFilename)
							} else {
								mustBeEqual(t, string(implData), implBuf.String())
							}
						} else {
							implFile, testFile := gostringenum.GenerateWithTestFile(fileSet, pkgs, pkgPath, typeName)
							var implBuf bytes.Buffer
							mustNotError(t, format.Node(&implBuf, token.NewFileSet(), implFile))
							var testBuf bytes.Buffer
							mustNotError(t, format.Node(&testBuf, token.NewFileSet(), testFile))
							if recordSnapshots {
								mustNotError(t, ioutil.WriteFile(implFilename, implBuf.Bytes(), 0644))
								t.Errorf("recorded new snapshot %s", implFilename)
								mustNotError(t, ioutil.WriteFile(testFilename, testBuf.Bytes(), 0644))
								t.Errorf("recorded new snapshot %s", testFilename)
							} else {
								mustBeEqual(t, string(implData), implBuf.String())
								mustBeEqual(t, string(testData), testBuf.String())
							}
						}
					})
				}
			})
		}
	})
}

func mustNotError(tb testing.TB, err error) {
	tb.Helper()
	if err != nil {
		tb.Fatal(err)
	}
}

func mustBeEqual(tb testing.TB, expected, received interface{}) {
	if expected != received {
		_, isLStr := expected.(string)
		_, isRStr := received.(string)
		if isLStr && isRStr {
			tb.Fatalf("expected:\n%s\nreceived:\n%s", expected, received)
		} else {
			tb.Fatalf("expected:\n%#v\nreceived:\n%#v", expected, received)
		}
	}
}
