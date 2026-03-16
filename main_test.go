package main

import (
	"testing"
)

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

