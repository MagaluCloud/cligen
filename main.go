package main

import (
	"fmt"
	"os"

	"cligen/commands"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "cligen",
		Short: "Gerador de código para CLI baseado no SDK",
		Long:  `Gerador de código que cria automaticamente o código fonte da CLI baseado no SDK.`,
	}

	// Adicionar comandos
	rootCmd.AddCommand(commands.CloneSDKCmd())
	rootCmd.AddCommand(commands.GenCLICmd())
	rootCmd.AddCommand(commands.GenCLISDKStructureCmd())
	rootCmd.AddCommand(commands.GenCLICodeCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
