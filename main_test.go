package main

import (
	"testing"
)

func TestExtractPath(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		want    string
		wantErr bool
	}{
		{
			name: "clone output",
			output: `     clone https://github.com/kotaoue/HandsOn-Docker -> /Users/kotaoue/ghq/github.com/kotaoue/HandsOn-Docker
       git clone --recursive https://github.com/kotaoue/HandsOn-Docker /Users/kotaoue/ghq/github.com/kotaoue/HandsOn-Docker
Cloning into '/Users/kotaoue/ghq/github.com/kotaoue/HandsOn-Docker'...
`,
			want: "/Users/kotaoue/ghq/github.com/kotaoue/HandsOn-Docker",
		},
		{
			name:   "exists output",
			output: "    exists /Users/kotaoue/ghq/github.com/kotaoue/HandsOn-Docker\n",
			want:   "/Users/kotaoue/ghq/github.com/kotaoue/HandsOn-Docker",
		},
		{
			name:    "unrecognized output",
			output:  "something unexpected\n",
			wantErr: true,
		},
		{
			name:    "empty output",
			output:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractPath(tt.output)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("extractPath() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseCloneLine(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		want  string
		found bool
	}{
		{
			name:  "valid clone line",
			line:  "clone https://github.com/kotaoue/HandsOn-Docker -> /Users/kotaoue/ghq/github.com/kotaoue/HandsOn-Docker",
			want:  "/Users/kotaoue/ghq/github.com/kotaoue/HandsOn-Docker",
			found: true,
		},
		{
			name:  "not a clone line",
			line:  "exists /some/path",
			found: false,
		},
		{
			name:  "clone line missing arrow",
			line:  "clone https://github.com/kotaoue/HandsOn-Docker",
			found: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseCloneLine(tt.line)
			if ok != tt.found {
				t.Errorf("parseCloneLine() found = %v, want %v", ok, tt.found)
				return
			}
			if got != tt.want {
				t.Errorf("parseCloneLine() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseExistsLine(t *testing.T) {
	tests := []struct {
		name  string
		line  string
		want  string
		found bool
	}{
		{
			name:  "valid exists line",
			line:  "exists /Users/kotaoue/ghq/github.com/kotaoue/HandsOn-Docker",
			want:  "/Users/kotaoue/ghq/github.com/kotaoue/HandsOn-Docker",
			found: true,
		},
		{
			name:  "not an exists line",
			line:  "clone something -> /path",
			found: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseExistsLine(tt.line)
			if ok != tt.found {
				t.Errorf("parseExistsLine() found = %v, want %v", ok, tt.found)
				return
			}
			if got != tt.want {
				t.Errorf("parseExistsLine() = %q, want %q", got, tt.want)
			}
		})
	}
}
