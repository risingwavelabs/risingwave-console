package config

import (
	"os"

	"github.com/cloudcarver/edc/conf"
)

type Pg struct {
	// (Required) The DSN (Data Source Name) for postgres database connection. If specified, Host, Port, User, Password, and Db settings will be ignored.
	DSN *string `yaml:"dsn,omitempty"`
}

type Jwt struct {
	// (Optional) The secret of the jwt. If not set, a random secret will be used.
	Secret string `yaml:"secret"`
}

type Root struct {
	// (Optional) The password of the root user, if not set, the default password is "123456"
	Password string `yaml:"password"`
}

type Config struct {
	// (Optional) The path of file to store the initialization data, if not set, skip the initialization
	Init string `yaml:"init,omitempty"`

	// (Optional) The host of the wavekit server, it is used in the API endpoint of the web UI.
	// If not set, the host will be localhost.
	Host string `yaml:"host,omitempty"`

	// (Optional) The port of the wavekit server, default is 8020
	Port int `yaml:"port,omitempty"`

	// The jwt configuration
	Jwt Jwt `yaml:"jwt,omitempty"`

	// The postgres configuration
	Pg Pg `yaml:"pg,omitempty"`

	// The root user configuration
	Root *Root `yaml:"root,omitempty"`

	// (Optional) Whether to disable internet access, default is false. If public internet is not allowed, set it to true. Then mount risectl files to <risectl dir>/<version>/risectl.
	NoInternet bool `yaml:"nointernet,omitempty"`

	// (Optional) The path of the directory to store the risectl files, default is "$HOME/.risectl"
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
