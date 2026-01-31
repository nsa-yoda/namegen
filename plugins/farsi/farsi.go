package farsi

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type farsiProfile struct{}

const PROFILE = "farsi"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p farsiProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Persian (Farsi) names (romanized, ASCII): realism blends curated lists with procedural syllables; deterministic with seed",
	}
}

// Curated given names (romanized; ASCII only).
var firstMale = []string{
	"Ali", "Reza", "Mohammad", "Hossein", "Mehdi", "Amir", "Saeed", "Morteza", "Hassan", "Javad",
	"Farhad", "Arash", "Kourosh", "Soroush", "Babak", "Shahin", "Kamran", "Navid", "Sina", "Yashar",
	"Ehsan", "Masoud", "Hamid", "Majid", "Behnam", "Pouya", "Ramin", "Payam", "Shahram", "Kian",
}

var firstFemale = []string{
	"Sara", "Maryam", "Fatemeh", "Zahra", "Neda", "Leila", "Mina", "Niloofar", "Shirin", "Parisa",
	"Golnaz", "Roya", "Elham", "Arezoo", "Samira", "Nazanin", "Hoda", "Mahtab", "Atena", "Darya",
	"Azadeh", "Ladan", "Yasaman", "Setareh", "Shadi", "Fereshteh", "Sahar", "Taraneh", "Shabnam", "Kiana",
}

var firstNeutral = []string{
	"Sara", "Neda", "Darya", "Sina", "Navid", "Kian", "Roya", "Shirin", "Ari", "Sam",
}

// Curated surnames (common Persian-family names; ASCII only).
var lastNames = []string{
	"Ahmadi", "Hosseini", "Mohammadi", "Rezaei", "Karimi", "Rahimi", "Ebrahimi", "Shirazi", "Tehrani", "Jafari",
	"Farhadi", "Khosravi", "Kazemi", "Ghasemi", "Soleimani", "Moradi", "Sadeghi", "Mahdavi", "Bakhtiari", "Zand",
	"Mehrabi", "Salehi", "Hedayati", "Rostami", "Shahbazi", "Abbasi", "Azimi", "Tavakoli", "Darvishi", "Nouri",
}

// Procedural building blocks (Persian-ish romanization, simplified).
var vowels = []string{"a", "e", "i", "o", "u", "aa", "ee", "oo", "ai", "ou"}
var onsets = []string{
	"", "",
	"b", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "r", "s", "t", "v", "y", "z",
	"sh", "ch", "kh", "gh", "zh",
	"br", "kr", "dr",
}
var codas = []string{"", "", "", "n", "m", "r", "l", "d", "t", "k", "sh"}

var givenEndingsMale = []string{"", "", "", "an", "ar", "ad", "id", "in", "shah"}
var givenEndingsFemale = []string{"", "", "", "a", "eh", "ieh", "naz", "gol"}
var givenEndingsNeutral = []string{"", "", "", "a", "an", "in"}

var surnameEndings = []string{"", "", "", "i", "ian", "zadeh", "pour", "nejad"}

func (p farsiProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
var Profile farsiProfile
