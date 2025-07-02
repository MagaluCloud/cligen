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
	genPackageCode(&sdkStructure)
	genServiceCode(&sdkStructure)
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
func genPackageCode(sdkStructure *sdk_structure.SDKStructure) error {
	for _, pkg := range sdkStructure.Packages {
		genPackageCodeRecursive(&pkg, nil)
	}

	return nil
}

func genPackageCodeRecursive(pkg *sdk_structure.Package, parentPkg *sdk_structure.Package) error {
	packageData := NewPackageGroupData()
	packageData.SetPackageName(pkg.Name)
	packageData.SetFunctionName(strutils.FirstUpper(pkg.Name))
	packageData.SetUseName(pkg.MenuName)
	packageData.SetDescriptions(defaultShortDesc, defaultLongDesc)
	packageData.SetGroupID(groupProducts)
	packageData.SetServiceParam(serviceParamPattern)
	packageData.AddImport(importSDK)

	if len(pkg.Services) > 0 {
		packageData.AddServiceInit(fmt.Sprintf("%sService := %sSdk.New(sdkCoreConfig)", pkg.Name, pkg.Name))
		for _, service := range pkg.Services {
			packageData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", pkg.Name, pkg.Name))
			packageData.AddImport(importCobra)
			if parentPkg != nil {
				packageData.AddImport(fmt.Sprintf("\"mgccli/cmd/gen/%s/%s/%s\"", strings.ToLower(parentPkg.Name), strings.ToLower(pkg.Name), strings.ToLower(service.Name)))
			} else {
				packageData.AddImport(fmt.Sprintf("\"mgccli/cmd/gen/%s/%s\"", strings.ToLower(pkg.Name), strings.ToLower(service.Name)))
			}
			packageData.AddSubCommand(service.Name, service.Name, fmt.Sprintf("%sService.%s()", pkg.Name, service.Name))
		}
	}

	if len(pkg.SubPkgs) > 0 {
		for _, subPkg := range pkg.SubPkgs {
			packageData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", subPkg.Name, subPkg.Name))
			packageData.AddImport(importCobra)
			packageData.AddImport(fmt.Sprintf("\"mgccli/cmd/gen/%s/%s\"", strings.ToLower(pkg.Name), strings.ToLower(subPkg.Name)))
			packageData.AddSubCommand(subPkg.Name, strutils.FirstUpper(subPkg.Name)+"Cmd", "sdkCoreConfig")
			genPackageCodeRecursive(&subPkg, pkg)
		}
	}
	dir := genDir
	if parentPkg != nil {
		dir = filepath.Join(dir, strings.ToLower(parentPkg.Name), strings.ToLower(pkg.Name), fmt.Sprintf("%s.go", pkg.Name))
	} else {
		dir = filepath.Join(dir, strings.ToLower(pkg.Name), fmt.Sprintf("%s.go", pkg.Name))
	}
	err := packageData.WriteGroupToFile(dir)
	if err != nil {
		return fmt.Errorf("erro ao escrever o arquivo %s.go para o pacote %s: %v", pkg.Name, pkg.Name, err)
	}
	return nil
}

func genServiceCode(sdkStructure *sdk_structure.SDKStructure) error {
	for _, pkg := range sdkStructure.Packages {
		genServiceCodeRecursive(&pkg, nil)
	}
	return nil
}

// func genServiceCodeRecursive(pkg *sdk_structure.Package, parentPkg *sdk_structure.Package) error {
// 	serviceData := NewPackageGroupData()
// 	serviceData.SetPackageName(service.Name)
// 	serviceData.SetFunctionName(service.Name)
// 	serviceData.SetUseName(service.Name)
// 	serviceData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", pkg.Name, pkg.Name))
// 	serviceData.AddImport(importCobra)
// 	serviceData.SetDescriptions(defaultShortDesc, defaultLongDesc)
// 	serviceData.SetServiceParam(fmt.Sprintf("%s %sSdk.%s", strutils.FirstLower(service.Interface), pkg.Name, service.Interface))

// }

func genServiceCodeRecursive(pkg *sdk_structure.Package, parentPkg *sdk_structure.Package) error {
	if len(pkg.Services) == 0 {
		return nil
	}
	if len(pkg.SubPkgs) > 0 {
		return nil
	}
	serviceData := NewPackageGroupData()
	serviceData.SetPackageName(pkg.Name)
	serviceData.SetFunctionName(pkg.Name)
	serviceData.SetUseName(pkg.Name)
	serviceData.AddImport(importCobra)
	serviceData.AddImport(importSDK)
	serviceData.SetDescriptions(defaultShortDesc, defaultLongDesc)
	serviceData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", pkg.Name, pkg.Name))
	serviceData.SetServiceParam("sdkCoreConfig *sdk.CoreClient")

	serviceData.AddServiceInit(fmt.Sprintf("%sService := %sSdk.New(sdkCoreConfig)", pkg.Name, pkg.Name))

	for _, service := range pkg.Services {

		serviceData.AddImport(fmt.Sprintf("\"mgccli/cmd/gen/%s/%s\"", strings.ToLower(pkg.Name), strings.ToLower(service.Name)))
		serviceData.AddSubCommand(service.Name, service.Name, fmt.Sprintf("%sService.%s()", pkg.Name, service.Name))
	}
	dir := genDir
	if parentPkg != nil {
		dir = filepath.Join(dir, strings.ToLower(parentPkg.Name), strings.ToLower(pkg.Name), fmt.Sprintf("%s.go", pkg.Name))
	} else {
		dir = filepath.Join(dir, strings.ToLower(pkg.Name), fmt.Sprintf("%s.go", pkg.Name))
	}

	err := serviceData.WriteGroupToFile(dir)
	if err != nil {
		return fmt.Errorf("erro ao escrever o arquivo %s.go para o serviço %s do pacote %s: %v", pkg.Name, pkg.Name, pkg.Name, err)
	}

	return nil
}

// 	// Usar mutex para operações thread-safe no rootGenData
// 	var mu sync.Mutex
// 	mu.Lock()
// 	rootGenData.AddSubCommand(pkg.Name, strutils.FirstUpper(pkg.Name)+"Cmd")
// 	rootGenData.AddImport(fmt.Sprintf("\"mgccli/cmd/gen/%s\"", strings.ToLower(pkg.Name)))
// 	mu.Unlock()

// 	// Se o pacote tem subpacotes (menu de agrupamento), processar os subpacotes
// 	if len(pkg.SubPkgs) > 0 {
// 		for subPkgName, subPkg := range pkg.SubPkgs {

// 			packageData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", subPkgName, subPkgName))
// 			packageData.AddImport(importCobra)
// 			packageData.AddImport(fmt.Sprintf("\"mgccli/cmd/gen/%s/%s\"", strings.ToLower(pkg.Name), strings.ToLower(subPkgName)))
// 			packageData.AddServiceInit(fmt.Sprintf("%sService := %sSdk.New(sdkCoreConfig)", subPkgName, subPkgName))
// 			for _, service := range subPkg.Services {
// 				// packageData.AddSubCommand(service.Name, service.Name, fmt.Sprintf("%sService.%s()", subPkgName, service.Name))
// 				packageData.AddSubCommand(subPkgName, strutils.ToPascalCase(subPkg.MenuName), fmt.Sprintf("%sService.%s()", subPkgName, service.Name))
// 			}

// 			// Gerar código para o subpacote
// 			if err := generateSubPackageCodeParallel(pkg.Name, &subPkg, *packageData); err != nil {
// 				return fmt.Errorf("erro ao gerar código do subpacote %s: %v", subPkgName, err)
// 			}
// 		}
// 	} else if len(pkg.Services) == 0 {
// 		// Se o pacote não tem serviços nem subpacotes, criar apenas o comando de agrupamento
// 		packageData.AddImport(importSDK)
// 		packageData.AddImport(importCobra)
// 		err := packageData.WriteGroupToFile(filepath.Join(genDir, strings.ToLower(pkg.Name), fmt.Sprintf("%s.go", pkg.Name)))
// 		if err != nil {
// 			return fmt.Errorf("erro ao escrever o arquivo %s.go para o pacote %s: %v", pkg.Name, pkg.Name, err)
// 		}
// 		return nil
// 	} else {
// 		// Processar serviços normalmente
// 		for _, service := range pkg.Services {
// 			packageData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", pkg.Name, pkg.Name))
// 			packageData.AddImport(importCobra)
// 			packageData.AddImport(fmt.Sprintf("\"mgccli/cmd/gen/%s/%s\"", strings.ToLower(pkg.Name), strings.ToLower(service.Name)))
// 			packageData.AddServiceInit(fmt.Sprintf("%sService := %sSdk.New(sdkCoreConfig)", pkg.Name, pkg.Name))
// 			packageData.AddSubCommand(service.Name, service.Name, fmt.Sprintf("%sService.%s()", pkg.Name, service.Name))

// 			if err := generateServiceCodeParallel(*pkg, &service, *packageData, ""); err != nil {
// 				return fmt.Errorf("erro ao gerar código do serviço %s: %v", service.Name, err)
// 			}
// 		}
// 	}

// 	packageData.AddImport(importSDK)
// 	err := packageData.WriteGroupToFile(filepath.Join(genDir, strings.ToLower(pkg.Name), fmt.Sprintf("%s.go", pkg.Name)))
// 	if err != nil {
// 		return fmt.Errorf("erro ao escrever o arquivo %s.go para o pacote %s: %v", pkg.Name, pkg.Name, err)
// 	}
// 	return nil
// }

// // generateServiceCodeParallel é a versão thread-safe da função generateServiceCode
// func generateServiceCodeParallel(parentPkg sdk_structure.Package, service *sdk_structure.Service, data PackageGroupData, isChieldOf string) error {
// 	genDirReplace := genDir
// 	if isChieldOf != "" {
// 		genDirReplace = path.Join(genDirReplace, isChieldOf)
// 	}
// 	dir := filepath.Join(genDirReplace, strings.ToLower(parentPkg.Name), strings.ToLower(service.Name))
// 	if err := os.MkdirAll(dir, 0755); err != nil {
// 		return fmt.Errorf("erro ao criar diretório %s: %v", dir, err)
// 	}

// 	serviceData := data.Copy()
// 	serviceData.SetPackageName(service.Name)
// 	serviceData.SetFunctionName(service.Name)
// 	serviceData.SetUseName(service.Name)
// 	filteredImports := make([]string, 0, len(serviceData.Imports))
// 	serviceDir := ""
// 	if isChieldOf != "" {
// 		serviceDir = fmt.Sprintf("mgccli/cmd/gen/%s/%s/%s", strings.ToLower(parentPkg.Name), strings.ToLower(isChieldOf), strings.ToLower(service.Name))
// 	} else {
// 		serviceDir = fmt.Sprintf("mgccli/cmd/gen/%s", strings.ToLower(parentPkg.Name))
// 	}

// 	for _, imp := range serviceData.Imports {
// 		if !strings.Contains(imp, serviceDir) && !strings.Contains(imp, fmt.Sprintf("mgccli/cmd/gen/%s/", strings.ToLower(parentPkg.Name))) {
// 			filteredImports = append(filteredImports, imp)
// 		}
// 	}
// 	serviceData.Imports = filteredImports
// 	serviceData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", parentPkg.Name, parentPkg.Name))
// 	serviceData.AddImport(importCobra)
// 	serviceData.SetDescriptions(defaultShortDesc, defaultLongDesc)
// 	serviceData.SetGroupID("")
// 	serviceData.SetServiceParam(fmt.Sprintf("%s %sSdk.%s", strutils.FirstLower(service.Interface), parentPkg.Name, service.Interface))

// 	for _, method := range service.Methods {
// 		productData := serviceData.Copy()
// 		productData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", parentPkg.Name, parentPkg.Name))
// 		productData.AddImport(importCobra)
// 		serviceData.AddCommand(method.Name, strutils.FirstLower(service.Interface))
// 		productData.AddCommand(method.Name, strutils.FirstLower(service.Interface))
// 		productData.SetServiceCall(fmt.Sprintf("%s.%s", strutils.FirstLower(service.Interface), method.Name))
// 		productData.SetFunctionName(method.Name)
// 		productData.SetUseName(method.Name)
// 		productData.SetDescriptions(defaultShortDesc, defaultLongDesc)

// 		for key, param := range method.Parameters {
// 			if key == "ctx" {
// 				productData.AddParam("context.Background()")
// 			} else {
// 				productData.AddParam(fmt.Sprintf("%s.%s{}", parentPkg.Name, param))
// 			}
// 		}
// 		err := productData.WriteProductToFile(filepath.Join(dir, fmt.Sprintf("%s.go", strings.ToLower(method.Name))))
// 		if err != nil {
// 			return fmt.Errorf("erro ao escrever o arquivo %s.go para o método %s do serviço %s do pacote %s: %v", method.Name, method.Name, service.Name, parentPkg.Name, err)
// 		}
// 	}

// 	err := serviceData.WriteServiceToFile(filepath.Join(dir, fmt.Sprintf("%s.go", strings.ToLower(service.Name))))
// 	if err != nil {
// 		return fmt.Errorf("erro ao escrever o arquivo %s.go para o serviço %s do pacote %s: %v", service.Name, service.Name, parentPkg.Name, err)
// 	}

// 	for _, subService := range service.SubServices {
// 		if err := generateServiceCodeParallel(parentPkg, &subService, data, isChieldOf); err != nil {
// 			return fmt.Errorf("erro ao gerar código do subserviço %s: %v", subService.Name, err)
// 		}
// 	}
// 	return nil
// }

// // generateSubPackageCodeParallel gera código para um subpacote dentro de um menu de agrupamento
// func generateSubPackageCodeParallel(parentPkgName string, subPkg *sdk_structure.Package, data PackageGroupData) error {
// 	dir := filepath.Join(genDir, strings.ToLower(parentPkgName), strings.ToLower(subPkg.Name))
// 	if err := os.MkdirAll(dir, 0755); err != nil {
// 		return fmt.Errorf("erro ao criar diretório %s: %v", dir, err)
// 	}

// 	subPackageData := NewPackageGroupData()
// 	subPackageData.SetPackageName(subPkg.Name)
// 	subPackageData.SetFunctionName(strutils.ToPascalCase(subPkg.MenuName))
// 	subPackageData.SetUseName(subPkg.MenuName)
// 	subPackageData.SetDescriptions(defaultShortDesc, defaultLongDesc)
// 	subPackageData.SetGroupID("")
// 	subPackageData.SetServiceParam(fmt.Sprintf("%sService *%sSdk.Client", strutils.FirstLower(subPkg.Name), subPkg.Name))

// 	subPackageData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", subPkg.Name, subPkg.Name))
// 	subPackageData.AddImport(importCobra)

// 	for _, service := range subPkg.Services {
// 		// Gerar código para o serviço diretamente no diretório do subpacote
// 		if err := generateServiceCodeParallel(*subPkg, &service, *subPackageData, parentPkgName); err != nil {
// 			return fmt.Errorf("erro ao gerar código do serviço %s do subpacote %s: %v", service.Name, subPkg.Name, err)
// 		}
// 	}

// 	err := subPackageData.WriteGroupToFile(filepath.Join(dir, fmt.Sprintf("%s.go", strings.ToLower(subPkg.Name))))
// 	if err != nil {
// 		return fmt.Errorf("erro ao escrever o arquivo %s.go para o subpacote %s: %v", subPkg.Name, subPkg.Name, err)
// 	}

// 	return nil
// }
