package config

import (
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/packages"
)

// Config representa a estrutura principal do arquivo config.yaml
type Config struct {
	CLIVersion          string  `json:"cli_version"`
	SDKBranch           string  `json:"sdk_branch"`
	SDKTag              string  `json:"sdk_tag"`
	TagOrBranchOrLatest string  `json:"tag_or_branch_or_latest"`
	ShowGitError        bool    `json:"show_git_error"`
	Menus               []*Menu `json:"menus"`
}

// Menu representa um menu principal
type Menu struct {
	ID               string               `json:"id,omitempty"`
	Name             string               `json:"name,omitempty"`
	Enabled          bool                 `json:"enabled,omitempty"`
	Description      string               `json:"description"`
	LongDescription  string               `json:"long_description"`
	SDKPackage       string               `json:"sdk_package,omitempty"`
	CliGroup         string               `json:"cli_group,omitempty"`
	Alias            []string             `json:"alias,omitempty"`
	Menus            []*Menu              `json:"menus,omitempty"`
	ServiceInterface string               `json:"service_interface,omitempty"`
	Methods          []*Method            `json:"methods,omitempty"`
	SDKFile          string               `json:"sdk_file,omitempty"`
	CustomFile       string               `json:"custom_file,omitempty"`
	IsGroup          bool                 `json:"is_group,omitempty"`
	ParentMenu       *Menu                `json:"-"`
	Pkgs             *packages.Package    `json:"-"`
	Fset             *token.FileSet       `json:"-"`
	MapFile          map[string]*ast.File `json:"-"`
}

// Method representa um método de um serviço
type Method struct {
	Description     string        `json:"description,omitempty"`
	LongDescription string        `json:"long_description,omitempty"`
	Name            string        `json:"name,omitempty"`
	Parameters      []Parameter   `json:"parameters,omitempty"`
	Returns         []Parameter   `json:"returns,omitempty"`
	Comments        string        `json:"comments,omitempty"`
	Confirmation    *Confirmation `json:"confirmation,omitempty"`
	IsService       bool          `json:"is_service,omitempty"`
	ServiceImport   string        `json:"service_import,omitempty"`
	SDKFile         string        `json:"sdk_file,omitempty"`
	CustomFile      string        `json:"custom_file,omitempty"`
}

// Parameter representa um parâmetro de método
type Parameter struct {
	Name            string               `json:"name"`
	Type            string               `json:"type"`
	Description     string               `json:"description"`
	IsPrimitive     bool                 `json:"is_primitive"`
	IsPointer       bool                 `json:"is_pointer"`
	IsOptional      bool                 `json:"is_optional"`
	IsArray         bool                 `json:"is_array"`
	IsPositional    bool                 `json:"is_positional"`
	PositionalIndex int                  `json:"positional_index"`
	Struct          map[string]Parameter `json:"struct,omitempty"`
	AliasType       string               `json:"alias_type"`
}

type Confirmation struct {
	Enabled *bool   `json:"enabled,omitempty"`
	Value   *string `json:"value,omitempty"`
	Field   *string `json:"field,omitempty"`
	Type    *string `json:"type,omitempty"`
	Text    *string `json:"text,omitempty"`
}
