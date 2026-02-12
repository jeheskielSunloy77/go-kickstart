package validate

import (
	"errors"
	"os"
	"path/filepath"
)

func ResolveDestination(arg string) (string, error) {
	if arg == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		return cwd, nil
	}
	abs, err := filepath.Abs(arg)
	if err != nil {
		return "", err
	}
	return abs, nil
}

func ResolveProjectDestination(baseArg, projectName string) (string, error) {
	base, err := ResolveDestination(baseArg)
	if err != nil {
		return "", err
	}
	base = filepath.Clean(base)
	if filepath.Base(base) == projectName {
		return base, nil
	}
	return filepath.Join(base, projectName), nil
}

func IsNonEmptyDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}
		return false, err
	}
	if !info.IsDir() {
		return false, errors.New("destination exists and is not a directory")
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}
	return len(entries) > 0, nil
}
