package app

import (
	"os"
	"path/filepath"
	"testing"
)

func initTempEnv(t *testing.T, version string) string {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	InitEnv(version)
	return home
}

func TestInitEnvCreatesDirectoriesAndSetsVersion(t *testing.T) {
	home := initTempEnv(t, "1.2.3")
	expectedAppDir := filepath.Join(home, ".gocron")
	if AppDir != expectedAppDir {
		t.Fatalf("expected AppDir %s, got %s", expectedAppDir, AppDir)
	}
	for _, dir := range []string{AppDir, ConfDir, LogDir} {
		if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
			t.Fatalf("expected directory %s to exist", dir)
		}
	}
	expectedVersion := ToNumberVersion("1.2.3")
	if VersionId != expectedVersion {
		t.Fatalf("expected VersionId %d, got %d", expectedVersion, VersionId)
	}
	if Installed {
		t.Fatal("app should not be marked installed without lock file")
	}
}

func TestCreateInstallLockAndIsInstalled(t *testing.T) {
	initTempEnv(t, "1.0.0")
	lockPath := filepath.Join(ConfDir, "install.lock")
	if IsInstalled() {
		t.Fatal("expected not installed before lock file exists")
	}
	if err := CreateInstallLock(); err != nil {
		t.Fatalf("CreateInstallLock failed: %v", err)
	}
	if _, err := os.Stat(lockPath); err != nil {
		t.Fatalf("install lock not created: %v", err)
	}
	if !IsInstalled() {
		t.Fatal("expected installed after lock creation")
	}
}

func TestUpdateVersionFileAndGetCurrentVersionId(t *testing.T) {
	initTempEnv(t, "1.0.0")
	VersionId = 789
	UpdateVersionFile()
	id := GetCurrentVersionId()
	if id != 789 {
		t.Fatalf("expected version id 789, got %d", id)
	}
}

func TestGetCurrentVersionIdWhenMissing(t *testing.T) {
	initTempEnv(t, "2.0.0")
	if id := GetCurrentVersionId(); id != 0 {
		t.Fatalf("expected 0 when version file missing, got %d", id)
	}
}

func TestToNumberVersion(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"v1.2.3", 123},
		{"1.2", 120},
		{"2.0.10", 2010},
	}
	for _, tt := range tests {
		got := ToNumberVersion(tt.input)
		if got != tt.want {
			t.Fatalf("ToNumberVersion(%s) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestCreateDirIfNotExists(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "dir")
	createDirIfNotExists(dir)
	if fi, err := os.Stat(dir); err != nil || !fi.IsDir() {
		t.Fatalf("expected directory %s to exist", dir)
	}
}
