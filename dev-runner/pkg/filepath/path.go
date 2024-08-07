package filepath

import (
	"fmt"
	"os"
	"path/filepath"
)

func IsExists(name string) bool {
	_, err := os.Stat(name)
	if err != nil {
		return false
	}
	return true
}

func IsDir(name string) bool {
	stat, err := os.Stat(name)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func IsFile(name string) bool {
	stat, err := os.Stat(name)
	if err != nil {
		return false
	}
	return !stat.IsDir()
}

func MakePaths(paths ...string) (err error) {
	for _, path := range paths {
		if IsExists(path) && IsDir(path) {
			continue
		}

		if IsExists(path) && !IsDir(path) {
			return fmt.Errorf("path '%s' already exists and it is not directtory", path)
		}

		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("cannot create directory '%s': %w", path, err)
		}
	}
	return nil
}

func MakeFiles(paths ...string) (err error) {
	for _, path := range paths {
		err = MakeFile(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func MakeFile(path string) (err error) {
	if IsFile(path) {
		return nil
	}

	parent := filepath.Dir(path)
	err = MakePaths(parent)
	if err != nil {
		return fmt.Errorf("cannot create parent directory '%s' for file '%s': %w", parent, parent, err)
	}

	var f *os.File
	f, err = os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create empty file '%s': %w", path, err)
	}
	defer func() { _ = f.Close() }()

	return nil
}
