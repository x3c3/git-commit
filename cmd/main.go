package main

import (
	"fmt"
	"os"

	"bufio"

	"github.com/pkg/errors"
	"livingit.de/code/git-commit/cmd/helper"
	"livingit.de/code/git-commit/cmd/methods"
	"livingit.de/code/git-commit/hook"
)

func main() {
	methods.PrintVersion()

	if len(os.Args) < 2 {
		methods.Help()
		os.Exit(0)
		return
	}

	if os.Args[1] == "install" {
		os.Exit(methods.InstallHook())
		return
	}

	if os.Args[1] == "uninstall" {
		os.Exit(methods.UninstallHook())
		return
	}

	validationResult := validateInput()
	if 0 != validationResult {
		os.Exit(validationResult)
		return
	}

	commitMessageFile := os.Args[1]
	config, err := hook.NewForVersion(commitMessageFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		return
	}

	commitFileContent, err := loadCommitMessageFile(commitMessageFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		return
	}

	ok, err := config.Validate(commitFileContent)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		return
	}

	if !ok {
		os.Exit(1)
		return
	}

	os.Exit(0)
}

// validateInput checks program oarameters when running
// as a hook
func validateInput() int {
	commitMessageFile := os.Args[1]
	if commitMessageFile == "" {
		fmt.Fprintln(os.Stderr, errors.New("no commit message file passed as parameter 1"))
		return 1
	}

	if !helper.FileExists(commitMessageFile) {
		fmt.Fprintln(os.Stderr, errors.New("passed commit message file not found"))
		return 1
	}
	return 0
}

// loadCommitMessageFile reads the commit message and returns it as a
// array containing the lines
func loadCommitMessageFile(commitMessageFile string) ([]string, error) {
	file, err := os.Open(commitMessageFile)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	commitFileContent := make([]string, 0)
	for scanner.Scan() {
		commitFileContent = append(commitFileContent, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return commitFileContent, nil
}
