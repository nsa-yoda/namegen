package spanish

import (
	"fmt"
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type spanishProfile struct{}

func init() {
	api.RegisterProfile("spanish", Profile)
}

func (p spanishProfile) Info() map[string]string {
	return map[string]string{
		"name":  "spanish",
		"notes": "Spanish-like generator with -ez / -es suffixes and vowel harmony",
	}
}

func (p spanishProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
	// Deterministic when cfg.Seed != 0
	r := api.NewRand(cfg)

	vowels := []string{"a", "e", "i", "o", "u"}
	consonants := []string{"b", "c", "d", "f", "g", "h", "j", "l", "m", "n", "p", "r", "s", "t", "v", "z"}
	frags := []string{"mar", "ana", "ro", "el", "carlos", "iza", "al"}
	randPick := func(arr []string) string { return arr[r.Intn(len(arr))] }

	syl := func() string {
		if r.Intn(100) < 60 {
			return randPick(consonants) + randPick(vowels)
		}
		return randPick(vowels) + randPick(consonants)
	}

	numSyl := 2 + r.Intn(2)
	if cfg.Realism < 40 {
		numSyl = 1 + r.Intn(3)
	}

	first := ""
	// blend fragments according to realism
	for i := 0; i < numSyl; i++ {
		if cfg.Realism > 60 && r.Intn(100) < cfg.Realism/2 {
			first += randPick(frags)
		} else {
			first += syl()
		}
	}

	last := ""
	if cfg.IncludeLast {
		for i := 0; i < 1+r.Intn(2); i++ {
			last += syl()
		}
		if cfg.Family == "spanish" || cfg.Family == "" {
			if cfg.Realism > 50 && r.Intn(100) < 60 {
				suffs := []string{"ez", "es", "ado", "ias"}
				last += randPick(suffs)
			}
		}
	}

	first = strings.Title(first)
	last = strings.Title(last)
	return api.NameResult{First: first, Last: last}, nil
}

// Profile is the core exported symbol
var Profile spanishProfile

func main() {
	fmt.Println("spanish plugin")
}
