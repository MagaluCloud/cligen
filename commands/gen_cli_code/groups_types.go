package gen_cli_code

import (
	"bytes"
	_ "embed"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"text/template"
)

//go:embed package_group.template
var packageGroupTemplate string

//go:embed service_group.template
var serviceGroupTemplate string

//go:embed product.template
var productTemplate string

//go:embed rootgen.template
var rootGenTemplate string

// PackageGroupData representa os dados necessários para gerar um arquivo de grupo de comandos
type PackageGroupData struct {
	// Informações básicas do pacote
	PackageName string `json:"package_name"`

	// Imports necessários para o arquivo
	Imports []string `json:"imports"`

	// Informações da função principal
	FunctionName string `json:"function_name"`
	ServiceParam string `json:"service_param"`
	ServiceInit  string `json:"service_init"`

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
		Imports:     []string{},
		SubCommands: []SubCommandData{},
		GroupID:     "", // Opcional
	}
}

// AddImport adiciona um import à lista de imports
func (pgd *PackageGroupData) AddImport(importPath string) {
	if slices.Contains(pgd.Imports, importPath) {
		return
	}
	pgd.Imports = append(pgd.Imports, importPath)
}

// AddCommand adiciona um comando ao grupo
func (pgd *PackageGroupData) AddCommand(functionName, serviceCall string) {
	cmd := CommandData{
		FunctionName: functionName,
		ServiceCall:  serviceCall,
	}
	pgd.Commands = append(pgd.Commands, cmd)
}

// AddSubCommand adiciona um subcomando ao grupo
func (pgd *PackageGroupData) AddSubCommand(packageName, functionName, serviceCall string) {
	subCmd := SubCommandData{
		PackageName:  strings.ToLower(packageName),
		FunctionName: functionName,
		ServiceCall:  serviceCall,
	}
	pgd.SubCommands = append(pgd.SubCommands, subCmd)
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
	pgd.UseName = strings.ToLower(useName)
}

// SetServiceInit define o código para inicializar o serviço
func (pgd *PackageGroupData) SetServiceInit(serviceInit string) {
	pgd.ServiceInit = serviceInit
}

// WriteToFile escreve os dados no arquivo
func (pgd *PackageGroupData) WriteGroupToFile(filePath string) error {
	if pgd.GenerateGroup {
		return nil
	}

	tmpl, err := template.New("package_group").Parse(packageGroupTemplate)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, pgd)
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(filePath), 0755)
	pgd.GenerateGroup = true
	return os.WriteFile(filePath, buf.Bytes(), 0644)

}

// WriteToFile escreve os dados no arquivo
func (pgd *PackageGroupData) WriteServiceToFile(filePath string) error {
	if pgd.GenerateGroup {
		return nil
	}

	tmpl, err := template.New("package_group").Parse(serviceGroupTemplate)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, pgd)
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(filePath), 0755)
	pgd.GenerateGroup = true
	return os.WriteFile(filePath, buf.Bytes(), 0644)

}

// WriteToFile escreve os dados no arquivo
func (pgd *PackageGroupData) WriteProductToFile(filePath string) error {
	if pgd.GenerateGroup {
		return nil
	}

	tmpl, err := template.New("package_group").Parse(productTemplate)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, pgd)
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(filePath), 0755)
	pgd.GenerateGroup = true
	return os.WriteFile(filePath, buf.Bytes(), 0644)

}
func (pgd *PackageGroupData) Copy() PackageGroupData {
	return *pgd
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
		Imports:     []string{},
		SubCommands: []RootSubCommandData{},
	}
}

// AddImport adiciona um import à lista de imports
func (rgd *RootGenData) AddImport(importPath string) {
	if slices.Contains(rgd.Imports, importPath) {
		return
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
	tmpl, err := template.New("root_gen").Parse(rootGenTemplate)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, rgd)
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(filePath), 0755)
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}
