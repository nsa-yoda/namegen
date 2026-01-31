package yoruba

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type yorubaProfile struct{}

const PROFILE = "yoruba"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p yorubaProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Yoruba names (ASCII): realism blends curated lists with procedural syllables; deterministic with seed",
	}
}

// Yoruba names often have meaningful compounds. Romanization varies; we keep ASCII.
var givenMale = []string{
	"Oladele", "Oluwaseun", "Oluwatobi", "Olamide", "Olawale", "Adewale", "Adekunle", "Adebayo", "Adeyemi", "Babajide",
	"Kayode", "Kehinde", "Taiwo", "Segun", "Tunde", "Femi", "Seyi", "Kunle", "Bode", "Dayo",
}

var givenFemale = []string{
	"Yetunde", "Oluwafunke", "Oluwatoyin", "Olamide", "Olayinka", "Aderonke", "Adebimpe", "Bolanle", "Temitope", "Funmilayo",
	"Folake", "Sade", "Kehinde", "Taiwo", "Bimpe", "Tola", "Dami", "Morayo", "Simisola", "Bisi",
}

var givenNeutral = []string{
	"Olamide", "Kehinde", "Taiwo", "Temitope", "Dami", "Seyi", "Tola", "Morayo", "Simisola", "Bode",
}

var surnames = []string{
	"Adeyemi", "Adebayo", "Adewale", "Ogunleye", "Olawale", "Oluwole", "Ojo", "Balogun", "Olawuyi",
	"Akinyemi", "Akinwale", "Olatunji", "Olaoye", "Adeniran", "Adesina",
}

// Yoruba-ish syllables (ASCII only; includes common digraphs like gb).
var onsets = []string{
	"", "",
	"b", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "r", "s", "t", "w", "y",
	"gb",
	"ol", "ad", "ak", "og", "oy",
}

var vowels = []string{
	"a", "e", "i", "o", "u",
	"aa", "ee", "ii", "oo", "ai", "ei", "oi", "ua",
}

var codas = []string{
	"", "", "", "",
	"n", "m", "r", "l", "s",
}

var givenEndings = []string{"", "", "", "de", "mi", "se", "to", "bo", "ye", "ni"}
var surnameEndings = []string{"", "", "", "yemi", "bayo", "wale", "tunde", "kunle", "tobi"}

func (p yorubaProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
		// Yoruba names can be 3-4 syllables; bias slightly longer.
		n := 3
		if realism < 40 {
			n = 2 + r.Intn(3) // 2..4
		} else if r.Intn(100) < 25 {
			n = 2 + r.Intn(3)
		}
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString(genSyl())
		}
		if r.Intn(100) < 50 {
			b.WriteString(api.PickRand(givenEndings, r))
		}
		return b.String()
	}

	genSurnameProcedural := func() string {
		n := 2 + r.Intn(2) // 2..3
		if realism < 40 {
			n = 2 + r.Intn(3) // 2..4
		}
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString(genSyl())
		}
		if r.Intn(100) < 55 {
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

var Profile yorubaProfile
