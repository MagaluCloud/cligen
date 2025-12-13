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
			dir := genDir
			if parentPkg != nil {
				dir = filepath.Join(dir, strings.ToLower(parentPkg.Name), strings.ToLower(pkg.Name), strings.ToLower(service.Name))
			} else {
				dir = filepath.Join(dir, strings.ToLower(pkg.Name), strings.ToLower(service.Name))
			}
			err := genService(&service, pkg, parentPkg, dir)
			if err != nil {
				return err
			}
			for _, subService := range service.SubServices {
				dir := filepath.Join(dir, strings.ToLower(subService.Name))
				err := genService(&subService, pkg, parentPkg, dir)
				if err != nil {
					return err
				}
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

func genService(service *sdk_structure.Service, pkg *sdk_structure.Package, parentPkg *sdk_structure.Package, dir string) error {
	serviceData := NewPackageGroupData()
	serviceData.SetPackageName(service.Name)
	serviceData.SetFunctionName(service.Name)
	serviceData.SetUseName(service.Name)
	serviceData.AddImport(importCobra)
	serviceData.SetServiceParam(fmt.Sprintf("%s %sSdk.%s", strutils.FirstLower(service.Interface), pkg.Name, service.Interface))
	for _, method := range service.Methods {
		serviceData.SetDescriptions(service.Description, service.LongDescription)
		serviceData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", pkg.Name, pkg.Name))
		if !method.IsService {
			serviceData.AddCommand(method.Name, strutils.FirstLower(service.Interface))
		}
		if method.IsService {
			serviceData.AddImport(fmt.Sprintf("\"github.com/magaluCloud/mgccli/cmd/gen/%s/%s/%s\"", pkg.Name, strings.ToLower(service.Name), strings.ToLower(method.ServiceImport)))
			serviceData.AddCommand(fmt.Sprintf("%s.%sCmd", strings.ToLower(method.ServiceImport), method.ServiceImport), fmt.Sprintf("%s.%s()", strutils.FirstLower(service.Interface), method.Name))

		}
	}

	err := serviceData.WriteServiceToFile(filepath.Join(dir, fmt.Sprintf("%s.go", strings.ToLower(service.Name))))
	if err != nil {
		return fmt.Errorf("erro ao escrever o arquivo %s.go para o serviço %s do pacote %s: %v", pkg.Name, pkg.Name, pkg.Name, err)
	}
	return nil
}
