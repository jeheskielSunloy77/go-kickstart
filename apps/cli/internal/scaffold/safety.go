package scaffold

import (
	"errors"

	"github.com/jeheskielSunloy77/go-kickstart/apps/cli/internal/validate"
)

func EnsureSafeDestination(path string, allowOverwrite bool) error {
	nonEmpty, err := validate.IsNonEmptyDir(path)
	if err != nil {
		return err
	}
	if nonEmpty && !allowOverwrite {
		return errors.New("destination directory is not empty")
	}
	return nil
}
