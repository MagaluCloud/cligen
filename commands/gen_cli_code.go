package commands

import (
	genclicode "cligen/commands/gen_cli_code"

	"github.com/spf13/cobra"
)

func GenCLICodeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "gen-cli-code",
		Short: "Gerar o código da CLI",
		Run: func(cmd *cobra.Command, args []string) {
			genclicode.GenCliCode()
		},
	}
}
