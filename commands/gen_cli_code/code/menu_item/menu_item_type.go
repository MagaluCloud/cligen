package menu_item

import (
	"fmt"
	"slices"
	"strings"

	"github.com/magaluCloud/cligen/config"
	"github.com/magaluCloud/cligen/file_utils"
)

type menuItemType struct {
	SaveToFile            string                  `json:"save_to_file_path"`
	PackageName           string                  `json:"package_name"`
	Imports               []string                `json:"imports"`
	FunctionName          string                  `json:"function_name"`
	ServiceParam          string                  `json:"service_param"`
	UseName               string                  `json:"use_name"`
	ShortDescription      string                  `json:"short_description"`
	LongDescription       string                  `json:"long_description"`
	CobraFlagsDefinition  []CobraFlagsDefinition  `json:"cobra_flags_definition"`
	CobraFlagsCreation    []string                `json:"cobra_flags_creation"`
	CobraFlagsAssign      []string                `json:"cobra_flags_assign"`
	CobraStructInitialize []string                `json:"cobra_struct_initialize"`
	Confirmation          string                  `json:"confirmation"`
	AssignResult          string                  `json:"assign_result"`
	PrintResult           string                  `json:"print_result"`
	ServiceCall           string                  `json:"service_call"`
	ErrorResult           string                  `json:"error_result"`
	PositionalArgs        []string                `json:"positional_args"`
	ServiceSDKParamCreate []ServiceSDKParamCreate `json:"service_sdk_param_create"`
	ServiceSDKParam       string                  `json:"service_sdk_param"`
	params                []string                `json:"-"`
}

type CobraFlagsDefinition struct {
	ParamName   string
	Name        string
	FlagType    string
	param       config.Parameter
	parents     []string
	cobraVar    string
	cobraAssign string
}

type ServiceSDKParamCreate struct {
	ParamName string
	Name      string
	ParamType string
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
	AddCobraFlagsDefinition(cobraFlagsDefinition CobraFlagsDefinition)
	GetCobraFlagsDefinition() []CobraFlagsDefinition
	AddServiceSDKParamCreate(serviceSDKParamCreate ServiceSDKParamCreate)
	AddCobraFlagsCreation(cobraFlagsCreation string)
	AddCobraFlagsAssign(cobraFlagsAssign string)
	AddCobraStructInitialize(cobraStructInitialize string)
	SetConfirmation(confirmation string)
	SetAssignResult(assignResult string)
	SetPrintResult(printResult string)
	SetServiceCall(serviceCall string)
	AddPositionalArgs(positionalArgs []string)
	SetErrorResult(errorResult string)
	SetServiceSDKParam(serviceSDKParam string)
	AddParam(param string)
	GetParams() []string
	Save() error
}

func NewMenuItem() MenuItem {
	return &menuItemType{}
}

func (m *menuItemType) SetErrorResult(errorResult string) {
	m.ErrorResult = errorResult
}

func (m *menuItemType) AddParam(param string) {
	m.params = append(m.params, param)
}

func (m *menuItemType) GetParams() []string {
	return m.params
}

func (m *menuItemType) AddImport(importt string) {
	if slices.Contains(m.Imports, importt) {
		return
	}
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

func (m *menuItemType) AddServiceSDKParamCreate(serviceSDKParamCreate ServiceSDKParamCreate) {
	m.ServiceSDKParamCreate = append(m.ServiceSDKParamCreate, serviceSDKParamCreate)
	slices.SortFunc(m.ServiceSDKParamCreate, func(a, b ServiceSDKParamCreate) int {
		return strings.Compare(a.Name, b.Name)
	})
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

func (m *menuItemType) AddCobraFlagsDefinition(cobraFlagsDefinition CobraFlagsDefinition) {
	m.CobraFlagsDefinition = append(m.CobraFlagsDefinition, cobraFlagsDefinition)

	slices.SortFunc(m.CobraFlagsDefinition, func(a, b CobraFlagsDefinition) int {
		return strings.Compare(a.Name, b.Name)
	})
}

func (m *menuItemType) GetCobraFlagsDefinition() []CobraFlagsDefinition {
	return m.CobraFlagsDefinition
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

func (m *menuItemType) AddCobraStructInitialize(cobraStructInitialize string) {
	if cobraStructInitialize == "" {
		return
	}
	m.CobraStructInitialize = append(m.CobraStructInitialize, cobraStructInitialize)
	slices.Sort(m.CobraStructInitialize)
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
