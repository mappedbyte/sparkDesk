package main

import "net/http"

type Config struct {
	Cookie  string `yaml:"cookie"`
	Fd      string `yaml:"fd"`
	GtToken string `yaml:"GtToken"`
}

type SparkWeb struct {
	Config      Config
	ChatId      string
	ChatHistory []map[string]string
	HttpClient  http.Client
}
