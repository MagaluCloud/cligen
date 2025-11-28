package gen_cli_code

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/magaluCloud/cligen/commands/sdk_structure"
	strutils "github.com/magaluCloud/cligen/str_utils"
)

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
			serviceData.SetUseName(service.Name)
			serviceData.AddImport(importCobra)
			serviceData.SetServiceParam(fmt.Sprintf("%s %sSdk.%s", strutils.FirstLower(service.Interface), pkg.Name, service.Interface))
			for _, method := range service.Methods {
				serviceData.SetDescriptions(service.Description, service.LongDescription)
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
