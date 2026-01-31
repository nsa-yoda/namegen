package hebrew

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type hebrewProfile struct{}

const PROFILE = "hebrew"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p hebrewProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Hebrew names (romanized, ASCII): realism blends curated lists with procedural syllables; deterministic with seed",
	}
}

// Curated given names (romanized; ASCII only).
var firstMale = []string{
	"David", "Daniel", "Yosef", "Moshe", "Avi", "Ariel", "Eitan", "Noam", "Omer", "Itai",
	"Yonatan", "Natan", "Shlomo", "Yitzhak", "Yaakov", "Gideon", "Uri", "Amir", "Eli", "Shai",
	"Lev", "Asher", "Hillel", "Nadav", "Baruch", "Elazar", "Ze'ev", "Reuven", "Shimon", "Yoav",
}

var firstFemale = []string{
	"Sarah", "Rivka", "Leah", "Rachel", "Miriam", "Hannah", "Noa", "Yael", "Tamar", "Avigail",
	"Shira", "Michal", "Noga", "Eden", "Lior", "Adi", "Maya", "Tal", "Orly", "Naama",
	"Esther", "Hadassah", "Chaya", "Tzipora", "Ofra", "Dana", "Gali", "Roni", "Batya", "Nitzan",
}

var firstNeutral = []string{
	"Noam", "Ariel", "Adi", "Tal", "Lior", "Eden", "Roni", "Nitzan", "Shai", "Maya",
}

// Curated surnames (common in Israeli / Jewish contexts; ASCII only).
var lastNames = []string{
	"Cohen", "Levi", "Mizrahi", "Peretz", "Biton", "Dahan", "Katz", "Shapiro", "Friedman", "Rosenberg",
	"Goldberg", "Weiss", "Klein", "Golan", "Barak", "BenAmi", "BenDavid", "BenHaim", "Azoulay", "Amar",
	"Dayan", "Halevi", "Navon", "Sharabi", "Ohayon", "Sasson", "Segal", "Gross", "Edelstein", "Rabin",
}

// Procedural building blocks (Hebrew-ish romanization, simplified).
var vowels = []string{"a", "e", "i", "o", "u", "ai", "ei", "ia", "oa"}
var onsets = []string{
	"", "", // vowel-start sometimes
	"b", "d", "g", "h", "k", "l", "m", "n", "p", "r", "s", "t", "v", "y", "z",
	"sh", "ch", "tz", "kh",
}
var codas = []string{"", "", "", "n", "m", "r", "l", "t", "k", "sh"}

var givenEndingsMale = []string{"", "", "", "el", "an", "am", "on", "ai"}
var givenEndingsFemale = []string{"", "", "", "a", "ah", "el", "it", "ya"}
var givenEndingsNeutral = []string{"", "", "", "a", "el", "on"}

var surnameEndings = []string{"", "", "", "man", "berg", "stein", "son", "i"}

func (p hebrewProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
	r := api.NewRand(cfg)
	caser := cases.Title(language.Und)

	realism := cfg.Realism
	if realism < 0 {
		realism = 0
	}
	if realism > 100 {
		realism = 100
	}

	// same curve as other profiles
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
		numSyl := 2 + r.Intn(2)
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
		if chooseFromReal() {
			last = api.PickRand(lastNames, r)
		} else {
			last = caser.String(genSurnameProcedural())
		}
	}

	first = caser.String(first)
	last = caser.String(last)
	return api.NameResult{First: first, Last: last}, nil
}

// Profile is the core exported symbol
var Profile hebrewProfile
