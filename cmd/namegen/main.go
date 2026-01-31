package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/nsa-yoda/namegen/api"
	_ "github.com/nsa-yoda/namegen/plugins/amharic"
	_ "github.com/nsa-yoda/namegen/plugins/arabic"
	_ "github.com/nsa-yoda/namegen/plugins/aramaic"
	_ "github.com/nsa-yoda/namegen/plugins/baltic"
	_ "github.com/nsa-yoda/namegen/plugins/celtic"
	_ "github.com/nsa-yoda/namegen/plugins/chinese"
	_ "github.com/nsa-yoda/namegen/plugins/english"
	_ "github.com/nsa-yoda/namegen/plugins/farsi"
	_ "github.com/nsa-yoda/namegen/plugins/filipino"
	_ "github.com/nsa-yoda/namegen/plugins/french"
	_ "github.com/nsa-yoda/namegen/plugins/germanic"
	_ "github.com/nsa-yoda/namegen/plugins/greek"
	_ "github.com/nsa-yoda/namegen/plugins/hawaiian"
	_ "github.com/nsa-yoda/namegen/plugins/hebrew"
	_ "github.com/nsa-yoda/namegen/plugins/hindi"
	_ "github.com/nsa-yoda/namegen/plugins/igbo"
	_ "github.com/nsa-yoda/namegen/plugins/indonesian"
	_ "github.com/nsa-yoda/namegen/plugins/italian"
	_ "github.com/nsa-yoda/namegen/plugins/japanese"
	_ "github.com/nsa-yoda/namegen/plugins/kazakh"
	_ "github.com/nsa-yoda/namegen/plugins/korean"
	_ "github.com/nsa-yoda/namegen/plugins/malay"
	_ "github.com/nsa-yoda/namegen/plugins/maori"
	_ "github.com/nsa-yoda/namegen/plugins/nahuatl"
	_ "github.com/nsa-yoda/namegen/plugins/nordic"
	_ "github.com/nsa-yoda/namegen/plugins/portuguese"
	_ "github.com/nsa-yoda/namegen/plugins/samoan"
	_ "github.com/nsa-yoda/namegen/plugins/slavic"
	_ "github.com/nsa-yoda/namegen/plugins/spanish"
	_ "github.com/nsa-yoda/namegen/plugins/swahili"
	_ "github.com/nsa-yoda/namegen/plugins/tamil"
	_ "github.com/nsa-yoda/namegen/plugins/thai"
	_ "github.com/nsa-yoda/namegen/plugins/turkish"
	_ "github.com/nsa-yoda/namegen/plugins/uzbek"
	_ "github.com/nsa-yoda/namegen/plugins/vietnamese"
	_ "github.com/nsa-yoda/namegen/plugins/yoruba"
)

// defaultFallbackGenerator
const defaultFallbackGenerator = "english"

func main() {
	// CLI flags
	mode := flag.String("mode", "english", "Mode/profile name (compiled-in). If not found, English is used by default.")
	includeLast := flag.Bool("l", false, "Include last name")
	reverse := flag.Bool("r", false, "Reverse order (last first)")
	gender := flag.String("gender", "neutral", "Gender: male|female|neutral")
	family := flag.String("family", "", "Family override for surname rules (e.g., japan, nordic, spanish)")
	realism := flag.Int("realism", 50, "Realism 0..100 (0 fictional phonotactics, 100 real-looking names)")
	seed := flag.Int64("s", 0, "Seed (0 or omit for random)")
	count := flag.Int("c", 1, "Number of names to generate, 1 by default or omitted")
	devMode := flag.Bool("d", false, "Development mode")
	flag.Parse()

	cfg := api.ProfileConfig{
		Count:       *count,
		Mode:        *mode,
		Seed:        *seed,
		Realism:     *realism,
		Gender:      *gender,
		Family:      *family,
		IncludeLast: *includeLast,
		Reverse:     *reverse,
		DevMode:     *devMode,
	}

	if *devMode {
		b, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			log.Printf("devMode: could not marshal config: %v\n", err)
		} else {
			fmt.Printf("Dev Mode Active - Config:\n%s\n", string(b))
		}
	}

	// Load our chosen profile (compiled-in registry)
	profile, err := api.GetProfile(*mode)
	if err != nil {
		log.Printf("profile not found for mode %q — using builtin fallback %q\n", *mode, defaultFallbackGenerator)
		profile, err = api.GetProfile(defaultFallbackGenerator)
		if err != nil {
			log.Fatalf("builtin fallback profile not found for mode %q\n", defaultFallbackGenerator)
		}
	}

	n := cfg.Count
	if n <= 0 {
		n = 1
	}

	for i := 0; i < n; i++ {
		// Generate and print
		res, err := profile.Generate(cfg)
		if err != nil {
			log.Fatalf("generate failed: %v", err)
		}

		if cfg.IncludeLast {
			if cfg.Reverse {
				fmt.Printf("%s %s\n", res.Last, res.First)
			} else {
				fmt.Printf("%s %s\n", res.First, res.Last)
			}
		} else {
			fmt.Println(res.First)
		}
	}
}
