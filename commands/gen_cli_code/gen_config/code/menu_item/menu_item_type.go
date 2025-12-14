package menu_item

import (
	"fmt"
	"slices"
	"strings"

	"github.com/magaluCloud/cligen/file_utils"
)

type menuItemType struct {
	SaveToFile            string   `json:"save_to_file_path"`
	PackageName           string   `json:"package_name"`
	Imports               []string `json:"imports"`
	FunctionName          string   `json:"function_name"`
	ServiceParam          string   `json:"service_param"`
	UseName               string   `json:"use_name"`
	ShortDescription      string   `json:"short_description"`
	LongDescription       string   `json:"long_description"`
	CobraFlagsDefinition  []string `json:"cobra_flags_definition"`
	CobraFlagsCreation    []string `json:"cobra_flags_creation"`
	CobraFlagsAssign      []string `json:"cobra_flags_assign"`
	CobraFlagsRequired    []string `json:"cobra_flags_required"`
	CobraStructInitialize []string `json:"cobra_struct_initialize"`
	CobraArrayParse       []string `json:"cobra_array_parse"`
	Confirmation          string   `json:"confirmation"`
	AssignResult          string   `json:"assign_result"`
	PrintResult           string   `json:"print_result"`
	ServiceCall           string   `json:"service_call"`
	PositionalArgs        []string `json:"positional_args"`
	ServiceSDKParamCreate []string `json:"service_sdk_param_create"`
	ServiceSDKParam       string   `json:"service_sdk_param"`
}

type MenuItem interface {
	SetPathSaveToFile(pathToSave string)
	SetPackageName(packageName string)
	SetFunctionName(functionName string)
	SetServiceParam(serviceParam string)
	SetUseName(useName string)
	SetShortDescription(shortDescription string)
	SetLongDescription(longDescription string)
	AddImport(importt string)
	AddCobraFlagsDefinition(cobraFlagsDefinition string)
	AddCobraFlagsCreation(cobraFlagsCreation string)
	AddCobraFlagsAssign(cobraFlagsAssign string)
	AddCobraFlagsRequired(cobraFlagsRequired string)
	AddCobraStructInitialize(cobraStructInitialize string)
	AddCobraArrayParse(cobraArrayParse []string)
	SetConfirmation(confirmation string)
	SetAssignResult(assignResult string)
	SetPrintResult(printResult string)
	SetServiceCall(serviceCall string)
	AddPositionalArgs(positionalArgs []string)
	AddServiceSDKParamCreate(serviceSDKParamCreate string)
	SetServiceSDKParam(serviceSDKParam string)
	Save() error
}

func NewMenuItem() MenuItem {
	return &menuItemType{}
}

func (m *menuItemType) AddImport(importt string) {
	m.Imports = append(m.Imports, importt)
	slices.Sort(m.Imports)
}

func (m *menuItemType) SetServiceSDKParam(serviceSDKParam string) {
	m.ServiceSDKParam = serviceSDKParam
}

func (m *menuItemType) SetPathSaveToFile(saveToFilePath string) {
	m.SaveToFile = saveToFilePath
}

func (m *menuItemType) AddPositionalArgs(positionalArgs []string) {
	m.PositionalArgs = append(m.PositionalArgs, positionalArgs...)
	slices.Sort(m.PositionalArgs)
}

func (m *menuItemType) AddServiceSDKParamCreate(serviceSDKParamCreate string) {
	if serviceSDKParamCreate == "" {
		return
	}
	m.ServiceSDKParamCreate = append(m.ServiceSDKParamCreate, serviceSDKParamCreate)
	slices.Sort(m.ServiceSDKParamCreate)
}

func (m *menuItemType) Save() error {
	if m.SaveToFile == "" {
		return fmt.Errorf("save to file path is not set")
	}
	return file_utils.WriteTemplateToFile(menuItemTmpl, m, m.SaveToFile)
}

func (m *menuItemType) SetPackageName(packageName string) {
	m.PackageName = strings.ToLower(packageName)
}

func (m *menuItemType) SetFunctionName(functionName string) {
	m.FunctionName = functionName
}

func (m *menuItemType) SetServiceParam(serviceParam string) {
	m.ServiceParam = serviceParam
}

func (m *menuItemType) SetUseName(useName string) {
	m.UseName = strings.ToLower(useName)
}

func (m *menuItemType) SetShortDescription(shortDescription string) {
	m.ShortDescription = shortDescription
}

func (m *menuItemType) SetLongDescription(longDescription string) {
	m.LongDescription = longDescription
}

func (m *menuItemType) AddCobraFlagsDefinition(cobraFlagsDefinition string) {
	if cobraFlagsDefinition == "" {
		return
	}
	m.CobraFlagsDefinition = append(m.CobraFlagsDefinition, cobraFlagsDefinition)
	slices.Sort(m.CobraFlagsDefinition)
}

func (m *menuItemType) AddCobraFlagsCreation(cobraFlagsCreation string) {
	if cobraFlagsCreation == "" {
		return
	}
	m.CobraFlagsCreation = append(m.CobraFlagsCreation, cobraFlagsCreation)
	slices.Sort(m.CobraFlagsCreation)
}

func (m *menuItemType) AddCobraFlagsAssign(cobraFlagsAssign string) {
	if cobraFlagsAssign == "" {
		return
	}
	m.CobraFlagsAssign = append(m.CobraFlagsAssign, cobraFlagsAssign)
	slices.Sort(m.CobraFlagsAssign)
}

func (m *menuItemType) AddCobraFlagsRequired(cobraFlagsRequired string) {
	if cobraFlagsRequired == "" {
		return
	}
	m.CobraFlagsRequired = append(m.CobraFlagsRequired, cobraFlagsRequired)
	slices.Sort(m.CobraFlagsRequired)
}

func (m *menuItemType) AddCobraStructInitialize(cobraStructInitialize string) {
	if cobraStructInitialize == "" {
		return
	}
	m.CobraStructInitialize = append(m.CobraStructInitialize, cobraStructInitialize)
	slices.Sort(m.CobraStructInitialize)
}

func (m *menuItemType) AddCobraArrayParse(cobraArrayParse []string) {
	m.CobraArrayParse = append(m.CobraArrayParse, cobraArrayParse...)
	slices.Sort(m.CobraArrayParse)
}

func (m *menuItemType) SetConfirmation(confirmation string) {
	m.Confirmation = confirmation
}

func (m *menuItemType) SetAssignResult(assignResult string) {
	m.AssignResult = assignResult
}

func (m *menuItemType) SetPrintResult(printResult string) {
	m.PrintResult = printResult
}

func (m *menuItemType) SetServiceCall(serviceCall string) {
	m.ServiceCall = serviceCall
}
