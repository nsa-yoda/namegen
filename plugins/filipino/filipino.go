package filipino

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type filipinoProfile struct{}

const PROFILE = "filipino"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p filipinoProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Filipino names: realism blends curated lists (Tagalog/Spanish-influenced) with procedural syllables; deterministic with seed",
	}
}

// Curated given names commonly used in the Philippines (mix of Tagalog, Spanish, and modern).
var firstMale = []string{
	"Juan", "Jose", "Antonio", "Miguel", "Andres", "Ramon", "Ricardo", "Eduardo", "Fernando", "Manuel",
	"Paolo", "Marco", "Carlo", "Enrique", "Gabriel", "Angelo", "Noel", "Renato", "Emilio", "Vicente",
	"Arnel", "Danilo", "Ernesto", "Isko", "Jun", "Junjun", "Nico", "Rafael", "Roberto", "Tomas",
}

var firstFemale = []string{
	"Maria", "Ana", "Carmen", "Isabel", "Teresa", "Rosa", "Elena", "Patricia", "Cristina", "Sofia",
	"Paula", "Andrea", "Daniela", "Carla", "Angelica", "Maricel", "May", "Mae", "Joy", "Grace",
	"Liza", "Lea", "Nena", "Nenita", "Charo", "Nora", "Regina", "Victoria", "Yvonne", "Michelle",
}

var firstNeutral = []string{
	"Alex", "Jamie", "Jordan", "Sam", "Taylor", "Rene", "Noel", "Angel", "Rio", "Ariel",
}

// Curated surnames common in the Philippines (Spanish influence + local).
var lastNames = []string{
	"Santos", "Reyes", "Cruz", "Bautista", "Gonzales", "Garcia", "Aquino", "Ramos", "Mendoza", "Torres",
	"Flores", "Rivera", "Castillo", "Navarro", "Domingo", "Villanueva", "Dela Cruz", "Del Rosario", "Salazar", "De Guzman",
	"Mercado", "Valdez", "Fernandez", "Diaz", "Morales", "Hernandez", "Manalo", "Pascual", "Vergara", "Rosales",
}

// Procedural syllable building blocks.
// Keep it simple and readable: mostly open syllables, light consonant clusters.
var vowels = []string{"a", "e", "i", "o", "u"}

var onsets = []string{
	"b", "k", "d", "g", "h", "l", "m", "n", "p", "r", "s", "t", "w", "y",
	"ch", "ng", "sh",
}

// Tagalog often uses open syllables; codas are rarer.
var codas = []string{"", "", "", "", "n", "ng", "s", "r", "t"}

var givenEndings = []string{"", "", "", "a", "o", "i", "an", "en", "in"}
var surnameEndings = []string{"", "", "", "son", "san", "dez", "ez", "ano", "ista"}

func (p filipinoProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
		// Mostly CV(+optional coda), sometimes VCV.
		if r.Intn(100) < 75 {
			return api.PickRand(onsets, r) + api.PickRand(vowels, r) + api.PickRand(codas, r)
		}
		return api.PickRand(vowels, r) + api.PickRand(onsets, r) + api.PickRand(vowels, r)
	}

	genGivenProcedural := func() string {
		// 2–3 syllables; low realism allows 1–3
		numSyl := 2 + r.Intn(2) // 2..3
		if realism < 40 {
			numSyl = 1 + r.Intn(3) // 1..3
		}
		var b strings.Builder
		for i := 0; i < numSyl; i++ {
			b.WriteString(genSyl())
		}
		end := api.PickRand(givenEndings, r)
		if end != "" && !strings.HasSuffix(b.String(), end) {
			b.WriteString(end)
		}
		return b.String()
	}

	genSurnameProcedural := func() string {
		// 2 syllables typically
		numSyl := 2
		if realism < 30 && r.Intn(100) < 20 {
			numSyl = 1
		}
		var b strings.Builder
		for i := 0; i < numSyl; i++ {
			b.WriteString(genSyl())
		}
		// Surname endings are modest; more likely at higher realism.
		thr := 20
		if realism >= 80 {
			thr = 40
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
			// If Family override is something else, still produce Filipino-ish surname for now.
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
var Profile filipinoProfile
