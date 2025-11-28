package config

// Config representa a estrutura principal do arquivo config.yaml
type Config struct {
	Version string  `yaml:"version"`
	Menus   []*Menu `yaml:"menus"`
}

// Menu representa um menu principal
type Menu struct {
	Name            string   `yaml:"name,omitempty"`
	Enabled         bool     `yaml:"enabled,omitempty"`
	Description     string   `yaml:"description,omitempty"`
	LongDescription string   `yaml:"long_description,omitempty"`
	Level           int      `yaml:"level,omitempty"`
	SDKPackage      string   `yaml:"sdk_package,omitempty"`
	CliGroup        string   `yaml:"cli_group,omitempty"`
	Alias           []string `yaml:"alias,omitempty"`
	Menus           []*Menu  `yaml:"menus,omitempty"`
}

// type Parameter struct {
// 	Name            string               `json:"name" yaml:"parameter_name"`
// 	Type            string               `json:"type" yaml:"parameter_type"`
// 	Description     string               `json:"description" yaml:"parameter_description"`
// 	IsPrimitive     bool                 `json:"is_primitive" yaml:"parameter_is_primitive"`
// 	IsPointer       bool                 `json:"is_pointer" yaml:"parameter_is_pointer"`
// 	IsOptional      bool                 `json:"is_optional" yaml:"parameter_is_optional"`
// 	IsArray         bool                 `json:"is_array" yaml:"parameter_is_array"`
// 	IsPositional    bool                 `json:"is_positional" yaml:"parameter_is_positional"`
// 	PositionalIndex int                  `json:"positional_index" yaml:"parameter_positional_index"`
// 	Struct          map[string]Parameter `json:"struct,omitempty" yaml:"parameter_struct,omitempty"`
// 	AliasType       string               `json:"alias_type" yaml:"parameter_alias_type"`
// }

// // Method representa um método de um serviço
// type Method struct {
// 	Description string      `json:"description" yaml:"method_description"`
// 	Name        string      `json:"name" yaml:"method_name"`
// 	Parameters  []Parameter `json:"parameters" yaml:"method_parameters"` // nome -> tipo
// 	Returns     []Parameter `json:"returns" yaml:"method_returns"`       // nome -> tipo
// 	Comments    string      `json:"comments" yaml:"method_comments"`
// }
