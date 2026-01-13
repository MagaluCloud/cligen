package objects

import (
	"context"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/objects/common"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type versionsOptions struct {
	Dst string
}

func VersionsCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts versionsOptions

	cmd := &cobra.Command{
		Use:   "versions [dst]",
		Short: manager.T("cli.auth.object_storage.objects.versions.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runVersions(ctx, objectService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.objects.dst"))

	return cmd
}

func runVersions(ctx context.Context, objectService objSdk.ObjectService, args []string, opts versionsOptions, rawMode bool) error {
	path := opts.Dst

	if len(args) > 0 {
		path = args[0]
	}

	if path == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o caminho do objeto como argumento ou usar a flag --dst")

		return nil
	}

	bucketName, objectKey := common.ParseBucketNameAndObjectKey(path)

	versions, err := objectService.ListAllVersions(ctx, bucketName, objectKey)
	if err != nil {
		return err
	}

	beautiful.NewOutput(false).PrintData(versions)

	return nil
}
