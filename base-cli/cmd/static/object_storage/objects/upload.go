package objects

import (
	"context"
	"fmt"
	"net/http"
	"os"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/objects/common"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type uploadOptions struct {
	Dst          string
	Src          string
	StorageClass string
}

func UploadCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts uploadOptions

	cmd := &cobra.Command{
		Use:   "upload [src] [dst]",
		Short: manager.T("cli.auth.object_storage.objects.upload.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runUpload(ctx, objectService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.objects.upload.dst"))
	cmd.Flags().StringVar(&opts.Src, "src", "", manager.T("cli.auth.object_storage.objects.upload.src"))
	cmd.Flags().StringVar(&opts.StorageClass, "storage-class", "standard", manager.T("cli.auth.object_storage.objects.upload.storage_class"))

	return cmd
}

func runUpload(ctx context.Context, objectService objSdk.ObjectService, args []string, opts uploadOptions, rawMode bool) error {
	if objectService == nil {
		return nil
	}

	src := opts.Src

	if len(args) > 0 {
		src = args[0]
	}

	if src == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o caminho do arquivo que deseja fazer upload como argumento ou usar a flag --src")

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

	fileBytes, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("erro ao ler o arquivo: %w", err)
	}

	contentType := "application/octet-stream"
	if len(fileBytes) > 0 {
		contentType = http.DetectContentType(fileBytes)
	}

	err = objectService.Upload(ctx, bucketName, objectKey, fileBytes, contentType, opts.StorageClass)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "✓ Upload realizado com sucesso!")

	return nil
}
