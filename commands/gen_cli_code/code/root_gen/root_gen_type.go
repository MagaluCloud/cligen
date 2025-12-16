package root_gen

import (
	"fmt"
	"slices"
	"strings"

	"github.com/magaluCloud/cligen/file_utils"
)

type rootGenType struct {
	SubCommands []SubCommandType
	Imports     []string
	SaveToFile  string `json:"save_to_file_path"`
}

type SubCommandType struct {
	PackageName  string `json:"package_name"`
	FunctionName string `json:"function_name"`
}

type RootGen interface {
	SetPathSaveToFile(saveToFilePath string)
	AddSubCommand(subCommand SubCommandType)
	AddImport(importPath string)
	Save() error
}

func NewRootGen() RootGen {
	return &rootGenType{
		SubCommands: make([]SubCommandType, 0),
		Imports:     make([]string, 0),
		SaveToFile:  "",
	}
}

func (rg *rootGenType) Save() error {
	if rg.SaveToFile == "" {
		return fmt.Errorf("save to file path is not set")
	}
	return file_utils.WriteTemplateToFile(rootGenTmpl, rg, rg.SaveToFile)
}

func (rg *rootGenType) SetPathSaveToFile(saveToFilePath string) {
	rg.SaveToFile = saveToFilePath
}

func (rg *rootGenType) AddSubCommand(subCommand SubCommandType) {
	rg.SubCommands = append(rg.SubCommands, subCommand)

	slices.SortFunc(rg.SubCommands, func(a, b SubCommandType) int {
		return strings.Compare(a.PackageName, b.PackageName)
	})
	slices.SortFunc(rg.SubCommands, func(a, b SubCommandType) int {
		return strings.Compare(a.FunctionName, b.FunctionName)
	})
}

func (rg *rootGenType) AddImport(importPath string) {
	if rg.hasImport(importPath) {
		return
	}
	rg.Imports = append(rg.Imports, importPath)
	slices.Sort(rg.Imports)
}

func (rg *rootGenType) hasImport(importPath string) bool {
	for _, imp := range rg.Imports {
		if imp == importPath {
			return true
		}
	}
	return false
}
