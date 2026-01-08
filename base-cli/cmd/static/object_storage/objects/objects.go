package objects

import (
	"context"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

func ObjectsCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "objects",
		Short: manager.T("cli.auth.object_storage.objects.short"),
		Long:  manager.T("cli.auth.object_storage.objects.long"),
	}

	cmd.AddCommand(ListCommand(ctx, objectService))

	return cmd
}
