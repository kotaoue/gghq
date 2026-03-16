package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: gghq <repository>")
		os.Exit(1)
	}

	url := os.Args[1]
	path, err := getLocalPath(url)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(path)
}

func getLocalPath(url string) (string, error) {
	var stderr bytes.Buffer
	cmd := exec.Command("ghq", "get", url)
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("ghq get failed: %w\n%s", err, stderr.String())
	}

	return extractPath(stderr.String())
}

func extractPath(output string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(output))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if path, ok := parseCloneLine(line); ok {
			return path, nil
		}
		if path, ok := parseExistsLine(line); ok {
			return path, nil
		}
	}
	return "", fmt.Errorf("could not extract local path from ghq output:\n%s", output)
}

// parseCloneLine parses lines like:
//
//	clone https://github.com/user/repo -> /path/to/repo
func parseCloneLine(line string) (string, bool) {
	if !strings.HasPrefix(line, "clone ") {
		return "", false
	}
	parts := strings.SplitN(line, " -> ", 2)
	if len(parts) != 2 {
		return "", false
	}
	return strings.TrimSpace(parts[1]), true
}

// parseExistsLine parses lines like:
//
//	exists /path/to/repo
func parseExistsLine(line string) (string, bool) {
	const prefix = "exists "
	if !strings.HasPrefix(line, prefix) {
		return "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(line, prefix)), true
}
