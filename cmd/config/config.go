package config

import "log"

type ZSConfig struct {
	URL      string
	Admin    string
	Password string
}

type Config struct {
	Port      int
	Env       string
	Zs        ZSConfig
	IndexPath string
}

type Application struct {
	Config Config
	Logger *log.Logger
}
