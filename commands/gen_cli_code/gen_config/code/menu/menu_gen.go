package menu

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/magaluCloud/cligen/config"
	strutils "github.com/magaluCloud/cligen/str_utils"
)

//go:embed menu.template
var menuTemplate string

var menuTmpl *template.Template

const (
	genDir      = "base-cli-gen/cmd/gen"
	importCobra = "\"github.com/spf13/cobra\""
)

func init() {
	var err error
	menuTmpl, err = template.New("module").Parse(menuTemplate)
	if err != nil {
		panic(err)
	}
}

func GenerateMenu(cfg *config.Config) error {
	for _, menu := range cfg.Menus {
		for _, submenu := range menu.Menus {
			GenMenu(cfg, submenu)
		}
	}
	return nil
}

func GenMenu(cfg *config.Config, submenu *config.Menu) error {
	menuData := NewMenu()
	parents := FindParents(cfg, submenu.ParentMenuID)

	parents = append(parents, strings.ToLower(submenu.Name))
	parentsPath := strings.Join(parents, "/")
	menuData.SetPathSaveToFile(filepath.Join(genDir, parentsPath, fmt.Sprintf("%s.go", strings.ToLower(submenu.Name))))

	nameFromParent, sdkPackageFromParent := FindSDKPackageFromParents(cfg, submenu.ParentMenuID)
	menuData.SetServiceParam(fmt.Sprintf("%s %s.%s", strutils.FirstLower(submenu.ServiceInterface), nameFromParent, submenu.ServiceInterface))
	menuData.AddImport(fmt.Sprintf("\"%s\"", sdkPackageFromParent))

	menuData.SetPackageName(submenu.Name)
	menuData.SetFunctionName(strutils.FirstUpper(submenu.Name))
	menuData.SetUseName(submenu.Name)
	menuData.SetAliases(submenu.Alias...)
	menuData.AddImport(importCobra)
	menuData.SetShortDescription(submenu.Description)
	menuData.SetLongDescription(submenu.LongDescription)
	menuData.SetGroupID(submenu.CliGroup)

	for _, method := range submenu.Methods {
		menuData.AddCommand(CommandType{
			FunctionName: strutils.FirstUpper(method.Name),
			ServiceCall:  strutils.FirstLower(submenu.ServiceInterface),
		})
	}
	for _, ssmenu := range submenu.Menus {
		menuData.AddServiceInit(fmt.Sprintf("%s.%sCmd(ctx, cmd, %s)", strings.ToLower(ssmenu.Name), ssmenu.Name, strutils.FirstLower(submenu.ServiceInterface)))
		menuData.AddImport(fmt.Sprintf("\"github.com/magaluCloud/mgccli/cmd/gen/%s/%s\"", parentsPath, strings.ToLower(ssmenu.Name)))
		GenMenu(cfg, ssmenu)
	}
	menuData.Save()

	return nil
}

func FindParents(cfg *config.Config, menuID string) []string {
	parents := []string{}
	menu := FindMenuByID(cfg.Menus, menuID)
	if menu != nil {
		parents = append(parents, strings.ToLower(menu.Name))
		if menu.ParentMenuID != "" {
			parents = append(FindParents(cfg, menu.ParentMenuID), parents...)
		}
	}
	return parents
}

func FindMenuByID(menus []*config.Menu, id string) *config.Menu {
	for _, menu := range menus {
		if menu.ID == id {
			return menu
		}
		if len(menu.Menus) > 0 {
			menu := FindMenuByID(menu.Menus, id)
			if menu != nil {
				return menu
			}
		}
	}
	return nil
}

func FindSDKPackageFromParents(cfg *config.Config, menuID string) (name string, sdkPackage string) {
	menu := FindMenuByID(cfg.Menus, menuID)
	if menu == nil {
		return "", ""
	}
	if menu.SDKPackage != "" {
		return menu.Name, menu.SDKPackage
	}
	if menu.ParentMenuID != "" {
		return FindSDKPackageFromParents(cfg, menu.ParentMenuID)
	}
	return "", ""
}
