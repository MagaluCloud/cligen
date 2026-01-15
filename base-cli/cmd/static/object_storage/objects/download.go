package objects

import (
	"context"
	"fmt"
	"io"
	"os"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/objects/common"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type downloadOptions struct {
	Src        string
	ObjVersion string
	Dst        string
}

type downloadReturn struct {
	Dst string `json:"dst"`
	Src string `json:"src"`
}

func DownloadCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts downloadOptions

	cmd := &cobra.Command{
		Use:   "download [src]",
		Short: manager.T("cli.auth.object_storage.objects.download.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runDownload(ctx, objectService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Src, "src", "", manager.T("cli.auth.object_storage.objects.dst"))
	cmd.Flags().StringVar(&opts.ObjVersion, "obj-version", "", manager.T("cli.auth.object_storage.objects.obj_version"))
	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.objects.download.dst"))

	return cmd
}

func runDownload(ctx context.Context, objectService objSdk.ObjectService, args []string, opts downloadOptions, rawMode bool) error {
	if objectService == nil {
		return nil
	}

	src := opts.Src

	if len(args) > 0 {
		src = args[0]
	}

	if src == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o caminho do objeto como argumento ou usar a flag --src")

		return nil
	}

	bucketName, objectKey := common.ParseBucketNameAndObjectKey(src)

	var downloadOptions *objSdk.DownloadStreamOptions

	if opts.ObjVersion != "" {
		downloadOptions = &objSdk.DownloadStreamOptions{VersionID: opts.ObjVersion}
	}

	objectReader, err := objectService.DownloadStream(ctx, bucketName, objectKey, downloadOptions)
	if err != nil {
		return fmt.Errorf("failed to get object content: %w", err)
	}

	dst, err := common.GetFileDst(opts.Dst, objectKey)
	if err != nil {
		return err
	}

	file, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %q: %w", dst, err)
	}
	defer file.Close()

	success := false
	defer func() {
		file.Close()
		if !success {
			_ = os.Remove(dst)
		}
	}()

	if _, err := io.Copy(file, objectReader); err != nil {
		return fmt.Errorf("failed to write object content to file %q: %w", dst, err)
	}

	success = true

	beautiful.NewOutput(rawMode).PrintData(downloadReturn{
		Dst: dst,
		Src: src,
	})

	return nil
}
