package celtic

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type celticProfile struct{}

const PROFILE = "celtic"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p celticProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Celtic-inspired (Irish/Scottish/Welsh) names (ASCII): curated + procedural fallback; deterministic",
	}
}

// Curated: common Irish/Scottish/Welsh given names (ASCII only; no accents).
var givenMale = []string{
	"Sean", "Liam", "Conor", "Ciaran", "Eoin", "Niall", "Fionn", "Declan", "Ronan", "Cormac",
	"Aidan", "Patrick", "Donal", "Darragh", "Colm", "Padraig", "Gavin", "Owain", "Rhys", "Dylan",
	"Alasdair", "Callum", "Ewan", "Angus", "Fergus",
}

var givenFemale = []string{
	"Siobhan", "Aoife", "Niamh", "Saoirse", "Orla", "Maeve", "Deirdre", "Brigid", "Grainne", "Aisling",
	"Ciara", "Eimear", "Fiona", "Mairead", "Roisin", "Keira", "Erin", "Bronagh", "Catriona", "Gwen",
	"Rhian", "Sian", "Eleri", "Megan", "Bethan",
}

var givenNeutral = []string{
	"Rowan", "Morgan", "Rory", "Erin", "Gavin", "Fiona", "Rhys", "Dylan", "Aidan", "Maeve",
}

// Curated surnames + patronymic prefixes (Mac/Mc/O'/ap/fitz).
var surnames = []string{
	"Murphy", "Kelly", "OBrien", "ONeill", "Byrne", "Ryan", "Walsh", "Sullivan", "Doyle", "McCarthy",
	"MacLeod", "MacDonald", "Campbell", "Stewart", "Fraser", "Sinclair", "MacKenzie", "Douglas",
	"Jones", "Evans", "Williams", "Davies", "Morgan", "Thomas",
}

var patronymicPrefixes = []string{"Mac", "Mc", "O", "Fitz", "Ap"}

// Procedural blocks (Celtic-ish phonotactics; simplified).
var onsets = []string{
	"", "",
	"b", "c", "d", "f", "g", "h", "k", "l", "m", "n", "p", "r", "s", "t", "v", "w", "y",
	"br", "cr", "dr", "fr", "gr", "tr",
	"cl", "gl", "pl", "sl",
	"ch",
}

var vowels = []string{
	"a", "e", "i", "o", "u", "y",
	"ae", "ai", "ao", "ea", "ei", "eo", "ia", "ie", "io", "oa", "oi", "ou", "ua", "ui",
}

var codas = []string{
	"", "", "", "",
	"n", "m", "r", "l", "s", "t", "d", "g", "k",
	"nn", "ll", "rr",
	"ch", "sh",
}

var givenEndings = []string{"", "", "", "an", "en", "in", "on", "ach", "aidh", "wyn", "wen"}
var surnameEndings = []string{"", "", "", "son", "ley", "lan", "nan", "don", "more", "ford"}

func (p celticProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
		n := 2 + r.Intn(2) // 2..3
		if realism < 40 {
			n = 1 + r.Intn(3)
		}
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString(genSyl())
		}
		if r.Intn(100) < 40 {
			b.WriteString(api.PickRand(givenEndings, r))
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
		if r.Intn(100) < 45 {
			b.WriteString(api.PickRand(surnameEndings, r))
		}
		return b.String()
	}

	// First
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

	// Last
	last := ""
	if cfg.IncludeLast {
		if chooseFromReal() {
			// Some chance to fabricate a patronymic: Prefix + CuratedSurname (no punctuation)
			if r.Intn(100) < 35 {
				pfx := api.PickRand(patronymicPrefixes, r)
				base := api.PickRand(surnames, r)
				base = strings.ReplaceAll(base, " ", "")
				// "O" is typically "O" + base without apostrophe in ASCII mode.
				last = pfx + base
			} else {
				last = api.PickRand(surnames, r)
			}
		} else {
			last = caser.String(genSurnameProcedural())
		}
	}

	return api.NameResult{First: caser.String(first), Last: caser.String(last)}, nil
}

var Profile celticProfile
