package config

import (
	"os"

	"github.com/cloudcarver/edc/conf"
)

type Pg struct {
	// The DSN (Data Source Name) for postgres database connection. If specified, Host, Port, User, Password, and Db settings will be ignored.
	DSN *string `yaml:"dsn,omitempty"`

	// The host of the postgres database
	Host string `yaml:"host"`

	// The user of the postgres database
	User string `yaml:"user"`

	// The password of the postgres database
	Password string `yaml:"password"`

	// The database of the postgres database
	Db string `yaml:"db"`

	// The port of the postgres database
	Port int `yaml:"port"`
}

type Jwt struct {
	// The secret of the jwt. If not set, a random secret will be used.
	Secret string `yaml:"secret"`
}

type Root struct {
	// The password of the root user, if not set, the default password is "123456"
	Password string `yaml:"password"`
}

type Config struct {
	// The path of file to store the initialization data
	Init string `yaml:"init,omitempty"`

	// The port of the wavekit server
	Port int `yaml:"port,omitempty"`

	// The jwt configuration
	Jwt Jwt `yaml:"jwt,omitempty"`

	// The postgres configuration
	Pg Pg `yaml:"pg,omitempty"`

	// The root user configuration
	Root *Root `yaml:"root,omitempty"`

	// Whether to disable internet access
	NoInternet bool `yaml:"nointernet,omitempty"`

	// The path of the directory to store the risectl files
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
