// Package all registers all built-in namegen profiles.
//
// Importing this package has side effects: it blank-imports every built-in
// plugin so their init() functions run and they register themselves with
// the api registry.
//
// Usage:
//
//	import (
//		_ "github.com/nsa-yoda/namegen/all"
//		"github.com/nsa-yoda/namegen/api"
//	)
package all

import (
	// Blank imports force plugin packages to init(), when adding a new imoprt it must be a blank import
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
