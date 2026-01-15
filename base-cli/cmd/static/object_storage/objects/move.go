package objects

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/objects/common"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type moveOptions struct {
	Src string
	Dst string
}

func MoveCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts moveOptions

	cmd := &cobra.Command{
		Use:   "move [src] [dst]",
		Short: manager.T("cli.auth.object_storage.objects.move.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runMove(ctx, objectService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.objects.move.dst"))
	cmd.Flags().StringVar(&opts.Src, "src", "", manager.T("cli.auth.object_storage.objects.move.src"))

	return cmd
}

func runMove(ctx context.Context, objectService objSdk.ObjectService, args []string, opts moveOptions, rawMode bool) error {
	if objectService == nil {
		return nil
	}

	src := opts.Src
	if len(args) > 0 {
		src = args[0]
	}
	if src == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o caminho do arquivo que deseja mover como argumento ou usar a flag --src")

		return nil
	}

	dst := opts.Dst
	if len(args) > 1 {
		dst = args[1]
	}
	if dst == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o caminho de destino como argumento ou usar a flag --dst")

		return nil
	}

	bucketNameSrc, objectKeySrc := common.ParseBucketNameAndObjectKey(src)
	bucketNameDst, objectKeyDst := common.ParseBucketNameAndObjectKey(dst)

	srcIsRemote := isRemote(src)
	dstIsRemote := isRemote(dst)

	srcConfig := objSdk.CopySrcConfig{
		BucketName: bucketNameSrc,
		ObjectKey:  objectKeySrc,
	}
	dstConfig := objSdk.CopyDstConfig{
		BucketName: bucketNameDst,
		ObjectKey:  objectKeyDst,
	}

	if !srcIsRemote && !dstIsRemote {
		return fmt.Errorf("operation not supported, this command cannot be used to move a local source to a local destination")
	} else if srcIsRemote && dstIsRemote {

		err := moveRemote(ctx, objectService, srcConfig, dstConfig)
		if err != nil {
			return err
		}
	} else if srcIsRemote {
		err := moveRemoteLocal(ctx, objectService, srcConfig, dst)
		if err != nil {
			return err
		}
	} else {
		err := moveLocalRemote(ctx, objectService, src, dstConfig)
		if err != nil {
			return err
		}
	}

	fmt.Fprintln(os.Stderr, "✓ Movido com sucesso!")

	return nil
}

func moveRemote(ctx context.Context, objectService objSdk.ObjectService, src objSdk.CopySrcConfig, dst objSdk.CopyDstConfig) error {
	err := objectService.Copy(ctx, src, dst)
	if err != nil {
		return fmt.Errorf("erro ao copiar o objeto para o destino")
	}

	err = objectService.Delete(ctx, src.BucketName, src.ObjectKey, nil)
	if err != nil {
		return fmt.Errorf("erro ao deletar o objeto de origem")
	}

	return nil
}

func moveLocalRemote(ctx context.Context, objectService objSdk.ObjectService, src string, dst objSdk.CopyDstConfig) error {
	srcAbs, err := filepath.Abs(src)
	if err != nil {
		return fmt.Errorf("error to get absolute representation of path: %w", err)
	}

	fileBytes, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("erro ao ler o arquivo: %w", err)
	}

	contentType := "application/octet-stream"
	if len(fileBytes) > 0 {
		contentType = http.DetectContentType(fileBytes)
	}

	err = objectService.Upload(ctx, dst.BucketName, dst.ObjectKey, fileBytes, contentType, nil)
	if err != nil {
		return err
	}

	err = os.Remove(srcAbs)
	if err != nil {
		return fmt.Errorf("error to delete the source: %w", err)
	}

	return nil
}

func moveRemoteLocal(ctx context.Context, objectService objSdk.ObjectService, src objSdk.CopySrcConfig, dst string) error {
	objectReader, err := objectService.DownloadStream(ctx, src.BucketName, src.ObjectKey, nil)
	if err != nil {
		return fmt.Errorf("failed to get object content: %w", err)
	}

	fileDst, err := common.GetFileDst(dst, src.ObjectKey)
	if err != nil {
		return err
	}

	file, err := os.Create(fileDst)
	if err != nil {
		return fmt.Errorf("failed to create destination file %q: %w", fileDst, err)
	}
	defer file.Close()

	success := false
	defer func() {
		file.Close()
		if !success {
			_ = os.Remove(fileDst)
		}
	}()

	if _, err := io.Copy(file, objectReader); err != nil {
		return fmt.Errorf("failed to write object content to file %q: %w", fileDst, err)
	}

	success = true

	err = objectService.Delete(ctx, src.BucketName, src.ObjectKey, nil)
	if err != nil {
		return fmt.Errorf("erro ao deletar o objeto de origem")
	}

	return nil
}

func isRemote(path string) bool {
	pathUrl, err := url.Parse(path)
	if err != nil {
		return false
	}

	return pathUrl.Scheme == "s3"
}
