package web

import "embed"

// Dist 是前端构建产物（dist/）。
//
//go:embed all:dist
var Dist embed.FS
