package tamil

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type tamilProfile struct{}

const PROFILE = "tamil"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p tamilProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Tamil-inspired names: realism blends curated lists with procedural syllables; deterministic with seed",
	}
}

// Curated given names commonly used among Tamil speakers (romanized; ASCII only).
// (Not exhaustive; expand anytime.)
var firstMale = []string{
	"Arun", "Karthik", "Vijay", "Ajith", "Suresh", "Ramesh", "Prakash", "Ganesh", "Hari", "Kumar",
	"Murugan", "Senthil", "Saravanan", "Madhan", "Naveen", "Sathish", "Dinesh", "Rajesh", "Bala", "Venkatesh",
	"Anand", "Shankar", "Sekar", "Mani", "Gopi", "Sivakumar", "Kathir", "Ravi", "Subash", "Thiru",
}

var firstFemale = []string{
	"Anjali", "Lakshmi", "Meena", "Priya", "Divya", "Kavitha", "Nandhini", "Deepa", "Revathi", "Sindhu",
	"Shalini", "Saranya", "Keerthi", "Aishwarya", "Pavithra", "Geetha", "Uma", "Mahalakshmi", "Sangeetha", "Vaishnavi",
	"Janani", "Ranjani", "Thenmozhi", "Vidhya", "Malathi", "Padma", "Sujatha", "Anitha", "Swathi", "Radhika",
}

var firstNeutral = []string{
	"Kiran", "Arun", "Naveen", "Anand", "Hari", "Mani", "Devi", "Sasi", "Bala", "Ravi",
}

// Curated surnames / family identifiers.
// Note: Tamil naming conventions vary widely (patronymics, initials, place names). For generator purposes
// we provide common-style surnames as a stand-in.
var lastNames = []string{
	"Iyer", "Iyengar", "Pillai", "Nadar", "Gounder", "Thevar", "Chettiar", "Mudaliar", "Naicker", "Reddy",
	"Menon", "Krishnan", "Subramanian", "Narayanan", "Raman", "Sundaram", "Srinivasan", "Venkatesan", "Balakrishnan", "Chandrasekar",
	"Rajendran", "Shanmugam", "Arumugam", "Kumar", "Anand", "Mohan", "Raghavan", "Varadarajan", "Sivakumar", "Murthy",
}

// Procedural building blocks to produce Tamil-ish romanized phonotactics.
// Keep it readable in ASCII; lean into common syllables like ka/tha/na/ra/ma/sa and endings like -an/-ar/-am/-i.
var vowels = []string{"a", "aa", "i", "ii", "u", "uu", "e", "ee", "o", "oo"}

var onsets = []string{
	"", "", // allow vowel-start sometimes
	"k", "g", "c", "j", "t", "th", "d", "n", "p", "b", "m", "y", "r", "l", "v", "s", "h",
	"kr", "pr", "tr", "sr",
	"ch", "sh",
}

var codas = []string{"", "", "", "n", "m", "r", "l", "y", "s", "k"}

var givenEndingsMale = []string{"", "", "", "an", "ar", "am", "esh", "kumar"}
var givenEndingsFemale = []string{"", "", "", "a", "i", "ini", "laxmi", "devi"}
var givenEndingsNeutral = []string{"", "", "", "a", "i", "an"}

var surnameEndings = []string{"", "", "", "an", "ar", "am", "iah", "appa"}

func (p tamilProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
			// If Family override is something else, still produce Tamil-ish surname for now.
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
var Profile tamilProfile
