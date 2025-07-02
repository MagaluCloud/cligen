package gen_cli_code

import (
	"bytes"
	_ "embed"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	strutils "cligen/str_utils"
)

//go:embed package_group.template
var packageGroupTemplate string

//go:embed service_group.template
var serviceGroupTemplate string

//go:embed product.template
var productTemplate string

//go:embed rootgen.template
var rootGenTemplate string

//go:embed sub_package_group.template
var subPackageGroupTemplate string

// Templates pré-compilados para melhor performance
var (
	packageGroupTmpl    *template.Template
	serviceGroupTmpl    *template.Template
	productTmpl         *template.Template
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

	// Informações do comando
	UseName          string `json:"use_name"`
	ShortDescription string `json:"short_description"`
	LongDescription  string `json:"long_description"`
	GroupID          string `json:"group_id,omitempty"`

	// Subcomandos que serão adicionados ao grupo
	SubCommands []SubCommandData `json:"sub_commands"`
	Commands    []CommandData    `json:"commands"`

	// Controle de geração
	GenerateGroup bool `json:"generate_group"`
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
		Imports:     make([]string, 0, 10),
		SubCommands: make([]SubCommandData, 0, 5),
		Commands:    make([]CommandData, 0, 5),
		Params:      make([]string, 0, 5),
		GroupID:     "",
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

// SetUseName define o nome de uso do comando
func (pgd *PackageGroupData) SetUseName(useName string) {
	pgd.UseName = strings.ToLower(strutils.ToSnakeCase(useName, "-"))
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
