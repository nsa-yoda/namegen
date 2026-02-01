[![Test & CI](https://github.com/nsa-yoda/namegen/actions/workflows/makefile.yml/badge.svg)](https://github.com/nsa-yoda/namegen/actions/workflows/makefile.yml)
[![Build](https://github.com/nsa-yoda/namegen/actions/workflows/go.yml/badge.svg)](https://github.com/nsa-yoda/namegen/actions/workflows/go.yml)
[![Codacy Security](https://github.com/nsa-yoda/namegen/actions/workflows/codacy.yml/badge.svg)](https://github.com/nsa-yoda/namegen/actions/workflows/codacy.yml)
[![CodeQL Advanced](https://github.com/nsa-yoda/namegen/actions/workflows/codeql.yml/badge.svg)](https://github.com/nsa-yoda/namegen/actions/workflows/codeql.yml)

# namegen

A fast, deterministic name generator written in Go.  

`namegen` generates realistic-looking or fictional names across many language 
families using **compiled-in profiles** (no runtime plugins). Output can be 
made fully reproducible via seeds, making it suitable for games, testing, 
world-building, and data generation.

---

## Features

- 🌍 Many language profiles (English, Japanese, Spanish, Nordic, Slavic, Tamil, Nahuatl, and more)
- 🎲 Deterministic randomness via seed (`-s`)
- 🎚 Realism control (`-realism 0..100`)
- 🧑 Gender hints (`male`, `female`, `neutral`)
- 👨‍👩‍👧 Optional surnames by default ( turn them on with `-l`)
- 🔁 Reverse order (last name first)
- 🔢 Batch generation (`-c`)
- 🛠 Dev mode whih prints resolved config (`-d`)
- 🚀 Single static binary - no CGO, no `.so` plugins

---

## Project layout

- `cmd/namegen/` – CLI entrypoint (imports all compiled-in profiles)
- `api/` – profile interface, deterministic RNG helpers, shared utilities
- `plugins/<name>/` – profiles (each registers itself via `init()`)

---

## Build

From repo root:

```bash
make build
```

Output:

```bash
$ ls .
bin/namegen
```

## Run

Show help:
```bash
$ ./bin/namegen -h
```

Examples:

```bash
# default (english):
./bin/namegen

# choose a profile:
./bin/namegen -mode japanese

# include last name:
./bin/namegen -mode spanish -l

# reverse (last first):
./bin/namegen -mode spanish -l -r

# generate 10 names:
./bin/namegen -mode english -l -c 10

# deterministic (repeatable):
./bin/namegen -mode english -l -s 42 -realism 80

# gender:
./bin/namegen -mode english -gender female -l -realism 90

# prints a list of all available profiles
./bin/namegen -p 

# dev mode prints the config JSON used:
./bin/namegen -mode japanese -l -c 10 -d
```


## Flags

| Flag                              | Meaning                                                            |
|-----------------------------------|--------------------------------------------------------------------|
| `-mode <name>`                    | Profile/mode name (compiled-in). Defaults to english.              |
| `-l`                              | Include last name                                                  |
| `-r`                              | Reverse output order (last first)                                  |
| `-gender <male, female, neutral>` | Gender hint passed to profile                                      |
| `-family <key>`                   | Optional “family override” (profiles may interpret it differently) |
| `-realism 0...100`                | 0 = fictional phonotactics, 100 = curated/real-looking             |
| `-s <seed>`                       | Seed (0 / omit = random each run)                                  |
| `-c <count>`                      | Number of names to generate                                        |
| `-d`                              | Dev mode: prints config JSON                                       |
| `-p`                              | List all avilable profiles                                         | 

## Available profiles

These are compiled into the binary via blank imports in `cmd/namegen/main.go`

- amharic
- arabic
- aramaic
- baltic
- celtic
- chinese
- english
- farsi
- filipino
- french
- germanic
- greek
- hawaiian
- hebrew
- hindi
- igbo
- indonesian
- italian
- japanese
- kazakh
- korean
- malay
- maori
- nahuatl
- nordic
- portuguese
- samoan
- slavic
- spanish
- swahili
- tamil
- thai
- turkish
- uzbek
- vietnamese
- yoruba

If you run a mode that doesn't exist, the CLI falls back to english.

## How it works

The CLI passes `api.ProfileConfig` to the selected profile.
- Profiles call `api.NewRand(cfg)` which:
  - returns a deterministic RNG when `cfg.Seed != 0`
  - returns a time-seeded RNG when `cfg.Seed == 0`
- Profiles should use `api.PickRand(slice, r)` to select items.

Rule of thumb: If you want reproducibility, always pass `-s <seed>`

## Writing a new profile (compiled-in)

1. Create a folder:

```bash
$ mkdir -p plugins/myprofile
```

2. Add `plugins/myprofile/myprofile.go` and paste this minimal skeleton:

```go 
package myprofile

import (
	"github.com/nsa-yoda/namegen/api"
)

type myProfile struct{}

const PROFILE = "myprofile"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p myProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "What this profile does",
	}
}

func (p myProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
	r := api.NewRand(cfg)

	first := api.PickRand([]string{"Ada", "Grace", "Linus"}, r)
	last := ""
	if cfg.IncludeLast {
		last = api.PickRand([]string{"Lovelace", "Hopper", "Torvalds"}, r)
	}

	return api.NameResult{First: first, Last: last}, nil
}

var Profile myProfile
```

3. Compile it into the binary by adding a blank import to `cmd/namegen/main.go`:

```go 
_ "github.com/nsa-yoda/namegen/plugins/myprofile"
```

4. Build and run:

```bash
$ make build
$ ./bin/namegen -mode myprofile -l -s 123 -c 5
```

## Notes / gotchas

- realism is profile-specific behavior. Every profile implements its own blend between curated names and procedural phonotactics.
- ASCII output by design is the defualt. Some languages normally use diacritics or special punctuation; these profiles intentionally keep output ASCII-friendly.
- Profile not found fallback. If -mode isn't registered, the CLI logs a warning and uses english.


## License

This project is licensed under the MIT license, see LICENSE file

