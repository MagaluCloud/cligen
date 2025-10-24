package gen_cli_code

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"

	strutils "github.com/magaluCloud/cligen/str_utils"
)

//go:embed package_group.template
var packageGroupTemplate string

//go:embed service_group.template
var serviceGroupTemplate string

//go:embed product.template
var productTemplate string

//go:embed product_custom.template
var productCustomTemplate string

//go:embed rootgen.template
var rootGenTemplate string

//go:embed sub_package_group.template
var subPackageGroupTemplate string

// Templates pré-compilados para melhor performance
var (
	packageGroupTmpl    *template.Template
	serviceGroupTmpl    *template.Template
	productTmpl         *template.Template
	productCustomTmpl   *template.Template
	rootGenTmpl         *template.Template
	subPackageGroupTmpl *template.Template
)

func init() {
	var err error
	packageGroupTmpl, err = template.New("package_group").Parse(packageGroupTemplate)
	if err != nil {
		panic(err)
	}
	serviceGroupTmpl, err = template.New("service_group").Parse(serviceGroupTemplate)
	if err != nil {
		panic(err)
	}
	productTmpl, err = template.New("product").Parse(productTemplate)
	if err != nil {
		panic(err)
	}
	productCustomTmpl, err = template.New("product_custom").Parse(productCustomTemplate)
	if err != nil {
		panic(err)
	}
	rootGenTmpl, err = template.New("rootgen").Parse(rootGenTemplate)
	if err != nil {
		panic(err)
	}
	subPackageGroupTmpl, err = template.New("sub_package_group").Parse(subPackageGroupTemplate)
	if err != nil {
		panic(err)
	}
}

// PackageGroupData representa os dados necessários para gerar um arquivo de grupo de comandos
type PackageGroupData struct {
	FileID        string `json:"file_id"`
	CustomFile    string `json:"custom_file"`
	HasCustomFile bool   `json:"has_custom_file"`
	CustomContent string `json:"custom_content"`

	// Informações básicas do pacote
	PackageName string `json:"package_name"`

	// Imports necessários para o arquivo
	Imports []string `json:"imports"`

	// Informações da função principal
	FunctionName string   `json:"function_name"`
	ServiceParam string   `json:"service_param"`
	ServiceInit  []string `json:"service_init"`

	// Parâmetros da função
	Params []string `json:"params"`

	// Código da função
	ServiceCall string `json:"service_call"`

	// Parâmetros do serviço
	ServiceSDKParam       string   `json:"service_sdk_param"`
	ServiceSDKParamType   []string `json:"service_sdk_param_type"`
	ServiceSDKParamCreate []string `json:"service_sdk_param_create"`

	// Informações do comando
	UseName          string   `json:"use_name"`
	Aliases          []string `json:"aliases"`
	ShortDescription string   `json:"short_description"`
	LongDescription  string   `json:"long_description"`
	GroupID          string   `json:"group_id,omitempty"`

	AllowPositionalArgs bool `json:"allow_positional_args"`

	// Subcomandos que serão adicionados ao grupo
	SubCommands []SubCommandData `json:"sub_commands"`
	Commands    []CommandData    `json:"commands"`

	// Controle de geração
	GenerateGroup bool `json:"generate_group"`

	// Cobra flags definition
	CobraFlagsDefinition  []string `json:"cobra_flags_definition"`
	CobraFlagsCreation    []string `json:"cobra_flags_creation"`
	CobraFlagsAssign      []string `json:"cobra_flags_assign"`
	PositionalArgs        string   `json:"positional_args"`
	CobraFlagsRequired    []string `json:"cobra_flags_required"`
	CobraStructInitialize []string `json:"cobra_struct_initialize"`
	CobraArrayParse       []string `json:"cobra_array_parse"`

	// Used chars
	UsedChars []string `json:"used_chars"`

	// PrintResult
	AssignResult string `json:"assign_result"`
	PrintResult  string `json:"print_result"`
}

// SubCommandData representa um subcomando que será adicionado ao grupo
type SubCommandData struct {
	PackageName  string `json:"package_name"`
	FunctionName string `json:"function_name"`
	ServiceCall  string `json:"service_call"`
}

type CommandData struct {
	FunctionName string `json:"function_name"`
	ServiceCall  string `json:"service_call"`
}

// TemplateData representa os dados completos para renderização do template
type TemplateData struct {
	PackageGroup PackageGroupData `json:"package_group"`
}

// NewPackageGroupData cria uma nova instância de PackageGroupData com valores padrão
func NewPackageGroupData() *PackageGroupData {
	return &PackageGroupData{
		Imports:     []string{},
		SubCommands: []SubCommandData{},
		Commands:    []CommandData{},
		Params:      []string{},
		GroupID:     "",
		UsedChars:   []string{},
	}
}

func (pgd *PackageGroupData) SetFileID(fileID string) {
	pgd.FileID = fileID
	pgd.CustomFile = strings.Replace(fileID, "base-cli-gen", "base-cli-custom", 1)

	if _, err := os.Stat(pgd.CustomFile); err == nil {
		pgd.HasCustomFile = true
		content, err := os.ReadFile(pgd.CustomFile)
		if err != nil {
			panic(err)
		}
		pgd.CustomContent = string(content)
		return
	}

	if os.Getenv("GEN_CUSTOM_FILE") == "true" {
		os.MkdirAll(filepath.Dir(pgd.CustomFile), 0755)
		_, err := os.Create(pgd.CustomFile + ".keep")
		if err != nil {
			panic(err)
		}
	}
}

// AddImport adiciona um import à lista de imports (evita duplicatas)
func (pgd *PackageGroupData) AddImport(importPath string) {
	for _, imp := range pgd.Imports {
		if imp == importPath {
			return
		}
	}
	pgd.Imports = append(pgd.Imports, importPath)
}

// AddCommand adiciona um comando ao grupo (evita duplicatas)
func (pgd *PackageGroupData) AddCommand(functionName, serviceCall string) {
	for _, cmd := range pgd.Commands {
		if cmd.FunctionName == functionName && cmd.ServiceCall == serviceCall {
			return
		}
	}
	pgd.Commands = append(pgd.Commands, CommandData{
		FunctionName: functionName,
		ServiceCall:  serviceCall,
	})
}

// AddSubCommand adiciona um subcomando ao grupo (evita duplicatas)
func (pgd *PackageGroupData) AddSubCommand(packageName, functionName, serviceCall string) {
	for _, subCmd := range pgd.SubCommands {
		if subCmd.PackageName == strings.ToLower(packageName) && subCmd.FunctionName == functionName && subCmd.ServiceCall == serviceCall {
			return
		}
	}
	pgd.SubCommands = append(pgd.SubCommands, SubCommandData{
		PackageName:  strings.ToLower(packageName),
		FunctionName: functionName,
		ServiceCall:  serviceCall,
	})
}

// SetGroupID define o ID do grupo (usado para agrupamento na CLI)
func (pgd *PackageGroupData) SetGroupID(groupID string) {
	pgd.GroupID = groupID
}

// SetDescriptions define as descrições do comando
func (pgd *PackageGroupData) SetDescriptions(short, long string) {
	pgd.ShortDescription = short
	pgd.LongDescription = long
}

// SetServiceParam define o parâmetro do serviço
func (pgd *PackageGroupData) SetServiceParam(serviceParam string) {
	pgd.ServiceParam = serviceParam
}

// SetFunctionName define o nome da função
func (pgd *PackageGroupData) SetFunctionName(functionName string) {
	pgd.FunctionName = functionName
}

// SetPackageName define o nome do pacote
func (pgd *PackageGroupData) SetPackageName(packageName string) {
	pgd.PackageName = strings.ToLower(packageName)
}

var notAllowedPositionalArgs = []string{"create"}

// SetUseName define o nome de uso do comando
func (pgd *PackageGroupData) SetUseName(useName string) {
	pgd.UseName = strings.ToLower(strutils.ToSnakeCase(useName, "-"))
	pgd.AllowPositionalArgs = true
	if slices.Contains(notAllowedPositionalArgs, pgd.UseName) {
		pgd.AllowPositionalArgs = false
	}
}

func (pgd *PackageGroupData) AppendPositionalArgs(positionalArgs string) bool {
	if pgd.AllowPositionalArgs {
		pgd.UseName = fmt.Sprintf("%s [%s]", pgd.UseName, positionalArgs)
		return true
	}
	pgd.AllowPositionalArgs = false
	return pgd.AllowPositionalArgs
}

func (pgd *PackageGroupData) LoadCustomUse() {

}

// SetAliases define os aliases do comando
func (pgd *PackageGroupData) SetAliases(aliases []string) {
	pgd.Aliases = aliases
}

// SetServiceInit define o código para inicializar o serviço
func (pgd *PackageGroupData) AddServiceInit(serviceInit string) {
	pgd.ServiceInit = append(pgd.ServiceInit, serviceInit)
}

// WriteToFile escreve os dados no arquivo
func (pgd *PackageGroupData) WriteGroupToFile(filePath string) error {
	if pgd.GenerateGroup {
		return nil
	}

	buf := bytes.NewBuffer(nil)
	err := packageGroupTmpl.Execute(buf, pgd)
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(filePath), 0755)
	pgd.GenerateGroup = true
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

// WriteToFile escreve os dados no arquivo
func (pgd *PackageGroupData) WriteServiceToFile(filePath string) error {
	buf := bytes.NewBuffer(nil)
	err := serviceGroupTmpl.Execute(buf, pgd)
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(filePath), 0755)
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

func (pgd *PackageGroupData) WriteSubPackageToFile(filePath string) error {
	buf := bytes.NewBuffer(nil)
	err := subPackageGroupTmpl.Execute(buf, pgd)
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(filePath), 0755)
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

func (pgd *PackageGroupData) WriteProductCustomToFile(filePath string) error {
	buf := bytes.NewBuffer(nil)
	err := productCustomTmpl.Execute(buf, pgd)
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(filePath), 0755)
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

// WriteToFile escreve os dados no arquivo
func (pgd *PackageGroupData) WriteProductToFile(filePath string) error {
	buf := bytes.NewBuffer(nil)
	err := productTmpl.Execute(buf, pgd)
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(filePath), 0755)
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

func (pgd *PackageGroupData) Copy() PackageGroupData {
	// Cria uma cópia profunda para evitar problemas de referência
	copied := *pgd

	// Copia os slices de forma eficiente
	if len(pgd.Imports) > 0 {
		copied.Imports = make([]string, len(pgd.Imports))
		copy(copied.Imports, pgd.Imports)
	}

	if len(pgd.SubCommands) > 0 {
		copied.SubCommands = make([]SubCommandData, len(pgd.SubCommands))
		copy(copied.SubCommands, pgd.SubCommands)
	}

	if len(pgd.Commands) > 0 {
		copied.Commands = make([]CommandData, len(pgd.Commands))
		copy(copied.Commands, pgd.Commands)
	}

	if len(pgd.Params) > 0 {
		copied.Params = make([]string, len(pgd.Params))
		copy(copied.Params, pgd.Params)
	}

	return copied
}

func (pdg *PackageGroupData) AddParam(param string) {
	pdg.Params = append(pdg.Params, param)
}

func (pdg *PackageGroupData) SetServiceCall(serviceCall string) {
	pdg.ServiceCall = serviceCall
}

// RootGenData representa os dados necessários para gerar o arquivo root_gen.go
type RootGenData struct {
	// Imports necessários para o arquivo
	Imports []string `json:"imports"`

	// Subcomandos que serão adicionados ao root
	SubCommands []RootSubCommandData `json:"sub_commands"`
}

// RootSubCommandData representa um subcomando que será adicionado ao root
type RootSubCommandData struct {
	Package string `json:"package"`
	Command string `json:"command"`
}

// NewRootGenData cria uma nova instância de RootGenData com valores padrão
func NewRootGenData() *RootGenData {
	return &RootGenData{
		Imports:     make([]string, 0, 10),
		SubCommands: make([]RootSubCommandData, 0, 10),
	}
}

// AddImport adiciona um import à lista de imports (evita duplicatas)
func (rgd *RootGenData) AddImport(importPath string) {
	for _, imp := range rgd.Imports {
		if imp == importPath {
			return
		}
	}
	rgd.Imports = append(rgd.Imports, importPath)
}

func (pgd *PackageGroupData) SetServiceSDKParam(param string) {
	pgd.ServiceSDKParam = param
}

func (pgd *PackageGroupData) AddServiceSDKParamType(paramType string) {
	pgd.ServiceSDKParamType = append(pgd.ServiceSDKParamType, paramType)
}

func (pgd *PackageGroupData) AddServiceSDKParamCreate(paramCreate string) {
	pgd.ServiceSDKParamCreate = append(pgd.ServiceSDKParamCreate, paramCreate)
}

func (pgd *PackageGroupData) AddCobraFlagsDefinition(cobraFlagsDefinition string) {
	pgd.CobraFlagsDefinition = append(pgd.CobraFlagsDefinition, cobraFlagsDefinition)
}

func (pgd *PackageGroupData) AddCobraFlagsCreation(cobraFlagsCreation string) {
	pgd.CobraFlagsCreation = append(pgd.CobraFlagsCreation, cobraFlagsCreation)
}

func (pgd *PackageGroupData) AddCobraFlagsAssign(cobraFlagsAssign string) {
	pgd.CobraFlagsAssign = append(pgd.CobraFlagsAssign, cobraFlagsAssign)
}

func (pgd *PackageGroupData) AddPositionalArgs(positionalArgs string) {
	pgd.PositionalArgs = positionalArgs
}

func (pgd *PackageGroupData) AddCobraFlagsRequired(cobraFlagsRequired string) {
	pgd.CobraFlagsRequired = append(pgd.CobraFlagsRequired, cobraFlagsRequired)
}

func (pgd *PackageGroupData) AddCobraArrayParse(cobraArrayParse string) {
	pgd.CobraArrayParse = append(pgd.CobraArrayParse, cobraArrayParse)
}

func (pgd *PackageGroupData) AddAssignResult(assignResult string) {
	pgd.AssignResult = assignResult
}

func (pgd *PackageGroupData) AddPrintResult(printResult string) {
	pgd.PrintResult = printResult
}

func (pgd *PackageGroupData) AddCobraStructInitialize(cobraStructInitialize string) {
	for _, c := range pgd.CobraStructInitialize {
		if c == cobraStructInitialize {
			return
		}
	}
	pgd.CobraStructInitialize = append(pgd.CobraStructInitialize, cobraStructInitialize)
}

// AddSubCommand adiciona um subcomando ao root
func (rgd *RootGenData) AddSubCommand(packageName, commandName string) {
	subCmd := RootSubCommandData{
		Package: strings.ToLower(packageName),
		Command: commandName,
	}
	rgd.SubCommands = append(rgd.SubCommands, subCmd)
}

// WriteRootGenToFile escreve os dados do root_gen.go no arquivo
func (rgd *RootGenData) WriteRootGenToFile(filePath string) error {
	buf := bytes.NewBuffer(nil)
	err := rootGenTmpl.Execute(buf, rgd)
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(filePath), 0755)
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}
