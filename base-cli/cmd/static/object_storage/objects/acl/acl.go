package acl

import (
	"context"

	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

func AclCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "acl",
		Short: manager.T("cli.auth.object_storage.objects.acl.short"),
	}

	cmd.AddCommand(GetCommand(ctx))
	cmd.AddCommand(SetCommand(ctx))

	return cmd
}
