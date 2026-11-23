package persist

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

type testConfig struct {
	Name string `json:"name"`
}

func TestPushPull(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	p := Expect[testConfig]("example")

	if err := p.Push(testConfig{Name: "grant"}); err != nil {
		t.Fatal(err)
	}

	got, err := p.Pull()
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "grant" {
		t.Fatalf("pull: got %+v want name=%q", got, "grant")
	}

	if err := p.Push(testConfig{Name: ""}); err != nil {
		t.Fatal(err)
	}

	after, err := p.Pull()
	if err != nil {
		t.Fatal(err)
	}
	if after.Name != "" {
		t.Fatalf("push with zero value: got %+v", after)
	}
}

func TestPushWritesJSON(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	p := Expect[testConfig]("example")
	if err := p.Push(testConfig{Name: "grant"}); err != nil {
		t.Fatal(err)
	}

	path, err := p.path()
	if err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var got testConfig
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Name != "grant" {
		t.Fatalf("file contents: got %+v", got)
	}
}

// Pull on a missing file returns the zero value and no error.
func TestPullMissingFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	p := Expect[testConfig]("absent")

	got, err := p.Pull()
	if err != nil {
		t.Fatalf("pull on missing file should not error: %v", err)
	}
	if got.Name != "" {
		t.Fatalf("pull on missing file: got %+v want zero value", got)
	}
}

// Expect with a default returns the default when the file is missing.
func TestExpectDefault(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	p := Expect[testConfig]("example", testConfig{Name: "default"})

	got, err := p.Pull()
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "default" {
		t.Fatalf("pull with default: got %+v want name=%q", got, "default")
	}

	if err := p.Push(testConfig{Name: "grant"}); err != nil {
		t.Fatal(err)
	}

	got, err = p.Pull()
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "grant" {
		t.Fatalf("pull after push: got %+v want name=%q", got, "grant")
	}
}

// Deleting the file after a Push makes Pull return the constructor default
// again, not the last written value.
func TestPullReturnsDefaultAfterFileDeleted(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	p := Expect[testConfig]("example", testConfig{Name: "default"})

	if err := p.Push(testConfig{Name: "grant"}); err != nil {
		t.Fatal(err)
	}

	got, err := p.Pull()
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "grant" {
		t.Fatalf("pull after push: got %+v want name=%q", got, "grant")
	}

	path, err := p.path()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}

	got, err = p.Pull()
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "default" {
		t.Fatalf("pull after file deleted: got %+v want name=%q", got, "default")
	}
}

// A name containing a path separator is rejected by Pull and Push.
func TestNameWithPathSeparator(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	for _, name := range []string{"bad/name", `bad\name`} {
		p := Expect[testConfig](name)
		if _, err := p.Pull(); err == nil {
			t.Errorf("Pull with name %q should error", name)
		}
		if err := p.Push(testConfig{Name: "grant"}); err == nil {
			t.Errorf("Push with name %q should error", name)
		}
	}
}

// A corrupt file should still error on Pull.
func TestPullCorruptFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	p := Expect[testConfig]("corrupt")
	path, err := p.path()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := p.Pull(); err == nil {
		t.Fatal("pull on corrupt file should error")
	}
}
