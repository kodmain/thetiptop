package assets

import "embed"

// Utilisez `go:embed` pour embarquer les fichiers de template.
//
//go:embed mails/*
var Mails embed.FS
