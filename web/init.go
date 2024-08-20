package web

import "embed"

//go:embed template/*
var ViewTemplates embed.FS

//go:embed static/*
var StaticContent embed.FS
