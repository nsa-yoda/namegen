package amharic

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type amharicProfile struct{}

const PROFILE = "amharic"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p amharicProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Amharic/Ethiopian names (ASCII): patronymic-style, curated + procedural; deterministic",
	}
}

// Ethiopian names usually don't have surnames in the Western sense;
// we still generate a second name when includeLast is true.
var givenMale = []string{
	"Abebe", "Bekele", "Dawit", "Tesfaye", "Kebede", "Getachew", "Yohannes", "Mulugeta", "Solomon", "Alemayehu",
	"Biruk", "Girma", "Haile", "Mengistu", "Tadesse", "Fikru", "Eshetu", "Seifu", "Addisu", "Zerihun",
}

var givenFemale = []string{
	"Almaz", "Hanna", "Selam", "Mulu", "Meseret", "Tigist", "Rahel", "Saba", "Aster", "Genet",
	"Wubit", "Eden", "Liya", "Frehiwot", "Biruktawit", "Yeshi", "Marta", "Tsedey", "Yodit", "Mekdes",
}

var givenNeutral = []string{
	"Selam", "Biruk", "Mulu", "Genet", "Eden", "Liya", "Saba", "Haile", "Solomon", "Addisu",
}

// Used as second names / patronymics
var surnames = []string{
	"Bekele", "Tesfaye", "Kebede", "Abebe", "Getachew", "Alemayehu", "Girma", "Haile", "Mengistu", "Tadesse",
}

// Amharic romanized phonotactics
var onsets = []string{
	"", "",
	"b", "d", "f", "g", "h", "k", "l", "m", "n", "p", "r", "s", "t", "w", "y", "z",
	"ch", "sh",
}

var vowels = []string{
	"a", "e", "i", "o", "u",
	"aa", "ee", "ie",
}

var codas = []string{
	"", "", "", "",
	"n", "m", "r", "l", "t",
}

var givenEndings = []string{"", "", "", "e", "u", "a", "ye"}
var surnameEndings = []string{"", "", "", "ye", "w", "e"}

func (p amharicProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
		n := 2 + r.Intn(2)
		if realism < 40 {
			n = 1 + r.Intn(3)
		}
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString(genSyl())
		}
		if r.Intn(100) < 45 {
			b.WriteString(api.PickRand(givenEndings, r))
		}
		return b.String()
	}

	genSurnameProcedural := func() string {
		n := 2
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString(genSyl())
		}
		if r.Intn(100) < 50 {
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

var Profile amharicProfile
