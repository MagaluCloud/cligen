package objects

import (
	"context"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/objects/acl"
	objectlock "github.com/magaluCloud/mgccli/cmd/static/object_storage/objects/object_lock"
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
	cmd.AddCommand(DeleteCommand(ctx, objectService))
	cmd.AddCommand(DownloadCommand(ctx, objectService))
	cmd.AddCommand(HeadCommand(ctx, objectService))
	cmd.AddCommand(PresignCommand(ctx, objectService))
	cmd.AddCommand(VersionsCommand(ctx, objectService))
	cmd.AddCommand(UploadCommand(ctx, objectService))
	cmd.AddCommand(CopyCommand(ctx, objectService))
	cmd.AddCommand(PublicURLCommand(ctx))

	cmd.AddCommand(objectlock.ObjectLockCommand(ctx, objectService))
	cmd.AddCommand(acl.AclCommand(ctx))

	return cmd
}
