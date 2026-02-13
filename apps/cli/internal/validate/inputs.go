package validate

import (
	"errors"
	"regexp"
	"strings"
)

var modulePathRe = regexp.MustCompile(`^[a-zA-Z0-9_.-]+(/[a-zA-Z0-9_.-]+)+$`)

func ProjectName(name string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("project name is required")
	}
	if strings.ContainsAny(name, `/\`) || strings.ContainsAny(name, "\t\r\n") {
		return errors.New("project name must not contain path separators")
	}
	return nil
}

func ModulePath(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("module path is required")
	}
	if !modulePathRe.MatchString(path) {
		return errors.New("module path must look like domain.com/owner/name")
	}
	return nil
}
