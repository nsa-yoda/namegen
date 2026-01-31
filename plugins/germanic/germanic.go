package germanic

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type germanicProfile struct{}

const PROFILE = "germanic"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p germanicProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Germanic names: realism blends curated lists (German/Scandinavian/Old Norse-ish) with procedural syllables; deterministic with seed",
	}
}

// Curated given names (ASCII only; expand anytime).
var firstMale = []string{
	"Erik", "Karl", "Lars", "Sven", "Bjorn", "Leif", "Nils", "Oskar", "Otto", "Felix",
	"Hans", "Johan", "Jonas", "Magnus", "Henrik", "Rolf", "Ulf", "Gunnar", "Harald", "Sigurd",
	"Dietrich", "Heinrich", "Konrad", "Wilhelm", "Friedrich", "Johann", "Anders", "Hakon", "Einar", "Ragnar",
}

var firstFemale = []string{
	"Anna", "Elsa", "Ingrid", "Freya", "Astrid", "Sigrid", "Helga", "Greta", "Klara", "Maja",
	"Ida", "Lina", "Karin", "Hilda", "Brunhild", "Hedwig", "Gertrud", "Johanna", "Frida", "Solveig",
	"Liv", "Nora", "Emilia", "Matilda", "Hanna", "Lotte", "Sabine", "Anneliese", "Hildegard", "Kristin",
}

var firstNeutral = []string{
	"Alex", "Robin", "Kim", "Sascha", "Noa", "Nika", "Jules", "Toni", "Mika", "Lenn",
}

// Curated surnames (mix of German/Scandinavian style; ASCII only).
var lastNames = []string{
	"Muller", "Schmidt", "Schneider", "Fischer", "Weber", "Meyer", "Wagner", "Becker", "Hoffmann", "Schulz",
	"Koch", "Bauer", "Richter", "Klein", "Wolf", "Neumann", "Schroder", "Braun", "Kruger", "Jensen",
	"Hansen", "Olsen", "Lindberg", "Lund", "Berg", "Bergstrom", "Nygaard", "Dahl", "Soderberg", "Johansson",
}

// Procedural building blocks (Germanic-ish phonotactics; simple ASCII).
var vowels = []string{"a", "e", "i", "o", "u", "ae", "oe"}

var onsets = []string{
	"b", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "r", "s", "t", "v", "w",
	"br", "dr", "fr", "gr", "kr", "pr", "tr",
	"sk", "st", "sp", "sn", "sm", "sl", "sw",
	"ch", "sch",
}

var codas = []string{"", "", "", "n", "r", "s", "t", "d", "k", "l", "m", "ng"}

var givenEndingsMale = []string{"", "", "", "er", "ar", "rik", "ulf", "mund", "son"}
var givenEndingsFemale = []string{"", "", "", "a", "e", "hild", "gund", "lind", "borg"}
var givenEndingsNeutral = []string{"", "", "", "en", "in", "e"}

var surnameEndings = []string{"", "", "", "son", "sen", "berg", "strom", "mann", "wald", "heim", "gaard"}

func (p germanicProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
			// If Family override is something else, still produce Germanic-ish surname for now.
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
var Profile germanicProfile
