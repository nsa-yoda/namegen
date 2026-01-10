package api

// ProfileConfig holds runtime options the main binary passes to the plugin.
type ProfileConfig struct {
	Seed        int64  // 0 for random
	Realism     int    // 0..100
	Gender      string // "male", "female", "neutral"
	Family      string // optional family override like "japan", "nordic", etc.
	IncludeLast bool   // -l flag
	Reverse     bool   // -r flag
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
