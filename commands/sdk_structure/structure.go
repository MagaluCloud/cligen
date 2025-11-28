package sdk_structure

import (
	"github.com/magaluCloud/cligen/config"
)

func GenCliSDKStructure(config *config.Config) (SDKStructure, error) {
	sdkStructure := &SDKStructure{
		Packages: make(map[string]Package),
	}

	for _, menu := range config.Menus {
		if menu.Enabled {
			processMenu(menu, sdkStructure)
		}
	}

	return *sdkStructure, nil
}
