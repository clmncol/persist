package persist

import (
	"encoding/json"
	"os"
	"testing"
)

type testConfig struct {
	Name string `json:"name"`
}

func TestPushPullPop(t *testing.T) {
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

	got, err = p.Pop()
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "grant" {
		t.Fatalf("pop: got %+v want name=%q", got, "grant")
	}

	after, err := p.Pull()
	if err != nil {
		t.Fatal(err)
	}
	if after.Name != "" {
		t.Fatalf("pop should reset file to zero value, got %+v", after)
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

// Pop on a missing file returns the zero value and creates the file.
func TestPopMissingFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	p := Expect[testConfig]("absent")

	got, err := p.Pop()
	if err != nil {
		t.Fatalf("pop on missing file should not error: %v", err)
	}
	if got.Name != "" {
		t.Fatalf("pop on missing file: got %+v want zero value", got)
	}

	path, err := p.path()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("pop should have created the file at %s", path)
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

// A corrupt file should still error on Pull.
func TestPullCorruptFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	p := Expect[testConfig]("corrupt")
	path, err := p.path()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("not json"), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := p.Pull(); err == nil {
		t.Fatal("pull on corrupt file should error")
	}
}