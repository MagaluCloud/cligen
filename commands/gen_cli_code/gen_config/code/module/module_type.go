package module

import (
	"fmt"
	"slices"
	"strings"

	"github.com/magaluCloud/cligen/file_utils"
)

type moduleType struct {
	PackageName      string           `json:"package_name"`
	Imports          []string         `json:"imports"`
	FunctionName     string           `json:"function_name"`
	ServiceParam     string           `json:"service_param"`
	UseName          string           `json:"use_name"`
	ShortDescription string           `json:"short_description"`
	LongDescription  string           `json:"long_description"`
	Aliases          []string         `json:"aliases"`
	GroupID          string           `json:"group_id"`
	ServiceInit      []string         `json:"service_init"`
	SubCommands      []subCommandType `json:"sub_commands"`
	ServiceCall      string           `json:"service_call"`
	SaveToFile       string           `json:"save_to_file_path"`
}

type subCommandType struct {
	PackageName  string `json:"package_name"`
	FunctionName string `json:"function_name"`
	ServiceCall  string `json:"service_call"`
}

type Module interface {
	SetPathSaveToFile(saveToFilePath string)
	AddSubCommand(subCommand subCommandType)
	AddImport(importPath string)
	SetPackageName(packageName string)
	SetFunctionName(functionName string)
	SetUseName(useName string)
	SetShortDescription(shortDescription string)
	SetLongDescription(longDescription string)
	AddAliases(alias ...string)
	AddAlias(alias string)
	SetGroupID(groupID string)
	AddServiceInit(serviceInit string)
	SetServiceParam(serviceParam string)

	Save() error
}

func NewModule() Module {
	return &moduleType{
		SubCommands: make([]subCommandType, 0),
		Imports:     make([]string, 0),
		SaveToFile:  "",
	}
}

func (m *moduleType) Save() error {
	if m.SaveToFile == "" {
		return fmt.Errorf("save to file path is not set")
	}
	return file_utils.WriteTemplateToFile(moduleTmpl, m, m.SaveToFile)
}

func (m *moduleType) SetPathSaveToFile(saveToFilePath string) {
	m.SaveToFile = saveToFilePath
}

func (m *moduleType) AddSubCommand(subCommand subCommandType) {
	m.SubCommands = append(m.SubCommands, subCommand)
	slices.SortFunc(m.SubCommands, func(a, b subCommandType) int {
		return strings.Compare(a.PackageName, b.PackageName)
	})
	slices.SortFunc(m.SubCommands, func(a, b subCommandType) int {
		return strings.Compare(a.FunctionName, b.FunctionName)
	})
	slices.SortFunc(m.SubCommands, func(a, b subCommandType) int {
		return strings.Compare(a.ServiceCall, b.ServiceCall)
	})
}

func (m *moduleType) AddImport(importPath string) {
	m.Imports = append(m.Imports, importPath)
}

func (m *moduleType) SetPackageName(packageName string) {
	m.PackageName = strings.ToLower(packageName)
}

func (m *moduleType) SetFunctionName(functionName string) {
	m.FunctionName = functionName
}

func (m *moduleType) SetServiceParam(serviceParam string) {
	m.ServiceParam = serviceParam
}

func (m *moduleType) SetUseName(useName string) {
	m.UseName = strings.ToLower(useName)
}

func (m *moduleType) SetShortDescription(shortDescription string) {
	m.ShortDescription = shortDescription
}

func (m *moduleType) SetLongDescription(longDescription string) {
	m.LongDescription = longDescription
}

func (m *moduleType) AddAliases(alias ...string) {
	m.Aliases = append(m.Aliases, alias...)
	slices.Sort(m.Aliases)
}

func (m *moduleType) AddAlias(alias string) {
	m.Aliases = append(m.Aliases, alias)
	slices.Sort(m.Aliases)
}

func (m *moduleType) SetGroupID(groupID string) {
	m.GroupID = groupID
}

func (m *moduleType) AddServiceInit(serviceInit string) {
	m.ServiceInit = append(m.ServiceInit, serviceInit)
	slices.Sort(m.ServiceInit)
}
