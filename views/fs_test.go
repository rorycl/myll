package views

import (
	"embed"
	_ "embed"
	"fmt"
	"io/fs"
	"testing"

	"testing/fstest"
)

//go:embed templates/*html
var tpls embed.FS

func TestFS(t *testing.T) {

	tests := []struct {
		name          string
		fS            fs.FS
		path          string
		inDevelopment bool
	}{
		{
			name:          "t1",
			fS:            tpls,
			path:          "templates",
			inDevelopment: false,
		},
		{
			name:          "t2",
			fS:            nil,
			path:          "templates",
			inDevelopment: true,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("test_%d", i), func(t *testing.T) {
			f, err := newFileSystem(tt.name, tt.fS, tt.path, tt.inDevelopment)
			if err != nil {
				t.Fatal(err)
			}
			/*
				de, err := fs.ReadDir(f.fS, ".")
				if err != nil {
					t.Fatal(err)
				}
				fmt.Println(de)
			*/
			if err := fstest.TestFS(f.fS, "home.html"); err != nil {
				t.Fatal(err)
			}

		})
	}
}
