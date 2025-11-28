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
	PositionalIndex     int  `json:"positional_index"`

	// Subcomandos que serão adicionados ao grupo
	SubCommands []SubCommandData `json:"sub_commands"`
	Commands    []CommandData    `json:"commands"`

	// Controle de geração
	GenerateGroup bool `json:"generate_group"`

	// Cobra flags definition
	CobraFlagsDefinition  []string       `json:"cobra_flags_definition"`
	CobraFlagsCreation    []string       `json:"cobra_flags_creation"`
	CobraFlagsAssign      []string       `json:"cobra_flags_assign"`
	PositionalArgs        map[int]string `json:"positional_args"`
	CobraFlagsRequired    []string       `json:"cobra_flags_required"`
	CobraStructInitialize []string       `json:"cobra_struct_initialize"`
	CobraArrayParse       []string       `json:"cobra_array_parse"`

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

	if info, err := os.Stat(pgd.CustomFile); err == nil && !info.IsDir() {
		pgd.HasCustomFile = true
		content, err := os.ReadFile(pgd.CustomFile)
		if err != nil {
			panic(fmt.Sprintf("erro ao ler arquivo customizado %s: %v", pgd.CustomFile, err))
		}
		pgd.CustomContent = string(content)
		return
	}

	if os.Getenv("GEN_CUSTOM_FILE") == "true" {
		dir := filepath.Dir(pgd.CustomFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			panic(fmt.Sprintf("erro ao criar diretório %s: %v", dir, err))
		}
		if _, err := os.Create(pgd.CustomFile + ".keep"); err != nil {
			panic(fmt.Sprintf("erro ao criar arquivo .keep: %v", err))
		}
	}
}

// AddImport adiciona um import à lista de imports (evita duplicatas)
func (pgd *PackageGroupData) AddImport(importPath string) {
	if pgd.hasImport(importPath) {
		return
	}
	pgd.Imports = append(pgd.Imports, importPath)
	slices.Sort(pgd.Imports)
}

// hasImport verifica se um import já existe na lista
func (pgd *PackageGroupData) hasImport(importPath string) bool {
	for _, imp := range pgd.Imports {
		if imp == importPath {
			return true
		}
	}
	return false
}

// AddCommand adiciona um comando ao grupo (evita duplicatas)
func (pgd *PackageGroupData) AddCommand(functionName, serviceCall string) {
	if pgd.hasCommand(functionName, serviceCall) {
		return
	}
	pgd.Commands = append(pgd.Commands, CommandData{
		FunctionName: functionName,
		ServiceCall:  serviceCall,
	})
	slices.SortFunc(pgd.Commands, func(a, b CommandData) int {
		return strings.Compare(a.FunctionName, b.FunctionName)
	})
	slices.SortFunc(pgd.Commands, func(a, b CommandData) int {
		return strings.Compare(a.ServiceCall, b.ServiceCall)
	})
}

// hasCommand verifica se um comando já existe no grupo
func (pgd *PackageGroupData) hasCommand(functionName, serviceCall string) bool {
	for _, cmd := range pgd.Commands {
		if cmd.FunctionName == functionName && cmd.ServiceCall == serviceCall {
			return true
		}
	}
	return false
}

// AddSubCommand adiciona um subcomando ao grupo (evita duplicatas)
func (pgd *PackageGroupData) AddSubCommand(packageName, functionName, serviceCall string) {
	lowerPkgName := strings.ToLower(packageName)
	if pgd.hasSubCommand(lowerPkgName, functionName, serviceCall) {
		return
	}
	pgd.SubCommands = append(pgd.SubCommands, SubCommandData{
		PackageName:  lowerPkgName,
		FunctionName: functionName,
		ServiceCall:  serviceCall,
	})

	slices.SortFunc(pgd.SubCommands, func(a, b SubCommandData) int {
		return strings.Compare(a.PackageName, b.PackageName)
	})
	slices.SortFunc(pgd.SubCommands, func(a, b SubCommandData) int {
		return strings.Compare(a.FunctionName, b.FunctionName)
	})
	slices.SortFunc(pgd.SubCommands, func(a, b SubCommandData) int {
		return strings.Compare(a.ServiceCall, b.ServiceCall)
	})
}

// hasSubCommand verifica se um subcomando já existe no grupo
func (pgd *PackageGroupData) hasSubCommand(packageName, functionName, serviceCall string) bool {
	for _, subCmd := range pgd.SubCommands {
		if subCmd.PackageName == packageName && subCmd.FunctionName == functionName && subCmd.ServiceCall == serviceCall {
			return true
		}
	}
	return false
}

// SetGroupID define o ID do grupo (usado para agrupamento na CLI)
func (pgd *PackageGroupData) SetGroupID(groupID string) {
	pgd.GroupID = groupID
}

// SetDescriptions define as descrições do comando
func (pgd *PackageGroupData) SetDescriptions(description, longDescription string) {
	pgd.ShortDescription = description
	pgd.LongDescription = longDescription
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

// SetUseName define o nome de uso do comando
func (pgd *PackageGroupData) SetUseName(useName string) {
	pgd.UseName = strings.ToLower(strutils.ToSnakeCase(useName, "-"))
}

var notAllowedPositionalArgs = []string{"create"}

func (pgd *PackageGroupData) CanAddPositionalArgs(positionalArgs string) bool {
	return !slices.Contains(notAllowedPositionalArgs, strings.ToLower(positionalArgs))
}

func (pgd *PackageGroupData) SetAllowedPositionalArgs() {
	pgd.AllowPositionalArgs = true
}

func (pgd *PackageGroupData) SetNotAllowedPositionalArgs() {
	pgd.AllowPositionalArgs = false
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
	slices.Sort(pgd.ServiceInit)

}

// WriteToFile escreve os dados no arquivo
func (pgd *PackageGroupData) WriteGroupToFile(filePath string) error {
	if pgd.GenerateGroup {
		return nil
	}

	pgd.GenerateGroup = true
	return pgd.writeTemplateToFile(packageGroupTmpl, filePath)
}

// WriteToFile escreve os dados no arquivo
func (pgd *PackageGroupData) WriteServiceToFile(filePath string) error {
	return pgd.writeTemplateToFile(serviceGroupTmpl, filePath)
}

func (pgd *PackageGroupData) WriteSubPackageToFile(filePath string) error {
	return pgd.writeTemplateToFile(subPackageGroupTmpl, filePath)
}

func (pgd *PackageGroupData) WriteProductCustomToFile(filePath string) error {
	return pgd.writeTemplateToFile(productCustomTmpl, filePath)
}

// WriteToFile escreve os dados no arquivo
func (pgd *PackageGroupData) WriteProductToFile(filePath string) error {
	return pgd.writeTemplateToFile(productTmpl, filePath)
}

// writeTemplateToFile é uma função auxiliar que escreve um template no arquivo
func (pgd *PackageGroupData) writeTemplateToFile(tmpl *template.Template, filePath string) error {
	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, pgd); err != nil {
		return fmt.Errorf("erro ao executar template: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório: %w", err)
	}

	if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("erro ao escrever arquivo: %w", err)
	}

	return nil
}

func (pgd *PackageGroupData) Copy() PackageGroupData {
	copied := *pgd

	// Copia os slices de forma eficiente usando copy
	copied.Imports = make([]string, len(pgd.Imports))
	copy(copied.Imports, pgd.Imports)

	copied.SubCommands = make([]SubCommandData, len(pgd.SubCommands))
	copy(copied.SubCommands, pgd.SubCommands)

	copied.Commands = make([]CommandData, len(pgd.Commands))
	copy(copied.Commands, pgd.Commands)

	copied.Params = make([]string, len(pgd.Params))
	copy(copied.Params, pgd.Params)

	// Copia outros slices que podem ter sido inicializados
	copied.ServiceInit = append([]string(nil), pgd.ServiceInit...)
	copied.ServiceSDKParamType = append([]string(nil), pgd.ServiceSDKParamType...)
	copied.ServiceSDKParamCreate = append([]string(nil), pgd.ServiceSDKParamCreate...)
	copied.CobraFlagsDefinition = append([]string(nil), pgd.CobraFlagsDefinition...)
	copied.CobraFlagsCreation = append([]string(nil), pgd.CobraFlagsCreation...)
	copied.CobraFlagsAssign = append([]string(nil), pgd.CobraFlagsAssign...)
	copied.CobraFlagsRequired = append([]string(nil), pgd.CobraFlagsRequired...)
	copied.CobraStructInitialize = append([]string(nil), pgd.CobraStructInitialize...)
	copied.CobraArrayParse = append([]string(nil), pgd.CobraArrayParse...)
	copied.UsedChars = append([]string(nil), pgd.UsedChars...)
	copied.Aliases = append([]string(nil), pgd.Aliases...)

	return copied
}

func (pdg *PackageGroupData) AddParam(param string) {
	pdg.Params = append(pdg.Params, param)
	slices.Sort(pdg.Params)
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
	if rgd.hasImport(importPath) {
		return
	}
	rgd.Imports = append(rgd.Imports, importPath)
	slices.Sort(rgd.Imports)
}

// hasImport verifica se um import já existe na lista
func (rgd *RootGenData) hasImport(importPath string) bool {
	for _, imp := range rgd.Imports {
		if imp == importPath {
			return true
		}
	}
	return false
}

func (pgd *PackageGroupData) SetServiceSDKParam(param string) {
	pgd.ServiceSDKParam = param
}

func (pgd *PackageGroupData) AddServiceSDKParamType(paramType string) {
	pgd.ServiceSDKParamType = append(pgd.ServiceSDKParamType, paramType)
	slices.Sort(pgd.ServiceSDKParamType)
}

func (pgd *PackageGroupData) AddServiceSDKParamCreate(paramCreate string) {
	pgd.ServiceSDKParamCreate = append(pgd.ServiceSDKParamCreate, paramCreate)
	slices.Sort(pgd.ServiceSDKParamCreate)
}

func (pgd *PackageGroupData) AddCobraFlagsDefinition(cobraFlagsDefinition string) {
	pgd.CobraFlagsDefinition = append(pgd.CobraFlagsDefinition, cobraFlagsDefinition)
	slices.Sort(pgd.CobraFlagsDefinition)
}

func (pgd *PackageGroupData) AddCobraFlagsCreation(cobraFlagsCreation string) {
	pgd.CobraFlagsCreation = append(pgd.CobraFlagsCreation, cobraFlagsCreation)
	slices.Sort(pgd.CobraFlagsCreation)
}

func (pgd *PackageGroupData) AddCobraFlagsAssign(cobraFlagsAssign string) {
	pgd.CobraFlagsAssign = append(pgd.CobraFlagsAssign, cobraFlagsAssign)
	slices.Sort(pgd.CobraFlagsAssign)

}

func (pgd *PackageGroupData) AddPositionalArgs(positionalArgs string) int {
	if pgd.PositionalArgs == nil {
		pgd.PositionalArgs = make(map[int]string)
	}
	pgd.PositionalArgs[pgd.PositionalIndex] = positionalArgs

	pgd.PositionalIndex = pgd.PositionalIndex + 1

	return pgd.PositionalIndex - 1

}

func (pgd *PackageGroupData) AddCobraFlagsRequired(cobraFlagsRequired string) {
	pgd.CobraFlagsRequired = append(pgd.CobraFlagsRequired, cobraFlagsRequired)
	slices.Sort(pgd.CobraFlagsRequired)
}

func (pgd *PackageGroupData) AddCobraArrayParse(cobraArrayParse string) {
	pgd.CobraArrayParse = append(pgd.CobraArrayParse, cobraArrayParse)
	slices.Sort(pgd.CobraArrayParse)
}

func (pgd *PackageGroupData) AddAssignResult(assignResult string) {
	pgd.AssignResult = assignResult
}

func (pgd *PackageGroupData) AddPrintResult(printResult string) {
	pgd.PrintResult = printResult
}

func (pgd *PackageGroupData) AddCobraStructInitialize(cobraStructInitialize string) {
	if pgd.hasCobraStructInitialize(cobraStructInitialize) {
		return
	}
	pgd.CobraStructInitialize = append(pgd.CobraStructInitialize, cobraStructInitialize)
	slices.Sort(pgd.CobraStructInitialize)
}

// hasCobraStructInitialize verifica se uma inicialização de struct já existe
func (pgd *PackageGroupData) hasCobraStructInitialize(value string) bool {
	for _, c := range pgd.CobraStructInitialize {
		if c == value {
			return true
		}
	}
	return false
}

// AddSubCommand adiciona um subcomando ao root
func (rgd *RootGenData) AddSubCommand(packageName, commandName string) {
	subCmd := RootSubCommandData{
		Package: strings.ToLower(packageName),
		Command: commandName,
	}
	rgd.SubCommands = append(rgd.SubCommands, subCmd)
	slices.SortFunc(rgd.SubCommands, func(a, b RootSubCommandData) int {
		return strings.Compare(a.Package, b.Package)
	})
	slices.SortFunc(rgd.SubCommands, func(a, b RootSubCommandData) int {
		return strings.Compare(a.Command, b.Command)
	})

}

// WriteRootGenToFile escreve os dados do root_gen.go no arquivo
func (rgd *RootGenData) WriteRootGenToFile(filePath string) error {
	buf := bytes.NewBuffer(nil)
	if err := rootGenTmpl.Execute(buf, rgd); err != nil {
		return fmt.Errorf("erro ao executar template: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório: %w", err)
	}

	if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("erro ao escrever arquivo: %w", err)
	}

	return nil
}
