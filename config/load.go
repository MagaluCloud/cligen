package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig() (*Config, error) {

	yamlFile, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("erro ao ler o arquivo de configuração: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		return nil, fmt.Errorf("erro ao deserializar o arquivo de configuração: %w", err)
	}

	return &config, nil
}

func (c *Config) SaveConfig() error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("erro ao serializar o arquivo de configuração: %w", err)
	}
	err = os.WriteFile("config/config_new.yaml", data, 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar o arquivo de configuração: %w", err)
	}
	return nil

}
