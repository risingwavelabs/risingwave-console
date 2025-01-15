package config

import (
	"os"

	"github.com/cloudcarver/edc/conf"
)

type Pg struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Db       string `yaml:"db"`
	Port     int    `yaml:"port"`
	// the path of directory to store migration files
	Migration string `yaml:"migration"`
}

type Jwt struct {
	Secret string `yaml:"secret"`
}

type Root struct {
	Password string `yaml:"password"`
}

type Config struct {
	Init string `yaml:"init,omitempty"`

	Port int `yaml:"port,omitempty"`

	Jwt Jwt `yaml:"jwt,omitempty"`
	Pg  Pg  `yaml:"pg,omitempty"`

	Root *Root `yaml:"root,omitempty"`

	NoInternet bool   `yaml:"nointernet,omitempty"`
	RisectlDir string `yaml:"risectldir,omitempty"`
}

func NewConfig() (*Config, error) {
	c := &Config{}
	if err := conf.FetchConfig((func() string {
		if _, err := os.Stat("config.yaml"); err != nil {
			return ""
		}
		return "config.yaml"
	})(), "WK_", c); err != nil {
		return nil, err
	}
	return c, nil
}
