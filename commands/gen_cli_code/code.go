package gen_cli_code

import (
	"os"
	"path/filepath"

	"github.com/magaluCloud/cligen/commands/gen_cli_code/code"
	"github.com/magaluCloud/cligen/commands/gen_cli_code/generate_config"
	"github.com/magaluCloud/cligen/commands/gen_cli_code/manipulate_config"
)

const (
	genDir = "base-cli-gen/cmd/gen"
)

func GenConfig() {
	generate_config.Run()
}

func Manipulate() {
	manipulate_config.StartServer("9080")
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
