package objectlock

import (
	"context"
	"fmt"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/objects/common"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type getOptions struct {
	Dst string
}

func GetCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts getOptions

	cmd := &cobra.Command{
		Use:   "get [dst]",
		Short: manager.T("cli.auth.object_storage.objects.object_lock.get.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runGet(ctx, objectService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.objects.dst"))

	return cmd
}

func runGet(ctx context.Context, objectService objSdk.ObjectService, args []string, opts getOptions, rawMode bool) error {
	if objectService == nil {
		return nil
	}

	path := opts.Dst

	if len(args) > 0 {
		path = args[0]
	}

	if path == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o caminho do objeto como argumento ou usar a flag --dst")

		return nil
	}

	bucketName, objectKey := common.ParseBucketNameAndObjectKey(path)

	info, err := objectService.GetObjectLockInfo(ctx, bucketName, objectKey)
	if err != nil {
		return fmt.Errorf("failed to get object lock: %w", err)
	}

	beautiful.NewOutput(rawMode).PrintData(info)

	return nil
}
