package config

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/magaluCloud/mgccli/cmd/common/config"
	"github.com/spf13/cobra"
)

func List(config config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Listar configurações",
		Long:  `Listar configurações`,
		Run: func(cmd *cobra.Command, args []string) {

			configMap, err := config.List()
			if err != nil {
				fmt.Println("Erro ao listar configurações:", err)
				return
			}
			for key, value := range configMap {
				fmt.Printf("%s: %v\n", color.BlueString(key), color.YellowString(fmt.Sprintf("%v", value)))
			}
		},
	}
	return cmd
}
