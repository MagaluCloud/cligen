package gen_cli_code

import (
	"os"
	"path/filepath"

	"github.com/magaluCloud/cligen/commands/gen_cli_code/gen_config/code"
	"github.com/magaluCloud/cligen/commands/gen_cli_code/gen_config/generate"
	"github.com/magaluCloud/cligen/commands/gen_cli_code/gen_config/manipulate"
)

const (
	genDir              = "base-cli-gen/cmd/gen"
	importCobra         = "\"github.com/spf13/cobra\""
	importSDK           = "sdk \"github.com/MagaluCloud/mgc-sdk-go/client\""
	serviceParamPattern = "sdkCoreConfig sdk.CoreClient"
)

func GenConfig() {
	generate.Run()
}

func Manipulate() {
	manipulate.StartServer("8080")
}

func GenCliCode() {

	cleanDir(genDir)
	code.Run()

}

func cleanDir(dir string) {
	toRemove := filepath.Clean(dir)
	if _, err := os.Stat(toRemove); err == nil {
		os.RemoveAll(toRemove)
	}
}
