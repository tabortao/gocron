//go:build windows
// +build windows

package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestExecShellWithQuotes(t *testing.T) {
	tempDir := t.TempDir()
	srcDir := filepath.Join(tempDir, "My Documents")
	dstDir := filepath.Join(tempDir, "Backup Folder")
	if err := os.MkdirAll(srcDir, 0o755); err != nil {
		t.Fatalf("mkdir srcDir failed: %v", err)
	}
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		t.Fatalf("mkdir dstDir failed: %v", err)
	}
	srcFile := filepath.Join(srcDir, "report.txt")
	if err := os.WriteFile(srcFile, []byte("hello"), 0o644); err != nil {
		t.Fatalf("write srcFile failed: %v", err)
	}
	dstFile := filepath.Join(dstDir, "report.txt")
	dstFile2 := filepath.Join(dstDir, "report2.txt")
	newDir := filepath.Join(tempDir, "John Doe", "Projects")

	tests := []struct {
		name    string
		command string
		wantErr bool
	}{
		{
			name:    "Simple command without quotes",
			command: "echo hello",
			wantErr: false,
		},
		{
			name:    "Command with double quotes",
			command: `dir "C:\Program Files"`,
			wantErr: false,
		},
		{
			name:    "Copy command with quoted paths",
			command: fmt.Sprintf(`copy "%s" "%s"`, srcFile, dstFile),
			wantErr: false,
		},
		{
			name:    "Mkdir with quoted path",
			command: fmt.Sprintf(`mkdir "%s"`, newDir),
			wantErr: false,
		},
		{
			name:    "Command with HTML entity quotes",
			command: fmt.Sprintf(`copy &quot;%s&quot; &quot;%s&quot;`, srcFile, dstFile2),
			wantErr: false, // HTML实体会被CleanHTMLEntities清理，所以应该成功
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			output, err := ExecShell(ctx, tt.command)
			t.Logf("Command: %s\nOutput: %s\nError: %v", tt.command, output, err)

			if (err != nil) != tt.wantErr {
				t.Errorf("ExecShell() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
