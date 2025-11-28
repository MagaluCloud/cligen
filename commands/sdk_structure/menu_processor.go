package sdk_structure

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/magaluCloud/cligen/config"
	strutils "github.com/magaluCloud/cligen/str_utils"
)

func processMenu(menu *config.Menu, sdkStructure *SDKStructure) {
	processMenuRecursive(menu, "", sdkStructure)
}

func processMenuRecursive(menu *config.Menu, parentPath string, sdkStructure *SDKStructure) {
	if len(menu.Menus) > 0 && menu.Level == 0 {
		groupPkg := Package{
			MenuName:        menu.Name,
			Name:            menu.Name,
			Description:     menu.Description,
			LongDescription: menu.Description,
			Aliases:         menu.Alias,
			GroupID:         menu.CliGroup,
			Services:        []Service{},
			SubPkgs:         make(map[string]Package),
		}

		currentPath := menu.Name
		if parentPath != "" {
			currentPath = filepath.Join(parentPath, menu.Name)
		}

		for _, submenu := range menu.Menus {
			if submenu.SDKPackage != "" {
				subPkg := genCliCodeFromSDK(submenu)
				subPkg.MenuName = submenu.Name
				groupPkg.SubPkgs[submenu.SDKPackage] = subPkg
			} else if len(submenu.Menus) > 0 {
				subGroupPkg := Package{
					MenuName:        submenu.Name,
					Name:            submenu.Name,
					Aliases:         submenu.Alias,
					Description:     submenu.Description,
					LongDescription: submenu.LongDescription,
					Services:        []Service{},
					GroupID:         menu.CliGroup,
					SubPkgs:         make(map[string]Package),
				}

				for _, subSubmenu := range submenu.Menus {

					if subSubmenu.SDKPackage != "" {
						subSubPkg := genCliCodeFromSDK(subSubmenu)
						subSubPkg.MenuName = subSubmenu.Name
						subGroupPkg.SubPkgs[subSubmenu.SDKPackage] = subSubPkg
					} else if len(subSubmenu.Menus) > 0 {
						processMenuRecursive(subSubmenu, filepath.Join(currentPath, submenu.Name), sdkStructure)
					}
				}

				groupPkg.SubPkgs[submenu.Name] = subGroupPkg
			}
		}

		if parentPath == "" {
			sdkStructure.Packages[menu.Name] = groupPkg
		} else {
			packageKey := filepath.Join(parentPath, menu.Name)
			sdkStructure.Packages[packageKey] = groupPkg
		}
	} else if menu.SDKPackage != "" {
		pkg := genCliCodeFromSDK(menu)
		pkg.MenuName = menu.Name
		pkg.GroupID = menu.CliGroup
		if pkg.Description == "" && pkg.LongDescription != "" {
			strs := strings.Split(pkg.LongDescription, "\n")
			for _, str := range strs {
				if str != "" {
					str = strings.Replace(str, "Package ", "", 1)
					str = strutils.FirstUpper(str)
					pkg.Description = str
					break
				}
			}
		}
		if parentPath == "" {
			sdkStructure.Packages[menu.SDKPackage] = pkg
		} else {
			packageKey := filepath.Join(parentPath, menu.SDKPackage)
			sdkStructure.Packages[packageKey] = pkg
		}
	} else {
		fmt.Printf("⚠️  Menu '%s' não tem submenus nem SDK Package (menu vazio)\n", menu.Name)
	}
}
