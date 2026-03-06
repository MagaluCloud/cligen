package module

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/magaluCloud/cligen/config"
	strutils "github.com/magaluCloud/cligen/str_utils"
)

//go:embed module.template
var moduleTemplate string

var moduleTmpl *template.Template

const (
	genDir              = "base-cli-gen/cmd/gen"
	importCobra         = "\"github.com/spf13/cobra\""
	importSDK           = "sdk \"github.com/MagaluCloud/mgc-sdk-go/client\""
	serviceParamPattern = "sdkCoreConfig sdk.CoreClient"
)

func init() {
	var err error
	moduleTmpl, err = template.New("module").Parse(moduleTemplate)
	if err != nil {
		panic(err)
	}
}

func GenerateModule(cfg *config.Config) error {
	for _, menu := range cfg.Menus {
		if menu.IsGroup {
			GenGroupModule(menu)
			continue
		}
		GenModule(menu)
	}

	return nil
}

func GenModule(menu *config.Menu) error {
	module := NewModule()
	moduleServiceName := fmt.Sprintf("%sService", strings.ToLower(menu.SDKName))

	sdkImportName := fmt.Sprintf("%sSdk", moduleServiceName)
	module.AddImport(fmt.Sprintf("%s \"%s\"", sdkImportName, menu.SDKPackage))
	module.SetPathSaveToFile(filepath.Join(genDir, strings.ToLower(menu.SDKName), fmt.Sprintf("%s.go", strings.ToLower(menu.SDKName))))

	module.AddServiceInit(fmt.Sprintf("%s := %s.New(&sdkCoreConfig)", moduleServiceName, sdkImportName))

	for _, submenu := range menu.Menus {
		module.AddSubCommand(subCommandType{
			PackageName:  strings.ToLower(submenu.SDKName),
			FunctionName: strutils.FirstUpper(submenu.SDKName) + "Cmd",
			ServiceCall:  fmt.Sprintf("%s.%s()", moduleServiceName, submenu.SDKName),
		})
		module.AddImport(fmt.Sprintf("\"github.com/magaluCloud/mgccli/cmd/gen/%s/%s\"", strings.ToLower(menu.SDKName), strings.ToLower(submenu.SDKName)))
	}
	module.SetPackageName(strings.ToLower(menu.SDKName))
	module.SetFunctionName(strutils.FirstUpper(menu.SDKName) + "Cmd")
	module.SetUseName(menu.CliName)
	module.SetShortDescription(menu.Description)
	module.SetLongDescription(menu.LongDescription)
	module.AddAliases(menu.Alias...)
	module.SetGroupID(menu.CliGroup)
	module.AddImport(importCobra)
	module.AddImport(importSDK)
	module.SetServiceParam(serviceParamPattern)
	module.Save()
	return nil
}

func GenGroupModule(menu *config.Menu) error {
	module := NewModule()
	module.SetPathSaveToFile(filepath.Join(genDir, strings.ToLower(menu.SDKName), fmt.Sprintf("%s.go", strings.ToLower(menu.SDKName))))

	for _, submenu := range menu.Menus {

		if len(submenu.Menus) == 1 {
			// vamos suprimir um nivel de menu
			ssmenu := submenu.Menus[0]
			moduleServiceName := fmt.Sprintf("%sService", strings.ToLower(submenu.SDKName))
			sdkImportName := fmt.Sprintf("%sSdk", moduleServiceName)
			module.AddImport(fmt.Sprintf("%s \"%s\"", sdkImportName, submenu.SDKPackage))

			module.AddServiceInit(fmt.Sprintf("%s := %s.New(&sdkCoreConfig)", moduleServiceName, sdkImportName))

			module.AddSubCommand(subCommandType{
				PackageName:  strings.ToLower(submenu.SDKName),
				FunctionName: strutils.FirstUpper(submenu.SDKName) + "Cmd",
				ServiceCall:  fmt.Sprintf("%s.%s()", moduleServiceName, ssmenu.SDKName),
			})
			module.AddImport(fmt.Sprintf("\"github.com/magaluCloud/mgccli/cmd/gen/%s/%s\"", strings.ToLower(menu.SDKName), strings.ToLower(submenu.SDKName)))

			continue
		}

	}
	module.SetPackageName(strings.ToLower(menu.SDKName))
	module.SetFunctionName(strutils.FirstUpper(menu.SDKName) + "Cmd")
	module.SetUseName(menu.CliName)
	module.SetShortDescription(menu.Description)
	module.SetLongDescription(menu.LongDescription)
	module.AddAliases(menu.Alias...)
	module.SetGroupID(menu.CliGroup)
	module.AddImport(importCobra)
	module.AddImport(importSDK)
	module.SetServiceParam(serviceParamPattern)
	module.Save()
	return nil
}
