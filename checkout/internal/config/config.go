package config

import (
	"os"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type ConfigStruct struct {
	Token    string `yaml:"token"`
	Services struct {
		LOMS           string `yaml:"loms"`
		ProductService string `yaml:"productservice"`
	} `yaml:"services"`
	GRPCPort    int    `yaml:"grpc_port"`
	HTTPPort    int    `yaml:"http_port"`
	DatabaseURL string `yaml:"database_url"`
}

var ConfigData ConfigStruct

func Init() error {
	rawYAML, err := os.ReadFile("config.yml")
	if err != nil {
		return errors.WithMessage(err, "reading config file")
	}

	err = yaml.Unmarshal(rawYAML, &ConfigData)
	if err != nil {
		return errors.WithMessage(err, "parsing yaml")
	}

	return nil
}
