package samoan

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type samoanProfile struct{}

const PROFILE = "samoan"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p samoanProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Samoan-inspired names (ASCII): curated + phonotactic procedural fallback; deterministic",
	}
}

var givenMale = []string{
	"Tui", "Mika", "Sione", "Ioane", "Manu", "Peni", "Luka", "Iosefa", "Tavita", "Kelepi",
	"Faafoi", "Afa", "Toa", "Pita", "Tama", "Fetu", "Leota", "Faatoia", "Atoa", "Malie",
}

var givenFemale = []string{
	"Lupe", "Mele", "Lina", "Sala", "Fia", "Tala", "Sina", "Lagi", "Manu", "Tia",
	"Leilani", "Malia", "Fetu", "Alofa", "Saia", "Ava", "Tasi", "Moe", "Eseta", "Nia",
}

var givenNeutral = []string{
	"Manu", "Tala", "Fetu", "Ava", "Tui", "Lagi", "Tasi", "Nia", "Mika", "Tama",
}

var surnames = []string{
	"Tuimalealiifano", "Tuilagi", "Faumuina", "Malietoa", "Saelua", "Fepuleai", "Leota", "Toleafoa", "Tufuga", "Aiono",
}

// Samoan-like phonotactics: mostly open syllables (C)V; "ng" appears in Polynesian.
var onsets = []string{
	"", "",
	"f", "g", "h", "k", "l", "m", "n", "p", "r", "s", "t", "v",
	"ng",
}

var vowels = []string{
	"a", "e", "i", "o", "u",
	"ai", "au", "ei", "ia", "io", "oa", "oi", "ou", "ua", "ui",
}

var codas = []string{
	"", "", "", "",
}

var givenEndings = []string{"", "", "", "a", "i", "o", "u"}
var surnameEndings = []string{"", "", "", "toga", "lani", "mana", "toa"}

func (p samoanProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
		n := 3
		if realism < 40 {
			n = 2 + r.Intn(3)
		} else if r.Intn(100) < 25 {
			n = 2 + r.Intn(3)
		}
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString(genSyl())
		}
		if r.Intn(100) < 35 {
			b.WriteString(api.PickRand(givenEndings, r))
		}
		return b.String()
	}

	genSurnameProcedural := func() string {
		n := 4
		if realism < 40 {
			n = 3 + r.Intn(3) // 3..5
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
			first = api.PickRand(givenNeutral, r)
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

var Profile samoanProfile
