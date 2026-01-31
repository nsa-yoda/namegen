package baltic

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type balticProfile struct{}

const PROFILE = "baltic"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p balticProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Baltic-inspired (Lithuanian/Latvian) names (ASCII): curated + procedural fallback; deterministic",
	}
}

// Curated (ASCII; no diacritics).
var givenMale = []string{
	"Jonas", "Marius", "Tomas", "Darius", "Mindaugas", "Vytautas", "Paulius", "Andrius", "Rokas", "Lukas",
	"Martynas", "Arnas", "Gintaras", "Saulius", "Kestas", "Edgaras", "Karolis", "Domantas", "Justas", "Ignas",
}

var givenFemale = []string{
	"Aiste", "Ruta", "Egle", "Ieva", "Lina", "Rasa", "Jurate", "Vaida", "Gabriele", "Monika",
	"Kristina", "Inga", "Dovile", "Laura", "Milda", "Greta", "Aurelija", "Simona", "Viktorija", "Edita",
}

var givenNeutral = []string{
	"Ruta", "Lina", "Laura", "Monika", "Simona", "Tomas", "Lukas", "Rokas", "Marius", "Greta",
}

// Curated Baltic-ish surnames (ASCII).
var surnames = []string{
	"Kazlauskas", "Petrauskas", "Jankauskas", "Stankevicius", "Zukauskas", "Vaitkus", "Butkus", "Kavaliauskas",
	"Berzins", "Kalnins", "Ozols", "Liepa", "Jansons", "Krumins", "Balodis",
}

var maleSurnameEndings = []string{"as", "is", "us", "aitis", "enas", "onis"}
var femaleSurnameEndings = []string{"a", "e", "iene", "yte", "aite", "ute"}

// Procedural blocks
var onsets = []string{
	"", "",
	"b", "d", "g", "j", "k", "l", "m", "n", "p", "r", "s", "t", "v", "z",
	"br", "dr", "gr", "kr", "pr", "tr",
	"st", "sk", "sp",
	"dz",
}

var vowels = []string{
	"a", "e", "i", "o", "u", "y",
	"ai", "ei", "ie", "uo", "au", "ia",
}

var codas = []string{
	"", "", "", "",
	"n", "m", "r", "l", "s", "t", "d", "k", "g",
	"ns", "rs", "lis", "tis",
}

var givenEndings = []string{"", "", "", "as", "is", "us", "a", "e", "ius"}
var surnameEndingsNeutral = []string{"", "", "", "as", "is", "us", "ins", "aus", "aitis"}

func (p balticProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
		return api.PickRand(onsets, r) + api.PickRand(vowels, r) + api.PickRand(codas, r)
	}

	genGivenProcedural := func() string {
		n := 2
		if realism < 40 {
			n = 1 + r.Intn(3)
		} else if r.Intn(100) < 25 {
			n = 2 + r.Intn(2)
		}
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString(genSyl())
		}
		if r.Intn(100) < 55 {
			b.WriteString(api.PickRand(givenEndings, r))
		}
		return b.String()
	}

	genSurnameProcedural := func() string {
		n := 2 + r.Intn(2) // 2..3
		if realism < 40 {
			n = 1 + r.Intn(3)
		}
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString(genSyl())
		}

		// add a gender-ish surname ending sometimes
		if r.Intn(100) < 70 {
			switch cfg.Gender {
			case "male":
				b.WriteString(api.PickRand(maleSurnameEndings, r))
			case "female":
				b.WriteString(api.PickRand(femaleSurnameEndings, r))
			default:
				b.WriteString(api.PickRand(surnameEndingsNeutral, r))
			}
		}
		return b.String()
	}

	first := ""
	if chooseFromReal() {
		switch cfg.Gender {
		case "male":
			first = api.PickRand(givenMale, r)
		case "female":
			first = api.PickRand(givenFemale, r)
		default:
			roll := r.Intn(100)
			if roll < 60 {
				first = api.PickRand(givenNeutral, r)
			} else if roll < 80 {
				first = api.PickRand(givenMale, r)
			} else {
				first = api.PickRand(givenFemale, r)
			}
		}
	} else {
		first = caser.String(genGivenProcedural())
	}

	last := ""
	if cfg.IncludeLast {
		if chooseFromReal() {
			last = api.PickRand(surnames, r)
		} else {
			last = caser.String(genSurnameProcedural())
		}
	}

	return api.NameResult{First: caser.String(first), Last: caser.String(last)}, nil
}

var Profile balticProfile
