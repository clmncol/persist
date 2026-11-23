# Persist

Save one Go value as JSON between runs. Use it for CLI preferences or small configuration structs.

```sh
go get github.com/clmncol/persist
```

```go
type Config struct {
	Profile string
	Verbose bool
}

store := persist.Expect("mycli-config", Config{Profile: "default"})
cfg, err := store.Pull()
if err != nil {
	return err
}

cfg.Profile = "work" // Keep the other settings.
return store.Push(cfg)
```

See [example_test.go](example_test.go) for complete struct and single-value examples.

`Pull` returns the supplied default (or the type’s zero value) when the file is missing. `Push` replaces `expect/<name>.expect.json` under `os.UserConfigDir()` using a temp file and rename. Names can’t contain `/` or `\`. There’s no locking.
