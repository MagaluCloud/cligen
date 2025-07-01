package gen_cli_code

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"cligen/commands/sdk_structure"
	strutils "cligen/str_utils"
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

	rootGenData := NewRootGenData()
	rootGenData.AddImport(importSDK)
	rootGenData.AddImport(importCobra)

	// Canal para receber erros das goroutines
	errChan := make(chan error, len(sdkStructure.Packages))
	var wg sync.WaitGroup

	// Processar pacotes em paralelo
	for _, pkg := range sdkStructure.Packages {
		wg.Add(1)
		go func(pkg sdk_structure.Package) {
			defer wg.Done()
			log.Printf("📦 Processando pacote: %s", pkg.Name)
			if err := genPackageCodeParallel(&pkg, rootGenData); err != nil {
				errChan <- fmt.Errorf("erro ao processar pacote %s: %v", pkg.Name, err)
			}
		}(pkg)
	}

	// Aguardar todas as goroutines terminarem
	wg.Wait()
	close(errChan)

	// Verificar se houve erros
	for err := range errChan {
		log.Fatalf("❌ Erro na geração paralela: %v", err)
	}

	if err := rootGenData.WriteRootGenToFile(filepath.Join(genDir, "root_gen.go")); err != nil {
		log.Fatalf("Erro ao escrever o arquivo root_gen.go: %v", err)
	}

	log.Printf("✅ Geração do CLI concluída com sucesso")
}

func cleanDir(dir string) {
	toRemove := filepath.Clean(dir)
	if _, err := os.Stat(toRemove); err == nil {
		os.RemoveAll(toRemove)
	}
}

// genPackageCodeParallel é a versão thread-safe da função genPackageCode
func genPackageCodeParallel(pkg *sdk_structure.Package, rootGenData *RootGenData) error {
	packageData := NewPackageGroupData()
	packageData.SetPackageName(pkg.Name)
	packageData.SetFunctionName(strutils.FirstUpper(pkg.Name))
	packageData.SetUseName(pkg.MenuName)
	packageData.SetDescriptions(defaultShortDesc, defaultLongDesc)
	packageData.SetGroupID(groupProducts)
	packageData.SetServiceParam(serviceParamPattern)

	// Usar mutex para operações thread-safe no rootGenData
	var mu sync.Mutex
	mu.Lock()
	rootGenData.AddSubCommand(pkg.Name, strutils.FirstUpper(pkg.Name)+"Cmd")
	rootGenData.AddImport(fmt.Sprintf("\"mgccli/cmd/gen/%s\"", strings.ToLower(pkg.Name)))
	mu.Unlock()

	// Se o pacote tem subpacotes (menu de agrupamento), processar os subpacotes
	if len(pkg.SubPkgs) > 0 {
		for subPkgName, subPkg := range pkg.SubPkgs {
			packageData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", subPkgName, subPkgName))
			packageData.AddImport(importCobra)
			packageData.AddImport(fmt.Sprintf("\"mgccli/cmd/gen/%s/%s\"", strings.ToLower(pkg.Name), strings.ToLower(subPkgName)))
			packageData.AddServiceInit(fmt.Sprintf("%sService := %sSdk.New(sdkCoreConfig)", subPkgName, subPkgName))
			packageData.AddSubCommand(subPkgName, strutils.ToPascalCase(subPkg.MenuName), fmt.Sprintf("%sService", subPkgName))

			// Gerar código para o subpacote
			if err := generateSubPackageCodeParallel(pkg.Name, &subPkg, *packageData); err != nil {
				return fmt.Errorf("erro ao gerar código do subpacote %s: %v", subPkgName, err)
			}
		}
	} else if len(pkg.Services) == 0 {
		// Se o pacote não tem serviços nem subpacotes, criar apenas o comando de agrupamento
		packageData.AddImport(importSDK)
		packageData.AddImport(importCobra)
		err := packageData.WriteGroupToFile(filepath.Join(genDir, strings.ToLower(pkg.Name), fmt.Sprintf("%s.go", pkg.Name)))
		if err != nil {
			return fmt.Errorf("erro ao escrever o arquivo %s.go para o pacote %s: %v", pkg.Name, pkg.Name, err)
		}
		return nil
	} else {
		// Processar serviços normalmente
		for _, service := range pkg.Services {
			packageData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", pkg.Name, pkg.Name))
			packageData.AddImport(importCobra)
			packageData.AddImport(fmt.Sprintf("\"mgccli/cmd/gen/%s/%s\"", strings.ToLower(pkg.Name), strings.ToLower(service.Name)))
			packageData.AddServiceInit(fmt.Sprintf("%sService := %sSdk.New(sdkCoreConfig)", pkg.Name, pkg.Name))
			packageData.AddSubCommand(service.Name, service.Name, fmt.Sprintf("%sService.%s()", pkg.Name, service.Name))

			if err := generateServiceCodeParallel(*pkg, &service, *packageData); err != nil {
				return fmt.Errorf("erro ao gerar código do serviço %s: %v", service.Name, err)
			}
		}
	}

	packageData.AddImport(importSDK)
	err := packageData.WriteGroupToFile(filepath.Join(genDir, strings.ToLower(pkg.Name), fmt.Sprintf("%s.go", pkg.Name)))
	if err != nil {
		return fmt.Errorf("erro ao escrever o arquivo %s.go para o pacote %s: %v", pkg.Name, pkg.Name, err)
	}
	return nil
}

// generateServiceCodeParallel é a versão thread-safe da função generateServiceCode
func generateServiceCodeParallel(parentPkg sdk_structure.Package, service *sdk_structure.Service, data PackageGroupData) error {
	dir := filepath.Join(genDir, strings.ToLower(parentPkg.Name), strings.ToLower(service.Name))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório %s: %v", dir, err)
	}

	serviceData := data.Copy()
	serviceData.SetPackageName(service.Name)
	serviceData.SetFunctionName(service.Name)
	serviceData.SetUseName(service.Name)
	filteredImports := make([]string, 0, len(serviceData.Imports))
	serviceDir := fmt.Sprintf("mgccli/cmd/gen/%s/%s", strings.ToLower(parentPkg.Name), strings.ToLower(service.Name))
	for _, imp := range serviceData.Imports {
		if !strings.Contains(imp, serviceDir) && !strings.Contains(imp, fmt.Sprintf("mgccli/cmd/gen/%s/", strings.ToLower(parentPkg.Name))) {
			filteredImports = append(filteredImports, imp)
		}
	}
	serviceData.Imports = filteredImports
	serviceData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", parentPkg.Name, parentPkg.Name))
	serviceData.AddImport(importCobra)
	serviceData.SetDescriptions(defaultShortDesc, defaultLongDesc)
	serviceData.SetGroupID("")
	serviceData.SetServiceParam(fmt.Sprintf("%s %sSdk.%s", strutils.FirstLower(service.Interface), parentPkg.Name, service.Interface))

	for _, method := range service.Methods {
		productData := serviceData.Copy()
		productData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", parentPkg.Name, parentPkg.Name))
		productData.AddImport(importCobra)
		serviceData.AddCommand(method.Name, strutils.FirstLower(service.Interface))
		productData.AddCommand(method.Name, strutils.FirstLower(service.Interface))
		productData.SetServiceCall(fmt.Sprintf("%s.%s", strutils.FirstLower(service.Interface), method.Name))
		productData.SetFunctionName(method.Name)
		productData.SetUseName(method.Name)
		productData.SetDescriptions(defaultShortDesc, defaultLongDesc)

		for key, param := range method.Parameters {
			if key == "ctx" {
				productData.AddParam("context.Background()")
			} else {
				productData.AddParam(fmt.Sprintf("%s.%s{}", parentPkg.Name, param))
			}
		}
		err := productData.WriteProductToFile(filepath.Join(dir, fmt.Sprintf("%s.go", strings.ToLower(method.Name))))
		if err != nil {
			return fmt.Errorf("erro ao escrever o arquivo %s.go para o método %s do serviço %s do pacote %s: %v", method.Name, method.Name, service.Name, parentPkg.Name, err)
		}
	}

	err := serviceData.WriteServiceToFile(filepath.Join(dir, fmt.Sprintf("%s.go", strings.ToLower(service.Name))))
	if err != nil {
		return fmt.Errorf("erro ao escrever o arquivo %s.go para o serviço %s do pacote %s: %v", service.Name, service.Name, parentPkg.Name, err)
	}

	for _, subService := range service.SubServices {
		if err := generateServiceCodeParallel(parentPkg, &subService, data); err != nil {
			return fmt.Errorf("erro ao gerar código do subserviço %s: %v", subService.Name, err)
		}
	}
	return nil
}

// generateSubPackageCodeParallel gera código para um subpacote dentro de um menu de agrupamento
func generateSubPackageCodeParallel(parentPkgName string, subPkg *sdk_structure.Package, data PackageGroupData) error {
	dir := filepath.Join(genDir, strings.ToLower(parentPkgName), strings.ToLower(subPkg.Name))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório %s: %v", dir, err)
	}

	subPackageData := NewPackageGroupData()
	subPackageData.SetPackageName(subPkg.Name)
	subPackageData.SetFunctionName(strutils.ToPascalCase(subPkg.MenuName))
	subPackageData.SetUseName(subPkg.MenuName)
	subPackageData.SetDescriptions(defaultShortDesc, defaultLongDesc)
	subPackageData.SetGroupID("")
	subPackageData.SetServiceParam(fmt.Sprintf("%sService *%sSdk.Client", strutils.FirstLower(subPkg.Name), subPkg.Name))

	subPackageData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", subPkg.Name, subPkg.Name))
	subPackageData.AddImport(importCobra)

	for _, service := range subPkg.Services {
		// Gerar código para o serviço diretamente no diretório do subpacote
		if err := generateServiceCodeParallel(*subPkg, &service, *subPackageData); err != nil {
			return fmt.Errorf("erro ao gerar código do serviço %s do subpacote %s: %v", service.Name, subPkg.Name, err)
		}
	}

	err := subPackageData.WriteGroupToFile(filepath.Join(dir, fmt.Sprintf("%s.go", strings.ToLower(subPkg.Name))))
	if err != nil {
		return fmt.Errorf("erro ao escrever o arquivo %s.go para o subpacote %s: %v", subPkg.Name, subPkg.Name, err)
	}

	return nil
}
