package persist_test

import (
	"github.com/clmncol/persist"
)

// A full config example.
func Example() {
	type Config struct {
		Profile string `json:"profile"`
		Verbose bool   `json:"verbose"`
	}

	store := persist.Expect("mycli-config", Config{Profile: "default"})

	cfg, err := store.Pull()
	if err != nil {
		panic(err)
	}

	// Change one preference without losing the others.
	cfg.Profile = "work"
	if err := store.Push(cfg); err != nil {
		panic(err)
	}
}

// A single value example.
func Example_singleValue() {
	// Inside your command:
	p := persist.Expect[string]("mycli-profile", "default")
	profile, err := p.Pull()
	if err != nil {
		panic(err)
	}
	_ = profile // use profile, then save the new selection.
	if err := p.Push("work"); err != nil {
		panic(err)
	}
}
