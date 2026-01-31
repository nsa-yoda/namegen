package english

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type englishProfile struct{}

const PROFILE = "english"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p englishProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "English names: realism blends real lists with procedural syllables; deterministic with seed",
	}
}

// Small curated lists (expand anytime).
// Intentionally mixed: classic + modern + neutral-ish.
var firstMale = []string{
	"James", "John", "Robert", "Michael", "William", "David", "Richard", "Joseph", "Thomas", "Charles",
	"Daniel", "Matthew", "Anthony", "Mark", "Paul", "Steven", "Andrew", "Joshua", "Kevin", "Brian",
	"Nathan", "Ryan", "Ethan", "Noah", "Liam", "Logan", "Lucas", "Benjamin", "Henry", "Jack",
	"Oliver", "Leo", "Miles", "Caleb", "Aaron", "Adam", "Jason", "Sean", "Kyle", "Eric",
}

var firstFemale = []string{
	"Mary", "Patricia", "Jennifer", "Linda", "Elizabeth", "Barbara", "Susan", "Jessica", "Sarah", "Karen",
	"Nancy", "Lisa", "Margaret", "Betty", "Sandra", "Ashley", "Kimberly", "Emily", "Donna", "Michelle",
	"Amanda", "Melissa", "Stephanie", "Rebecca", "Laura", "Hannah", "Olivia", "Sophia", "Ava", "Isabella",
	"Mia", "Amelia", "Grace", "Chloe", "Ella", "Lily", "Zoe", "Nora", "Lucy", "Claire",
}

var firstNeutral = []string{
	"Alex", "Jordan", "Taylor", "Morgan", "Casey", "Riley", "Jamie", "Quinn", "Avery", "Parker",
	"Reese", "Rowan", "Skyler", "Cameron", "Hayden", "Emerson", "Sage", "Finley", "Dakota", "Harper",
}

var lastNames = []string{
	"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez",
	"Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson", "Martin",
	"Lee", "Perez", "Thompson", "White", "Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson",
	"Walker", "Young", "Allen", "King", "Wright", "Scott", "Torres", "Nguyen", "Hill", "Flores",
	"Green", "Adams", "Nelson", "Baker", "Hall", "Rivera", "Campbell", "Mitchell", "Carter", "Roberts",
}

// Conservative mutation: small “English-feeling” tweaks for variety.
// Only applied at higher realism and low probability.
func mutateEnglish(r api.RandLike, s string) string {
	// Conservative mutations to add variety without drifting into gibberish.
	// Intended mostly for procedurally-generated strings.
	if len(s) < 3 {
		return s
	}

	lower := strings.ToLower(s)

	// 1) Occasionally normalize a few common digraphs.
	if r.Intn(100) < 20 {
		repls := [][2]string{
			{"ph", "f"},
			{"ck", "k"},
			{"qu", "kw"},
			{"ae", "e"},
			{"oe", "e"},
		}
		pair := repls[r.Intn(len(repls))]
		lower = strings.ReplaceAll(lower, pair[0], pair[1])
	}

	// 2) Occasionally drop a doubled vowel/consonant sequence.
	if r.Intn(100) < 30 {
		for _, pat := range []string{"aa", "ee", "ii", "oo", "uu", "tt", "ll", "rr", "ss", "nn"} {
			if strings.Contains(lower, pat) {
				lower = strings.ReplaceAll(lower, pat, pat[:1])
			}
		}
	}

	// 3) Occasionally add a soft ending.
	if r.Intn(100) < 15 {
		suffs := []string{"", "e", "y", "a", "n", "s"}
		suf := suffs[r.Intn(len(suffs))]
		if suf != "" && !strings.HasSuffix(lower, suf) {
			lower += suf
		}
	}

	// 4) Occasionally remove a trailing vowel to make it feel more surname-like.
	if r.Intn(100) < 10 {
		if strings.HasSuffix(lower, "a") || strings.HasSuffix(lower, "e") || strings.HasSuffix(lower, "i") ||
			strings.HasSuffix(lower, "o") || strings.HasSuffix(lower, "u") {
			if len(lower) > 3 {
				lower = lower[:len(lower)-1]
			}
		}
	}

	// 5) Very rarely swap two adjacent letters.
	if r.Intn(100) < 5 && len(lower) >= 4 {
		idx := 1 + r.Intn(len(lower)-2)
		b := []byte(lower)
		b[idx], b[idx+1] = b[idx+1], b[idx]
		lower = string(b)
	}

	return lower
}

func (p englishProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
	r := api.NewRand(cfg)
	caser := cases.Title(language.English)

	// Procedural building blocks (kept from your original approach)
	vowels := []string{"a", "e", "i", "o", "u"}
	consonants := []string{"b", "c", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "r", "s", "t", "v", "w", "y", "z"}
	realFragments := []string{"el", "ric", "mar", "an", "beth", "ron", "ly", "ton", "den", "ley", "gar", "wyn", "la", "li", "jo", "na", "mi", "sa"}

	genSyl := func(pat string) string {
		var b strings.Builder
		for _, ch := range pat {
			if ch == 'C' {
				b.WriteString(api.PickRand(consonants, r))
			} else {
				b.WriteString(api.PickRand(vowels, r))
			}
		}
		return b.String()
	}

	genProceduralFirst := func() string {
		// syllables influenced by realism and gender
		numSyl := 1 + r.Intn(3)
		if cfg.Realism > 70 {
			numSyl = 2 + r.Intn(2)
		}
		first := ""
		for i := 0; i < numSyl; i++ {
			pat := "CV"
			if cfg.Gender == "male" {
				if r.Intn(100) < 40 {
					pat = "CVC"
				}
			} else if cfg.Gender == "female" {
				if r.Intn(100) < 30 {
					pat = "V"
				}
			} else {
				if r.Intn(100) < 30 {
					pat = "CVC"
				}
			}
			// realism: inject fragments sometimes
			if cfg.Realism > 60 && r.Intn(100) < cfg.Realism/2 {
				first += api.PickRand(realFragments, r)
			} else {
				first += genSyl(pat)
			}
		}
		return first
	}

	genProceduralLast := func() string {
		last := ""
		parts := 1 + r.Intn(2)
		for i := 0; i < parts; i++ {
			last += genSyl("CVC")
		}
		// suffixes — make less aggressive at high realism
		if cfg.Family == PROFILE || cfg.Family == "" {
			roll := r.Intn(100)
			// At realism 100, allow suffix sometimes, but not constantly.
			threshold := 20
			if cfg.Realism < 60 {
				threshold = 60
			} else if cfg.Realism < 80 {
				threshold = 40
			}
			if roll < threshold {
				suffs := []string{"son", "ford", "wood", "well", "shire", "field", "stone", "brook"}
				last += api.PickRand(suffs, r)
			}
		}
		return last
	}

	// --- Realism blending strategy ---
	// realism in [0..100]
	realism := cfg.Realism
	if realism < 0 {
		realism = 0
	}
	if realism > 100 {
		realism = 100
	}

	// probability (0..100) to use “real list” instead of procedural
	// Make it ramp hard after 60 and very strong after 80.
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

	// First name selection
	first := ""
	if chooseFromReal() {
		switch cfg.Gender {
		case "male":
			first = api.PickRand(firstMale, r)
		case "female":
			first = api.PickRand(firstFemale, r)
		default:
			// neutral: mix neutral list plus a bit of male/female
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
		first = caser.String(genProceduralFirst())
	}

	// Last name selection
	last := ""
	if cfg.IncludeLast {
		if chooseFromReal() {
			last = api.PickRand(lastNames, r)
		} else {
			last = caser.String(genProceduralLast())
		}
	}

	// Light mutation (very conservative) at high realism, low probability.
	// This keeps results from repeating when Count is large.
	if realism >= 90 && r.Intn(100) < 8 {
		// Simple tweak: if procedural was used, occasionally soften double letters.
		// (Keep it minimal to avoid gibberish.)
		first = strings.ReplaceAll(first, "aa", "a")
		first = strings.ReplaceAll(first, "ee", "e")
		first = strings.ReplaceAll(first, "ii", "i")
		first = strings.ReplaceAll(first, "oo", "o")
		first = strings.ReplaceAll(first, "uu", "u")
		last = strings.ReplaceAll(last, "aa", "a")
		last = strings.ReplaceAll(last, "ee", "e")
		last = strings.ReplaceAll(last, "ii", "i")
		last = strings.ReplaceAll(last, "oo", "o")
		last = strings.ReplaceAll(last, "uu", "u")
	}

	// Ensure proper casing if we generated procedurally
	first = caser.String(first)
	last = caser.String(last)

	return api.NameResult{First: first, Last: last}, nil
}

// Profile is the core exported symbol
var Profile englishProfile
