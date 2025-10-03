package commands

import (
	"log"

	"github.com/magaluCloud/cligen/commands/sdk_structure"

	"github.com/spf13/cobra"
)

func GenCLISDKStructureCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "sdk-structure",
		Short: "Printa a estrutura do SDK da CLI",
		Run: func(cmd *cobra.Command, args []string) {
			sdkStructure, err := sdk_structure.GenCliSDKStructure()
			if err != nil {
				log.Fatalf("Erro ao gerar a estrutura do SDK: %v", err)
			}
			sdk_structure.PrintSDKStructure(&sdkStructure)
		},
	}
}
