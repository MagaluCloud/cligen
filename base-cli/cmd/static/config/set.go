package config

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/magaluCloud/mgccli/cmd/common/config"
	"github.com/spf13/cobra"
)

type SetOptions struct {
	Key   string
	Value string
}

func Set(config config.Config) *cobra.Command {
	var opts SetOptions

	cmd := &cobra.Command{
		Use:   "set [key] [value]",
		Short: "Definir configurações",
		Long:  `Definir configurações`,
		Run: func(cmd *cobra.Command, args []string) {
			key := opts.Key
			value := opts.Value

			if len(args) > 0 {
				key = args[0]
			}
			if len(args) > 1 {
				value = args[1]
			}

			if key == "" || value == "" {
				fmt.Println("Erro: chave e/ou valor não especificados")
				return
			}

			err := config.Set(key, value)
			if err != nil {
				fmt.Println("Erro ao definir configuração:", err)
				return
			}
			fmt.Printf("%s: %v\n", color.BlueString(key), color.YellowString(value))
		},
	}

	cmd.Flags().StringVar(&opts.Key, "key", "", "Nome da configuração desejada (required)")
	cmd.Flags().StringVar(&opts.Value, "value", "", "Valor da configuração desejada (required)")

	return cmd
}
