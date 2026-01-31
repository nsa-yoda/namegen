package english

import (
	"fmt"
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type englishProfile struct {
	// could hold precomputed frequency tables
}

func init() {
	api.RegisterProfile("english", Profile)
}

func (p englishProfile) Info() map[string]string {
	return map[string]string{
		"name":  "english",
		"notes": "English-like generator with suffixes and realism blending",
	}
}

func (p englishProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
	// Deterministic when cfg.Seed != 0
	r := api.NewRand(cfg)

	vowels := []string{"a", "e", "i", "o", "u"}
	consonants := []string{"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "r", "s", "t", "v", "w", "y", "z"}
	realFragments := []string{"el", "ric", "mar", "an", "beth", "ron", "ly", "ton", "den", "ley", "gar", "wyn"}

	randPick := func(arr []string) string { return arr[r.Intn(len(arr))] }
	genSyl := func(pat string) string {
		var b strings.Builder
		for _, ch := range pat {
			if ch == 'C' {
				b.WriteString(randPick(consonants))
			} else {
				b.WriteString(randPick(vowels))
			}
		}
		return b.String()
	}

	// determine syllables influenced by realism and gender
	numSyl := 1 + r.Intn(3)
	if cfg.Realism > 70 {
		numSyl = 2 + r.Intn(2)
	}
	// gender bias: male slightly more CVC, female more CV/V patterns
	first := ""
	for i := 0; i < numSyl; i++ {
		pat := "CV"
		if cfg.Gender == "male" {
			if r.Intn(100) < 40 {
				pat = "CVC"
			}
		} else if cfg.Gender == "female" {
			if r.Intn(100) < 30 {
				pat = "V"
			}
		} else {
			if r.Intn(100) < 30 {
				pat = "CVC"
			}
		}
		// realism: inject real fragments at high realism
		if cfg.Realism > 60 && r.Intn(100) < cfg.Realism/2 {
			first += randPick(realFragments)
		} else {
			first += genSyl(pat)
		}
	}

	// surname logic
	last := ""
	if cfg.IncludeLast {
		parts := 1 + r.Intn(2)
		for i := 0; i < parts; i++ {
			last += genSyl("CVC")
		}
		// apply family override / common suffixes
		if cfg.Family == "english" || cfg.Family == "" {
			if cfg.Realism > 50 && r.Intn(100) < 60 {
				suffs := []string{"son", "ford", "wood", "well", "shire"}
				last += randPick(suffs)
			}
		}
	}

	// capitalization
	caser := cases.Title(language.English)
	first = caser.String(first)
	last = caser.String(last)

	return api.NameResult{First: first, Last: last}, nil
}

// Profile is the core exported symbol
var Profile englishProfile

func main() {
	// plugin main should be empty; kept to satisfy "package main" build target.
	fmt.Println("english plugin")
}
