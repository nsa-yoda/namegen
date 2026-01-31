package greek

import (
	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type greekProfile struct{}

const PROFILE = "greek"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p greekProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Greek names (romanized, ASCII)",
	}
}

var firstMale = []string{
	"Yannis", "Nikos", "Giorgos", "Dimitris", "Kostas", "Panagiotis",
	"Alexandros", "Stavros", "Christos", "Theodoros",
}

var firstFemale = []string{
	"Maria", "Eleni", "Katerina", "Sofia", "Anna", "Georgia",
	"Dimitra", "Ioanna", "Christina", "Eirini",
}

var firstNeutral = []string{
	"Alexis", "Niko", "Ari", "Danae",
}

var lastNames = []string{
	"Papadopoulos", "Nikolaidis", "Georgiou", "Dimitriou",
	"Christou", "Ioannou", "Kostopoulos", "Vasiliadis",
	"Panagiotou", "Theodorou",
}

var onsets = []string{
	"k", "g", "d", "t", "p", "m", "n", "l", "r", "s", "v",
	"ch", "th", "ps", "x",
	"", "",
}

var vowels = []string{
	"a", "e", "i", "o", "u", "ai", "ei", "oi", "ou",
}

var codas = []string{
	"", "", "", "s", "n", "r",
}

func (p greekProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
	r := api.NewRand(cfg)
	caser := cases.Title(language.Und)

	useReal := r.Intn(100) < cfg.Realism+15

	gen := func() string {
		return api.PickRand(onsets, r) +
			api.PickRand(vowels, r) +
			api.PickRand(codas, r)
	}

	first := ""
	if useReal {
		switch cfg.Gender {
		case "male":
			first = api.PickRand(firstMale, r)
		case "female":
			first = api.PickRand(firstFemale, r)
		default:
			first = api.PickRand(firstNeutral, r)
		}
	} else {
		first = caser.String(gen() + gen())
	}

	last := ""
	if cfg.IncludeLast {
		if useReal {
			last = api.PickRand(lastNames, r)
		} else {
			last = caser.String(gen() + gen() + "s")
		}
	}

	return api.NameResult{First: first, Last: last}, nil
}

var Profile greekProfile
