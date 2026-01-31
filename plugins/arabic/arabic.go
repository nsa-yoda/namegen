package arabic

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type arabicProfile struct{}

const PROFILE = "arabic"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p arabicProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Arabic names: realism blends curated transliterated lists with procedural syllables; deterministic with seed",
	}
}

// Curated transliterated lists (expand anytime).
var firstMale = []string{
	"Muhammad", "Ahmed", "Ali", "Omar", "Hassan", "Hussein", "Yusuf", "Ibrahim", "Abdullah", "Khalid",
	"Mahmoud", "Tariq", "Bilal", "Mustafa", "Sami", "Nabil", "Fadi", "Rami", "Zaid", "Hamza",
	"Amir", "Salim", "Karim", "Jamal", "Faisal", "Adel", "Ismail", "Marwan", "Anas", "Samir",
}

var firstFemale = []string{
	"Fatima", "Aisha", "Maryam", "Layla", "Noor", "Sara", "Hana", "Zainab", "Amal", "Salma",
	"Yasmin", "Rania", "Noura", "Huda", "Dalia", "Lina", "Nadia", "Iman", "Reem", "Samar",
	"Farah", "Mona", "Mariam", "Aya", "Leena", "Hiba", "Nada", "Jana", "Ruqayya", "Sumaya",
}

var firstNeutral = []string{
	"Noor", "Iman", "Rami", "Sami", "Salam", "Hadi", "Zain", "Amin", "Jude", "Rayan",
}

var lastNames = []string{
	"Almasri", "Alharbi", "Alsayed", "Haddad", "Nassar", "Khatib", "Salem", "Farah", "Yousef", "Hamdan",
	"Abbas", "Khalil", "Mansour", "Najjar", "Amin", "Sharif", "Bakri", "Qasim", "Saeed", "Fahmy",
	"Aziz", "Hussein", "Mahmoud", "Taha", "Darwish", "Sabbagh", "Zahran", "Fadel", "Ghanem", "Rashid",
}

// Procedural building blocks for Arabic-ish transliteration.
// Keep it simple: mostly CV/CVC with some common clusters/digraphs.
var vowels = []string{"a", "i", "u", "e", "o"}

// Common onsets, including transliterated digraphs.
// (We keep it readable and not too harsh.)
var onsets = []string{
	"b", "t", "th", "j", "h", "kh", "d", "dh", "r", "z", "s", "sh",
	"f", "q", "k", "l", "m", "n", "w", "y", "g",
}

// Codas are often empty; sometimes n/r/l/d/s to feel name-like.
var codas = []string{"", "", "", "", "n", "r", "l", "d", "s", "m"}

// Endings to nudge outputs into more “name-ish” shapes.
var givenEndings = []string{"", "", "", "a", "ah", "an", "in", "un", "i", "y"}
var surnameEndings = []string{"", "", "", "i", "iy", "awi", "ani", "ari", "ullah", "uddin"}

func (p arabicProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
		// Mostly onset+vowel(+optional coda), sometimes vowel+onset+vowel for variety.
		if r.Intn(100) < 75 {
			return api.PickRand(onsets, r) + api.PickRand(vowels, r) + api.PickRand(codas, r)
		}
		return api.PickRand(vowels, r) + api.PickRand(onsets, r) + api.PickRand(vowels, r)
	}

	genGivenProcedural := func() string {
		// 2–4 syllables; lower realism sometimes 1–3
		numSyl := 2 + r.Intn(3) // 2..4
		if realism < 40 {
			numSyl = 1 + r.Intn(3) // 1..3
		}

		var b strings.Builder
		for i := 0; i < numSyl; i++ {
			b.WriteString(genSyl())
		}

		// soft ending
		end := api.PickRand(givenEndings, r)
		if end != "" && !strings.HasSuffix(b.String(), end) {
			b.WriteString(end)
		}

		// slight gender bias at higher realism: more 'a/ah' endings for female,
		// more 'i/in' endings for male (very light touch).
		if realism >= 70 {
			if cfg.Gender == "female" && r.Intn(100) < 20 {
				s := b.String()
				if !strings.HasSuffix(s, "a") && !strings.HasSuffix(s, "ah") {
					b.WriteString("a")
				}
			}
			if cfg.Gender == "male" && r.Intn(100) < 15 {
				s := b.String()
				if strings.HasSuffix(s, "a") {
					b.Reset()
					b.WriteString(strings.TrimSuffix(s, "a"))
					b.WriteString("i")
				}
			}
		}

		return b.String()
	}

	genSurnameProcedural := func() string {
		// 2–3 syllables
		numSyl := 2 + r.Intn(2) // 2..3
		var b strings.Builder
		for i := 0; i < numSyl; i++ {
			b.WriteString(genSyl())
		}

		// surname endings slightly more likely at higher realism
		thr := 20
		if realism >= 80 {
			thr = 45
		} else if realism >= 60 {
			thr = 30
		}
		if r.Intn(100) < thr {
			b.WriteString(api.PickRand(surnameEndings, r))
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
			if roll < 55 {
				first = api.PickRand(firstNeutral, r)
			} else if roll < 78 {
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
			// If Family override is something else, still produce Arabic-ish surname for now.
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
var Profile arabicProfile
