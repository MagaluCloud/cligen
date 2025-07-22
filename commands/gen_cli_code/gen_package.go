package gen_cli_code

import (
	"cligen/commands/sdk_structure"
	strutils "cligen/str_utils"
	"fmt"
	"path/filepath"
	"strings"
)

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
	packageData.SetGroupID(groupProducts)
	packageData.SetPackageName(pkg.Name)
	packageData.SetFunctionName(strutils.FirstUpper(pkg.Name))
	packageData.SetUseName(pkg.MenuName)
	packageData.SetAliases(pkg.Aliases)
	packageData.AddImport(importCobra)
	packageData.AddImport(importSDK)
	packageData.SetDescriptions(pkg.Description, pkg.LongDescription)
	packageData.AddImport(fmt.Sprintf("%sSdk \"github.com/MagaluCloud/mgc-sdk-go/%s\"", pkg.Name, pkg.Name))
	packageData.SetServiceParam(serviceParamPattern)

	packageData.AddServiceInit(fmt.Sprintf("%sService := %sSdk.New(&sdkCoreConfig)", pkg.Name, pkg.Name))

	for _, service := range pkg.Services {
		packageData.AddImport(fmt.Sprintf("\"gfcli/cmd/gen/%s/%s\"", strings.ToLower(pkg.Name), strings.ToLower(service.Name)))
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
