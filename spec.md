# Persist
_simple, uncomplicated persistent configs in go_

Basically I want this to be something you import with a expected config shape, then either push, pull, or pop a single value for it. 

I think this is a basic requirement of many apps where viper is too complicated but its a repeated enough pattern to need a lib


maybe something like this:
```go
pkg main

import "persist"

type Config struct {
    Name string `json:"name"`
}

func main() {
    persist.Expect("example", Config) // Establishes a new file "example.persist" which is a json encoded Config in XCD_CONFIG

    err := persist.Push(Config{Name: "grant"}) // Overwrites example.presist with this config

    cfg, err := persist.Pull() // Copies persist file struct to cfg

    cfg, err := persist.Pop() // Moves persist file struct to cfg and reinitializes a empty struct

    }

```
