package sdk_structure

import (
	"log"

	"github.com/magaluCloud/cligen/config"
)

func GenCliSDKStructure() (SDKStructure, error) {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	sdkStructure := &SDKStructure{
		Packages: make(map[string]Package),
	}

	// Processar menus principais e seus submenus
	for _, menu := range config.Menus {
		processMenu(menu, sdkStructure)
	}

	return *sdkStructure, nil
}
