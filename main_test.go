package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"
)

// TestHelperProcess is not a real test – it is invoked as a subprocess by
// sequentialFakeExecCommand to simulate an external command.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	fmt.Fprint(os.Stdout, os.Getenv("GGHQ_MOCK_STDOUT"))
	code, _ := strconv.Atoi(os.Getenv("GGHQ_MOCK_EXIT"))
	os.Exit(code)
}

type fakeResponse struct {
	stdout   string
	exitCode int
}

// sequentialFakeExecCommand returns a replacement for execCommand that returns
// successive responses from the provided list on each call.
func sequentialFakeExecCommand(responses ...fakeResponse) func(string, ...string) *exec.Cmd {
	calls := 0
	return func(name string, args ...string) *exec.Cmd {
		idx := calls
		calls++
		var resp fakeResponse
		if idx < len(responses) {
			resp = responses[idx]
		}
		cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--")
		cmd.Env = append(os.Environ(),
			"GO_WANT_HELPER_PROCESS=1",
			"GGHQ_MOCK_STDOUT="+resp.stdout,
			"GGHQ_MOCK_EXIT="+strconv.Itoa(resp.exitCode),
		)
		return cmd
	}
}

func TestRepoPath(t *testing.T) {
	tests := []struct {
		name       string
		repository string
		want       string
	}{
		{
			name:       "https URL",
			repository: "https://github.com/kotaoue/HandsOn-Docker",
			want:       "github.com/kotaoue/HandsOn-Docker",
		},
		{
			name:       "https URL with .git suffix",
			repository: "https://github.com/kotaoue/HandsOn-Docker.git",
			want:       "github.com/kotaoue/HandsOn-Docker",
		},
		{
			name:       "http URL",
			repository: "http://github.com/kotaoue/HandsOn-Docker",
			want:       "github.com/kotaoue/HandsOn-Docker",
		},
		{
			name:       "git@ SSH URL",
			repository: "git@github.com:kotaoue/HandsOn-Docker.git",
			want:       "github.com/kotaoue/HandsOn-Docker",
		},
		{
			name:       "git@ SSH URL without .git",
			repository: "git@github.com:kotaoue/HandsOn-Docker",
			want:       "github.com/kotaoue/HandsOn-Docker",
		},
		{
			name:       "git:// URL",
			repository: "git://github.com/kotaoue/HandsOn-Docker.git",
			want:       "github.com/kotaoue/HandsOn-Docker",
		},
		{
			name:       "short form user/repo",
			repository: "kotaoue/HandsOn-Docker",
			want:       "kotaoue/HandsOn-Docker",
		},
		{
			name:       "short form with .git",
			repository: "kotaoue/HandsOn-Docker.git",
			want:       "kotaoue/HandsOn-Docker",
		},
		{
			name:       "empty string",
			repository: "",
			want:       "",
		},
		{
			name:       "repo name only",
			repository: "HandsOn-Docker",
			want:       "HandsOn-Docker",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := repoPath(tt.repository)
			if got != tt.want {
				t.Errorf("repoPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestListPath(t *testing.T) {
	tests := []struct {
		name     string
		repo     string
		stdout   string
		exitCode int
		want     string
		wantErr  bool
	}{
		{
			name:     "success",
			repo:     "kotaoue/HandsOn-Docker",
			stdout:   "/home/user/ghq/github.com/kotaoue/HandsOn-Docker",
			exitCode: 0,
			want:     "/home/user/ghq/github.com/kotaoue/HandsOn-Docker",
		},
		{
			name:     "multiple results returns first",
			repo:     "kotaoue/HandsOn-Docker",
			stdout:   "/path/first\n/path/second",
			exitCode: 0,
			want:     "/path/first",
		},
		{
			name:     "empty output",
			repo:     "kotaoue/HandsOn-Docker",
			stdout:   "",
			exitCode: 0,
			wantErr:  true,
		},
		{
			name:     "command failure",
			repo:     "kotaoue/HandsOn-Docker",
			stdout:   "",
			exitCode: 1,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execCommand = sequentialFakeExecCommand(fakeResponse{tt.stdout, tt.exitCode})
			defer func() { execCommand = exec.Command }()

			got, err := listPath(tt.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("listPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("listPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGetLocalPath(t *testing.T) {
	tests := []struct {
		name      string
		repo      string
		responses []fakeResponse
		want      string
		wantErr   bool
	}{
		{
			name: "success",
			repo: "kotaoue/HandsOn-Docker",
			responses: []fakeResponse{
				{stdout: "", exitCode: 0},
				{stdout: "/home/user/ghq/github.com/kotaoue/HandsOn-Docker", exitCode: 0},
			},
			want: "/home/user/ghq/github.com/kotaoue/HandsOn-Docker",
		},
		{
			name: "ghq get fails",
			repo: "kotaoue/HandsOn-Docker",
			responses: []fakeResponse{
				{stdout: "error output", exitCode: 1},
			},
			wantErr: true,
		},
		{
			name: "ghq list fails",
			repo: "kotaoue/HandsOn-Docker",
			responses: []fakeResponse{
				{stdout: "", exitCode: 0},
				{stdout: "", exitCode: 1},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execCommand = sequentialFakeExecCommand(tt.responses...)
			defer func() { execCommand = exec.Command }()

			got, err := getLocalPath(tt.repo)
			if (err != nil) != tt.wantErr {
				t.Errorf("getLocalPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("getLocalPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		responses []fakeResponse
		wantErr   bool
	}{
		{
			name:    "no args returns usage error",
			args:    []string{},
			wantErr: true,
		},
		{
			name: "success",
			args: []string{"kotaoue/HandsOn-Docker"},
			responses: []fakeResponse{
				{stdout: "", exitCode: 0},
				{stdout: "/home/user/ghq/github.com/kotaoue/HandsOn-Docker", exitCode: 0},
			},
			wantErr: false,
		},
		{
			name: "getLocalPath error propagates",
			args: []string{"kotaoue/HandsOn-Docker"},
			responses: []fakeResponse{
				{stdout: "error", exitCode: 1},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.responses) > 0 {
				execCommand = sequentialFakeExecCommand(tt.responses...)
				defer func() { execCommand = exec.Command }()
			}

			err := run(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

