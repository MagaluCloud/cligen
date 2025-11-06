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

	log.Printf("🔧 Iniciando geração do CLI com %d pacotes", len(sdkStructure.Packages))
	cleanDir(genDir)
	genGoModFile()
	generateRootCode(&sdkStructure)
	genMainPackageCode(&sdkStructure)
	genPackageCode(&sdkStructure)
	genServiceCode(&sdkStructure)
	genProductCode(&sdkStructure)

}

func cleanDir(dir string) {
	toRemove := filepath.Clean(dir)
	if _, err := os.Stat(toRemove); err == nil {
		os.RemoveAll(toRemove)
	}
}

func genMainPackageCode(sdkStructure *sdk_structure.SDKStructure) error {

	for _, pkg := range sdkStructure.Packages {
		genMainPackageCodeRecursive(&pkg, nil)
	}

	return nil
}

func genMainPackageCodeRecursive(pkg *sdk_structure.Package, parentPkg *sdk_structure.Package) error {
	mainPackageData := NewPackageGroupData()
	mainPackageData.SetPackageName(pkg.Name)
	mainPackageData.SetFunctionName(strutils.FirstUpper(pkg.Name))
	mainPackageData.SetUseName(pkg.MenuName)
	mainPackageData.SetDescriptions(pkg.Description, "defaultLongDesc 1")
	mainPackageData.SetServiceParam(serviceParamPattern)
	mainPackageData.AddImport(importSDK)
	mainPackageData.SetGroupID(groupProducts)

	if len(pkg.Services) > 0 {
		setupMainPackageServices(mainPackageData, pkg, parentPkg)
	}

	if len(pkg.SubPkgs) > 0 {
		if err := setupMainPackageSubPackages(mainPackageData, pkg, parentPkg); err != nil {
			return err
		}
	}

	filePath := buildMainPackageFilePath(pkg, parentPkg)
	var err error
	if parentPkg != nil {
		err = mainPackageData.WriteSubPackageToFile(filePath)
	} else {
		err = mainPackageData.WriteGroupToFile(filePath)
	}

	if err != nil {
		return fmt.Errorf("erro ao escrever o arquivo %s.go para o pacote %s: %w", pkg.Name, pkg.Name, err)
	}
	return nil
}

func setupMainPackageServices(pkgData *PackageGroupData, pkg *sdk_structure.Package, parentPkg *sdk_structure.Package) {
	pkgData.AddServiceInit(fmt.Sprintf("%sService := %sSdk.New(&sdkCoreConfig)", pkg.Name, pkg.Name))
	pkgData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", pkg.Name, pkg.Name))
	pkgData.AddImport(importCobra)

	for _, service := range pkg.Services {
		serviceImport := buildServiceImportPath(parentPkg, pkg.Name, service.Name)
		pkgData.AddImport(serviceImport)
		pkgData.AddSubCommand(service.Name, service.Name, fmt.Sprintf("%sService.%s()", pkg.Name, service.Name))
	}
}

func setupMainPackageSubPackages(pkgData *PackageGroupData, pkg *sdk_structure.Package, parentPkg *sdk_structure.Package) error {
	pkgData.AddImport(importCobra)

	for _, subPkg := range pkg.SubPkgs {
		subPkgImport := fmt.Sprintf("\"github.com/magaluCloud/mgccli/cmd/gen/%s/%s\"", strings.ToLower(pkg.Name), strings.ToLower(subPkg.Name))
		pkgData.AddImport(subPkgImport)
		pkgData.AddSubCommand(subPkg.Name, strutils.FirstUpper(subPkg.Name), "sdkCoreConfig")
		if err := genMainPackageCodeRecursive(&subPkg, pkg); err != nil {
			return err
		}
	}
	return nil
}

func buildServiceImportPath(parentPkg *sdk_structure.Package, pkgName, serviceName string) string {
	if parentPkg != nil {
		return fmt.Sprintf("\"github.com/magaluCloud/mgccli/cmd/gen/%s/%s/%s\"",
			strings.ToLower(parentPkg.Name), strings.ToLower(pkgName), strings.ToLower(serviceName))
	}
	return fmt.Sprintf("\"github.com/magaluCloud/mgccli/cmd/gen/%s/%s\"",
		strings.ToLower(pkgName), strings.ToLower(serviceName))
}

func buildMainPackageFilePath(pkg *sdk_structure.Package, parentPkg *sdk_structure.Package) string {
	if parentPkg != nil {
		return filepath.Join(genDir, strings.ToLower(parentPkg.Name), strings.ToLower(pkg.Name), fmt.Sprintf("%s.go", pkg.Name))
	}
	return filepath.Join(genDir, strings.ToLower(pkg.Name), fmt.Sprintf("%s.go", pkg.Name))
}
