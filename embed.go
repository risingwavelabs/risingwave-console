package root

import "embed"

//go:embed web/out/*
var StaticFiles embed.FS

//go:embed sql/migrations/*
var Migrations embed.FS
