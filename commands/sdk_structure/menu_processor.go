package sdk_structure

import (
	"cligen/config"
	"fmt"
	"path/filepath"
	"strings"
)

// processMenu processa um menu e seus submenus recursivamente
func processMenu(menu config.Menu, sdkStructure *SDKStructure) {
	processMenuRecursive(menu, "", sdkStructure)
}

// processMenuRecursive processa um menu e seus submenus recursivamente com suporte a hierarquia
func processMenuRecursive(menu config.Menu, parentPath string, sdkStructure *SDKStructure) {
	fmt.Printf("🔄 Processando menu: %s (caminho pai: %s)\n", menu.Name, parentPath)

	// Se o menu tem submenus, criar um pacote de agrupamento
	if len(menu.Menus) > 0 {
		fmt.Printf("📁 Menu '%s' é um agrupador com %d submenus\n", menu.Name, len(menu.Menus))

		// Criar um pacote vazio para o menu de agrupamento
		groupPkg := Package{
			MenuName:        menu.Name,
			Name:            menu.Name,
			Description:     menu.Description,
			LongDescription: "menu.LongDescription 1",
			Aliases:         menu.Alias,
			Services:        []Service{},
			SubPkgs:         make(map[string]Package),
		}

		// Construir o caminho atual para este menu
		currentPath := menu.Name
		if parentPath != "" {
			currentPath = filepath.Join(parentPath, menu.Name)
		}

		fmt.Printf("📍 Caminho atual para menu '%s': %s\n", menu.Name, currentPath)

		// Adicionar subpacotes para cada submenu
		for _, submenu := range menu.Menus {
			fmt.Printf("  🔍 Processando submenu: %s\n", submenu.Name)

			if submenu.SDKPackage != "" {
				fmt.Printf("  📦 Submenu '%s' tem SDK Package: %s\n", submenu.Name, submenu.SDKPackage)
				// Para menus filhos, o diretório será dentro do diretório pai
				subPkg := genCliCodeFromSDK(submenu)
				subPkg.MenuName = submenu.Name
				groupPkg.SubPkgs[submenu.SDKPackage] = subPkg
			} else if len(submenu.Menus) > 0 {
				fmt.Printf("  📁 Submenu '%s' é um agrupador com %d sub-submenus\n", submenu.Name, len(submenu.Menus))
				// Se o submenu também tem submenus, processar recursivamente
				// Criar um subpacote de agrupamento
				subGroupPkg := Package{
					MenuName:        submenu.Name,
					Name:            submenu.Name,
					Aliases:         submenu.Alias,
					Description:     "submenu.Description",
					LongDescription: "submenu.LongDescription 2",
					Services:        []Service{},
					SubPkgs:         make(map[string]Package),
				}

				// Processar submenus do submenu
				for _, subSubmenu := range submenu.Menus {
					fmt.Printf("    🔍 Processando sub-submenu: %s\n", subSubmenu.Name)

					if subSubmenu.SDKPackage != "" {
						fmt.Printf("    📦 Sub-submenu '%s' tem SDK Package: %s\n", subSubmenu.Name, subSubmenu.SDKPackage)
						// Para sub-submenus, o diretório será dentro do diretório do submenu pai
						subSubPkg := genCliCodeFromSDK(subSubmenu)
						subSubPkg.MenuName = subSubmenu.Name
						subGroupPkg.SubPkgs[subSubmenu.SDKPackage] = subSubPkg
					} else if len(subSubmenu.Menus) > 0 {
						fmt.Printf("    📁 Sub-submenu '%s' é um agrupador com %d sub-sub-submenus\n", subSubmenu.Name, len(subSubmenu.Menus))
						// Recursão para níveis mais profundos
						processMenuRecursive(subSubmenu, filepath.Join(currentPath, submenu.Name), sdkStructure)
					}
				}

				groupPkg.SubPkgs[submenu.Name] = subGroupPkg
			}
		}

		// Adicionar o pacote ao nível apropriado
		if parentPath == "" {
			// Menu principal - adicionar diretamente ao SDKStructure
			fmt.Printf("✅ Adicionando menu principal '%s' ao SDKStructure\n", menu.Name)
			sdkStructure.Packages[menu.Name] = groupPkg
		} else {
			// Submenu - adicionar ao pacote pai
			// Nota: Aqui precisamos adicionar ao pacote pai correto
			// Por enquanto, vamos adicionar diretamente ao SDKStructure com um nome único
			packageKey := filepath.Join(parentPath, menu.Name)
			fmt.Printf("✅ Adicionando submenu '%s' ao SDKStructure com chave: %s\n", menu.Name, packageKey)
			sdkStructure.Packages[packageKey] = groupPkg
		}
	} else if menu.SDKPackage != "" {
		fmt.Printf("📦 Menu '%s' tem SDK Package: %s\n", menu.Name, menu.SDKPackage)
		pkg := genCliCodeFromSDK(menu)
		pkg.MenuName = menu.Name
		if pkg.Description == "" {
			description := []string{}
			checkpoint := false
			for _, service := range pkg.Services {
				description = append(description, service.Name)
				if len(description) == 6 || len(strings.Join(description, ", ")) > 70 {
					checkpoint = true
					break
				}
			}
			pkg.LongDescription = "menu.LongDescription 3"
			if checkpoint {
				pkg.Description = strings.Join(description, ", ") + "..."
			} else {
				pkg.Description = strings.Join(description, ", ") + "."
			}
		}

		// Adicionar ao nível apropriado
		if parentPath == "" {
			// Menu principal
			fmt.Printf("✅ Adicionando menu principal com SDK '%s' ao SDKStructure\n", menu.SDKPackage)
			sdkStructure.Packages[menu.SDKPackage] = pkg
		} else {
			// Submenu - adicionar com nome único
			packageKey := filepath.Join(parentPath, menu.SDKPackage)
			fmt.Printf("✅ Adicionando submenu com SDK '%s' ao SDKStructure com chave: %s\n", menu.SDKPackage, packageKey)
			sdkStructure.Packages[packageKey] = pkg
		}
	} else {
		fmt.Printf("⚠️  Menu '%s' não tem submenus nem SDK Package (menu vazio)\n", menu.Name)
	}
}
