package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

var execCommand = exec.Command

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: gghq <repository>")
	}

	path, err := getLocalPath(args[0])
	if err != nil {
		return err
	}

	fmt.Println(path)
	return nil
}

func getLocalPath(repository string) (string, error) {
	out, err := execCommand("ghq", "get", repository).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ghq get failed: %w\n%s", err, out)
	}

	return listPath(repository)
}

// listPath uses `ghq list --full-path` to retrieve the local path for the given repository.
func listPath(repository string) (string, error) {
	out, err := execCommand("ghq", "list", "--full-path", repoPath(repository)).Output()
	if err != nil {
		return "", fmt.Errorf("ghq list failed: %w", err)
	}

	path := strings.TrimSpace(string(out))
	if path == "" {
		return "", fmt.Errorf("could not find local path for %s", repository)
	}

	// Return the first match in case there are multiple results.
	path, _, _ = strings.Cut(path, "\n")
	return path, nil
}

// repoPath extracts the "host/user/repo" identifier from a repository URL or short form,
// which is the query format expected by `ghq list`.
func repoPath(repository string) string {
	// https:// or http:// URL
	if strings.HasPrefix(repository, "https://") || strings.HasPrefix(repository, "http://") {
		if u, err := url.Parse(repository); err == nil {
			return u.Host + strings.TrimSuffix(u.Path, ".git")
		}
	}

	// git@host:user/repo.git (SCP-style SSH)
	if i := strings.Index(repository, "@"); i != -1 {
		rest := repository[i+1:]
		host, path, _ := strings.Cut(rest, ":")
		return host + "/" + strings.TrimSuffix(path, ".git")
	}

	// git://host/user/repo.git
	if strings.HasPrefix(repository, "git://") {
		return strings.TrimSuffix(strings.TrimPrefix(repository, "git://"), ".git")
	}

	// Already in short form: user/repo or host/user/repo
	return strings.TrimSuffix(repository, ".git")
}
