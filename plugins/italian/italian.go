package italian

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type italianProfile struct{}

const PROFILE = "italian"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p italianProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Italian names: realism blends curated lists with procedural syllables; deterministic with seed",
	}
}

// Curated given names.
var firstMale = []string{
	"Marco", "Luca", "Matteo", "Giovanni", "Francesco", "Alessandro", "Andrea", "Giorgio", "Paolo", "Stefano",
	"Roberto", "Davide", "Simone", "Federico", "Riccardo", "Antonio", "Giuseppe", "Salvatore", "Vincenzo", "Nicola",
	"Enrico", "Fabio", "Daniele", "Massimo", "Leonardo", "Emanuele", "Pietro", "Filippo", "Michele", "Claudio",
}

var firstFemale = []string{
	"Giulia", "Sofia", "Martina", "Francesca", "Chiara", "Alice", "Elena", "Valentina", "Sara", "Laura",
	"Federica", "Alessia", "Giorgia", "Silvia", "Elisa", "Paola", "Roberta", "Claudia", "Maria", "Anna",
	"Beatrice", "Camilla", "Arianna", "Lucia", "Ilaria", "Simona", "Caterina", "Serena", "Emanuela", "Cristina",
}

var firstNeutral = []string{
	"Andrea", "Gabriele", "Alex", "Noa", "Sasha", "Giovi", "Dani", "Vale", "Nico", "Rene",
}

// Curated surnames.
var lastNames = []string{
	"Rossi", "Russo", "Ferrari", "Esposito", "Bianchi", "Romano", "Colombo", "Ricci", "Marino", "Greco",
	"Bruno", "Gallo", "Conti", "Costa", "Giordano", "Mancini", "Rizzo", "Lombardi", "Moretti", "Barbieri",
	"Fontana", "Santoro", "Mariani", "Rinaldi", "Caruso", "Ferrara", "Gatti", "Longo", "Martinelli", "Leone",
}

// Procedural building blocks (Italian-ish, vowel-forward).
var vowels = []string{"a", "e", "i", "o", "u", "ai", "ei", "ia", "io", "ua"}
var onsets = []string{
	"b", "c", "d", "f", "g", "l", "m", "n", "p", "r", "s", "t", "v", "z",
	"br", "cr", "dr", "fr", "gr", "pr", "tr",
	"ch", "gh", "gl", "gn", "sc", "sp", "st",
	"", "", // allow vowel-start sometimes
}
var codas = []string{"", "", "", "n", "l", "r", "s", "t"}

var givenEndingsMale = []string{"", "", "", "o", "i", "e", "ino", "etto", "one"}
var givenEndingsFemale = []string{"", "", "", "a", "ia", "ina", "etta", "ella"}
var givenEndingsNeutral = []string{"", "", "", "a", "e", "i"}

var surnameEndings = []string{"", "", "", "i", "o", "a", "ini", "etti", "elli", "one", "aro"}

func (p italianProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
var Profile italianProfile
