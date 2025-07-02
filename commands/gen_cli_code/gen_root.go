package gen_cli_code

import (
	"cligen/commands/sdk_structure"
	strutils "cligen/str_utils"
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

func generateRootCode(sdkStructure *sdk_structure.SDKStructure) error {
	rootGenData := NewRootGenData()
	rootGenData.AddImport(importSDK)
	rootGenData.AddImport(importCobra)
	for _, pkg := range sdkStructure.Packages {
		rootGenData.AddSubCommand(pkg.Name, strutils.FirstUpper(pkg.Name)+"Cmd")
		rootGenData.AddImport(fmt.Sprintf("\"mgccli/cmd/gen/%s\"", strings.ToLower(pkg.Name)))
	}
	if err := rootGenData.WriteRootGenToFile(filepath.Join(genDir, "root_gen.go")); err != nil {
		log.Fatalf("Erro ao escrever o arquivo root_gen.go: %v", err)
	}
	return nil
}
