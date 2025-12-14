package root_gen

import (
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/magaluCloud/cligen/config"
	strutils "github.com/magaluCloud/cligen/str_utils"
)

const genDir = "base-cli-gen/cmd/gen"

//go:embed root_gen.template
var rootGenTemplate string

var rootGenTmpl *template.Template

func init() {
	var err error
	rootGenTmpl, err = template.New("rootgen").Parse(rootGenTemplate)
	if err != nil {
		panic(err)
	}
}

func GenerateRootGen(cfg *config.Config) error {
	rootGen := NewRootGen()
	rootGen.SetPathSaveToFile(filepath.Join(genDir, "root_gen.go"))

	if len(cfg.Menus) > 0 {
		for _, menu := range cfg.Menus {
			rootGen.AddSubCommand(SubCommandType{
				PackageName:  menu.Name,
				FunctionName: strutils.FirstUpper(menu.Name) + "Cmd",
			})
			rootGen.AddImport(fmt.Sprintf("\"github.com/magaluCloud/mgccli/cmd/gen/%s\"", strings.ToLower(menu.Name)))
		}
	}

	rootGen.Save()
	return nil
}
