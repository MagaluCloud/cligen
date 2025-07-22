package sdk_structure

import (
	"cligen/config"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// genCliCodeFromSDK processa um menu e gera código CLI baseado no SDK
func genCliCodeFromSDK(menu config.Menu) Package {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Erro ao obter diretório atual: %v", err)
	}

	// Construir o caminho do SDK baseado na hierarquia
	sdkDir := filepath.Join(dir, "tmp-sdk", menu.SDKPackage)
	fmt.Printf("🔍 Procurando SDK em diretório principal: %s\n", sdkDir)

	pkg := Package{
		MenuName:        menu.SDKPackage,
		Name:            menu.SDKPackage,
		Description:     menu.Description,
		LongDescription: "menu.LongDescription 4",
		Aliases:         menu.Alias,
		Services:        []Service{},
		SubPkgs:         make(map[string]Package),
	}

	// Verificar se o diretório do SDK existe
	if _, err := os.Stat(sdkDir); os.IsNotExist(err) {
		// Se o diretório não existe, retornar um pacote vazio (para menus de agrupamento)
		fmt.Printf("⚠️  Diretório do SDK não encontrado: %s (menu de agrupamento)\n", sdkDir)
		return pkg
	}

	fmt.Printf("✅ Diretório do SDK encontrado: %s\n", sdkDir)

	files, err := os.ReadDir(sdkDir)
	if err != nil {
		log.Fatalf("Erro ao ler diretório do SDK: %v", err)
	}

	fmt.Printf("📄 Total de arquivos no diretório: %d\n", len(files))

	for _, file := range files {
		if file.Name() == "client.go" {
			fmt.Printf("🔧 Processando arquivo client.go em: %s\n", sdkDir)
			services := genCliCodeFromClient(&pkg, sdkDir, filepath.Join(sdkDir, file.Name()))
			pkg.Services = services
			fmt.Printf("✅ Processados %d serviços do pacote %s\n", len(services), menu.SDKPackage)
		}
	}

	return pkg
}
