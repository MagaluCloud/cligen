package code

import (
	"fmt"

	"github.com/magaluCloud/cligen/commands/gen_cli_code/gen_config/code/module"
	"github.com/magaluCloud/cligen/commands/gen_cli_code/gen_config/code/root_gen"
	"github.com/magaluCloud/cligen/config"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("erro ao carregar configuração: %w", err))
	}

	root_gen.GenerateRootGen(cfg)
	module.GenerateModule(cfg)
	fmt.Println(cfg.SDKTag)
}
