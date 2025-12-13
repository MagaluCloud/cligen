package commands

import (
	"context"
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
		GenerateConfigCmd(),
		ManipulateConfigCmd(),
	}
}

// ManipulateConfigCmd retorna o comando para manipular a configuração
func ManipulateConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "manipulate-config",
		Short: "Manipula a configuração",
		Run: func(cmd *cobra.Command, args []string) {
			gen_cli_code.Manipulate()
		},
	}
}

// WriteConfigCmd retorna o comando para escrever a configuração
func GenerateConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "generate-config",
		Short: "Escreve a configuração",
		Run: func(cmd *cobra.Command, args []string) {
			gen_cli_code.GenConfig()
		},
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

			ctx := context.Background()
			sdkStructure, err := sdk_structure.GenCliSDKStructure(ctx, cfg)
			if err != nil {
				log.Fatalf("Erro ao gerar a estrutura do SDK: %v", err)
			}
			sdk_structure.PrintSDKStructure(&sdkStructure)
		},
	}
}
