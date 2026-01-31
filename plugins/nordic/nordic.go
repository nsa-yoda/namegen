package nordic

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type nordicProfile struct{}

const PROFILE = "nordic"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p nordicProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Nordic names: realism blends curated Scandinavian lists with procedural syllables; deterministic with seed",
	}
}

// Curated Scandinavian given names (ASCII only; expand anytime).
var firstMale = []string{
	"Erik", "Karl", "Lars", "Sven", "Bjorn", "Leif", "Nils", "Oskar", "Otto", "Felix",
	"Hans", "Johan", "Jonas", "Magnus", "Henrik", "Rolf", "Ulf", "Gunnar", "Harald", "Sigurd",
	"Anders", "Hakon", "Einar", "Ragnar", "Stellan", "Torbjorn", "Mikkel", "Kristian", "Mats", "Kjell",
}

var firstFemale = []string{
	"Anna", "Elsa", "Ingrid", "Freya", "Astrid", "Sigrid", "Helga", "Greta", "Klara", "Maja",
	"Ida", "Lina", "Karin", "Hilda", "Frida", "Solveig", "Liv", "Nora", "Emilia", "Matilda",
	"Hanna", "Lotte", "Saga", "Tove", "Sanna", "Eira", "Alva", "Linnea", "Agnes", "Kristin",
}

var firstNeutral = []string{
	"Alex", "Robin", "Kim", "Noa", "Mika", "Lenn", "Toni", "Jules", "Nika", "Elli",
}

// Curated Nordic-style surnames (ASCII; mix of Swedish/Norwegian/Danish patterns).
var lastNames = []string{
	"Johansson", "Andersson", "Karlsson", "Nilsson", "Larsson", "Olsson", "Persson", "Svensson", "Gustafsson", "Pettersson",
	"Hansen", "Jensen", "Nielsen", "Olsen", "Lund", "Dahl", "Berg", "Lindberg", "Lindstrom", "Bergstrom",
	"Nygaard", "Skov", "Haugland", "Solberg", "Sandberg", "Lind", "Holm", "Ekberg", "Soderberg", "Thorsen",
}

// Procedural building blocks (Nordic-ish phonotactics; simple ASCII).
var vowels = []string{"a", "e", "i", "o", "u", "y", "ae", "oe"}

var onsets = []string{
	"b", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "r", "s", "t", "v", "w",
	"bj", "dj", "fj", "gj", "hj", "kj", "lj", "mj", "nj", "rj", "sj", "tj", "vj",
	"sk", "st", "sp", "sn", "sm", "sl", "sv",
	"tr", "dr", "br", "gr", "kr", "fr",
}

var codas = []string{"", "", "", "n", "r", "s", "t", "d", "k", "l", "m", "ng"}

var givenEndingsMale = []string{"", "", "", "er", "ar", "rik", "ulf", "vald", "son"}
var givenEndingsFemale = []string{"", "", "", "a", "e", "hild", "frid", "borg", "dis"}
var givenEndingsNeutral = []string{"", "", "", "en", "in", "e"}

var surnameEndings = []string{"", "", "", "son", "sen", "berg", "strom", "lund", "holm", "gaard", "vik"}

func (p nordicProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
	r := api.NewRand(cfg)
	caser := cases.Title(language.Und)

	// clamp realism
	realism := cfg.Realism
	if realism < 0 {
		realism = 0
	}
	if realism > 100 {
		realism = 100
	}

	// Probability to use curated lists vs procedural (same curve as other profiles)
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
		// Mostly onset+vowel(+optional coda), sometimes vowel+onset+vowel.
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

		// Ending by gender (light touch)
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

		// Surname endings more likely at higher realism
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

	// ---- First name selection ----
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

	// ---- Last name selection ----
	last := ""
	if cfg.IncludeLast {
		if cfg.Family == "" || strings.EqualFold(cfg.Family, PROFILE) {
			if chooseFromReal() {
				last = api.PickRand(lastNames, r)
			} else {
				last = caser.String(genSurnameProcedural())
			}
		} else {
			// If Family override is something else, still produce Nordic-ish surname for now.
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
var Profile nordicProfile
