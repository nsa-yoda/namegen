package turkish

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type turkishProfile struct{}

const PROFILE = "turkish"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p turkishProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Turkish names (ASCII): realism blends curated lists with procedural syllables; deterministic with seed",
	}
}

// Curated given names (ASCII; diacritics removed, e.g., Ş->S, ğ->g, ı->i, ö->o, ü->u, ç->c).
var firstMale = []string{
	"Mehmet", "Mustafa", "Ahmet", "Ali", "Emre", "Murat", "Yusuf", "Osman", "Hasan", "Huseyin",
	"Kerem", "Can", "Burak", "Omer", "Eren", "Serkan", "Cem", "Kaan", "Baris", "Deniz",
	"Onur", "Ibrahim", "Halil", "Suleyman", "Fatih", "Sinan", "Cenk", "Umut", "Tolga", "Taylan",
}

var firstFemale = []string{
	"Ayse", "Fatma", "Emine", "Zeynep", "Elif", "Merve", "Seda", "Esra", "Ebru", "Ceren",
	"Selin", "Derya", "Deniz", "Buse", "Gul", "Asli", "Hande", "Yasemin", "Aylin", "Melis",
	"Sibel", "Sevgi", "Nazan", "Tugce", "Ece", "Pinar", "Aysun", "Gizem", "Nazli", "Damla",
}

var firstNeutral = []string{
	"Deniz", "Can", "Eren", "Umut", "Derya", "Onur", "Ece", "Naz", "Miran", "Cem",
}

var lastNames = []string{
	"Yilmaz", "Kaya", "Demir", "Sahin", "Celik", "Yildiz", "Aydin", "Ozdemir", "Arslan", "Dogan",
	"Kilic", "Koc", "Aslan", "Yavuz", "Ozturk", "Erdogan", "Polat", "Aksoy", "Gunes", "Bulut",
	"Kaplan", "Karaca", "Toprak", "Tas", "Tekin", "Ekinci", "Eren", "Kurt", "Yalcin", "Sari",
}

// Procedural building blocks (Turkish-ish, vowel harmony not enforced; ASCII).
var vowels = []string{"a", "e", "i", "o", "u", "ai", "ei", "ia", "io", "ua"}
var onsets = []string{
	"", "",
	"b", "c", "d", "f", "g", "h", "k", "l", "m", "n", "p", "r", "s", "t", "v", "y", "z",
	"ch", "sh",
	"kr", "tr", "pr", "gr",
}
var codas = []string{"", "", "", "n", "m", "r", "l", "k", "t", "s"}

var givenEndingsMale = []string{"", "", "", "han", "can", "em", "er", "in", "an"}
var givenEndingsFemale = []string{"", "", "", "a", "e", "in", "nur", "sel", "gul"}
var givenEndingsNeutral = []string{"", "", "", "a", "e", "in", "er"}

var surnameEndings = []string{"", "", "", "oglu", "soy", "li", "er", "ci"}

func (p turkishProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
		numSyl := 2 + r.Intn(2)
		if realism < 40 {
			numSyl = 1 + r.Intn(3)
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
var Profile turkishProfile
