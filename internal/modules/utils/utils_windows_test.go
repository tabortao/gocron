//go:build windows
// +build windows

package utils

import (
	"context"
	"testing"
	"time"
)

func TestExecShellWithQuotes(t *testing.T) {
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
			command: `copy "C:\My Documents\report.docx" "D:\Backup"`,
			wantErr: false,
		},
		{
			name:    "Mkdir with quoted path",
			command: `mkdir "C:\Users\John Doe\Projects"`,
			wantErr: false,
		},
		{
			name:    "Command with HTML entity quotes",
			command: `copy &quot;C:\My Documents\report.docx&quot; &quot;D:\Backup&quot;`,
			wantErr: true, // This should fail because &quot; is not valid in cmd
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
