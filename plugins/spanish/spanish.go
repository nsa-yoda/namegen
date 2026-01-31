package spanish

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type spanishProfile struct{}

const PROFILE = "spanish"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p spanishProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Spanish names: realism blends curated lists with procedural syllables; deterministic with seed",
	}
}

// Curated lists (expand anytime).
var firstMale = []string{
	"Juan", "Jose", "Carlos", "Luis", "Javier", "Miguel", "Antonio", "Manuel", "Francisco", "Pedro",
	"Sergio", "Diego", "Rafael", "Fernando", "Alejandro", "Pablo", "Andres", "Ricardo", "Roberto", "Alberto",
	"Mario", "Raul", "Hector", "Emilio", "Eduardo", "Jorge", "Victor", "Adrian", "Ivan", "Oscar",
}

var firstFemale = []string{
	"Maria", "Ana", "Carmen", "Isabel", "Laura", "Elena", "Sofia", "Lucia", "Paula", "Marta",
	"Patricia", "Claudia", "Andrea", "Raquel", "Sara", "Julia", "Natalia", "Silvia", "Rosa", "Teresa",
	"Beatriz", "Irene", "Noelia", "Cristina", "Alicia", "Monica", "Daniela", "Carolina", "Veronica", "Adriana",
}

var firstNeutral = []string{
	"Alex", "Cruz", "Angel", "Noa", "Ariel", "Dani", "Gael", "Andrea", "Sam", "Rene",
}

var lastNames = []string{
	"Garcia", "Gonzalez", "Rodriguez", "Fernandez", "Lopez", "Martinez", "Sanchez", "Perez", "Gomez", "Martin",
	"Jimenez", "Ruiz", "Hernandez", "Diaz", "Moreno", "Munoz", "Alvarez", "Romero", "Alonso", "Gutierrez",
	"Navarro", "Torres", "Dominguez", "Vazquez", "Ramos", "Gil", "Serrano", "Blanco", "Molina", "Morales",
}

// Procedural building blocks to produce Spanish-ish phonotactics.
var vowels = []string{"a", "e", "i", "o", "u"}

// Common onsets (include Spanish digraphs).
var onsets = []string{
	"b", "c", "d", "f", "g", "h", "j", "l", "m", "n", "p", "r", "s", "t", "v", "z",
	"ch", "ll", "rr",
}

// Common consonant endings (codas). Empty strings keep many syllables open.
var codas = []string{"", "", "", "n", "s", "r", "l", "d"}

// Endings to nudge names into more recognizable shapes.
var givenEndings = []string{"", "", "", "a", "o", "ia", "io", "el", "in"}
var surnameEndings = []string{"", "", "", "ez", "es", "ado", "era", "ero", "osa", "illo"}

func (p spanishProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
	r := api.NewRand(cfg)
	caser := cases.Title(language.Spanish)

	// clamp realism
	realism := cfg.Realism
	if realism < 0 {
		realism = 0
	}
	if realism > 100 {
		realism = 100
	}

	// Probability to use curated lists vs procedural
	// Ramp hard after 60 and very strong after 80.
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
		// Mostly CV(+optional coda), sometimes VCV
		if r.Intn(100) < 70 {
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

		end := api.PickRand(givenEndings, r)
		if end != "" && !strings.HasSuffix(b.String(), end) {
			b.WriteString(end)
		}

		// small gender nuance at higher realism
		if cfg.Gender == "female" && realism >= 70 && r.Intn(100) < 20 {
			if !strings.HasSuffix(b.String(), "a") {
				b.WriteString("a")
			}
		}
		if cfg.Gender == "male" && realism >= 70 && r.Intn(100) < 20 {
			if strings.HasSuffix(b.String(), "a") {
				s := b.String()
				b.Reset()
				b.WriteString(strings.TrimSuffix(s, "a"))
				b.WriteString("o")
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

		// Spanish-ish patronymic endings become more likely at higher realism
		thr := 25
		if realism >= 80 {
			thr = 55
		} else if realism >= 60 {
			thr = 40
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
		// honor Family override, but keep Spanish-ish behavior if unset or spanish
		if cfg.Family == "" || strings.EqualFold(cfg.Family, PROFILE) {
			if chooseFromReal() {
				last = api.PickRand(lastNames, r)
			} else {
				last = caser.String(genSurnameProcedural())
			}
		} else {
			// If Family override is something else, still produce Spanish-ish surname for now
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
var Profile spanishProfile
