package commands

import (
	"fmt"
	"log"

	"github.com/magaluCloud/cligen/commands/gen_cli_code"
	"github.com/magaluCloud/cligen/commands/sdk_structure"
	"github.com/magaluCloud/cligen/config"

	"github.com/spf13/cobra"
)

// AllCommands retorna todos os comandos disponíveis
func AllCommands() []*cobra.Command {
	return []*cobra.Command{
		CloneSDKCmd(),
		GenCLICmd(),
		GenCLISDKStructureCmd(),
		GenCLICodeCmd(),
	}
}

// GenCLICodeCmd retorna o comando para gerar o código da CLI
func GenCLICodeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "gen-cli-code",
		Short: "Gerar o código da CLI",
		Run: func(cmd *cobra.Command, args []string) {
			gen_cli_code.GenCliCode()
		},
	}
}

// GenCLISDKStructureCmd retorna o comando para imprimir a estrutura do SDK
func GenCLISDKStructureCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sdk-structure",
		Short: "Printa a estrutura do SDK da CLI",
		Run: func(cmd *cobra.Command, args []string) {

			cfg, err := config.LoadConfig()
			if err != nil {
				panic(fmt.Errorf("erro ao carregar configuração: %w", err))
			}

			sdkStructure, err := sdk_structure.GenCliSDKStructure(cfg)
			if err != nil {
				log.Fatalf("Erro ao gerar a estrutura do SDK: %v", err)
			}
			sdk_structure.PrintSDKStructure(&sdkStructure)
		},
	}
}
