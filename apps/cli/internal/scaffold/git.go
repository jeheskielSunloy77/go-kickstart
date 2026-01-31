package scaffold

import (
	"os/exec"
)

func InitGitRepo(path string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = path
	if err := cmd.Run(); err != nil {
		return err
	}
	add := exec.Command("git", "add", ".")
	add.Dir = path
	if err := add.Run(); err != nil {
		return err
	}
	commit := exec.Command("git", "commit", "-m", "chore: initial commit")
	commit.Dir = path
	if err := commit.Run(); err != nil {
		return err
	}
	return nil
}
