package main

import "embed"

//go:embed templates
var embeddedTemplates embed.FS

//go:embed build/assets
var embeddedAssets embed.FS