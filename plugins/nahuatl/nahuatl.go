package nahuatl

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type nahuatlProfile struct{}

const PROFILE = "nahuatl"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p nahuatlProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Nahuatl-inspired names: realism blends curated Nahuatl-style transliterations with procedural syllables; deterministic with seed",
	}
}

// Curated Nahuatl-inspired / Nahuatl-origin names in common Latin transliteration.
// (Not exhaustive; expand anytime.)
var firstMale = []string{
	"Cuauhtemoc", "Tenoch", "Nezahualcoyotl", "Itzcoatl", "Moctezuma", "Cuitlahuac", "Tlahuicole", "Axayacatl",
	"Tizoc", "Ahuizotl", "Xolotl", "Tlaloc", "Ocelotl", "Yaotl", "Mictlantecuhtli", "Huitzilihuitl",
	"Chimalpopoca", "Totoquihuatzin", "Ixtlilxochitl", "Xochipilli",
}

var firstFemale = []string{
	"Xochitl", "Citlali", "Izel", "Malinalli", "Metztli", "Yaretzi", "Xilonen", "Tonantzin",
	"Chalchiuhtlicue", "Tlaltecuhtli", "Xochiquetzal", "Cihuacoatl", "Mecatl", "Ilancueitl", "Atotoztli", "Zyanya",
	"Nahuatl", "Yolotzin", "Teyacapan", "Ixtli",
}

var firstNeutral = []string{
	"Xochitl", "Citlali", "Izel", "Metztli", "Yaotl", "Ocelotl", "Tenoch", "Xolotl", "Yolotzin", "Tlaloc",
}

// "Last name" style elements. Historically, Nahua naming traditions differ from modern surname usage;
// these are Nahuatl-style epithets/constructs for generator purposes.
var lastNames = []string{
	"Xochitlal", "Cuauhtli", "Ocelotzin", "Yolotzin", "Tepetl", "Tlalli", "Tonal", "Miztli",
	"Itzcuintli", "Chalchiuh", "Cihuatl", "Tecuhtli", "Popoca", "Tzompantli", "Cempoal", "Acatl",
	"Tletl", "Atl", "Coatl", "Xihuitl",
}

// Procedural building blocks to produce Nahuatl-ish phonotactics.
// Keep it readable in ASCII and lean into signature clusters (tl, tz, hu, cu, x, ch).
var vowels = []string{"a", "e", "i", "o", "u", "ia", "oa"}

var onsets = []string{
	"", "", "", // allow vowel-start
	"t", "k", "p", "m", "n", "l", "y", "w", "h",
	"x", "ch", "tz", "tl", "cu", "hu", "qu",
	"te", "to", "ta", // common starter feel
}

var codas = []string{"", "", "", "l", "n", "m", "t", "k", "tl", "tz"}

var givenEndings = []string{"", "", "", "tl", "tli", "tzin", "yotl", "coatl", "tecuhtli"}
var surnameEndings = []string{"", "", "", "tzin", "yotl", "tl", "tli", "co", "pan", "tlan"}

func (p nahuatlProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
		// Mostly onset+vowel(+optional coda)
		s := api.PickRand(onsets, r) + api.PickRand(vowels, r)
		if r.Intn(100) < 35 {
			s += api.PickRand(codas, r)
		}
		return s
	}

	genGivenProcedural := func() string {
		// 2–4 syllables; low realism allows 1–4
		numSyl := 2 + r.Intn(3) // 2..4
		if realism < 35 {
			numSyl = 1 + r.Intn(4) // 1..4
		}

		var b strings.Builder
		for i := 0; i < numSyl; i++ {
			b.WriteString(genSyl())
		}

		// Add a characteristic ending more often at higher realism
		thr := 20
		if realism >= 80 {
			thr = 55
		} else if realism >= 60 {
			thr = 40
		}
		if r.Intn(100) < thr {
			end := api.PickRand(givenEndings, r)
			if end != "" && !strings.HasSuffix(b.String(), end) {
				b.WriteString(end)
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

		thr := 25
		if realism >= 80 {
			thr = 55
		} else if realism >= 60 {
			thr = 40
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
			// If Family override is something else, still produce Nahuatl-ish surname for now.
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
var Profile nahuatlProfile
