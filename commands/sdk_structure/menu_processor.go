package sdk_structure

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/magaluCloud/cligen/config"
)

func processMenu(ctx context.Context, menu *config.Menu, sdkStructure *SDKStructure) {
	processMenuRecursive(ctx, menu, "", sdkStructure)
}

func processMenuRecursive(ctx context.Context, menu *config.Menu, parentPath string, sdkStructure *SDKStructure) {
	if len(menu.Menus) > 0 && menu.Name == "profile" {
		groupPkg := Package{
			MenuName:        menu.Name,
			Name:            menu.Name,
			Description:     menu.Description,
			LongDescription: menu.LongDescription,
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
				subPkg := genCliCodeFromSDK(ctx, submenu)
				subPkg.MenuName = submenu.Name
				subPkg.Description = submenu.Description
				subPkg.LongDescription = submenu.LongDescription
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
						subSubPkg := genCliCodeFromSDK(ctx, subSubmenu)
						subSubPkg.MenuName = subSubmenu.Name
						subSubPkg.Description = subSubmenu.Description
						subSubPkg.LongDescription = subSubmenu.LongDescription
						subGroupPkg.SubPkgs[subSubmenu.SDKPackage] = subSubPkg
					} else if len(subSubmenu.Menus) > 0 {
						processMenuRecursive(ctx, subSubmenu, filepath.Join(currentPath, submenu.Name), sdkStructure)
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
		pkg := genCliCodeFromSDK(ctx, menu)
		pkg.MenuName = menu.Name
		pkg.GroupID = menu.CliGroup
		pkg.Description = menu.Description
		pkg.LongDescription = menu.LongDescription

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
