package french

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type frenchProfile struct{}

const PROFILE = "french"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p frenchProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "French names (ASCII): realism blends curated lists with procedural syllables; deterministic with seed",
	}
}

// Curated given names (ASCII only; accents removed).
var firstMale = []string{
	"Jean", "Pierre", "Louis", "Michel", "Andre", "Paul", "Jacques", "Henri", "Luc", "Thomas",
	"Antoine", "Nicolas", "Julien", "Mathieu", "Hugo", "Arthur", "Guillaume", "Alexandre", "Victor", "Sebastien",
	"Maxime", "Theo", "Romain", "Damien", "Laurent", "Olivier", "Francois", "Benjamin", "Gabriel", "Etienne",
}

var firstFemale = []string{
	"Marie", "Anne", "Sophie", "Camille", "Julie", "Claire", "Isabelle", "Nathalie", "Helene", "Pauline",
	"Charlotte", "Emma", "Lea", "Manon", "Chloe", "Sarah", "Alice", "Juliette", "Celine", "Amandine",
	"Elise", "Margaux", "Aurelie", "Valerie", "Mathilde", "Audrey", "Lucie", "Noemie", "Ines", "Gabrielle",
}

var firstNeutral = []string{
	"Camille", "Alex", "Charlie", "Noa", "Sacha", "Lou", "Morgan", "Remy", "Jules", "Andrea",
}

// Curated surnames (ASCII; accents removed).
var lastNames = []string{
	"Martin", "Bernard", "Thomas", "Petit", "Robert", "Richard", "Durand", "Dubois", "Moreau", "Laurent",
	"Simon", "Michel", "Lefevre", "Garcia", "Roux", "David", "Bertrand", "Morel", "Fournier", "Girard",
	"Bonnet", "Dupont", "Lambert", "Fontaine", "Rousseau", "Vincent", "Muller", "Leroy", "Faure", "Andre",
}

var vowels = []string{"a", "e", "i", "o", "u", "y", "ai", "au", "ei", "eu", "ou", "oi", "ui"}
var onsets = []string{
	"b", "c", "d", "f", "g", "h", "j", "l", "m", "n", "p", "r", "s", "t", "v",
	"br", "bl", "cr", "cl", "dr", "fr", "fl", "gr", "gl", "pr", "pl", "tr",
	"ch", "gn", "ph", "qu",
	"", "", // allow vowel-start sometimes
}
var codas = []string{"", "", "", "n", "m", "r", "s", "t", "l", "d", "x"}

var givenEndingsMale = []string{"", "", "", "e", "el", "en", "ier", "ois", "on", "in"}
var givenEndingsFemale = []string{"", "", "", "e", "elle", "ine", "ette", "ane", "ie", "a"}
var givenEndingsNeutral = []string{"", "", "", "e", "i", "en"}

var surnameEndings = []string{"", "", "", "eau", "et", "ier", "in", "on", "ard", "oux", "ois"}

func (p frenchProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
	r := api.NewRand(cfg)
	caser := cases.Title(language.Und)

	realism := cfg.Realism
	if realism < 0 {
		realism = 0
	}
	if realism > 100 {
		realism = 100
	}

	useRealPct := 0
	switch {
	case realism >= 95:
		useRealPct = 95
	case realism >= 90:
		useRealPct = 90
	case realism >= 80:
		useRealPct = 80
	case realism >= 70:
		useRealPct = 55
	case realism >= 60:
		useRealPct = 35
	case realism >= 40:
		useRealPct = 20
	default:
		useRealPct = 5
	}
	chooseFromReal := func() bool { return r.Intn(100) < useRealPct }

	genSyl := func() string {
		if r.Intn(100) < 75 {
			return api.PickRand(onsets, r) + api.PickRand(vowels, r) + api.PickRand(codas, r)
		}
		return api.PickRand(vowels, r) + api.PickRand(onsets, r) + api.PickRand(vowels, r)
	}

	genGivenProcedural := func() string {
		numSyl := 2 + r.Intn(2) // 2..3
		if realism < 40 {
			numSyl = 1 + r.Intn(3) // 1..3
		}

		var b strings.Builder
		for i := 0; i < numSyl; i++ {
			b.WriteString(genSyl())
		}

		switch cfg.Gender {
		case "male":
			end := api.PickRand(givenEndingsMale, r)
			if end != "" && !strings.HasSuffix(b.String(), end) {
				b.WriteString(end)
			}
		case "female":
			end := api.PickRand(givenEndingsFemale, r)
			if end != "" && !strings.HasSuffix(b.String(), end) {
				b.WriteString(end)
			}
		default:
			end := api.PickRand(givenEndingsNeutral, r)
			if end != "" && !strings.HasSuffix(b.String(), end) {
				b.WriteString(end)
			}
		}
		return b.String()
	}

	genSurnameProcedural := func() string {
		numSyl := 2 + r.Intn(2) // 2..3
		var b strings.Builder
		for i := 0; i < numSyl; i++ {
			b.WriteString(genSyl())
		}

		thr := 20
		if realism >= 80 {
			thr = 45
		} else if realism >= 60 {
			thr = 30
		}
		if r.Intn(100) < thr {
			end := api.PickRand(surnameEndings, r)
			if end != "" && !strings.HasSuffix(b.String(), end) {
				b.WriteString(end)
			}
		}
		return b.String()
	}

	first := ""
	if chooseFromReal() {
		switch cfg.Gender {
		case "male":
			first = api.PickRand(firstMale, r)
		case "female":
			first = api.PickRand(firstFemale, r)
		default:
			roll := r.Intn(100)
			if roll < 60 {
				first = api.PickRand(firstNeutral, r)
			} else if roll < 80 {
				first = api.PickRand(firstMale, r)
			} else {
				first = api.PickRand(firstFemale, r)
			}
		}
	} else {
		first = caser.String(genGivenProcedural())
	}

	last := ""
	if cfg.IncludeLast {
		if cfg.Family == "" || strings.EqualFold(cfg.Family, PROFILE) {
			if chooseFromReal() {
				last = api.PickRand(lastNames, r)
			} else {
				last = caser.String(genSurnameProcedural())
			}
		} else {
			if chooseFromReal() {
				last = api.PickRand(lastNames, r)
			} else {
				last = caser.String(genSurnameProcedural())
			}
		}
	}

	first = caser.String(first)
	last = caser.String(last)
	return api.NameResult{First: first, Last: last}, nil
}

// Profile is the core exported symbol
var Profile frenchProfile
