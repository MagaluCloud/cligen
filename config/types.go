package config

// Config representa a estrutura principal do arquivo config.yaml
type Config struct {
	Menus []Menu `yaml:"menus"`
}

// Menu representa um menu principal
type Menu struct {
	Name       string   `yaml:"name,omitempty"`
	SDKPackage string   `yaml:"sdk_package,omitempty"`
	Alias      []string `yaml:"alias,omitempty"`
	Menus      []Menu   `yaml:"menus,omitempty"`
}
