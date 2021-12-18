package assets

import "embed"

//go:embed *.html css/*
var Content embed.FS
