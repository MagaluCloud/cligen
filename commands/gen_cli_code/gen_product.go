package gen_cli_code

import (
	"cligen/commands/sdk_structure"
	strutils "cligen/str_utils"
	"fmt"
	"path/filepath"
	"strings"
)

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
				productData.SetServiceCall(fmt.Sprintf("%s.%s", strutils.FirstLower(service.Interface), method.Name))

				var params []string
				for key, param := range method.Parameters {
					productData.AddServiceSDKParamCreate(fmt.Sprintf("var %s %s", key, param))
					params = append(params, key)
				}
				productData.SetServiceSDKParam(strings.Join(params, ", "))

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
