package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Debug            bool
	HTTPPort         int    `split_words:"true" default:"8080"`
	DBHost           string `split_words:"true"`
	DBUser           string `split_words:"true"`
	DBPass           string `split_words:"true"`
	DBName           string `split_words:"true"`
	DBMigrationsPath string `split_words:"true"`
}

func loadConfig() (*config, error) {
	var cfg config
	err := envconfig.Process("", &cfg)
	return &cfg, err
}

func (c *config) connectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", c.DBUser, c.DBPass, c.DBHost, c.DBName)
}
