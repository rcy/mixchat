package files

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
)

type Files struct {
	directory string
}

func (f *Files) filename(key string) string {
	return f.directory + "/" + key
}

func MustInit(directory string) *Files {
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Initialized file based storage %s\n", directory)
	return &Files{directory: directory}
}

func (f *Files) Put(ctx context.Context, key string, data []byte) error {
	filename := f.filename(key)

	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	if err := os.WriteFile(filename, data, 0666); err != nil {
		return err
	}

	return nil
}

func (f *Files) Get(ctx context.Context, key string) ([]byte, error) {
	return os.ReadFile(f.filename(key))
}

func (f *Files) URI(key string) string {
	return f.filename(key)
}
