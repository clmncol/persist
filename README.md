# Persist

Save one Go value as JSON between runs. Use it for CLI preferences, like remembering the last selected profile.

```sh
go get github.com/clmncol/persist
```

```go
import "github.com/clmncol/persist"

// Inside your command:
p := persist.Expect[string]("mycli-profile", "default")
profile, err := p.Pull()
if err != nil {
	panic(err)
}
// Use profile, then save the new selection.
if err := p.Push("work"); err != nil {
	panic(err)
}
```

`Pull` returns the supplied default (or the type’s zero value) if the file is missing. `Push` replaces `expect/<name>.expect.json` under `os.UserConfigDir()` via a temp file and rename. Names can’t contain `/` or `\`. No locking.
