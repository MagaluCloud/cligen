package config

import (
	"encoding/json"
	"fmt"
	"os"
)

func LoadConfig() (*Config, error) {

	jsonFile, err := os.ReadFile("config/config.json")
	if err != nil {
		return nil, fmt.Errorf("erro ao ler o arquivo de configuração: %w", err)
	}

	var config Config
	if err := json.Unmarshal(jsonFile, &config); err != nil {
		return nil, fmt.Errorf("erro ao deserializar o arquivo de configuração: %w", err)
	}

	return &config, nil
}

func (c *Config) SaveConfig() error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar o arquivo de configuração: %w", err)
	}
	err = os.WriteFile("config/config.json", data, 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar o arquivo de configuração: %w", err)
	}
	return nil

}
