package code

import (
	"fmt"

	"github.com/magaluCloud/cligen/commands/gen_cli_code/code/gomod"
	"github.com/magaluCloud/cligen/commands/gen_cli_code/code/menu"
	"github.com/magaluCloud/cligen/commands/gen_cli_code/code/menu_item"
	"github.com/magaluCloud/cligen/commands/gen_cli_code/code/module"
	"github.com/magaluCloud/cligen/commands/gen_cli_code/code/root_gen"
	"github.com/magaluCloud/cligen/config"
)

func Run() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("erro ao carregar configuração: %w", err))
	}

	gomod.GenGoModFile(cfg)
	root_gen.GenerateRootGen(cfg)
	module.GenerateModule(cfg)
	menu.GenerateMenu(cfg)
	menu_item.GenerateMenuItem(cfg)
}
