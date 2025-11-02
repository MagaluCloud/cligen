package config

import (
	"fmt"
	"os"
	"path"

	"github.com/magaluCloud/mgccli/cmd/common/structs"
	"github.com/magaluCloud/mgccli/cmd/common/workspace"
	"gopkg.in/yaml.v3"
)

type Config interface {
	Get(name string) (any, error)
	Set(name string, value any) error
	Delete(name string) error
	List() (map[string]any, error)
	Write() error
}

type CliConfig struct {
	ChunkSize        int    `yaml:"chunk_size" default:"8" validate:"minimum=8,maximum=5120"`
	Workers          int    `yaml:"workers" default:"5" validate:"minimum=1"`
	DefaultOutput    string `yaml:"default_output,omitempty"`
	Region           string `yaml:"region" default:"br-se1"`
	Env              string `yaml:"env" default:"prod"`
	Debug            bool   `yaml:"debug,omitempty"`
	NoConfirm        bool   `yaml:"no_confirm,omitempty"`
	RawOutput        bool   `yaml:"raw_output,omitempty"`
	Lang             string `yaml:"lang" default:"en-US"`
	ServerURL        string `yaml:"server_url,omitempty"`
	VersionLastCheck string `yaml:"version_last_check,omitempty"`
}

type config struct {
	cliConfig CliConfig
	workspace workspace.Workspace
}

func NewConfig(workspace workspace.Workspace) Config {
	configFile := path.Join(workspace.Dir(), "cli.yaml")
	cliConfig, err := loadConfig(configFile)
	if err != nil {
		//TODO: Handle error
		panic(err)
	}
	return &config{workspace: workspace, cliConfig: cliConfig}
}

func (c *config) Get(name string) (any, error) {
	configMap, err := structs.StructToMap(c.cliConfig)
	if err != nil {
		return "", err
	}
	value, ok := configMap[name]
	if !ok {
		return "", fmt.Errorf("config %s not found", name)
	}
	return value, nil
}

func (c *config) Set(name string, value any) error {
	err := structs.Set(&c.cliConfig, name, value)
	if err != nil {
		return err
	}
	err = c.Write()
	if err != nil {
		return err
	}
	return nil
}

func (c *config) Delete(name string) error {
	return c.Set(name, nil)
}

func (c *config) Write() error {

	data, err := yaml.Marshal(c.cliConfig)
	if err != nil {
		return err
	}
	err = os.WriteFile(path.Join(c.workspace.Dir(), "cli.yaml"), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *config) List() (map[string]any, error) {
	configMap, err := structs.StructToMap(c.cliConfig)
	if err != nil {
		return nil, err
	}
	return configMap, nil
}

func loadConfig(configPath string) (CliConfig, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return structs.InitConfig[CliConfig](), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return CliConfig{}, err
	}

	var cliConfig CliConfig
	err = yaml.Unmarshal(data, &cliConfig)
	if err != nil {
		return CliConfig{}, err
	}
	return cliConfig, nil
}
