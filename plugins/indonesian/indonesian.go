package indonesian

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type indonesianProfile struct{}

const PROFILE = "indonesian"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p indonesianProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Indonesian names (ASCII): realism blends curated lists with procedural syllables; deterministic with seed",
	}
}

// Indonesia has many naming conventions; many people have a single name.
// We'll generate a given name (First) and optionally a surname-ish (Last).
var givenMale = []string{
	"Agus", "Budi", "Dedi", "Eko", "Hadi", "Indra", "Joko", "Rizki", "Rudi", "Slamet",
	"Yusuf", "Ahmad", "Fajar", "Bayu", "Dimas", "Arif", "Hendra", "Wahyu", "Putra", "Surya",
}

var givenFemale = []string{
	"Ayu", "Dewi", "Sari", "Wulan", "Rina", "Intan", "Putri", "Indah", "Lestari", "Ratna",
	"Sri", "Nia", "Maya", "Rani", "Tika", "Fitri", "Nabila", "Aisyah", "Nurlaila", "Kartika",
}

var givenNeutral = []string{
	"Maya", "Rizki", "Indra", "Ayu", "Sari", "Bayu", "Dimas", "Nia", "Rani", "Wahyu",
}

// Some common-ish family names / second names used in Indonesia (not universal).
var surnames = []string{
	"Wijaya", "Saputra", "Pratama", "Santoso", "Setiawan", "Siregar", "Hidayat", "Wibowo",
	"Mahendra", "Kusuma", "Firmansyah", "Permata", "Nugroho", "Gunawan", "Utami",
}

var onsets = []string{
	"", "",
	"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "r", "s", "t", "w", "y", "z",
	"ng", "ny", "sy", "kh",
	"pr", "tr", "kr", "br",
}
var vowels = []string{
	"a", "e", "i", "o", "u",
	"ai", "au", "ia", "ua", "ei",
}
var codas = []string{
	"", "", "", "",
	"n", "m", "ng", "r", "h", "t", "k", "s",
}

var givenEndings = []string{"", "", "", "an", "ah", "i", "u"}
var surnameEndings = []string{"", "", "", "wan", "man", "yah", "tama", "putra", "sari"}

func (p indonesianProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
		// 2 syllables common; allow 1-3.
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
		if r.Intn(100) < 30 {
			b.WriteString(api.PickRand(givenEndings, r))
		}
		return b.String()
	}

	genSurnameProcedural := func() string {
		n := 2 + r.Intn(2) // 2..3
		if realism < 40 {
			n = 1 + r.Intn(3) // 1..3
		}
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString(genSyl())
		}
		if r.Intn(100) < 45 {
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
	} else {
		// Many Indonesians have a single name; at mid realism, often omit last implicitly.
		// (No-op)
	}

	return api.NameResult{First: caser.String(first), Last: caser.String(last)}, nil
}

var Profile indonesianProfile
