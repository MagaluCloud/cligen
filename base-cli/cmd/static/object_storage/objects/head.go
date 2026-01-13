package objects

import (
	"context"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/objects/common"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type headOptions struct {
	Dst        string
	ObjVersion string
}

func HeadCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts headOptions

	cmd := &cobra.Command{
		Use:   "head [dst]",
		Short: manager.T("cli.auth.object_storage.objects.head.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runHead(ctx, objectService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.objects.dst"))
	cmd.Flags().StringVar(&opts.ObjVersion, "obj-version", "", manager.T("cli.auth.object_storage.objects.obj_version"))

	return cmd
}

func runHead(ctx context.Context, objectService objSdk.ObjectService, args []string, opts headOptions, rawMode bool) error {
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

	var headOptions *objSdk.MetadataOptions

	if opts.ObjVersion != "" {
		headOptions = &objSdk.MetadataOptions{VersionID: opts.ObjVersion}
	}

	metadata, err := objectService.Metadata(ctx, bucketName, objectKey, headOptions)
	if err != nil {
		return err
	}

	beautiful.NewOutput(rawMode).PrintData(metadata)

	return nil
}
