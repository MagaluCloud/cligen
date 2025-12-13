package gen_cli_code

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/magaluCloud/cligen/config"
	strutils "github.com/magaluCloud/cligen/str_utils"
)

func generateRootCode(sdkStructure *config.Config) error {
	rootGenData := NewRootGenData()
	rootGenData.AddImport(importSDK)
	rootGenData.AddImport(importCobra)
	for _, pkg := range sdkStructure.Menus {
		rootGenData.AddSubCommand(pkg.Name, strutils.FirstUpper(pkg.Name)+"Cmd")
		rootGenData.AddImport(fmt.Sprintf("\"github.com/magaluCloud/mgccli/cmd/gen/%s\"", strings.ToLower(pkg.Name)))
	}
	if err := rootGenData.WriteRootGenToFile(filepath.Join(genDir, "root_gen.go")); err != nil {
		log.Fatalf("Erro ao escrever o arquivo root_gen.go: %v", err)
	}
	return nil
}
