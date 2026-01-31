package api

import (
	"errors"
	"math/rand"
	"sort"
	"strings"
	"sync"
	"time"
)

// NewRand returns a deterministic RNG when cfg.Seed != 0.
// When cfg.Seed == 0, it returns a time-seeded RNG.
func NewRand(cfg ProfileConfig) *rand.Rand {
	if cfg.Seed != 0 {
		return rand.New(rand.NewSource(cfg.Seed))
	}
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}

var (
	profilesMu         sync.RWMutex
	profiles           = map[string]NameProfile{}
	ErrProfileNotFound = errors.New("profile not found")
)

// RegisterProfile registers a new profile
func RegisterProfile(name string, p NameProfile) {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		panic("api.RegisterProfile: empty name")
	}
	if p == nil {
		panic("api.RegisterProfile: nil profile")
	}
	profilesMu.Lock()
	defer profilesMu.Unlock()
	profiles[name] = p
}

// GetProfile returns the given profile as
func GetProfile(name string) (NameProfile, error) {
	name = strings.TrimSpace(strings.ToLower(name))
	profilesMu.RLock()
	defer profilesMu.RUnlock()

	p, ok := profiles[name]
	if !ok {
		return nil, ErrProfileNotFound
	}
	return p, nil
}

func ListProfiles() []string {
	profilesMu.RLock()
	defer profilesMu.RUnlock()
	out := make([]string, 0, len(profiles))
	for k := range profiles {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// ProfileConfig holds runtime options the main binary passes to the plugin.
type ProfileConfig struct {
	Count       int    `json:"count,omitempty"`
	Mode        string `json:"mode,omitempty"`
	Seed        int64  `json:"seed,omitempty"`        // 0 for random
	Realism     int    `json:"realism,omitempty"`     // 0..100
	Gender      string `json:"gender,omitempty"`      // "male", "female", "neutral"
	Family      string `json:"family,omitempty"`      // optional family override like "japan", "nordic", etc.
	IncludeLast bool   `json:"includeLast,omitempty"` // -l flag
	Reverse     bool   `json:"reverse,omitempty"`     // -r flag
	DevMode     bool   `json:"devMode,omitempty"`
}

// NameResult is returned by plugin when asked to generate a name.
type NameResult struct {
	First string
	Last  string // may be empty if plugin doesn't generate surnames
}

// NameProfile is the interface plugin must expose as a symbol (e.g. "Profile").
// Plugins should export a variable named "Profile" of this type.
type NameProfile interface {
	// Generate returns a NameResult obeying the provided ProfileConfig.
	Generate(cfg ProfileConfig) (NameResult, error)

	// Info returns human-readable metadata: supported family keys, language name, notes.
	Info() map[string]string
}
