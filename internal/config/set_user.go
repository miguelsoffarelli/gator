package config

import (
	"encoding/json"
	"os"
)

func (cfg *Config) SetUser(username string) error {
	cfg.Current_user_name = username

	err := write(*cfg)
	if err != nil {
		return err
	}

	return nil
}

func write(cfg Config) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	if err = os.WriteFile(configFilePath, data, 0644); err != nil {
		return err
	}

	return nil
}
