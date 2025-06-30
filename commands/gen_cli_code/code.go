package gen_cli_code

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"cligen/commands/sdk_structure"
	strutils "cligen/str_utils"
)

const (
	genDir           = "base-cli-gen/cmd/gen"
	defaultShortDesc = "todo"
	defaultLongDesc  = "todo2"
)

func GenCliCode() {
	sdkStructure, err := sdk_structure.GenCliSDKStructure()
	if err != nil {
		log.Fatalf("Erro ao gerar a estrutura do SDK: %v", err)
	}
	cleanDir(genDir)

	rootGenData := NewRootGenData()
	rootGenData.AddImport("sdk \"github.com/MagaluCloud/mgc-sdk-go/client\"")
	rootGenData.AddImport("\"github.com/spf13/cobra\"")

	for _, pkg := range sdkStructure.Packages {
		genPackageCode(&pkg, rootGenData)
	}

	rootGenData.WriteRootGenToFile(filepath.Join(genDir, "root_gen.go"))
}

func cleanDir(dir string) {
	toRemove := filepath.Clean(dir)
	if _, err := os.Stat(toRemove); err == nil {
		os.RemoveAll(toRemove)
	}
}

func genPackageCode(pkg *sdk_structure.Package, rootGenData *RootGenData) {
	packageData := NewPackageGroupData()
	packageData.SetPackageName(pkg.Name)
	packageData.SetFunctionName(strutils.FirstUpper(pkg.Name))
	packageData.SetUseName(pkg.Name)
	packageData.SetDescriptions(defaultShortDesc, defaultLongDesc)
	packageData.SetGroupID("products")
	packageData.SetServiceParam("sdkCoreConfig *sdk.CoreClient")

	rootGenData.AddSubCommand(pkg.Name, strutils.FirstUpper(pkg.Name)+"Cmd")
	rootGenData.AddImport(fmt.Sprintf("\"mgccli/cmd/gen/%s\"", strings.ToLower(pkg.Name)))

	for _, service := range pkg.Services {
		packageData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", pkg.Name, pkg.Name))
		packageData.AddImport("\"github.com/spf13/cobra\"")
		packageData.AddImport(fmt.Sprintf("\"mgccli/cmd/gen/%s/%s\"", strings.ToLower(pkg.Name), strings.ToLower(service.Name)))
		packageData.SetServiceInit(fmt.Sprintf("%sService := %sSdk.New(sdkCoreConfig)", pkg.Name, pkg.Name))
		packageData.AddSubCommand(service.Name, service.Name, fmt.Sprintf("%sService.%s()", pkg.Name, service.Name))

		generateServiceCode(*pkg, &service, *packageData)
	}
	packageData.AddImport("sdk \"github.com/MagaluCloud/mgc-sdk-go/client\"")
	err := packageData.WriteGroupToFile(filepath.Join(genDir, strings.ToLower(pkg.Name), fmt.Sprintf("%s.go", pkg.Name)))
	if err != nil {
		log.Fatalf("Erro ao escrever o arquivo %s: %v", pkg.Name, err)
	}
}

func generateServiceCode(parentPkg sdk_structure.Package, service *sdk_structure.Service, data PackageGroupData) {
	dir := filepath.Join(genDir, strings.ToLower(parentPkg.Name), strings.ToLower(service.Name))
	os.MkdirAll(dir, 0755)

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
	serviceData.AddImport("\"github.com/spf13/cobra\"")
	serviceData.SetDescriptions(defaultShortDesc, defaultLongDesc)
	serviceData.SetGroupID("")
	serviceData.SetServiceParam(fmt.Sprintf("%s %sSdk.%s", strutils.FirstLower(service.Interface), parentPkg.Name, service.Interface))

	for _, method := range service.Methods {
		productData := serviceData.Copy()
		productData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", parentPkg.Name, parentPkg.Name))
		productData.AddImport("\"github.com/spf13/cobra\"")
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
			log.Fatalf("Erro ao escrever o arquivo %s: %v", strings.ToLower(method.Name), err)
		}
	}

	err := serviceData.WriteServiceToFile(filepath.Join(dir, fmt.Sprintf("%s.go", strings.ToLower(service.Name))))
	if err != nil {
		log.Fatalf("Erro ao escrever o arquivo %s: %v", service.Name, err)
	}

	for _, subService := range service.SubServices {
		generateServiceCode(parentPkg, &subService, data)
	}
}
