package api

import (
	"math/rand"
	"time"
)

// RandLike is the minimal interface this project needs for deterministic randomness.
// *rand.Rand satisfies this.
type RandLike interface {
	Intn(n int) int
}

func PickRand[T any](arr []T, r RandLike) T {
	if len(arr) == 0 {
		panic("api.Pick: empty slice")
	}

	// Defensive: allow callers to pass nil; fall back to time-seeded RNG.
	if r == nil {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	return arr[r.Intn(len(arr))]
}

// NewRand returns a deterministic RNG when cfg.Seed != 0.
// When cfg.Seed == 0, it returns a time-seeded RNG.
func NewRand(cfg ProfileConfig) *rand.Rand {
	if cfg.Seed != 0 {
		return rand.New(rand.NewSource(cfg.Seed))
	}
	return rand.New(rand.NewSource(time.Now().UnixNano()))
}
