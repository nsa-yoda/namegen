package maori

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type maoriProfile struct{}

const PROFILE = "maori"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p maoriProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Maori-inspired names (ASCII): curated + phonotactic procedural fallback; deterministic",
	}
}

// Maori uses macrons in real orthography; we keep ASCII.
var givenMale = []string{
	"Wiremu", "Hemi", "Rangi", "Tama", "Hone", "Rawiri", "Tane", "Kauri", "Manu", "Aroha",
	"Ngata", "Kahu", "Koro", "Matiu", "Hori", "Timi", "Pita", "TeRangi", "Kingi", "Hoani",
}

var givenFemale = []string{
	"Aroha", "Anahera", "Mere", "Moana", "Hine", "Ria", "Kiri", "Rangi", "Wai", "Maia",
	"Marama", "Rere", "Ata", "Hera", "Mereana", "TeAroha", "Tia", "Kahurangi", "Manawa", "Hinemoa",
}

var givenNeutral = []string{
	"Aroha", "Moana", "Rangi", "Manu", "Maia", "Wai", "Ata", "Kauri", "Kahu", "Manawa",
}

var surnames = []string{
	"Ngata", "TeRangi", "TeAroha", "TeKahu", "TeWai", "Tame", "Ranginui", "Tukiri", "Kahukura", "Manawa",
}

// Maori phonotactics are very strict: (C)V with limited consonants; "ng", "wh" common.
var onsets = []string{
	"", "",
	"h", "k", "m", "n", "p", "r", "t", "w",
	"ng", "wh",
}

var vowels = []string{
	"a", "e", "i", "o", "u",
	"ai", "ae", "ao", "au", "ei", "io", "oa", "oi", "ou", "ua", "ui",
}

var codas = []string{
	"", "", "", "", // Maori syllables usually open
}

var givenEndings = []string{"", "", "", "a", "e", "i", "o", "u"}
var surnameEndings = []string{"", "", "", "nui", "rangi", "waka", "manawa"}

func (p maoriProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
		// Maori names often 2-4 syllables.
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
		if r.Intn(100) < 40 {
			b.WriteString(api.PickRand(givenEndings, r))
		}
		return b.String()
	}

	genSurnameProcedural := func() string {
		n := 3
		if realism < 40 {
			n = 2 + r.Intn(3)
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

var Profile maoriProfile
