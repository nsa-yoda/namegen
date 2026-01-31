package kazakh

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type kazakhProfile struct{}

const PROFILE = "kazakh"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p kazakhProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Kazakh/Central Asian inspired names (ASCII): curated + procedural fallback; deterministic",
	}
}

// Curated: common Kazakh given names (ASCII transliteration).
var givenMale = []string{
	"Alikhan", "Nursultan", "Arman", "Bekzat", "Dias", "Erlan", "Yerlan", "Serik", "Timur", "Aidar",
	"Kanat", "Daniyar", "Marat", "Nurbol", "Sanzhar", "Azamat", "Bolat", "Bauyrzhan", "Mukhtar", "Zhanibek",
}

var givenFemale = []string{
	"Aigul", "Aigerim", "Dana", "Dinara", "Gulnaz", "Madina", "Aruzhan", "Zarina", "Assel", "Aisulu",
	"Malika", "Kamila", "Amina", "Sholpan", "Saule", "Gulnara", "Aliya", "Ainur", "Zhanna", "Karlygash",
}

var givenNeutral = []string{
	"Dana", "Amina", "Timur", "Arman", "Madina", "Aliya", "Dias", "Ainur", "Zarina", "Azamat",
}

// Curated surnames and common suffix styles.
var surnames = []string{
	"Nurpeisov", "Suleimenov", "Kudaibergenov", "Kenzhebekov", "Serikov", "Tursunov", "Abdullayev", "Iskakov",
	"Zhaksylykov", "Omarov", "Akhmetov", "Beketov", "Zhaparov", "Sadykov", "Bekturov",
}

// Procedural blocks (Turkic-ish; simplified, ASCII only).
var onsets = []string{
	"", "",
	"b", "d", "g", "h", "j", "k", "l", "m", "n", "p", "q", "r", "s", "t", "v", "y", "z",
	"zh", "sh", "kh", "ch",
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

var givenEndingsMale = []string{"", "", "", "bek", "khan", "bay", "mir", "nur", "lan"}
var givenEndingsFemale = []string{"", "", "", "gul", "nur", "ai", "ana", "ya"}
var givenEndingsNeutral = []string{"", "", "", "nur", "ai", "an"}

var surnameEndings = []string{"ov", "ova", "ev", "eva", "bekov", "bayev", "uly", "kyzy"}

func (p kazakhProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
		n := 2 + r.Intn(2) // 2..3
		if realism < 40 {
			n = 1 + r.Intn(3)
		}
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString(genSyl())
		}

		// Frequently add a surname suffix.
		if r.Intn(100) < 80 {
			sfx := api.PickRand(surnameEndings, r)
			// avoid double-suffix look
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

var Profile kazakhProfile
