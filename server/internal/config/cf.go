package config

import (
	"fmt"
	"os"
)

func ErrEnvParamNotDefined(param string) error {
	return fmt.Errorf("Env param %v not defined", param)
}

type Config struct {
	Address string
}

func Get() (*Config, error) {
	cfg := &Config{}
	address := os.Getenv("ADDRESS")
	if address == "" {
		return cfg, ErrEnvParamNotDefined("ADDRESS")
	}
	cfg.Address = address
	return cfg, nil
}
