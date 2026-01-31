package portuguese

import (
	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type portugueseProfile struct{}

const PROFILE = "portuguese"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p portugueseProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Portuguese (Portugal/Brazil) names, ASCII",
	}
}

var firstMale = []string{
	"Joao", "Pedro", "Lucas", "Mateus", "Rafael", "Bruno", "Tiago", "Andre",
	"Diego", "Felipe", "Gustavo", "Carlos", "Daniel", "Eduardo", "Fernando",
}

var firstFemale = []string{
	"Maria", "Ana", "Beatriz", "Carla", "Patricia", "Juliana", "Fernanda",
	"Camila", "Renata", "Luciana", "Paula", "Daniela", "Larissa", "Bianca",
}

var firstNeutral = []string{
	"Ariel", "Alex", "Noa", "Dani", "Rene",
}

var lastNames = []string{
	"Silva", "Santos", "Oliveira", "Pereira", "Costa", "Rodrigues",
	"Alves", "Lima", "Gomes", "Ribeiro", "Carvalho", "Souza",
	"Martins", "Araujo", "Rocha",
}

var onsets = []string{
	"b", "c", "d", "f", "g", "l", "m", "n", "p", "r", "s", "t", "v",
	"br", "cr", "tr", "pr", "cl",
	"", "",
}

var vowels = []string{
	"a", "e", "i", "o", "u", "ai", "ei", "ao", "ou",
}

var codas = []string{
	"", "", "", "s", "r", "l", "m", "n",
}

func (p portugueseProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
	r := api.NewRand(cfg)
	caser := cases.Title(language.Und)

	useReal := r.Intn(100) < cfg.Realism+10

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
			last = caser.String(gen() + gen())
		}
	}

	return api.NameResult{First: first, Last: last}, nil
}

var Profile portugueseProfile
