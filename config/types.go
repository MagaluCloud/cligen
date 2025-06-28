package config

// Config representa a estrutura principal do arquivo config.yaml
type Config struct {
	Menus []Menu `yaml:"menus"`
}

// Menu representa um menu principal
type Menu struct {
	Name        string    `yaml:"name"`
	Alias       string    `yaml:"alias"`
	SDKPackage  string    `yaml:"sdk_package"`
	Description string    `yaml:"description"`
	Submenus    []Submenu `yaml:"submenus"`
}

// Submenu representa um submenu que pode ser recursivo
type Submenu struct {
	Name        string    `yaml:"name"`
	Description string    `yaml:"description"`
	Actions     []Action  `yaml:"actions"`
	Submenus    []Submenu `yaml:"submenus,omitempty"` // Recursivo - pode conter outros submenus
}

// Action representa uma ação disponível em um submenu
type Action struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}
