package objects

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/objects/common"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type uploadDirOptions struct {
	Dst          string
	Src          string
	StorageClass string
	Shallow      bool
	Filter       string
}

func UploadDirCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts uploadDirOptions

	cmd := &cobra.Command{
		Use:   "upload-dir [src] [dst]",
		Short: manager.T("cli.auth.object_storage.objects.upload_dir.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Filter == "help" {
				common.PrintFilterHelp()
				return nil
			}

			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runUploadDir(ctx, objectService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.objects.upload_dir.dst"))
	cmd.Flags().StringVar(&opts.Src, "src", "", manager.T("cli.auth.object_storage.objects.upload_dir.src"))
	cmd.Flags().StringVar(&opts.StorageClass, "storage-class", "standard", manager.T("cli.auth.object_storage.objects.upload_dir.storage_class"))
	cmd.Flags().StringVar(&opts.Filter, "filter", "", manager.T("cli.auth.object_storage.objects.upload_dir.filter"))
	cmd.Flags().BoolVar(&opts.Shallow, "shallow", false, manager.T("cli.auth.object_storage.objects.upload_dir.shallow"))

	return cmd
}

func runUploadDir(ctx context.Context, objectService objSdk.ObjectService, args []string, opts uploadDirOptions, rawMode bool) error {
	if objectService == nil {
		return nil
	}

	src := opts.Src

	if len(args) > 0 {
		src = args[0]
	}

	if src == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o caminho dos arquivos que deseja fazer upload como argumento ou usar a flag --src")

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

	bucketName, objectKey := common.ParseBucketNameAndObjectKey(path)

	uploadOpts := objSdk.UploadDirOptions{
		Shallow: opts.Shallow,
	}

	if opts.StorageClass != "" {
		uploadOpts.StorageClass = opts.StorageClass
	}

	if opts.Filter != "" {
		var filter *[]objSdk.FilterOptions
		if err := json.Unmarshal([]byte(opts.Filter), &filter); err != nil {
			return fmt.Errorf("--filter JSON inválido: %w", err)
		}

		uploadOpts.Filter = filter
	}

	fmt.Println("Uploading...")

	_, err := objectService.UploadDir(ctx, bucketName, objectKey, src, &uploadOpts)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "✓ Upload realizado com sucesso!")

	return nil
}
