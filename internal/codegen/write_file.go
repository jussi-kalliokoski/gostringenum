package codegen

import (
	"go/ast"
	"go/format"
	"go/token"
	"os"
)

// WriteASTFile writes provided AST to a file.
func WriteASTFile(filename string, astFile *ast.File) error {
	return writeASTFile(filename, astFile)
}

func writeASTFile(filename string, astFile *ast.File) (returnedErr error) {
	fd, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func() {
		if err := fd.Close(); err != nil && returnedErr == nil {
			returnedErr = err
		}
	}()
	if err := format.Node(fd, token.NewFileSet(), astFile); err != nil {
		return err
	}
	return nil
}
