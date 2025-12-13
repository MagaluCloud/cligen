package code

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/magaluCloud/cligen/config"
	strutils "github.com/magaluCloud/cligen/str_utils"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("erro ao carregar configuração: %w", err))
	}

	for _, pkg := range cfg.Menus {
		rootGen := NewRootGen()
		rootGen.SetPathSaveToFile(filepath.Join("./base-cli-gen", pkg.Name, "root_gen.go"))

		if len(pkg.Menus) > 0 {
			for _, subPkg := range pkg.Menus {
				rootGen.AddSubCommand(SubCommandType{
					PackageName:  subPkg.Name,
					FunctionName: strutils.FirstUpper(subPkg.Name) + "Cmd",
				})
			}
		}

		rootGen.AddImport(fmt.Sprintf("\"github.com/magaluCloud/mgccli/cmd/gen/%s\"", strings.ToLower(pkg.Name)))
		rootGen.Save()
	}

	fmt.Println(cfg.SDKTag)
}
