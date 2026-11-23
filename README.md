# Persist

_simple, uncomplicated persistent configs in go_

A tiny library for saving a single value to a JSON file in the user's config
directory. Import it with an expected config shape, then `Push`, `Pull`, or
`Pop` a single value for it.

A basic requirement of many apps where viper is too complicated, but it's a
repeated enough pattern to need a lib.

## Installation

```sh
go get github.com/clmncol/persist
```

## Usage

```go
package main

import (
	"fmt"

	"github.com/clmncol/persist"
)

type Config struct {
	Name string `json:"name"`
}

func main() {
	p := persist.Expect[Config]("example")

	// Overwrites the file with this value.
	if err := p.Push(Config{Name: "grant"}); err != nil {
		panic(err)
	}

	// Reads the file back. The file is left untouched.
	cfg, err := p.Pull()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)

	// Reads the value back, then resets the file to the zero value.
	cfg, err = p.Pop()
	if err != nil {
		panic(err)
	}
	fmt.Println(cfg)
}
```

## API

Values are stored as JSON in the `expect` subdirectory of
`os.UserConfigDir()` in a file named `<name>.expect.json`.

### `persist.Expect[T any](name string, defaults ...T) *Persist[T]`

Establishes a new persist handle for the given config name. The file is
created lazily on `Push` or `Pop`.

If a default is provided, it is returned by `Pull`/`Pop` when no file exists
yet — handy for first-run behavior.

```go
p := persist.Expect[Config]("example", Config{Name: "default"})
```

### `(*Persist[T]).Push(v T) error`

Overwrites the persist file with `v`.

### `(*Persist[T]).Pull() (T, error)`

Returns the value stored in the persist file. The file itself is left
untouched. If the file does not exist yet, returns the default (or the zero
value of `T`) with no error.

### `(*Persist[T]).Pop() (T, error)`

Returns the value stored in the persist file, then resets the file to the
zero value of `T` — effectively clearing the persisted value. If the file
does not exist, returns the default (or the zero value of `T`) and writes it
to disk, creating the file.

## Design notes

- **No locking.** Persist does not guarantee exclusive read/write access; it
  is not a replacement for a file-locking store. Viper doesn't guarantee it
  either.
- **JSON** is the only encoding. It's stdlib, human-readable, and
  cross-language — the opposite of "viper is too complicated."