package gen_cli_code

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/magaluCloud/cligen/commands/sdk_structure"
	strutils "github.com/magaluCloud/cligen/str_utils"
)

const (
	// Diretório onde os arquivos gerados serão salvos
	genDir = "base-cli-gen/cmd/gen"

	// Descrições padrão para comandos (placeholders)
	defaultShortDesc = "todo aaa"
	defaultLongDesc  = "todo2"

	// Nomes de grupos para agrupamento de comandos no CLI
	groupProducts = "products"
	groupSettings = "settings"
	groupOther    = "other"

	// Padrões de imports comuns
	importCobra = "\"github.com/spf13/cobra\""
	importSDK   = "sdk \"github.com/MagaluCloud/mgc-sdk-go/client\""

	// Padrões de parâmetros de serviço
	serviceParamPattern = "sdkCoreConfig sdk.CoreClient"
)

func GenCliCode() {
	sdkStructure, err := sdk_structure.GenCliSDKStructure()
	if err != nil {
		log.Fatalf("Erro ao gerar a estrutura do SDK: %v", err)
	}

	custom := NewCustom()
	err = custom.Load()
	if err != nil {
		panic(fmt.Errorf("erro ao carregar os comandos customizados: %v", err))
	}

	log.Printf("🔧 Iniciando geração do CLI com %d pacotes", len(sdkStructure.Packages))
	cleanDir(genDir)
	genGoModFile()
	generateRootCode(&sdkStructure)
	genMainPackageCode(custom, &sdkStructure)
	genPackageCode(custom, &sdkStructure)
	genServiceCode(custom, &sdkStructure)
	genProductCode(custom, &sdkStructure)
	err = custom.Write()
	if err != nil {
		panic(fmt.Errorf("erro ao escrever os comandos customizados: %v", err))
	}
}

func cleanDir(dir string) {
	toRemove := filepath.Clean(dir)
	if _, err := os.Stat(toRemove); err == nil {
		os.RemoveAll(toRemove)
	}
}

func genMainPackageCode(custom *CustomHeader, sdkStructure *sdk_structure.SDKStructure) error {

	for _, pkg := range sdkStructure.Packages {
		genMainPackageCodeRecursive(custom, &pkg, nil)
	}

	return nil
}

func genMainPackageCodeRecursive(custom *CustomHeader, pkg *sdk_structure.Package, parentPkg *sdk_structure.Package) error {
	mainPackageData := NewPackageGroupData(custom)
	mainPackageData.SetPackageName(pkg.Name)
	mainPackageData.SetFunctionName(strutils.FirstUpper(pkg.Name))
	mainPackageData.SetUseName(pkg.MenuName)
	mainPackageData.SetDescriptions(pkg.Description, "defaultLongDesc 1")
	mainPackageData.SetServiceParam(serviceParamPattern)
	mainPackageData.AddImport(importSDK)
	mainPackageData.SetGroupID(groupProducts)
	if len(pkg.Services) > 0 {
		mainPackageData.AddServiceInit(fmt.Sprintf("%sService := %sSdk.New(&sdkCoreConfig)", pkg.Name, pkg.Name))
		for _, service := range pkg.Services {
			mainPackageData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", pkg.Name, pkg.Name))
			mainPackageData.AddImport(importCobra)
			if parentPkg != nil {
				mainPackageData.AddImport(fmt.Sprintf("\"github.com/magaluCloud/mgccli/cmd/gen/%s/%s/%s\"", strings.ToLower(parentPkg.Name), strings.ToLower(pkg.Name), strings.ToLower(service.Name)))
			} else {
				mainPackageData.AddImport(fmt.Sprintf("\"github.com/magaluCloud/mgccli/cmd/gen/%s/%s\"", strings.ToLower(pkg.Name), strings.ToLower(service.Name)))
			}
			mainPackageData.AddSubCommand(service.Name, service.Name, fmt.Sprintf("%sService.%s()", pkg.Name, service.Name))
		}
	}

	if len(pkg.SubPkgs) > 0 {
		for _, subPkg := range pkg.SubPkgs {
			mainPackageData.AddImport(importCobra)
			mainPackageData.AddImport(fmt.Sprintf("\"github.com/magaluCloud/mgccli/cmd/gen/%s/%s\"", strings.ToLower(pkg.Name), strings.ToLower(subPkg.Name)))
			mainPackageData.AddSubCommand(subPkg.Name, strutils.FirstUpper(subPkg.Name), "sdkCoreConfig")
			genMainPackageCodeRecursive(custom, &subPkg, pkg)
		}
	}
	var err error
	dir := genDir
	if parentPkg != nil {
		dir = filepath.Join(dir, strings.ToLower(parentPkg.Name), strings.ToLower(pkg.Name), fmt.Sprintf("%s.go", pkg.Name))
		err = mainPackageData.WriteSubPackageToFile(dir)
	} else {
		dir = filepath.Join(dir, strings.ToLower(pkg.Name), fmt.Sprintf("%s.go", pkg.Name))
		err = mainPackageData.WriteGroupToFile(dir)
	}

	if err != nil {
		return fmt.Errorf("erro ao escrever o arquivo %s.go para o pacote %s: %v", pkg.Name, pkg.Name, err)
	}
	return nil
}
