package config

// Config representa a estrutura principal do arquivo config.yaml
type Config struct {
	CLIVersion          string  `yaml:"cli_version"`
	SDKBranch           string  `yaml:"sdk_branch"`
	SDKTag              string  `yaml:"sdk_tag"`
	TagOrBranchOrLatest string  `yaml:"tag_or_branch_or_latest"`
	ShowGitError        bool    `yaml:"show_git_error"`
	Menus               []*Menu `yaml:"menus"`
}

// Menu representa um menu principal
type Menu struct {
	Name            string        `yaml:"name,omitempty"`
	Enabled         bool          `yaml:"enabled,omitempty"`
	Description     string        `yaml:"description"`
	LongDescription string        `yaml:"long_description"`
	Level           int           `yaml:"level,omitempty"`
	SDKPackage      string        `yaml:"sdk_package,omitempty"`
	CliGroup        string        `yaml:"cli_group,omitempty"`
	Alias           []string      `yaml:"alias,omitempty"`
	Menus           []*Menu       `yaml:"menus,omitempty"`
	Confirmation    *Confirmation `yaml:"confirmation,omitempty"`
}

type Confirmation struct {
	Enabled *bool   `yaml:"enabled,omitempty"`
	Value   *string `yaml:"value,omitempty"`
	Field   *string `yaml:"field,omitempty"`
	Type    *string `yaml:"type,omitempty"`
	Text    *string `yaml:"text,omitempty"`
}

//type:
/*
simple-ask
type-value
just-type
*/
