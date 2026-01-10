package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/sphireinc/namegen/api"
)

type spanishProfile struct{}

func (p spanishProfile) Info() map[string]string {
	return map[string]string{"name": "spanish", "notes": "Spanish-like generator with -ez / -es suffixes and vowel harmony"}
}

func (p spanishProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
	var r *rand.Rand
	if cfg.Seed != 0 {
		r = rand.New(rand.NewSource(cfg.Seed + time.Now().UnixNano()))
	} else {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

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

var Profile spanishProfile

func main() {
	fmt.Println("spanish plugin")
}
