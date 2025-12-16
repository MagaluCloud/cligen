package commands

import (
	"github.com/magaluCloud/cligen/commands/gen_cli_code"

	"github.com/spf13/cobra"
)

// AllCommands retorna todos os comandos disponíveis
func AllCommands() []*cobra.Command {
	return []*cobra.Command{
		CloneSDKCmd(),
		GenCLICmd(),
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
