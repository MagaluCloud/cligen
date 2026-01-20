package objects

import (
	"context"
	"fmt"
	"os"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/objects/common"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type copyOptions struct {
	Dst          string
	Src          string
	StorageClass string
	ObjVersion   string
}

func CopyCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts copyOptions

	cmd := &cobra.Command{
		Use:   "copy [src] [dst]",
		Short: manager.T("cli.auth.object_storage.objects.copy.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runCopy(ctx, objectService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.objects.copy.dst"))
	cmd.Flags().StringVar(&opts.Src, "src", "", manager.T("cli.auth.object_storage.objects.copy.src"))
	cmd.Flags().StringVar(&opts.StorageClass, "storage-class", "", manager.T("cli.auth.object_storage.objects.copy.storage_class"))
	cmd.Flags().StringVar(&opts.ObjVersion, "obj-version", "", manager.T("cli.auth.object_storage.objects.obj_version"))

	return cmd
}

func runCopy(ctx context.Context, objectService objSdk.ObjectService, args []string, opts copyOptions, rawMode bool) error {
	if objectService == nil {
		return nil
	}

	src := opts.Src

	if len(args) > 0 {
		src = args[0]
	}

	if src == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o caminho do arquivo que deseja copiar como argumento ou usar a flag --src")

		return nil
	}

	path := opts.Dst

	if len(args) > 1 {
		path = args[1]
	}

	if path == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o caminho de destino no bucket como argumento ou usar a flag --dst")

		return nil
	}

	bucketNameSrc, objectKeySrc := common.ParseBucketNameAndObjectKey(src)
	bucketNameDst, objectKeyDst := common.ParseBucketNameAndObjectKey(path)

	srcCopy := objSdk.CopySrcConfig{
		BucketName: bucketNameSrc,
		ObjectKey:  objectKeySrc,
	}
	dstCopy := objSdk.CopyDstConfig{
		BucketName: bucketNameDst,
		ObjectKey:  objectKeyDst,
	}

	if opts.ObjVersion != "" {
		srcCopy.VersionID = opts.ObjVersion
	}

	if opts.StorageClass != "" {
		dstCopy.StorageClass = opts.StorageClass
	}

	err := objectService.Copy(ctx, srcCopy, dstCopy)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "✓ Cópia realizada com sucesso!")

	return nil
}
