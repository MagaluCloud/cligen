package commands

import (
	"cligen/commands/sdk_structure"

	"github.com/spf13/cobra"
)

func GenCLISDKStructureCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sdk-structure",
		Short: "Gerar a estrutura do SDK da CLI",
		Run: func(cmd *cobra.Command, args []string) {
			sdk_structure.GenCliSDKStructure()
		},
	}
}
