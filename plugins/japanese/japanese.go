package japanese

import (
	"strings"

	"github.com/nsa-yoda/namegen/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type japaneseProfile struct{}

const PROFILE = "japanese"

func init() {
	api.RegisterProfile(PROFILE, Profile)
}

func (p japaneseProfile) Info() map[string]string {
	return map[string]string{
		"name":  PROFILE,
		"notes": "Japanese names: realism blends curated romaji lists with kana-like procedural syllables; deterministic with seed",
	}
}

// Curated romaji lists (expand whenever you want).
// These are common/recognizable enough to feel “real” without being huge datasets.
var firstMale = []string{
	"Haruto", "Yuto", "Sota", "Yuki", "Koki", "Ren", "Kaito", "Takumi", "Daiki", "Ryota",
	"Yuma", "Riku", "Shota", "Tatsuya", "Kenta", "Keita", "Kazuki", "Shinji", "Hiroshi", "Taro",
	"Kenji", "Naoki", "Koji", "Masato", "Yusuke", "Hayato", "Shun", "Minato", "Itsuki", "Sora",
}

var firstFemale = []string{
	"Yui", "Aoi", "Sakura", "Hina", "Rin", "Mio", "Yuna", "Akari", "Hana", "Mei",
	"Nanami", "Rina", "Ayaka", "Haruka", "Miku", "Misaki", "Kaori", "Emi", "Nozomi", "Yoko",
	"Keiko", "Sachiko", "Naoko", "Maki", "Chihiro", "Reina", "Sumire", "Koharu", "Saki", "Natsumi",
}

var firstNeutral = []string{
	"Akira", "Hikaru", "Kaoru", "Makoto", "Nao", "Rei", "Ryo", "Sora", "Yu", "Haruka",
}

var lastNames = []string{
	"Sato", "Suzuki", "Takahashi", "Tanaka", "Watanabe", "Ito", "Yamamoto", "Nakamura", "Kobayashi", "Kato",
	"Yoshida", "Yamada", "Sasaki", "Yamaguchi", "Matsumoto", "Inoue", "Kimura", "Hayashi", "Shimizu", "Yamazaki",
	"Morita", "Okada", "Abe", "Fujita", "Ishikawa", "Hashimoto", "Ikeda", "Maeda", "Fukuda", "Ota",
}

// --- Procedural syllables ---
// Keep these “Japanese-feeling” (simple CV / some common clusters).
var vowels = []string{"a", "i", "u", "e", "o"}

var consonantOnsets = []string{
	"k", "s", "t", "n", "h", "m", "y", "r", "w",
	"g", "z", "d", "b", "p",
}

var clusters = []string{
	"ky", "gy",
	"sh", "ch", "j",
	"ny", "hy", "by", "py", "my", "ry",
	"ts",
}

// Some common endings to make results feel more name-like (still conservative).
var givenEndings = []string{"", "", "", "to", "ta", "ki", "shi", "ya", "na", "ko"} // ko can occur in female names
var surnameEndings = []string{"", "", "", "moto", "yama", "kawa", "zaki", "mura", "naka", "shita", "gawa"}

func (p japaneseProfile) Generate(cfg api.ProfileConfig) (api.NameResult, error) {
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
	// Ramp hard after 60 and very strong after 80.
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

	// ---- Procedural generators ----
	genCV := func() string {
		// 20% chance of cluster at higher realism; otherwise simple onset
		onset := api.PickRand(consonantOnsets, r)
		if realism >= 50 && r.Intn(100) < 20 {
			onset = api.PickRand(clusters, r)
		}
		return onset + api.PickRand(vowels, r)
	}

	genGivenProcedural := func() string {
		// 2–4 syllables normally; lower realism sometimes 1–3
		numSyl := 2 + r.Intn(3) // 2..4
		if realism < 40 {
			numSyl = 1 + r.Intn(3) // 1..3
		}

		var b strings.Builder
		for i := 0; i < numSyl; i++ {
			b.WriteString(genCV())
		}

		// occasional ending
		end := api.PickRand(givenEndings, r)
		if end != "" && !strings.HasSuffix(b.String(), end) {
			b.WriteString(end)
		}

		// small gender nuance: if explicitly female and high realism, bias toward "-ko" sometimes
		if cfg.Gender == "female" && realism >= 70 && r.Intn(100) < 15 {
			s := b.String()
			if !strings.HasSuffix(s, "ko") {
				b.Reset()
				b.WriteString(s)
				b.WriteString("ko")
			}
		}

		return b.String()
	}

	genSurnameProcedural := func() string {
		// 2–3 syllables surname-like
		numSyl := 2 + r.Intn(2) // 2..3
		var b strings.Builder
		for i := 0; i < numSyl; i++ {
			b.WriteString(genCV())
		}
		// add surname ending sometimes, more likely at higher realism
		if r.Intn(100) < 35 {
			end := api.PickRand(surnameEndings, r)
			if end != "" && !strings.HasSuffix(b.String(), end) {
				b.WriteString(end)
			}
		}
		return b.String()
	}

	// ---- Choose first name ----
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

	// ---- Choose last name ----
	last := ""
	if cfg.IncludeLast {
		// Allow cfg.Family override to force Japanese-style surname rules if you expand other modes later.
		if cfg.Family == "" || strings.EqualFold(cfg.Family, PROFILE) {
			if chooseFromReal() {
				last = api.PickRand(lastNames, r)
			} else {
				last = caser.String(genSurnameProcedural())
			}
		} else {
			// If Family override is set to something else, still produce a Japanese-ish surname
			// (keeps behavior stable, but you can refine later).
			if chooseFromReal() {
				last = api.PickRand(lastNames, r)
			} else {
				last = caser.String(genSurnameProcedural())
			}
		}
	}

	// final casing (safe if input already proper)
	first = caser.String(first)
	last = caser.String(last)

	return api.NameResult{First: first, Last: last}, nil
}

// Profile is the core exported symbol
var Profile japaneseProfile
