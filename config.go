package main

import (
	"github.com/BurntSushi/toml"
	"sort"
)

type config struct {
	ListenOn string

	TLS bool
	TLSPem string
	TLSKey string

	MaxFileSize int64
	Formats []string
}

var Config = config{
	ListenOn: ":80",

	MaxFileSize: 1_000_000,
	Formats: []string{"image/jpeg", "image/png", "image/webp"},
} // Default config, may be overwritten by toml file,


func loadConfig() (err error) {
	_, err = toml.DecodeFile("config.toml", &Config)

	if Config.Formats != nil {
		sort.Strings(Config.Formats)
	}
	return
}