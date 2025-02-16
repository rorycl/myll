package views

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

// fileSystem is a simple wrapper around an fs.FS
type fileSystem struct {
	name string
	fS   fs.FS
}

// newFileSystem returns a new fileSystem from either a path (in
// inDevelopment mode) or provided fs.FS depending.
func newFileSystem(name string, fS fs.FS, path string, inDevelopment bool) (*fileSystem, error) {
	if path == "" {
		return nil, errors.New("path needed for both embedded and inDevelopment modes")
	}
	f := &fileSystem{
		name: name,
	}
	if inDevelopment {
		dirOK := func(d string) bool {
			if _, err := os.Stat(d); os.IsNotExist(err) {
				return false
			}
			return true
		}
		if !dirOK(path) {
			return f, fmt.Errorf("path %s does not exist", path)
		}
		f.fS = os.DirFS(path)
		return f, nil
	}

	var err error
	f.fS, err = fs.Sub(fS, path)
	if err != nil {
		return nil, err
	}
	return f, nil
}
