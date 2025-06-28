package config

import (
	sdk "github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/spf13/cobra"
)

func ConfigCmd(sdkCoreConfig sdk.CoreClient, root *cobra.Command) {
	cmd := &cobra.Command{
		Use:     "config",
		Short:   "Configuração do CLI",
		Long:    `Configuração do CLI`,
		GroupID: "settings",
	}

	cmd.AddCommand(List())
	cmd.AddCommand(Delete())
	cmd.AddCommand(Get())
	cmd.AddCommand(Set())

	root.AddCommand(cmd)
}
