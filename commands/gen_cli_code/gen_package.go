package gen_cli_code

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/magaluCloud/cligen/config"
	strutils "github.com/magaluCloud/cligen/str_utils"
)

func genPackageCode(cfg *config.Config) error {
	for _, pkg := range cfg.Menus {
		genPackageCodeRecursive(pkg, nil)
	}
	return nil
}

func genPackageCodeRecursive(pkg *config.Menu, parentPkg *config.Menu) error {
	if len(pkg.Methods) == 0 {
		if len(pkg.Menus) > 0 {
			for _, subPkg := range pkg.Menus {
				genPackageCodeRecursive(subPkg, pkg)
			}
		}
		return nil
	}
	if len(pkg.Menus) > 0 {
		return nil
	}
	packageData := NewPackageGroupData()
	packageData.SetGroupID(pkg.CliGroup)
	packageData.SetPackageName(pkg.Name)
	packageData.SetFunctionName(strutils.FirstUpper(pkg.Name))
	packageData.SetUseName(pkg.Name)
	packageData.SetAliases(pkg.Alias)
	packageData.AddImport(importCobra)
	packageData.AddImport(importSDK)
	packageData.SetDescriptions(pkg.Description, pkg.LongDescription)
	if parentPkg != nil {
		packageData.AddImport(fmt.Sprintf("\"github.com/magaluCloud/mgccli/cmd/gen/%s\"", strings.ToLower(parentPkg.Name)))
	} else {
		packageData.AddImport(fmt.Sprintf("\"github.com/magaluCloud/mgccli/cmd/gen/%s\"", strings.ToLower(pkg.Name)))
	}
	packageData.SetServiceParam(serviceParamPattern)

	for _, method := range pkg.Methods {
		packageData.AddSubCommand(method.Name, method.Name, fmt.Sprintf("%sService.%s()", pkg.Name, method.Name))
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
