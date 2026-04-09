package providers

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCursorCLIProvider_buildArgs_permMode(t *testing.T) {
	base := NewCursorCLIProvider("agent")
	cases := []struct {
		name     string
		permMode string
		want     []string
	}{
		{
			name: "default force",
			want: []string{"--print", "--output-format", "json", "--model", "m", "--force", "--trust", "--workspace", "/w"},
		},
		{
			name:     "explicit force",
			permMode: "force",
			want:     []string{"--print", "--output-format", "json", "--model", "m", "--force", "--trust", "--workspace", "/w"},
		},
		{
			name:     "default mode",
			permMode: "default",
			want:     []string{"--print", "--output-format", "json", "--model", "m", "--trust", "--workspace", "/w"},
		},
		{
			name:     "sandbox",
			permMode: "sandbox",
			want:     []string{"--print", "--output-format", "json", "--model", "m", "--force", "--trust", "--sandbox", "enabled", "--workspace", "/w"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := base
			if tc.permMode != "" {
				p = NewCursorCLIProvider("agent", WithCursorCLIPermMode(tc.permMode))
			}
			got := p.buildArgs("m", "/w", false, "", "json")
			if strings.Join(got, " ") != strings.Join(tc.want, " ") {
				t.Fatalf("buildArgs: got %v want %v", got, tc.want)
			}
		})
	}
}

func TestPermModeFromCursorCLISettings(t *testing.T) {
	if got := PermModeFromCursorCLISettings([]byte(`{"perm_mode":"sandbox"}`)); got != "sandbox" {
		t.Fatalf("got %q", got)
	}
	if got := PermModeFromCursorCLISettings(nil); got != "" {
		t.Fatalf("got %q", got)
	}
}

func TestNormalizeCursorSessionLockKey(t *testing.T) {
	if got := normalizeCursorSessionLockKey(""); got != "default" {
		t.Fatalf("empty key: got %q", got)
	}
	if got := normalizeCursorSessionLockKey("   "); got != "default" {
		t.Fatalf("whitespace key: got %q", got)
	}
	if got := normalizeCursorSessionLockKey("session-1"); got != "session-1" {
		t.Fatalf("normal key: got %q", got)
	}
}

func TestCursorCLIProvider_ensureWorkDir_error(t *testing.T) {
	tmp := t.TempDir()
	blocker := filepath.Join(tmp, "blocker")
	if err := os.WriteFile(blocker, []byte("x"), 0600); err != nil {
		t.Fatalf("write blocker file: %v", err)
	}

	p := NewCursorCLIProvider("agent", WithCursorCLIWorkDir(blocker))
	if _, err := p.ensureWorkDir("session-a"); err == nil {
		t.Fatal("expected error when base workdir is a file")
	}
}
