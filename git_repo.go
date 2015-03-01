package main

import (
	"io"
	"os/exec"
)

type GitRepo struct {
	Path string
}

func (r *GitRepo) Exec(stdin io.Reader, args ...string) ([]byte, error) {
	args = append([]string{"--git-dir=" + r.Path + "/.git", "--work-tree=" + r.Path}, args...)
	return ExecGit(stdin, args...)
}

func (r *GitRepo) IsClean() (bool, error) {
	out, err := r.Exec(nil, "status", "--porcelain")
	if err != nil {
		return false, err
	}

	return len(out) == 0, nil
}

func ExecGit(stdin io.Reader, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Stdin = stdin
	return cmd.CombinedOutput()
}
