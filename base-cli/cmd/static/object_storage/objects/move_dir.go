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

type moveDirOptions struct {
	Src       string
	Dst       string
	BatchSize int
}

func MoveDirCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts moveDirOptions

	cmd := &cobra.Command{
		Use:   "move-dir [src] [dst]",
		Short: manager.T("cli.auth.object_storage.objects.move_dir.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runMoveDir(ctx, objectService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.objects.move_dir.dst"))
	cmd.Flags().StringVar(&opts.Src, "src", "", manager.T("cli.auth.object_storage.objects.move_dir.src"))
	cmd.Flags().IntVar(&opts.BatchSize, "batch-size", 1000, manager.T("cli.auth.object_storage.objects.move_dir.batch_size"))

	return cmd
}

func runMoveDir(ctx context.Context, objectService objSdk.ObjectService, args []string, opts moveDirOptions, rawMode bool) error {
	if objectService == nil {
		return nil
	}

	src := opts.Src
	if len(args) > 0 {
		src = args[0]
	}
	if src == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o caminho dos arquivos que deseja mover como argumento ou usar a flag --src")

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

	if !srcIsRemote && !dstIsRemote {
		return fmt.Errorf("operation not supported, this command cannot be used to move a local source to a local destination")
	}

	srcConfig := objSdk.CopyPath{
		BucketName: bucketNameSrc,
		ObjectKey:  objectKeySrc,
	}
	dstConfig := objSdk.CopyPath{
		BucketName: bucketNameDst,
		ObjectKey:  objectKeyDst,
	}

	fmt.Fprintln(os.Stderr, "Movendo...")

	if srcIsRemote && dstIsRemote {
		err := moveDirRemote(ctx, objectService, srcConfig, dstConfig, opts.BatchSize)
		if err != nil {
			return err
		}
	} else if srcIsRemote {
		err := moveDirRemoteLocal(ctx, objectService, srcConfig, dst, opts.BatchSize)
		if err != nil {
			return err
		}
	} else {
		err := moveDirLocalRemote(ctx, objectService, src, dstConfig, opts.BatchSize)
		if err != nil {
			return err
		}
	}

	fmt.Fprintln(os.Stderr, "✓ Movido com sucesso!")

	return nil
}

func moveDirRemote(ctx context.Context, objectService objSdk.ObjectService, src objSdk.CopyPath, dst objSdk.CopyPath, batchSize int) error {
	_, err := objectService.CopyAll(ctx, src, dst, &objSdk.CopyAllOptions{
		BatchSize: batchSize,
	})
	if err != nil {
		return fmt.Errorf("erro ao copiar os objetos para o destino: %w", err)
	}

	deleteAllOpts := &objSdk.DeleteAllOptions{BatchSize: &batchSize}

	if src.ObjectKey != "" {
		deleteAllOpts.ObjectKey = src.ObjectKey
	}

	_, err = objectService.DeleteAll(ctx, src.BucketName, deleteAllOpts)
	if err != nil {
		return fmt.Errorf("erro ao deletar os objetos de origem: %w", err)
	}

	return nil
}

func moveDirLocalRemote(ctx context.Context, objectService objSdk.ObjectService, src string, dst objSdk.CopyPath, batchSize int) error {
	_, err := objectService.UploadDir(ctx, dst.BucketName, dst.ObjectKey, src, &objSdk.UploadDirOptions{
		BatchSize: batchSize,
	})

	err = os.RemoveAll(src)
	if err != nil {
		return fmt.Errorf("erro ao deletar os objetos de origem: %w", err)
	}

	return nil
}

func moveDirRemoteLocal(ctx context.Context, objectService objSdk.ObjectService, src objSdk.CopyPath, dst string, batchSize int) error {
	_, err := objectService.DownloadAll(ctx, src.BucketName, dst, &objSdk.DownloadAllOptions{
		Prefix:    src.ObjectKey,
		BatchSize: batchSize,
	})
	if err != nil {
		return fmt.Errorf("erro ao fazer o download dos objetos de origem %w", err)
	}

	_, err = objectService.DeleteAll(ctx, src.BucketName, &objSdk.DeleteAllOptions{
		ObjectKey: src.ObjectKey,
		BatchSize: &batchSize,
	})
	if err != nil {
		return fmt.Errorf("erro ao deletar os objetos de origem: %w", err)
	}

	return nil
}
