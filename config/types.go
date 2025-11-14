package config

// Config representa a estrutura principal do arquivo config.yaml
type Config struct {
	Version string `yaml:"version"`
	Menus   []Menu `yaml:"menus"`
}

// Menu representa um menu principal
type Menu struct {
	Name            string   `yaml:"name,omitempty"`
	Description     string   `yaml:"description,omitempty"`
	LongDescription string   `yaml:"long_description,omitempty"`
	SDKPackage      string   `yaml:"sdk_package,omitempty"`
	CliGroup        string   `yaml:"cli_group,omitempty"`
	Alias           []string `yaml:"alias,omitempty"`
	Menus           []Menu   `yaml:"menus,omitempty"`
}
