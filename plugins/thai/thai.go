package thai

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type thaiProfile struct{}

const PROFILE = "thai"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p thaiProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Thai names (ASCII romanization): realism blends curated lists with procedural syllables; deterministic with seed",
	}
}

// Thai naming is complex; romanization varies. This is a lightweight generator.
var givenMale = []string{
	"Somchai", "Somsak", "Prasit", "Krit", "Niran", "Anan", "Kittisak", "Surasak", "Wichai", "Chaiwat",
	"Thanakorn", "Preecha", "Sakchai", "Teerapong", "Narong", "Suthipong", "Phakorn", "Kosin", "Chanon", "Tanin",
}

var givenFemale = []string{
	"Siri", "Suda", "Kanya", "Nok", "Pim", "Sunee", "Wipa", "Chompoo", "Natcha", "Ploy",
	"Patcharaporn", "Sudarat", "Kanchana", "Supaporn", "Woranuch", "Chanida", "Nanthita", "Ratri", "Mayuree", "Araya",
}

var givenNeutral = []string{
	"Nok", "Pim", "Siri", "Krit", "Niran", "Anan", "Ploy", "Natcha", "Chanon", "Araya",
}

// Thai surnames are often long/unique; we include some common-ish examples.
var surnames = []string{
	"Saetang", "Srisai", "Wongsa", "Rattanakorn", "Sukhum", "Kanchanapong", "Boonyarat", "Chantarangsu",
	"Phromma", "Sanguansak", "Kittipong", "Rattanapong", "Sukprasert", "Wattanakul", "Srisuk", "Wongchai",
}

var onsets = []string{
	"", "",
	"b", "d", "f", "g", "h", "k", "kh", "l", "m", "n", "p", "ph", "r", "s", "t", "th", "w", "y", "ch", "j",
	"kr", "pr", "tr", "pl", "kl",
}
var vowels = []string{
	"a", "aa", "e", "ee", "i", "ii", "o", "oo", "u", "uu",
	"ai", "ao", "ua", "ue", "ia",
}
var codas = []string{
	"", "", "", "",
	"n", "m", "ng", "t", "k", "p", "r", "l", "s",
}

var givenEndingsMale = []string{"", "", "", "chai", "sak", "pon", "wat", "korn"}
var givenEndingsFemale = []string{"", "", "", "rat", "porn", "nee", "da", "ya"}
var givenEndingsNeutral = []string{"", "", "", "n", "ng", "da"}

var surnameEndings = []string{"", "", "", "kul", "sak", "pong", "chai", "wat", "korn"}

func (p thaiProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
		// Thai surnames tend to be longer; bias 3-4 syllables.
		numSyl := 3 + r.Intn(2) // 3..4
		if realism < 40 {
			numSyl = 2 + r.Intn(3) // 2..4
		}
		var b strings.Builder
		for i := 0; i < numSyl; i++ {
			b.WriteString(genSyl())
		}

		thr := 25
		if realism >= 80 {
			thr = 50
		} else if realism >= 60 {
			thr = 35
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

var Profile thaiProfile
