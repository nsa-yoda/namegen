package uzbek

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type uzbekProfile struct{}

const PROFILE = "uzbek"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p uzbekProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Uzbek/Central Asian inspired names (ASCII): curated + procedural fallback; deterministic",
	}
}

var givenMale = []string{
	"Aziz", "Bekzod", "Jasur", "Sardor", "Rustam", "Shavkat", "Ulugbek", "Temur", "Akmal", "Dilshod",
	"Farrukh", "Kamol", "Bunyod", "Odil", "Asad", "Sherzod", "Islom", "Siroj", "Anvar", "Jamshid",
}

var givenFemale = []string{
	"Malika", "Dilnoza", "Madina", "Nigina", "Gulnora", "Zarina", "Aziza", "Saida", "Shahnoza", "Feruza",
	"Munisa", "Shirin", "Lola", "Gulbahor", "Sitora", "Nodira", "Nargiza", "Sevara", "Rayhona", "Laylo",
}

var givenNeutral = []string{
	"Aziz", "Aziza", "Madina", "Malika", "Dilshod", "Zarina", "Kamol", "Odil", "Shirin", "Anvar",
}

var surnames = []string{
	"Karimov", "Rakhimov", "Yusupov", "Abdullayev", "Ismailov", "Nazarov", "Tursunov", "Saidov", "Khodjaev", "Qodirov",
	"Aliyev", "Usmonov", "Soliev", "Mamatov", "Shukurov",
}

var onsets = []string{
	"", "",
	"b", "d", "f", "g", "h", "j", "k", "l", "m", "n", "p", "q", "r", "s", "t", "v", "y", "z",
	"sh", "ch", "kh", "zh",
	"br", "kr", "tr",
}

var vowels = []string{
	"a", "e", "i", "o", "u", "y",
	"ai", "au", "ia", "ie", "oi", "ua",
}

var codas = []string{
	"", "", "", "",
	"n", "m", "r", "l", "t", "k", "s",
	"ng",
}

var givenEndingsMale = []string{"", "", "", "bek", "jon", "mir", "khon", "dor", "shod"}
var givenEndingsFemale = []string{"", "", "", "a", "ya", "noza", "nora", "gul", "oy"}
var givenEndingsNeutral = []string{"", "", "", "a", "an", "bek", "oy"}

var surnameEndings = []string{"ov", "ova", "ev", "eva", "bekov", "bayev", "zoda"}

func (p uzbekProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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

	genSyl := func() string {
		return api.PickRand(onsets, r) + api.PickRand(vowels, r) + api.PickRand(codas, r)
	}

	genGivenProcedural := func() string {
		n := 2
		if realism < 40 {
			n = 1 + r.Intn(3)
		} else if r.Intn(100) < 30 {
			n = 2 + r.Intn(2)
		}
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString(genSyl())
		}
		switch cfg.Gender {
		case "male":
			b.WriteString(api.PickRand(givenEndingsMale, r))
		case "female":
			b.WriteString(api.PickRand(givenEndingsFemale, r))
		default:
			b.WriteString(api.PickRand(givenEndingsNeutral, r))
		}
		return b.String()
	}

	genSurnameProcedural := func() string {
		n := 2 + r.Intn(2)
		if realism < 40 {
			n = 1 + r.Intn(3)
		}
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString(genSyl())
		}
		if r.Intn(100) < 80 {
			sfx := api.PickRand(surnameEndings, r)
			if !strings.HasSuffix(b.String(), sfx) {
				b.WriteString(sfx)
			}
		}
		return b.String()
	}

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

	last := ""
	if cfg.IncludeLast {
		if chooseFromReal() {
			last = api.PickRand(surnames, r)
		} else {
			last = caser.String(genSurnameProcedural())
		}
	}

	return api.NameResult{First: caser.String(first), Last: caser.String(last)}, nil
}

var Profile uzbekProfile
