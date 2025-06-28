package gen_cli_code

import (
	"fmt"
	"log"

	"cligen/commands/sdk_structure"
)

func GenCliCode() {
	//read the sdk structure
	sdkStructure, err := sdk_structure.GenCliSDKStructure()
	if err != nil {
		log.Fatalf("Erro ao gerar a estrutura do SDK: %v", err)
	}
	for _, pkg := range sdkStructure.Packages {
		genPackageCode(&pkg)
	}
}

func genPackageCode(pkg *sdk_structure.Package) {
	for _, service := range pkg.Services {
		fmt.Println(service.Name)
	}
}
