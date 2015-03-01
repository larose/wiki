package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"
)

type GitStorage struct {
	mu            sync.Mutex
	pagesDir      string
	pageExtension string
	repo          *GitRepo
}

type DirtyWorkTree struct {
}

func (d *DirtyWorkTree) Error() string {
	return "Work tree is not clean"
}

func NewGitStorage(path string, pagesDir string, pageExtension string) *GitStorage {
	return &GitStorage{
		pagesDir:      pagesDir,
		pageExtension: pageExtension,
		repo: &GitRepo{
			Path: path,
		},
	}
}

func (s *GitStorage) ensureIsClean() error {
	if clean, err := s.repo.IsClean(); !clean {
		if err != nil {
			return err
		}
		return &DirtyWorkTree{}
	}
	return nil
}

func (s *GitStorage) DeletePage(title string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureIsClean(); err != nil {
		return err
	}

	filename := path.Join(s.pagesDir, title+s.pageExtension)

	if out, err := s.repo.Exec(nil, "rm", filename); err != nil {
		if len(out) > 0 {
			return errors.New(string(out))
		}
		return err
	}

	if out, err := s.repo.Exec(nil, "commit", "--allow-empty", "-m", "Delete "+title); err != nil {
		if len(out) > 0 {
			return errors.New(string(out))
		}
		return err
	}

	return nil
}

func (s *GitStorage) Diff(title string, body string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureIsClean(); err != nil {
		return nil, err
	}

	file, err := ioutil.TempFile(os.TempDir(), "")
	if err != nil {
		return nil, err
	}
	defer os.Remove(file.Name())

	err = ioutil.WriteFile(file.Name(), []byte(body), 0600)
	if err != nil {
		return nil, err
	}

	out, err := ExecGit(strings.NewReader(body), "diff", "--no-index", "--",
		path.Join(s.repo.Path, s.pagesDir, title+s.pageExtension), file.Name())

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitStatus := exitError.Sys().(syscall.WaitStatus).ExitStatus(); exitStatus == 1 {
				return out, nil
			}
		}

		if len(out) > 0 {
			return nil, errors.New(string(out))
		}

		return nil, err
	}

	return out, nil
}

func (s *GitStorage) History(title string) ([]Commit, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureIsClean(); err != nil {
		return nil, err
	}

	out, err := s.repo.Exec(nil, "log", "--name-status", "--", path.Join(s.pagesDir, title+s.pageExtension))

	if err != nil {
		if len(out) > 0 {
			return nil, errors.New(string(out))
		}
		return nil, err
	}

	commits := logParser(string(out), title)
	lines := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")

	revisions := make([]string, len(lines))
	for index, line := range lines {
		parts := strings.SplitN(line, " ", 2)
		revisions[index] = parts[0]
	}

	return commits, nil
}

func (s *GitStorage) Init() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	out, err := ExecGit(nil, "init", s.repo.Path)

	if err != nil {
		if len(out) > 0 {
			return errors.New(string(out))
		}
		return err
	}

	out, err = s.repo.Exec(nil, "show-ref", "--verify", "--quiet", "--", "refs/heads/master")

	if err == nil {
		return nil
	}

	if exitError, ok := err.(*exec.ExitError); ok {
		if exitStatus := exitError.Sys().(syscall.WaitStatus).ExitStatus(); exitStatus == 1 {
			out, err = s.repo.Exec(nil, "commit", "--allow-empty", "-m", "Initial commit")
			if err == nil {
				return nil
			}
			if len(out) > 0 {
				return errors.New(string(out))
			}
			return err
		}
	}
	return err
}

func (s *GitStorage) ListDeletedPages() ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureIsClean(); err != nil {
		return nil, err
	}

	var out []byte
	var err error

	if out, err = s.repo.Exec(nil, "log", "-z", "--name-status", "--oneline", "--", s.pagesDir); err != nil {
		if len(out) > 0 {
			return nil, errors.New(string(out))
		}
		return nil, err
	}

	byteTitles := bytes.Split(bytes.TrimSuffix(out, []byte{0}), []byte{0})

	titles := make(map[string]struct{})
	for _, title := range byteTitles {
		title := string(title)
		if !strings.HasPrefix(title, s.pagesDir+"/") || !strings.HasSuffix(title, s.pageExtension) {
			continue
		}
		title = strings.TrimPrefix(title, s.pagesDir+"/")
		title = strings.TrimSuffix(title, s.pageExtension)
		titles[title] = struct{}{}
	}

	var notDeletedTitles []string

	notDeletedTitles, err = s.listPages()
	if err != nil {
		return nil, err
	}

	for _, title := range notDeletedTitles {
		delete(titles, title)
	}

	_titles := make([]string, 0)

	for title, _ := range titles {
		_titles = append(_titles, title)
	}

	sort.Strings(_titles)

	return _titles, nil
}

func (s *GitStorage) ListPages() ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureIsClean(); err != nil {
		return nil, err
	}
	return s.listPages()
}

func (s *GitStorage) listPages() ([]string, error) {
	out, err := s.repo.Exec(nil, "ls-tree", "-z", "-r", "HEAD", "--name-only", "--", s.pagesDir)
	if err != nil {
		if len(out) > 0 {
			return nil, errors.New(string(out))
		}
		return nil, err
	}

	if len(out) == 0 {
		return nil, nil
	}

	byteTitles := bytes.Split(bytes.TrimSuffix(out, []byte{0}), []byte{0})

	titles := make([]string, 0)
	for _, title := range byteTitles {
		title := string(title)
		title = strings.TrimPrefix(title, s.pagesDir+"/")
		title = strings.TrimSuffix(title, s.pageExtension)
		titles = append(titles, title)
	}
	return titles, nil
}

func (s *GitStorage) PageBody(title string, revision string) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureIsClean(); err != nil {
		return nil, err
	}

	out, err := s.repo.Exec(nil, "cat-file", "-p", revision+":"+path.Join(s.pagesDir, title+s.pageExtension))

	if err != nil {
		if len(out) > 0 {
			return nil, errors.New(string(out))
		}
		return nil, err
	}
	return out, nil
}

func (s *GitStorage) Search(q string) ([]PageSearchResult, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureIsClean(); err != nil {
		return nil, err
	}

	out, err := s.repo.Exec(nil, "grep", "--ignore-case", "-e", q, "--", "HEAD", s.pagesDir)

	if err != nil {
		if msg, ok := err.(*exec.ExitError); ok {
			if exitCode := msg.Sys().(syscall.WaitStatus).ExitStatus(); exitCode == 1 {
				// No results
				return nil, nil
			}
		}

		if len(out) > 0 {
			return nil, errors.New(string(out))
		}
		return nil, err
	}

	results := make(map[string][]string)

	for _, line := range strings.Split(strings.TrimRight(string(out), "\n"), "\n") {
		tokens := strings.SplitN(line, ":", 2)
		pageName := strings.TrimPrefix(strings.TrimRight(tokens[0], pageExtension), pagesDir+"/")

		results[pageName] = append(results[pageName], tokens[1])
	}

	searchResults := make([]PageSearchResult, 0)

	for pageName, lines := range results {
		pageSearchResult := PageSearchResult{
			Title: pageName,
			Lines: lines,
		}
		searchResults = append(searchResults, pageSearchResult)
	}

	return searchResults, nil
}

func (s *GitStorage) SetPageBody(title string, body string, message string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.ensureIsClean(); err != nil {
		return err
	}

	filename := path.Join(s.repo.Path, s.pagesDir, title+s.pageExtension)
	dirName := filepath.Dir(filename)

	if err := os.MkdirAll(dirName, 0770); err != nil {
		return err
	}

	if err := ioutil.WriteFile(filename, []byte(body), 0660); err != nil {
		return err
	}

	out, err := s.repo.Exec(nil, "add", filename)
	if err != nil {
		if len(out) > 0 {
			return errors.New(string(out))
		}
		return err
	}

	out, err = s.repo.Exec(strings.NewReader(message), "commit", "--allow-empty", "-F", "-")
	if err != nil {
		if len(out) > 0 {
			return errors.New(string(out))
		}
		return err
	}

	return nil
}

func logParser(log string, title string) []Commit {
	lines := strings.Split(strings.TrimSuffix(log, "\n"), "\n")
	commits := make([]Commit, 0)

	var commit *Commit
	var message []string = make([]string, 0)

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "commit"):
			if commit != nil {
				commit.Message = strings.Join(message, " ")
				commits = append(commits, *commit)
			}
			message = make([]string, 0)
			commit = &Commit{
				ID: strings.TrimSpace(strings.TrimPrefix(line, "commit")),
			}
		case strings.HasPrefix(line, "Author:"):
		case strings.HasPrefix(line, "Date:"):
			commit.Date = strings.TrimSpace(strings.TrimPrefix(line, "Date:"))
		case strings.HasPrefix(line, "    "):
			message = append(message, strings.TrimPrefix(line, "    "))
		case strings.HasPrefix(line, "A"):
		case strings.HasPrefix(line, "M"):
		case strings.HasPrefix(line, "D"):
			commit.Delete = true
		}
	}

	if commit != nil {
		commit.Message = strings.Join(message, " ")
		commits = append(commits, *commit)
	}
	return commits
}
