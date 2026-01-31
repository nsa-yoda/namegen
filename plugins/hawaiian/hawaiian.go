package hawaiian

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type hawaiianProfile struct{}

const PROFILE = "hawaiian"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p hawaiianProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Hawaiian-inspired names (ASCII): curated + strict phonotactic procedural fallback; deterministic",
	}
}

// Real Hawaiian uses okina and kahako; we keep ASCII-only approximations.
var givenMale = []string{
	"Kai", "Keanu", "Koa", "Noa", "Ikaika", "Kekoa", "Makana", "Keoni", "Kaleo", "Kanani",
	"Maleko", "Kainoa", "Kimo", "Kekai", "Lono", "Keola", "Kekoa", "Makoa", "Nalu", "Kekai",
}

var givenFemale = []string{
	"Leilani", "Kalani", "Malia", "Noelani", "Nalani", "Keala", "Moana", "Anela", "Kiana", "Lani",
	"Makana", "Kailani", "Melia", "Alana", "Kapua", "Mahina", "Kalea", "Kamalani", "Nanea", "Kekepania",
}

var givenNeutral = []string{
	"Kai", "Kalani", "Noa", "Moana", "Makana", "Lani", "Nalu", "Keala", "Kaleo", "Mahina",
}

var surnames = []string{
	"Kamehameha", "Kalakaua", "Kealoha", "Kawika", "Kailani", "Makana", "Kaleo", "Kamaka", "Keoni", "Kahale",
}

// Hawaiian phonotactics are very strict: consonants {h,k,l,m,n,p,w} + vowels.
var onsets = []string{
	"", "",
	"h", "k", "l", "m", "n", "p", "w",
}

var vowels = []string{
	"a", "e", "i", "o", "u",
	"ai", "ae", "ao", "au", "ei", "io", "oa", "oi", "ou", "ua", "ui",
}

var codas = []string{
	"", "", "", "", // mostly open syllables
}

var givenEndings = []string{"", "", "", "a", "i", "o", "u"}
var surnameEndings = []string{"", "", "", "lani", "nui", "loa", "mano"}

func (p hawaiianProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
		// Hawaiian names often 2-4 syllables.
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
			n = 3 + r.Intn(3)
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

var Profile hawaiianProfile
