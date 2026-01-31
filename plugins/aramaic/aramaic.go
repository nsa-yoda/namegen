package aramaic

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type aramaicProfile struct{}

const PROFILE = "aramaic"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p aramaicProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Aramaic/Syriac-inspired names (ASCII romanization): curated + procedural fallback; deterministic with seed",
	}
}

// Note: This is a lightweight romanized set inspired by common Biblical/Syriac-era forms.
// ASCII only.
var givenMale = []string{
	"Bartholomew", "Thomas", "Yohannan", "Yeshua", "Shimon", "Yosef", "Yaqub", "Matthai", "Taddeus", "Philip",
	"Andreas", "Petros", "Paulos", "Barnaba", "Hanania", "Azaria", "Mishael", "Natan", "Eliya", "Gamaliel",
}

var givenFemale = []string{
	"Maryam", "Martha", "Hannah", "Sarah", "Rivqa", "Leah", "Rachel", "Elizabeth", "Salome", "Susanna",
	"Deborah", "Judith", "Tamar", "Dinah", "Esther", "Miriam", "Naomi", "Abigail", "Shifra", "Zipporah",
}

var givenNeutral = []string{
	"Shimon", "Yosef", "Hannah", "Maryam", "Eliya", "Natan", "Tamar", "Miriam", "Naomi", "Judith",
}

// Patronymic / clan-like / place-like endings (not truly “surnames” historically).
var surnames = []string{
	"Bar", "BarNatan", "BarYosef", "BarShimon", "Bethlehem", "Nazareth", "Ephesus", "Edessa", "Antioch", "Damascus",
	"HaLevi", "Cohen",
}

// Procedural building blocks (Semitic-ish romanization, simplified).
var onsets = []string{
	"", "",
	"b", "d", "g", "h", "k", "l", "m", "n", "p", "q", "r", "s", "t", "w", "y", "z",
	"sh", "ch", "kh", "th", "ts",
	"br", "tr", "kr",
}

var vowels = []string{
	"a", "e", "i", "o", "u",
	"aa", "ee", "ia", "ie", "oa", "ou",
}

var codas = []string{
	"", "", "", "",
	"n", "m", "r", "l", "t", "k", "sh", "th",
}

var givenEndingsMale = []string{"", "", "", "el", "an", "am", "on", "ya", "iah"}
var givenEndingsFemale = []string{"", "", "", "a", "ah", "el", "it", "ya"}
var givenEndingsNeutral = []string{"", "", "", "a", "el", "on"}

var surnameEndings = []string{"", "", "", "bar", "beth", "iya", "el", "an"}

func (p aramaicProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
	r := api.NewRand(cfg)
	caser := cases.Title(language.Und)

	realism := cfg.Realism
	if realism < 0 {
		realism = 0
	}
	if realism > 100 {
		realism = 100
	}

	// same curve you’re using everywhere else
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
		// Semitic-ish CV(C) feel, with occasional vowel-start.
		if r.Intn(100) < 75 {
			return api.PickRand(onsets, r) + api.PickRand(vowels, r) + api.PickRand(codas, r)
		}
		return api.PickRand(vowels, r) + api.PickRand(onsets, r) + api.PickRand(vowels, r)
	}

	genGivenProcedural := func() string {
		// often 2-3 syllables
		n := 2 + r.Intn(2) // 2..3
		if realism < 40 {
			n = 1 + r.Intn(3) // 1..3
		}
		var b strings.Builder
		for i := 0; i < n; i++ {
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
		// 2-3 syllables, sometimes with a suffix.
		n := 2 + r.Intn(2)
		if realism < 40 {
			n = 1 + r.Intn(3)
		}
		var b strings.Builder
		for i := 0; i < n; i++ {
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

	// ---- First (given) ----
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

	// ---- Last (surname/patronymic-ish) ----
	last := ""
	if cfg.IncludeLast {
		if chooseFromReal() {
			// With some chance, treat "Bar X" as a patronymic.
			if r.Intn(100) < 35 {
				// "Bar" + (a given name)
				var child string
				switch cfg.Gender {
				case "male":
					child = api.PickRand(givenMale, r)
				case "female":
					child = api.PickRand(givenFemale, r)
				default:
					child = api.PickRand(givenNeutral, r)
				}
				last = "Bar" + caser.String(child)
			} else {
				last = api.PickRand(surnames, r)
			}
		} else {
			last = caser.String(genSurnameProcedural())
		}
	}

	first = caser.String(first)
	last = caser.String(last)
	return api.NameResult{First: first, Last: last}, nil
}

var Profile aramaicProfile
