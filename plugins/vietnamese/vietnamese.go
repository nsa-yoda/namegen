package vietnamese

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type vietnameseProfile struct{}

const PROFILE = "vietnamese"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p vietnameseProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Vietnamese names (ASCII romanization): realism blends curated lists with procedural syllables; deterministic with seed",
	}
}

// Note: Vietnamese naming convention is typically Family (surname) + Middle + Given.
// This generator returns First + Last; here we treat \"First\" as given name and \"Last\" as surname,
// with an optional middle-like component folded into First at lower realism.

var givenMale = []string{
	"Anh", "Bao", "Binh", "Cuong", "Duc", "Hieu", "Hoang", "Hung", "Khanh", "Khoa",
	"Long", "Minh", "Nam", "Phuc", "Quan", "Son", "Tuan", "Viet", "Thanh", "Thien",
	"Dat", "Kiet", "Lam", "Luan", "Nghia", "Phu", "Tai", "Trung", "Vu", "Xuan",
}

var givenFemale = []string{
	"An", "Chi", "Diem", "Dung", "Giang", "Han", "Hanh", "Hoa", "Huong", "Lan",
	"Linh", "Mai", "My", "Nga", "Ngoc", "Nhi", "Phuong", "Quynh", "Thao", "Trang",
	"Thuy", "Tien", "Trinh", "Tuyet", "Vy", "Yen", "Ha", "Hien", "Kim", "Thao",
}

var givenNeutral = []string{
	"Anh", "Khanh", "Linh", "Minh", "An", "Chi", "Giang", "Ha", "My", "Vy",
}

// Very common Vietnamese surnames (ASCII).
var surnames = []string{
	"Nguyen", "Tran", "Le", "Pham", "Huynh", "Hoang", "Phan", "Vu", "Vo", "Dang",
	"Bui", "Do", "Ho", "Ngo", "Duong", "Ly", "Dinh", "Truong", "Ha", "Dao",
}

// Common middle names (often gendered, but used flexibly here).
var middles = []string{
	"Van", "Thi", "Huu", "Gia", "Quoc", "Duc", "Minh", "Ngoc", "Thanh", "Thuy",
}

var onsets = []string{
	"", "",
	"b", "c", "d", "g", "h", "k", "l", "m", "n", "p", "q", "r", "s", "t", "v", "x",
	"ch", "ng", "nh", "ph", "th", "tr",
}
var vowels = []string{
	"a", "e", "i", "o", "u", "y", "ai", "ao", "au", "ea", "eo", "ia", "ie", "oa", "oi", "oo", "ua", "uo",
}
var codas = []string{
	"", "", "", "", // many Vietnamese syllables are open in romanized text
	"n", "m", "ng", "nh", "t", "c", "p",
}

var givenEndings = []string{"", "", "", "h", "n", "t", "ng"}

func (p vietnameseProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
	r := api.NewRand(cfg)
	caser := cases.Title(language.Und)

	realism := cfg.Realism
	if realism < 0 {
		realism = 0
	}
	if realism > 100 {
		realism = 100
	}

	// same curve as your other profiles
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
		// Keep it compact: onset + vowel + optional coda.
		s := api.PickRand(onsets, r) + api.PickRand(vowels, r) + api.PickRand(codas, r)
		if r.Intn(100) < 25 {
			s += api.PickRand(givenEndings, r)
		}
		return s
	}

	genGivenProcedural := func() string {
		// Vietnamese given names are often 1 syllable; sometimes 2.
		n := 1
		if realism < 40 && r.Intn(100) < 40 {
			n = 2
		} else if r.Intn(100) < 20 {
			n = 2
		}
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString(genSyl())
		}
		return b.String()
	}

	// ---- Given name (First) ----
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

	// Optionally add a middle name at mid/low realism (or when realism is high but random says so).
	// We fold it into First as "First Middle" so output remains First/Last.
	if r.Intn(100) < 35 && realism < 85 {
		mid := api.PickRand(middles, r)
		first = caser.String(first) + " " + caser.String(mid)
	}

	// ---- Surname (Last) ----
	last := ""
	if cfg.IncludeLast {
		if chooseFromReal() {
			last = api.PickRand(surnames, r)
		} else {
			// Surname-like: usually 1 syllable but can be 2 in procedural mode
			n := 1
			if r.Intn(100) < 25 {
				n = 2
			}
			var b strings.Builder
			for i := 0; i < n; i++ {
				b.WriteString(genSyl())
			}
			last = caser.String(b.String())
		}
	}

	return api.NameResult{First: caser.String(first), Last: caser.String(last)}, nil
}

var Profile vietnameseProfile
