package sdk_structure

import (
	"context"

	"github.com/magaluCloud/cligen/config"
)

func GenCliSDKStructure(ctx context.Context, config *config.Config) (SDKStructure, error) {
	sdkStructure := &SDKStructure{
		Packages: make(map[string]Package),
	}

	for _, menu := range config.Menus {
		if menu.Enabled {
			processMenu(ctx, menu, sdkStructure)
		}
	}

	return *sdkStructure, nil
}
