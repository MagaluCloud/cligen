package auth

import (
	"context"

	"github.com/magaluCloud/mgccli/i18n"

	sdk "github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/spf13/cobra"
)

func AuthCmd(ctx context.Context, parent *cobra.Command, sdkCoreConfig sdk.CoreClient) {
	manager := i18n.GetInstance()
	cmd := &cobra.Command{
		Use:     "auth",
		Short:   manager.T("cli.auth.short"),
		Long:    manager.T("cli.auth.long"),
		Aliases: []string{"auth"},
		GroupID: "settings",
	}

	cmd.AddCommand(Login(ctx))
	parent.AddCommand(cmd)
}
