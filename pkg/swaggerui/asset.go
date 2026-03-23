package swaggerui

import "embed"

//go:embed asset/*
var SwaggerAsset embed.FS

//go:embed scala/*
var ScalarAsset embed.FS
