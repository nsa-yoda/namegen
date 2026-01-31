package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/nsa-yoda/namegen/api"
	_ "github.com/nsa-yoda/namegen/plugins/english"
	_ "github.com/nsa-yoda/namegen/plugins/japanese"
	_ "github.com/nsa-yoda/namegen/plugins/spanish"
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
	flag.Parse()

	cfg := api.ProfileConfig{
		Seed:        *seed,
		Realism:     *realism,
		Gender:      *gender,
		Family:      *family,
		IncludeLast: *includeLast,
		Reverse:     *reverse,
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
