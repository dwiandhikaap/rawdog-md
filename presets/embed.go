package presets

import "embed"

//go:embed all:basic
var BasicPreset embed.FS

//go:embed all:skeleton
var SkeletonPreset embed.FS

//go:embed all:docs
var DocsPreset embed.FS
