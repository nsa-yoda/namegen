package korean

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type koreanProfile struct{}

const PROFILE = "korean"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p koreanProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Korean names (romanized): realism blends curated lists with procedural syllables; deterministic with seed",
	}
}

// Curated given names (romanized; ASCII only).
// These are common-ish modern given names, not Hangul.
var firstMale = []string{
	"Minjun", "Seojun", "Jiho", "Joon", "Hyunwoo", "Taehyun", "Junho", "Donghyun", "Seungmin", "Jisung",
	"Hyun", "Sungmin", "Jinhyuk", "Jaehoon", "Wonjun", "Daehyun", "Kangmin", "Sangwoo", "Youngho", "Byungwoo",
	"Jaewon", "Seungwoo", "Kihyun", "Sungwoo", "Hyeonjin", "Seongho", "Jinwoo", "Kyungsoo", "Inho", "Gunwoo",
}

var firstFemale = []string{
	"Seoyeon", "Seoah", "Jiwon", "Soojin", "Hyejin", "Yuna", "Minseo", "Jiyeon", "Eunji", "Soyeon",
	"Hayoung", "Yeji", "Dahyun", "Seulgi", "Nayeon", "Jisoo", "Jieun", "Eunseo", "Chaeyoung", "Sumin",
	"Yejin", "Hana", "Hyerin", "Jimin", "Bomin", "Sora", "Yuri", "Sena", "Mina", "Euna",
}

var firstNeutral = []string{
	"Jiwon", "Jimin", "Hana", "Yuna", "Mina", "Yuri", "Sora", "Hyun", "Jun", "Eun",
}

// Curated surnames (romanized; common family names).
var lastNames = []string{
	"Kim", "Lee", "Park", "Choi", "Jung", "Kang", "Cho", "Yoon", "Jang", "Lim",
	"Han", "Oh", "Seo", "Shin", "Kwon", "Hwang", "Ahn", "Song", "Ryu", "Hong",
	"Yang", "Ko", "Moon", "Baek", "Heo", "Nam", "Jeon", "Bae", "No", "Min",
}

// Procedural building blocks (very simplified romanization).
// Korean romanized given names commonly combine 2 syllables.
var initials = []string{
	"g", "k", "n", "d", "t", "r", "m", "b", "p", "s", "j", "ch", "h",
	"kk", "tt", "pp", "ss", "jj", // stylized double consonants (rarely used here)
	"", "", "", // allow vowel-start sometimes
}

var vowels = []string{
	"a", "ae", "ya", "yae", "eo", "e", "yeo", "ye", "o", "wa", "wae", "oe",
	"u", "wo", "we", "wi", "yu", "eu", "ui", "i",
}

var finals = []string{
	"", "", "", "", // open syllables are common in romanized style names
	"n", "m", "ng", "k", "t", "l", "r", "s",
}

func (p koreanProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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

	genSyllable := func() string {
		ini := api.PickRand(initials, r)
		v := api.PickRand(vowels, r)
		fin := api.PickRand(finals, r)

		// Light cleanup: avoid awkward doubled-double initials too often.
		if (ini == "kk" || ini == "tt" || ini == "pp" || ini == "ss" || ini == "jj") && r.Intn(100) < 75 {
			ini = strings.TrimLeft(ini, "kptsj")
			if ini == "" {
				ini = "k"
			}
		}
		// Avoid empty + ui too often
		if ini == "" && v == "ui" && r.Intn(100) < 70 {
			v = "i"
		}
		return ini + v + fin
	}

	genGivenProcedural := func() string {
		// Typically 2 syllables; sometimes 3 at low realism.
		n := 2
		if realism < 40 && r.Intn(100) < 25 {
			n = 3
		}
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString(genSyllable())
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
		if cfg.Family == "" || strings.EqualFold(cfg.Family, PROFILE) {
			if chooseFromReal() {
				last = api.PickRand(lastNames, r)
			} else {
				// Korean surnames are usually one syllable; keep it short.
				s := genSyllable()
				// Force shorter-ish surname by trimming to first 2-5 chars
				if len(s) > 5 {
					s = s[:5]
				}
				last = caser.String(s)
			}
		} else {
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
var Profile koreanProfile
