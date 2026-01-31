package hindi

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type hindiProfile struct{}

const PROFILE = "hindi"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p hindiProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Hindi / North Indian names (romanized, ASCII)",
	}
}

var firstMale = []string{
	"Rahul", "Amit", "Vikram", "Arjun", "Rohit", "Suresh", "Anil", "Rajesh",
	"Manish", "Sanjay", "Deepak", "Kunal", "Nitin", "Ashok", "Pradeep",
	"Vijay", "Rakesh", "Sachin", "Anand", "Harish",
}

var firstFemale = []string{
	"Priya", "Anita", "Sunita", "Pooja", "Neha", "Kavita", "Ritu", "Asha",
	"Rekha", "Suman", "Meena", "Anjali", "Shilpa", "Nisha", "Seema",
	"Divya", "Kiran", "Jyoti", "Sarita", "Rashmi",
}

var firstNeutral = []string{
	"Kiran", "Ravi", "Aman", "Arya", "Nikhil", "Dev", "Shiv", "Rani",
}

var lastNames = []string{
	"Sharma", "Verma", "Gupta", "Singh", "Kumar", "Agarwal", "Mishra",
	"Yadav", "Chaudhary", "Patel", "Malhotra", "Kapoor", "Khanna",
	"Mehta", "Bansal", "Joshi", "Pandey", "Tiwari", "Goyal", "Jain",
}

var onsets = []string{
	"b", "bh", "d", "dh", "g", "gh", "k", "kh", "m", "n", "p", "ph",
	"r", "s", "sh", "t", "th", "v", "y", "ch", "j",
	"", "",
}

var vowels = []string{
	"a", "aa", "i", "ee", "u", "oo", "e", "ai", "o", "au",
}

var codas = []string{
	"", "", "", "n", "m", "r", "sh", "t", "k",
}

func (p hindiProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
	r := api.NewRand(cfg)
	caser := cases.Title(language.Und)

	realism := cfg.Realism
	if realism < 0 {
		realism = 0
	}
	if realism > 100 {
		realism = 100
	}

	realPct := 10 + realism
	if realPct > 95 {
		realPct = 95
	}
	useReal := func() bool { return r.Intn(100) < realPct }

	genSyl := func() string {
		return api.PickRand(onsets, r) +
			api.PickRand(vowels, r) +
			api.PickRand(codas, r)
	}

	genGiven := func() string {
		n := 2
		if realism < 40 && r.Intn(100) < 30 {
			n = 3
		}
		var b strings.Builder
		for i := 0; i < n; i++ {
			b.WriteString(genSyl())
		}
		return b.String()
	}

	first := ""
	if useReal() {
		switch cfg.Gender {
		case "male":
			first = api.PickRand(firstMale, r)
		case "female":
			first = api.PickRand(firstFemale, r)
		default:
			first = api.PickRand(firstNeutral, r)
		}
	} else {
		first = caser.String(genGiven())
	}

	last := ""
	if cfg.IncludeLast {
		if useReal() {
			last = api.PickRand(lastNames, r)
		} else {
			last = caser.String(genSyl() + genSyl())
		}
	}

	return api.NameResult{
		First: caser.String(first),
		Last:  caser.String(last),
	}, nil
}

var Profile hindiProfile
