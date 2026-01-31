package chinese

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type chineseProfile struct{}

const PROFILE = "chinese"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p chineseProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Chinese names (pinyin): realism blends curated lists with procedural pinyin syllables; deterministic with seed",
	}
}

// Curated pinyin given names (no tone marks for simplicity).
var firstMale = []string{
	"Wei", "Jie", "Jun", "Hao", "Ming", "Lei", "Qiang", "Bo", "Chen", "Feng",
	"Yu", "Peng", "Tao", "Yang", "Bin", "Guang", "Dong", "Chao", "Gang", "Sheng",
	"Zhi", "Heng", "Xiang", "Rui", "Yong", "Xuan", "Yifan", "Haoran", "Zhe", "Yuze",
}

var firstFemale = []string{
	"Mei", "Ling", "Yan", "Na", "Jing", "Xiu", "Hua", "Fang", "Ying", "Li",
	"Juan", "Min", "Qian", "Xue", "Xia", "Lan", "Ting", "Rong", "Xin", "Shan",
	"Yutong", "Yihan", "Zihan", "Ruoxi", "Xinyi", "Jia", "Yue", "Yuxi", "Kexin", "Meng",
}

var firstNeutral = []string{
	"Wei", "Yu", "Rui", "Xin", "Jia", "Yue", "Ming", "Yang", "Lin", "An",
}

// Curated pinyin surnames (common family names).
var lastNames = []string{
	"Wang", "Li", "Zhang", "Liu", "Chen", "Yang", "Huang", "Zhao", "Wu", "Zhou",
	"Xu", "Sun", "Ma", "Zhu", "Hu", "Guo", "He", "Gao", "Lin", "Luo",
	"Zheng", "Liang", "Xie", "Song", "Tang", "Han", "Feng", "Yu", "Dong", "Xiao",
}

// Pinyin syllable building blocks (simplified, no tones).
var initials = []string{
	"", "", "b", "p", "m", "f", "d", "t", "n", "l",
	"g", "k", "h", "j", "q", "x", "zh", "ch", "sh", "r",
	"z", "c", "s", "y", "w",
}

var finals = []string{
	"a", "ai", "an", "ang", "ao",
	"e", "ei", "en", "eng", "er",
	"i", "ia", "ian", "iang", "iao", "ie", "in", "ing", "iong", "iu",
	"o", "ong", "ou",
	"u", "ua", "uai", "uan", "uang", "ui", "un", "uo",
	"v", "ve", "van", "vn", // represent ü as v
}

// Common two-syllable given-name patterns are frequent; we keep optional 1-syllable too.
func (p chineseProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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

	// Probability to use curated lists vs procedural
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

	genSyllable := func() string {
		ini := api.PickRand(initials, r)
		fin := api.PickRand(finals, r)
		// Avoid impossible/awkward combos (very light filtering).
		// e.g., 'y' or 'w' often stand in for i/u glides; keep them only with certain finals.
		if (ini == "y" || ini == "w") && strings.HasPrefix(fin, "v") {
			ini = ""
		}
		// Avoid empty+er too often
		if ini == "" && fin == "er" && r.Intn(100) < 70 {
			fin = "e"
		}
		return ini + fin
	}

	genGivenProcedural := func() string {
		// Many given names are 2 syllables; allow 1 sometimes.
		n := 2
		if realism < 40 {
			if r.Intn(100) < 35 {
				n = 1
			}
		} else {
			if r.Intn(100) < 15 {
				n = 1
			}
		}

		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString(genSyllable())
		}
		return b.String()
	}

	// ---- Given name selection ----
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

	// ---- Surname selection ----
	last := ""
	if cfg.IncludeLast {
		if cfg.Family == "" || strings.EqualFold(cfg.Family, PROFILE) {
			if chooseFromReal() {
				last = api.PickRand(lastNames, r)
			} else {
				// Procedural surname: 1 syllable is most common; sometimes 2 for variety.
				n := 1
				if realism < 40 {
					if r.Intn(100) < 10 {
						n = 2
					}
				} else {
					if r.Intn(100) < 5 {
						n = 2
					}
				}
				var b strings.Builder
				for i := 0; i < n; i++ {
					b.WriteString(genSyllable())
				}
				last = caser.String(b.String())
			}
		} else {
			// If Family override is something else, still produce Chinese-ish surname for now.
			if chooseFromReal() {
				last = api.PickRand(lastNames, r)
			} else {
				last = caser.String(genSyllable())
			}
		}
	}

	first = caser.String(first)
	last = caser.String(last)
	return api.NameResult{First: first, Last: last}, nil
}

// Profile is the core exported symbol
var Profile chineseProfile
