package gen_cli_code

import (
	"cligen/commands/sdk_structure"
	strutils "cligen/str_utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	// Diretório onde os arquivos gerados serão salvos
	genDir = "base-cli-gen/cmd/gen"

	// Descrições padrão para comandos (placeholders)
	defaultShortDesc = "todo"
	defaultLongDesc  = "todo2"

	// Nomes de grupos para agrupamento de comandos no CLI
	groupProducts = "products"
	groupSettings = "settings"
	groupOther    = "other"

	// Padrões de imports comuns
	importCobra = "\"github.com/spf13/cobra\""
	importSDK   = "sdk \"github.com/MagaluCloud/mgc-sdk-go/client\""

	// Padrões de parâmetros de serviço
	serviceParamPattern = "sdkCoreConfig *sdk.CoreClient"
)

func GenCliCode() {
	sdkStructure, err := sdk_structure.GenCliSDKStructure()
	if err != nil {
		log.Fatalf("Erro ao gerar a estrutura do SDK: %v", err)
	}

	log.Printf("🔧 Iniciando geração do CLI com %d pacotes", len(sdkStructure.Packages))
	cleanDir(genDir)

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

func generateRootCode(sdkStructure *sdk_structure.SDKStructure) error {
	rootGenData := NewRootGenData()
	rootGenData.AddImport(importSDK)
	rootGenData.AddImport(importCobra)
	for _, pkg := range sdkStructure.Packages {
		rootGenData.AddSubCommand(pkg.Name, strutils.FirstUpper(pkg.Name)+"Cmd")
		rootGenData.AddImport(fmt.Sprintf("\"mgccli/cmd/gen/%s\"", strings.ToLower(pkg.Name)))
	}
	if err := rootGenData.WriteRootGenToFile(filepath.Join(genDir, "root_gen.go")); err != nil {
		log.Fatalf("Erro ao escrever o arquivo root_gen.go: %v", err)
	}
	return nil
}

// // genPackageCodeParallel é a versão thread-safe da função genPackageCode
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
	mainPackageData.SetDescriptions(defaultShortDesc, defaultLongDesc)
	mainPackageData.SetServiceParam(serviceParamPattern)
	mainPackageData.AddImport(importSDK)

	if len(pkg.Services) > 0 {
		if parentPkg == nil {
			mainPackageData.SetGroupID(groupProducts)
		}
		mainPackageData.AddServiceInit(fmt.Sprintf("%sService := %sSdk.New(sdkCoreConfig)", pkg.Name, pkg.Name))
		for _, service := range pkg.Services {
			mainPackageData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", pkg.Name, pkg.Name))
			mainPackageData.AddImport(importCobra)
			if parentPkg != nil {
				mainPackageData.AddImport(fmt.Sprintf("\"mgccli/cmd/gen/%s/%s/%s\"", strings.ToLower(parentPkg.Name), strings.ToLower(pkg.Name), strings.ToLower(service.Name)))
			} else {
				mainPackageData.AddImport(fmt.Sprintf("\"mgccli/cmd/gen/%s/%s\"", strings.ToLower(pkg.Name), strings.ToLower(service.Name)))
			}
			mainPackageData.AddSubCommand(service.Name, service.Name, fmt.Sprintf("%sService.%s()", pkg.Name, service.Name))
		}
	}

	if len(pkg.SubPkgs) > 0 {
		for _, subPkg := range pkg.SubPkgs {
			mainPackageData.AddImport(importCobra)
			mainPackageData.AddImport(fmt.Sprintf("\"mgccli/cmd/gen/%s/%s\"", strings.ToLower(pkg.Name), strings.ToLower(subPkg.Name)))
			mainPackageData.AddSubCommand(subPkg.Name, strutils.FirstUpper(subPkg.Name), "sdkCoreConfig")
			genMainPackageCodeRecursive(&subPkg, pkg)
		}
	}
	dir := genDir
	if parentPkg != nil {
		dir = filepath.Join(dir, strings.ToLower(parentPkg.Name), strings.ToLower(pkg.Name), fmt.Sprintf("%s.go", pkg.Name))
	} else {
		dir = filepath.Join(dir, strings.ToLower(pkg.Name), fmt.Sprintf("%s.go", pkg.Name))
	}
	err := mainPackageData.WriteGroupToFile(dir)
	if err != nil {
		return fmt.Errorf("erro ao escrever o arquivo %s.go para o pacote %s: %v", pkg.Name, pkg.Name, err)
	}
	return nil
}

func genPackageCode(sdkStructure *sdk_structure.SDKStructure) error {
	for _, pkg := range sdkStructure.Packages {
		genPackageCodeRecursive(&pkg, nil)
	}
	return nil
}

func genPackageCodeRecursive(pkg *sdk_structure.Package, parentPkg *sdk_structure.Package) error {
	if len(pkg.Services) == 0 {
		return nil
	}
	if len(pkg.SubPkgs) > 0 {
		return nil
	}
	packageData := NewPackageGroupData()
	packageData.SetPackageName(pkg.Name)
	packageData.SetFunctionName(strutils.FirstUpper(pkg.Name))
	packageData.SetUseName(pkg.Name)
	packageData.AddImport(importCobra)
	packageData.AddImport(importSDK)
	packageData.SetDescriptions(defaultShortDesc, defaultLongDesc)
	packageData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", pkg.Name, pkg.Name))
	packageData.SetServiceParam("sdkCoreConfig *sdk.CoreClient")

	packageData.AddServiceInit(fmt.Sprintf("%sService := %sSdk.New(sdkCoreConfig)", pkg.Name, pkg.Name))

	for _, service := range pkg.Services {
		packageData.AddImport(fmt.Sprintf("\"mgccli/cmd/gen/%s/%s\"", strings.ToLower(pkg.Name), strings.ToLower(service.Name)))
		packageData.AddSubCommand(service.Name, service.Name, fmt.Sprintf("%sService.%s()", pkg.Name, service.Name))
	}
	dir := genDir
	if parentPkg != nil {
		dir = filepath.Join(dir, strings.ToLower(parentPkg.Name), strings.ToLower(pkg.Name), fmt.Sprintf("%s.go", pkg.Name))
	} else {
		dir = filepath.Join(dir, strings.ToLower(pkg.Name), fmt.Sprintf("%s.go", pkg.Name))
	}

	err := packageData.WriteGroupToFile(dir)
	if err != nil {
		return fmt.Errorf("erro ao escrever o arquivo %s.go para o serviço %s do pacote %s: %v", pkg.Name, pkg.Name, pkg.Name, err)
	}

	return nil
}

func genServiceCode(sdkStructure *sdk_structure.SDKStructure) error {
	for _, pkg := range sdkStructure.Packages {
		genServiceCodeRecursive(&pkg, nil)
	}
	return nil
}

func genServiceCodeRecursive(pkg *sdk_structure.Package, parentPkg *sdk_structure.Package) error {
	if len(pkg.Services) > 0 {
		for _, service := range pkg.Services {
			serviceData := NewPackageGroupData()
			serviceData.SetPackageName(service.Name)
			serviceData.SetFunctionName(service.Name)
			serviceData.SetUseName(strutils.FirstLower(service.Name))
			serviceData.AddImport(importCobra)
			serviceData.SetServiceParam(fmt.Sprintf("%s %sSdk.%s", strutils.FirstLower(service.Interface), pkg.Name, service.Interface))
			for _, method := range service.Methods {
				serviceData.SetDescriptions(defaultShortDesc, defaultLongDesc)
				serviceData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", pkg.Name, pkg.Name))
				serviceData.AddCommand(method.Name, strutils.FirstLower(service.Interface))
			}
			dir := genDir
			if parentPkg != nil {
				dir = filepath.Join(dir, strings.ToLower(parentPkg.Name), strings.ToLower(pkg.Name), strings.ToLower(service.Name), fmt.Sprintf("%s.go", strings.ToLower(service.Name)))
			} else {
				dir = filepath.Join(dir, strings.ToLower(pkg.Name), strings.ToLower(service.Name), fmt.Sprintf("%s.go", strings.ToLower(service.Name)))
			}
			err := serviceData.WriteServiceToFile(dir)
			if err != nil {
				return fmt.Errorf("erro ao escrever o arquivo %s.go para o serviço %s do pacote %s: %v", pkg.Name, pkg.Name, pkg.Name, err)
			}
		}
	}
	if len(pkg.SubPkgs) > 0 {
		for _, subPkg := range pkg.SubPkgs {
			genServiceCodeRecursive(&subPkg, pkg)
		}
	}
	return nil
}

func genProductCode(sdkStructure *sdk_structure.SDKStructure) error {
	for _, pkg := range sdkStructure.Packages {
		genProductCodeRecursive(&pkg, nil)
	}
	return nil
}

func genProductCodeRecursive(pkg *sdk_structure.Package, parentPkg *sdk_structure.Package) error {
	if len(pkg.Services) > 0 {
		for _, service := range pkg.Services {
			for _, method := range service.Methods {
				productData := NewPackageGroupData()
				productData.SetPackageName(service.Name)
				productData.AddImport(importCobra)
				productData.SetServiceParam(fmt.Sprintf("%s %sSdk.%s", strutils.FirstLower(service.Interface), pkg.Name, service.Interface))
				productData.SetFunctionName(method.Name)
				productData.SetUseName(strutils.FirstLower(method.Name))
				productData.SetDescriptions(defaultShortDesc, defaultLongDesc)
				productData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", pkg.Name, pkg.Name))
				productData.AddCommand(method.Name, strutils.FirstLower(service.Interface))

				dir := genDir
				if parentPkg != nil {
					dir = filepath.Join(dir, strings.ToLower(parentPkg.Name), strings.ToLower(pkg.Name), strings.ToLower(service.Name), fmt.Sprintf("%s.go", strings.ToLower(method.Name)))
				} else {
					dir = filepath.Join(dir, strings.ToLower(pkg.Name), strings.ToLower(service.Name), fmt.Sprintf("%s.go", strings.ToLower(method.Name)))
				}
				err := productData.WriteProductToFile(dir)
				if err != nil {
					return fmt.Errorf("erro ao escrever o arquivo %s.go para o serviço %s do pacote %s: %v", pkg.Name, pkg.Name, pkg.Name, err)
				}
			}
		}
	}
	if len(pkg.SubPkgs) > 0 {
		for _, subPkg := range pkg.SubPkgs {
			genProductCodeRecursive(&subPkg, pkg)
		}
	}
	return nil
}
