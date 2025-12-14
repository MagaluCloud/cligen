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
	module.AddImport(fmt.Sprintf("\"%s\"", menu.SDKPackage))
	module.SetPathSaveToFile(filepath.Join(genDir, strings.ToLower(menu.Name), fmt.Sprintf("%s.go", strings.ToLower(menu.Name))))

	moduleServiceName := fmt.Sprintf("%sService", strings.ToLower(menu.Name))
	module.AddServiceInit(fmt.Sprintf("%s := %s.New(&sdkCoreConfig)", moduleServiceName, strings.ToLower(menu.Name)))

	for _, submenu := range menu.Menus {
		module.AddSubCommand(subCommandType{
			PackageName:  strings.ToLower(submenu.Name),
			FunctionName: strutils.FirstUpper(submenu.Name) + "Cmd",
			ServiceCall:  fmt.Sprintf("%s.%s()", moduleServiceName, submenu.Name),
		})
		module.AddImport(fmt.Sprintf("\"github.com/magaluCloud/mgccli/cmd/gen/%s/%s\"", strings.ToLower(menu.Name), strings.ToLower(submenu.Name)))
	}
	module.SetPackageName(strings.ToLower(menu.Name))
	module.SetFunctionName(strutils.FirstUpper(menu.Name) + "Cmd")
	module.SetUseName(menu.Name)
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
	module.SetPathSaveToFile(filepath.Join(genDir, strings.ToLower(menu.Name), fmt.Sprintf("%s.go", strings.ToLower(menu.Name))))

	for _, submenu := range menu.Menus {
		module.AddImport(fmt.Sprintf("\"%s\"", submenu.SDKPackage))

		moduleServiceName := fmt.Sprintf("%sService", strings.ToLower(submenu.Name))
		module.AddServiceInit(fmt.Sprintf("%s := %s.New(&sdkCoreConfig)", moduleServiceName, strings.ToLower(submenu.Name)))

		for _, ssmenu := range submenu.Menus {
			module.AddSubCommand(subCommandType{
				PackageName:  strings.ToLower(ssmenu.Name),
				FunctionName: strutils.FirstUpper(ssmenu.Name) + "Cmd",
				ServiceCall:  fmt.Sprintf("%s.%s()", moduleServiceName, ssmenu.Name),
			})
			module.AddImport(fmt.Sprintf("\"github.com/magaluCloud/mgccli/cmd/gen/%s/%s\"", strings.ToLower(menu.Name), strings.ToLower(ssmenu.Name)))
		}
	}
	module.SetPackageName(strings.ToLower(menu.Name))
	module.SetFunctionName(strutils.FirstUpper(menu.Name) + "Cmd")
	module.SetUseName(menu.Name)
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
