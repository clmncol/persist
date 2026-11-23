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