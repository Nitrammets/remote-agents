package launchd

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// fakeLaunchctl stands in for launchctl: it records its arguments and exits
// with status.
func fakeLaunchctl(t *testing.T, status string) (path, argsFile string) {
	t.Helper()
	dir := t.TempDir()
	path = filepath.Join(dir, "launchctl")
	argsFile = filepath.Join(dir, "args")
	script := "#!/bin/sh\nprintf '%s\\n' \"$@\" > " + argsFile + "\nexit " + status + "\n"
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	return path, argsFile
}

func TestLoadedAsksLaunchdForTheUsersJob(t *testing.T) {
	path, argsFile := fakeLaunchctl(t, "0")
	loaded, err := Launchctl{Path: path, UID: 501}.Loaded(context.Background(), "com.example.job")
	if err != nil || !loaded {
		t.Fatalf("Loaded = %v, %v", loaded, err)
	}
	args, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(args), "print\ngui/501/com.example.job\n"; got != want {
		t.Fatalf("launchctl args = %q, want %q", got, want)
	}
}

func TestNotLoadedWhenLaunchdHasNoSuchJob(t *testing.T) {
	path, _ := fakeLaunchctl(t, "113")
	loaded, err := Launchctl{Path: path, UID: 501}.Loaded(context.Background(), "com.example.job")
	if err != nil || loaded {
		t.Fatalf("Loaded = %v, %v", loaded, err)
	}
}

func TestLoadedFailsOnAnyOtherLaunchctlFailure(t *testing.T) {
	path, _ := fakeLaunchctl(t, "1")
	if _, err := (Launchctl{Path: path, UID: 501}).Loaded(context.Background(), "com.example.job"); err == nil {
		t.Fatal("Loaded succeeded, want an error")
	}
}
