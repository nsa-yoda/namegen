package igbo

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type igboProfile struct{}

const PROFILE = "igbo"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p igboProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Igbo names (ASCII): meaning-based compounds with procedural fallback; deterministic",
	}
}

// Igbo names are often meaningful phrases; many are gender-neutral.
var givenMale = []string{
	"Chinedu", "Emeka", "Ifeanyi", "Nnamdi", "Obinna", "Chukwudi", "Uche", "Ikenna", "Onyekachi", "Ifeoma",
	"Chibuike", "Somto", "Chima", "Okechukwu", "Chijioke", "Chukwuka", "Nwachukwu", "Chukwuma", "Uzoma", "Ifechukwu",
}

var givenFemale = []string{
	"Chiamaka", "Ngozi", "Ifunanya", "Nkiru", "Uju", "Chinwe", "Ifeoma", "Obiageli", "Chizoba", "Nkechi",
	"Amarachi", "Chisom", "Uchechi", "Nkiruka", "Onyinye", "Somadina", "Chinenye", "Chidimma", "Chinyere", "Nnenna",
}

var givenNeutral = []string{
	"Uche", "Chisom", "Somto", "Ifeoma", "Uzoma", "Onyekachi", "Ifeanyi", "Amarachi", "Somadina", "Chibuike",
}

var surnames = []string{
	"Okafor", "Okeke", "Nwoye", "Eze", "Obi", "Nnamdi", "Chukwu", "Anyanwu", "Okorie", "Onyekachi",
	"Nwankwo", "Nwoye", "Uche", "Nwafor", "Onyekwere",
}

// Igbo phonotactics (simple CV-heavy structure)
var onsets = []string{
	"", "",
	"b", "ch", "d", "f", "g", "h", "j", "k", "l", "m", "n", "nw", "ny", "p", "r", "s", "t", "w", "y", "z",
}

var vowels = []string{
	"a", "e", "i", "o", "u",
	"ia", "io", "ua", "ee",
}

var codas = []string{
	"", "", "", "",
	"m", "n", "r",
}

var givenEndings = []string{"", "", "", "chi", "ma", "na", "du", "ka"}
var surnameEndings = []string{"", "", "", "eze", "chukwu", "nna", "for"}

func (p igboProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
		return api.PickRand(onsets, r) +
			api.PickRand(vowels, r) +
			api.PickRand(codas, r)
	}

	genGivenProcedural := func() string {
		n := 3
		if realism < 40 {
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
		n := 2 + r.Intn(2)
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

var Profile igboProfile
