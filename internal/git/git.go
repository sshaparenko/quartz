package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var ErrUnrecognizablePath = errors.New("unrecognizable path")
var ErrInitializedGit = errors.New("your vault already has a git repo")

// Init initiazizes Git repo by provided path
func Init(path string) (err error) {

	var dirPath string

	defer func() {
		if err != nil {
			err = fmt.Errorf("in git.Init: %w", err)
		}
	}()

	switch {
	case strings.HasPrefix(path, ".") && len(path) == 1:
		dirPath, err = currentDirPath()
		if err != nil {
			return err
		}
	case strings.HasPrefix(path, "./"):
		dirPath, err = relativePath(path)
		if err != nil {
			return err
		}
	case strings.HasPrefix(path, "/"):
		dirPath, err = obsolutePath(path)
		if err != nil {
			return err
		}
	default:
		return ErrUnrecognizablePath
	}

	if err = initRepo(dirPath); err != nil {
		return err
	}
	return nil
}

// currentDirPath returns path of the current directory
// if this directory is an Obsidian vault
func currentDirPath() (string, error) {
	currentDir, _ := os.Getwd()
	filePath := currentDir + "/.obsidian"

	if err := checkVault(filePath); err != nil {
		return "", err
	}
	return currentDir, nil
}

// relativePath returns path of the directory specified
// by relative path if this directory is an Obsidian vault
func relativePath(path string) (string, error) {
	currentDir, _ := os.Getwd()
	vaultPath := currentDir + path[1:]
	filePath := vaultPath + ".obsidian"

	if err := checkVault(filePath); err != nil {
		return "", err
	}
	return path, nil
}

// relativePath returns path of the directory specified
// by obsolute path if this directory is an Obsidian vault
func obsolutePath(path string) (string, error) {
	filePath := path + ".obsidian"

	if err := checkVault(filePath); err != nil {
		return "", err
	}

	return path, nil
}

// checkVault cheks if directory specified by path
// is an Obsidian vault
func checkVault(path string) error {
	path = filepath.Clean(path)
	f, err := os.Open(path)
	defer func() {
		_ = f.Close()
	}()

	if err != nil {
		return errors.New("in git.checkVault: specified path is not an obsidian vault")
	}

	return nil
}

func checkRepo(path string) (err error) {
	path = filepath.Join(filepath.Clean(path), "/.git/info")
	f, err := os.Open(path)
	defer func() {
		_ = f.Close()
	}()

	if err != nil {
		return nil
	}

	return ErrInitializedGit
}

// initRepo initizalizes a Git repository by
// specified path
func initRepo(path string) error {

	if err := checkRepo(path); err != nil {
		return err
	}

	cmd := exec.Command("git", "init", path)
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Println("quartz initilized new git repository")
	return nil
}
