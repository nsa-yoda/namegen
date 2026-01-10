# Namegen


## Build 

Run this from project root

```bash
make build
```

This produces:
•	bin/namegen (main binary)
•	plugins/english.so, plugins/japanese.so, plugins/spanish.so

## Run

```bash
# simple:
./bin/namegen -mode english
# with last name:
./bin/namegen -mode japanese -l
# reverse name order:
./bin/namegen -mode spanish -l -r
# specify plugins dir
./bin/namegen -mode english -plugins ./plugins -l
# deterministically:
./bin/namegen -mode english -l -s 42 --realism 80
# gender:
./bin/namegen -mode english -gender male -l --realism 70
# custom family override (surname rules in plugin may react to this):
./bin/namegen -mode english -family nordic -l --realism 90
```

## How to write your own plugin

	•	Create a new directory plugins/<yourname>/.
	•	package main with a file exporting var Profile <type> that satisfies api.NameProfile.
	•	Build with:
go build -buildmode=plugin -o ../<yourname>.so from inside the plugin dir (or use Makefile).
•	The Profile variable must be the symbol Profile (case-sensitive).

Minimal plugin skeleton:

```go
package main

import (
	"example.com/namegen/api"
)

type myProfile struct{}

func (m myProfile) Info() map[string]string { return map[string]string{"name":"myprofile"} }
func (m myProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
	// generate
	return api.NameResult{First:"X", Last:"Y"}, nil
}

var Profile myProfile
```

## Important caveats & troubleshooting

1.	Go plugin limitations
	•	Go plugins (plugin package) only work on systems that support dlopen for Go plugin binaries (typically Linux). Plugins built on one OS/arch generally won’t work on another.
	•	Plugin and main binary must be built with the same Go toolchain version (e.g., both with Go 1.20). Mismatches can cause plugin.Open errors or type mismatches.

2. Type identity
The api package must be the same module import path and compiled version for both plugin and main. If you change module path or API types, rebuild both.

3. Permissions & Execution
•	Ensure .so files are in the -plugins directory and readable.
•	If a plugin fails to load, the binary falls back to builtin English generator.