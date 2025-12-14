package menu

import (
	"fmt"
	"slices"
	"strings"

	"github.com/magaluCloud/cligen/file_utils"
)

type menuType struct {
	SaveToFile       string        `json:"save_to_file_path"`
	PackageName      string        `json:"package_name"`
	Imports          []string      `json:"imports"`
	FunctionName     string        `json:"function_name"`
	ServiceParam     string        `json:"service_param"`
	UseName          string        `json:"use_name"`
	ShortDescription string        `json:"short_description"`
	LongDescription  string        `json:"long_description"`
	Aliases          []string      `json:"aliases"`
	GroupID          string        `json:"group_id"`
	ServiceInit      []string      `json:"service_init"`
	Commands         []CommandType `json:"commands"`
}

type CommandType struct {
	FunctionName string `json:"function_name"`
	ServiceCall  string `json:"service_call"`
}

type Menu interface {
	SetPathSaveToFile(saveToFilePath string)
	AddCommand(command CommandType)
	AddImport(importPath string)
	SetPackageName(packageName string)
	SetFunctionName(functionName string)
	SetUseName(useName string)
	SetShortDescription(shortDescription string)
	SetLongDescription(longDescription string)
	SetAliases(alias ...string)
	AddAlias(alias string)
	SetServiceParam(serviceParam string)
	SetGroupID(groupID string)
	AddServiceInit(serviceInit string)
	Save() error
}

func NewMenu() Menu {
	return &menuType{
		Commands: make([]CommandType, 0),
		Imports:  make([]string, 0),
	}
}

func (m *menuType) SetPathSaveToFile(saveToFilePath string) {
	m.SaveToFile = saveToFilePath
}

func (m *menuType) AddCommand(command CommandType) {
	m.Commands = append(m.Commands, command)
	slices.SortFunc(m.Commands, func(a, b CommandType) int {
		return strings.Compare(a.FunctionName, b.FunctionName)
	})
}

func (m *menuType) AddImport(importPath string) {
	m.Imports = append(m.Imports, importPath)
}

func (m *menuType) SetPackageName(packageName string) {
	m.PackageName = strings.ToLower(packageName)
}

func (m *menuType) SetFunctionName(functionName string) {
	m.FunctionName = functionName
}

func (m *menuType) SetUseName(useName string) {
	m.UseName = strings.ToLower(useName)
}

func (m *menuType) SetShortDescription(shortDescription string) {
	m.ShortDescription = shortDescription
}

func (m *menuType) SetLongDescription(longDescription string) {
	m.LongDescription = longDescription
}

func (m *menuType) AddAlias(alias string) {
	m.Aliases = append(m.Aliases, alias)
	slices.Sort(m.Aliases)
}

func (m *menuType) SetAliases(alias ...string) {
	m.Aliases = append(m.Aliases, alias...)
	slices.Sort(m.Aliases)
}

func (m *menuType) SetServiceParam(serviceParam string) {
	m.ServiceParam = serviceParam
}

func (m *menuType) SetGroupID(groupID string) {
	m.GroupID = groupID
}

func (m *menuType) AddServiceInit(serviceInit string) {
	m.ServiceInit = append(m.ServiceInit, serviceInit)
	slices.Sort(m.ServiceInit)
}

func (m *menuType) Save() error {
	if m.SaveToFile == "" {
		return fmt.Errorf("save to file path is not set")
	}
	return file_utils.WriteTemplateToFile(menuTmpl, m, m.SaveToFile)
}
