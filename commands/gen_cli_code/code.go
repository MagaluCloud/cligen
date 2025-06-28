package gen_cli_code

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"cligen/commands/sdk_structure"
)

const (
	genDir = "base-cli-gen/cmd/gen"
)

func GenCliCode() {
	sdkStructure, err := sdk_structure.GenCliSDKStructure()
	if err != nil {
		log.Fatalf("Erro ao gerar a estrutura do SDK: %v", err)
	}
	cleanDir(genDir)
	for _, pkg := range sdkStructure.Packages {
		genPackageCode(&pkg)
	}
}

func cleanDir(dir string) {
	toRemove := filepath.Clean(dir)
	if _, err := os.Stat(toRemove); err == nil {
		os.RemoveAll(toRemove)
	}
}

func genPackageCode(pkg *sdk_structure.Package) {
	for _, service := range pkg.Services {
		generateServiceCode(*pkg, &service)
	}
}

func generateServiceCode(parentPkg sdk_structure.Package, service *sdk_structure.Service) {
	dir := filepath.Join(genDir, parentPkg.Name, service.Name)
	os.MkdirAll(dir, 0755)
	fmt.Println(service.Name)
	for _, subService := range service.SubServices {
		generateServiceCode(parentPkg, &subService)
	}
}
