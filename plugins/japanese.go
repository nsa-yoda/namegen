package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/sphireinc/namegen/api"
)

type japaneseProfile struct{}

func (p japaneseProfile) Info() map[string]string {
	return map[string]string{"name": "japanese", "notes": "CV-heavy Japanese-like generator; includes typical suffixes"}
}

func (p japaneseProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
	// deterministic seeding
	var r *rand.Rand
	if cfg.Seed != 0 {
		r = rand.New(rand.NewSource(cfg.Seed + time.Now().UnixNano()))
	} else {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	vowels := []string{"a", "i", "u", "e", "o"}
	consonants := []string{"k", "s", "t", "n", "h", "m", "y", "r", "w", "g", "z"}
	clusters := []string{"ky", "sh", "ch", "ny", "ry"}
	randPick := func(arr []string) string { return arr[r.Intn(len(arr))] }

	genCV := func() string { return randPick(consonants) + randPick(vowels) }

	// Japanese tends to 2-4 syllables
	numSyl := 2 + r.Intn(3)
	if cfg.Realism < 30 {
		numSyl = 1 + r.Intn(3)
	}

	first := ""
	for i := 0; i < numSyl; i++ {
		if cfg.Realism > 70 && r.Intn(100) < 20 {
			first += randPick(clusters) + randPick(vowels)
		} else {
			first += genCV()
		}
	}

	last := ""
	if cfg.IncludeLast {
		// last names in Japanese often 2-3 moras; use constructed parts
		for i := 0; i < 2; i++ {
			last += genCV()
		}
		if cfg.Family == "japan" || cfg.Family == "" {
			if cfg.Realism > 60 && r.Intn(100) < 40 {
				suffs := []string{"moto", "suke", "naga", "shita", "gawa", "yama"}
				last += randPick(suffs)
			}
		}
	}

	first = strings.Title(first)
	last = strings.Title(last)

	return api.NameResult{First: first, Last: last}, nil
}

var Profile japaneseProfile

func main() {
	fmt.Println("japanese plugin")
}
