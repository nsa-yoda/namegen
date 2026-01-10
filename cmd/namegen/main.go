package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"plugin"
	"time"

	"github.com/nsa-yoda/namegen/api"
)

const defaultPluginDir = "./plugins" // runtime plugins path

// loadPlugin loads the .so plugin file and extracts the Profile symbol.
func loadPlugin(path string) (api.NameProfile, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}
	sym, err := p.Lookup("Profile")
	if err != nil {
		return nil, err
	}
	prof, ok := sym.(api.NameProfile)
	if !ok {
		return nil, fmt.Errorf("symbol Profile has wrong type in %s", path)
	}
	return prof, nil
}

// findPluginFile tries to map a mode/plugin name to a .so file in pluginDir
func findPluginFile(pluginDir, mode string) (string, error) {
	candidate := filepath.Join(pluginDir, mode+".so")
	if _, err := os.Stat(candidate); err == nil {
		return candidate, nil
	}
	// fallback: try mode as full path
	if _, err := os.Stat(mode); err == nil {
		return mode, nil
	}
	return "", fmt.Errorf("plugin not found: %s (checked %s)", mode, candidate)
}

// fallback builtin english generator (used when plugin unavailable)
type builtinEnglish struct{}

func (b builtinEnglish) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
	r := rand.New(rand.NewSource(cfg.Seed + time.Now().UnixNano()))
	vowels := []rune("aeiou")
	consonants := []rune("bcdfghjklmnpqrstvwxyz")
	pick := func(rs []rune) rune { return rs[r.Intn(len(rs))] }
	makeSyll := func(pat string) string {
		s := make([]rune, 0, len(pat))
		for _, ch := range pat {
			if ch == 'C' {
				s = append(s, pick(consonants))
			} else {
				s = append(s, pick(vowels))
			}
		}
		return string(s)
	}
	// simple heuristics based on realism
	numSyl := 1 + r.Intn(3)
	if cfg.Realism > 70 {
		numSyl = 2 + r.Intn(2)
	}
	first := ""
	for i := 0; i < numSyl; i++ {
		first += makeSyll("CV")
	}
	last := ""
	if cfg.IncludeLast {
		last = ""
		for i := 0; i < 1+r.Intn(2); i++ {
			last += makeSyll("CVC")
		}
		// apply english-ish suffixes if realism high
		if cfg.Realism > 60 {
			suffs := []string{"son", "ford", "wood", "well"}
			last += suffs[r.Intn(len(suffs))]
		}
	}
	// capitalize
	if len(first) > 0 {
		first = string([]rune(first[:1])) + first[1:]
	}
	if len(last) > 0 {
		last = string([]rune(last[:1])) + last[1:]
	}
	return api.NameResult{First: first, Last: last}, nil
}

func (b builtinEnglish) Info() map[string]string {
	return map[string]string{"name": "builtin-english", "notes": "Fallback builtin English-like generator"}
}

func main() {
	// CLI flags
	mode := flag.String("mode", "english", "Mode/plugin name (plugin file without .so). If plugin not found, builtin fallback used.")
	pluginDir := flag.String("plugins", defaultPluginDir, "Directory with plugin .so files")
	includeLast := flag.Bool("l", false, "Include last name")
	reverse := flag.Bool("r", false, "Reverse order (last first)")
	gender := flag.String("gender", "neutral", "Gender: male|female|neutral")
	family := flag.String("family", "", "Family override for surname rules (e.g., japan, nordic, spanish)")
	realism := flag.Int("realism", 50, "Realism 0..100 (0 fictional phonotactics, 100 real-looking names)")
	seed := flag.Int64("s", 0, "Seed (0 random)")
	flag.Parse()

	// Seed
	if *seed == 0 {
		rand.Seed(time.Now().UnixNano())
	} else {
		rand.Seed(*seed)
	}

	cfg := api.ProfileConfig{
		Seed:        *seed,
		Realism:     *realism,
		Gender:      *gender,
		Family:      *family,
		IncludeLast: *includeLast,
		Reverse:     *reverse,
	}

	// Attempt to load plugin
	var profile api.NameProfile
	pluginFile, err := findPluginFile(*pluginDir, *mode)
	if err == nil {
		prof, errPlugin := loadPlugin(pluginFile)
		if errPlugin != nil {
			log.Printf("failed to load plugin %s: %v — using builtin fallback\n", pluginFile, errPlugin)
			profile = builtinEnglish{}
		} else {
			profile = prof
			log.Printf("loaded plugin: %s\n", pluginFile)
		}
	} else {
		log.Printf("plugin not found for mode %s: %v — using builtin fallback\n", *mode, err)
		profile = builtinEnglish{}
	}

	// Generate and print
	res, err := profile.Generate(cfg)
	if err != nil {
		log.Fatalf("generate failed: %v", err)
	}

	if cfg.IncludeLast {
		if cfg.Reverse {
			fmt.Printf("%s %s\n", res.Last, res.First)
		} else {
			fmt.Printf("%s %s\n", res.First, res.Last)
		}
	} else {
		fmt.Println(res.First)
	}
}
