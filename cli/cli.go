package cli

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"github.com/jussi-kalliokoski/gostringenum/internal/codegen"
	"github.com/jussi-kalliokoski/gostringenum/internal/gostringenum"
)

var flagTestFile = flag.Bool("test-file", false, "Generates a _test.go to support the implementation.")

// Run the CLI program.
func Run() {
	// ignore -- flag from go run.
	if len(os.Args) > 1 && os.Args[1] == "--" {
		os.Args = os.Args[:1+copy(os.Args[1:], os.Args[2:])]
	}
	fmt.Println(os.Args)
	flag.Parse()
	args := flag.Args()
	if len(args) != 3 {
		help()
	}
	directory, typeName, filename := args[0], args[1], args[2]
	fileSet := token.NewFileSet()
	pkgs, err := parser.ParseDir(fileSet, directory, nil, 0)
	if err != nil {
		panic(err)
	}
	pkgPath, err := codegen.GetImportPath(directory)
	if err != nil {
		panic(err)
	}
	if *flagTestFile {
		implementationFile, testFile := gostringenum.GenerateWithTestFile(fileSet, pkgs, pkgPath, typeName)
		if err = codegen.WriteASTFile(filename, implementationFile); err != nil {
			panic(err)
		}
		testFilename := strings.TrimSuffix(filename, ".go") + "_test.go"
		if err = codegen.WriteASTFile(testFilename, testFile); err != nil {
			panic(err)
		}
	} else {
		implementationFile := gostringenum.Generate(fileSet, pkgs, pkgPath, typeName)
		if err = codegen.WriteASTFile(filename, implementationFile); err != nil {
			panic(err)
		}
	}
}

func help() {
	fmt.Println(strings.TrimSpace(`
Usage:
  gostringenum [flags] <directory> <name_of_type> <filename>

Flags:
  --help      Show this message.
  --test-file Generates a _test.go to support the implementation.
`))
	os.Exit(1)
}
