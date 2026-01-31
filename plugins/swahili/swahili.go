package swahili

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type swahiliProfile struct{}

const PROFILE = "swahili"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p swahiliProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Swahili names (ASCII): realism blends curated lists with procedural syllables; deterministic with seed",
	}
}

var givenMale = []string{
	"Juma", "Hassan", "Ali", "Said", "Bakari", "Hamisi", "Omari", "Salim", "Kassim", "Abdallah",
	"Daudi", "Musa", "Ismail", "Rashid", "Faraji", "Baraka", "Amani", "Shaban", "Azizi", "Idris",
}

var givenFemale = []string{
	"Asha", "Zainab", "Fatuma", "Halima", "Rehema", "Neema", "Zuri", "Safiya", "Subira", "Mariam",
	"Najma", "Amina", "Bahati", "Upendo", "Wema", "Imani", "Zawadi", "Nuru", "Siti", "Ruqayya",
}

var givenNeutral = []string{
	"Amani", "Baraka", "Imani", "Nuru", "Bahati", "Zawadi", "Neema", "Zuri", "Wema", "Rehema",
}

var surnames = []string{
	"Ali", "Hassan", "Said", "Abdallah", "Omari", "Bakari", "Juma", "Salim", "Kassim", "Musa",
	"Daudi", "Ismail", "Idris", "Rashid", "Azizi",
}

// Swahili-ish syllables
var onsets = []string{
	"", "",
	"b", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "r", "s", "t", "v", "w", "y", "z",
	"ch", "sh", "ng", "ny", "mw", "kw",
}
var vowels = []string{"a", "e", "i", "o", "u", "aa", "ee", "ia", "ua", "ai"}
var codas = []string{"", "", "", "", "n", "m", "ng", "r", "l", "t", "k", "s"}

var givenEndings = []string{"", "", "", "a", "i", "u", "ni", "ri"}
var surnameEndings = []string{"", "", "", "wa", "ani", "eni", "oni"}

func (p swahiliProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
			n = 1 + r.Intn(3)
		}
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString(genSyl())
		}
		if r.Intn(100) < 30 {
			b.WriteString(api.PickRand(givenEndings, r))
		}
		return b.String()
	}

	genSurnameProcedural := func() string {
		n := 2 + r.Intn(2)
		if realism < 40 {
			n = 1 + r.Intn(3)
		}
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString(genSyl())
		}
		if r.Intn(100) < 35 {
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

var Profile swahiliProfile
