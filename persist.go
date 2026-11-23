package persist

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Persist is a generic, file-backed persistent config value.
// Values are stored as JSON in os.UserConfigDir() in a file named
// "<name>.expect.json".
type Persist[T any] struct {
	Name  string
	Value T
}

// Expect establishes a new persist handle for the given config name.
// The file "<name>.expect.json" is created in os.UserConfigDir() on Push.
func Expect[T any](name string) *Persist[T] {
	return &Persist[T]{
		Name:  name,
		Value: *new(T),
	}
}

// Push overwrites the persist file with v.
func (p *Persist[T]) Push(v T) error {
	p.Value = v
	return p.write()
}

// Pull returns the value stored in the persist file.
// The file itself is left untouched.
func (p *Persist[T]) Pull() (T, error) {
	var v T
	path, err := p.path()
	if err != nil {
		return v, err
	}
	if err := readFile(path, &v); err != nil {
		return v, err
	}
	p.Value = v
	return v, nil
}

// Pop returns the value stored in the persist file and then resets the file
// to the zero value of T, effectively clearing the persisted value.
func (p *Persist[T]) Pop() (T, error) {
	v, err := p.Pull()
	if err != nil {
		return v, err
	}
	var zero T
	p.Value = zero
	if err := p.write(); err != nil {
		return v, err
	}
	return v, nil
}

func (p *Persist[T]) path() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, p.Name+".expect.json"), nil
}

func (p *Persist[T]) write() error {
	path, err := p.path()
	if err != nil {
		return err
	}
	data, err := json.Marshal(p.Value)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func readFile(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}