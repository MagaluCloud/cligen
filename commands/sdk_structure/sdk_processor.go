package sdk_structure

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/magaluCloud/cligen/config"
)

func genCliCodeFromSDK(ctx context.Context, menu *config.Menu) Package {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Erro ao obter diretório atual: %v", err)
	}

	sdkDir := filepath.Join(dir, "tmp-sdk", menu.SDKPackage)

	pkg := Package{
		MenuName:        menu.SDKPackage,
		Name:            menu.SDKPackage,
		Description:     menu.Description,
		LongDescription: menu.LongDescription,
		Aliases:         menu.Alias,
		Services:        []Service{},
		SubPkgs:         make(map[string]Package),
	}

	if _, err := os.Stat(sdkDir); os.IsNotExist(err) {
		fmt.Printf("⚠️  Diretório do SDK não encontrado: %s (menu de agrupamento)\n", sdkDir)
		return pkg
	}

	services := genCliCodeFromClient(ctx, menu, &pkg, sdkDir, filepath.Join(sdkDir, "client.go"))
	pkg.Services = services

	return pkg
}
