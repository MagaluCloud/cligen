package objects

import (
	"context"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

func ListCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "list",
		Short: manager.T("cli.auth.object_storage.objects.list.short"),
		Long:  manager.T("cli.auth.object_storage.objects.list.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runList(ctx, objectService, raw)
		},
	}

	return cmd
}

func runList(ctx context.Context, objectService objSdk.ObjectService, rawMode bool) error {
	return nil
}
