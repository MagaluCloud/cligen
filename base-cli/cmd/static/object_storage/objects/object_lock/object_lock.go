package objectlock

import (
	"context"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

func ObjectLockCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "object-lock",
		Short: manager.T("cli.auth.object_storage.objects.object_lock.short"),
	}

	cmd.AddCommand(SetCommand(ctx, objectService))

	return cmd
}
