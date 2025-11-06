package main

import (
	"fmt"
	"os"

	"github.com/magaluCloud/cligen/commands"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "cligen",
		Short: "Gerador de código para CLI baseado no SDK",
		Long:  `Gerador de código que cria automaticamente o código fonte da CLI baseado no SDK.`,
	}

	// Adicionar todos os comandos
	for _, cmd := range commands.AllCommands() {
		rootCmd.AddCommand(cmd)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
