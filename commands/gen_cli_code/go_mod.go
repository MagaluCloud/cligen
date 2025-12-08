package gen_cli_code

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"text/template"

	"github.com/magaluCloud/cligen/config"

	_ "embed"
)

//go:embed gomod.template
var goModTemplate string

var goModTmpl = template.Must(template.New("goMod").Parse(goModTemplate))

type GoModData struct {
	Version string
}

// WriteToFile escreve os dados no arquivo
func (gmd *GoModData) WriteGoModFile(filePath string) error {

	buf := bytes.NewBuffer(nil)
	err := goModTmpl.Execute(buf, gmd)
	if err != nil {
		return err
	}

	os.MkdirAll(filepath.Dir(filePath), 0755)
	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

func genGoModFile() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Erro ao carregar o arquivo de configuração: %v", err)
	}

	gmd := GoModData{
		Version: config.SDKTag,
	}

	gmd.WriteGoModFile(filepath.Join("base-cli-gen", "go.mod"))
}
