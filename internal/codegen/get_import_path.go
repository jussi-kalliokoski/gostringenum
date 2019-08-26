package codegen

import (
	"go/build"
	"path/filepath"
)

// GetImportPath gets the path for importing given directory.
func GetImportPath(directory string) (string, error) {
	abs, err := filepath.Abs(directory)
	if err != nil {
		return "", err
	}
	pkg, err := build.ImportDir(abs, build.ImportComment)
	if err != nil {
		return "", err
	}
	return pkg.ImportPath, nil
}
