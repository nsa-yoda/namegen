package malay

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type malayProfile struct{}

const PROFILE = "malay"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p malayProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Malay names (ASCII): realism blends curated lists with procedural syllables; deterministic with seed",
	}
}

// Malaysia naming varies (patronymics common, some family names).
// We'll generate a given name (First) and optionally a last/family (Last).
var givenMale = []string{
	"Ahmad", "Muhammad", "Hafiz", "Hakim", "Faiz", "Azlan", "Syafiq", "Firdaus", "Amir", "Farhan",
	"Imran", "Iskandar", "Razak", "Fikri", "Irfan", "Zul", "Zulkifli", "Aiman", "Adib", "Khairul",
}

var givenFemale = []string{
	"Nur", "Nurul", "Aisyah", "Siti", "Hannah", "Farah", "Nadia", "Aina", "Alya", "Balqis",
	"Syahirah", "Izzah", "Shahira", "Sofea", "Amira", "Maryam", "Husna", "Zara", "Najwa", "Diyana",
}

var givenNeutral = []string{
	"Nur", "Aiman", "Amir", "Nadia", "Farah", "Alya", "Hafiz", "Irfan", "Zara", "Najwa",
}

var surnames = []string{
	"Abdullah", "Ibrahim", "Ismail", "Hassan", "Hamid", "Rahman", "Razak", "Yusof", "Aziz", "Mahmud",
	"Zainal", "Mustafa", "Saleh", "Othman", "Kassim",
}

// Procedural romanized syllables (Malay-ish).
var onsets = []string{
	"", "",
	"b", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "r", "s", "t", "w", "y", "z",
	"ch", "sy", "kh",
	"br", "kr", "pr", "tr",
}
var vowels = []string{
	"a", "e", "i", "o", "u",
	"ai", "au", "ia", "ua", "ei",
}
var codas = []string{
	"", "", "", "",
	"n", "m", "ng", "r", "h", "t", "k", "s",
}

var givenEndingsMale = []string{"", "", "", "din", "man", "raf", "zul", "far"}
var givenEndingsFemale = []string{"", "", "", "ah", "a", "na", "ira", "nur"}
var givenEndingsNeutral = []string{"", "", "", "ah", "a", "an", "in"}

var surnameEndings = []string{"", "", "", "bin", "binti", "rahman", "din", "man"}

func (p malayProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
			n = 1 + r.Intn(3) // 1..3
		} else if r.Intn(100) < 20 {
			n = 1 + r.Intn(3)
		}
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString(genSyl())
		}
		switch cfg.Gender {
		case "male":
			b.WriteString(api.PickRand(givenEndingsMale, r))
		case "female":
			b.WriteString(api.PickRand(givenEndingsFemale, r))
		default:
			b.WriteString(api.PickRand(givenEndingsNeutral, r))
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
		if r.Intn(100) < 40 {
			b.WriteString(api.PickRand(surnameEndings, r))
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

var Profile malayProfile
