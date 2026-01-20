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

type copyAllOptions struct {
	Dst          string
	Src          string
	StorageClass string
	Filter       string
}

func CopyAllCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts copyAllOptions

	cmd := &cobra.Command{
		Use:   "copy-all [src] [dst]",
		Short: manager.T("cli.auth.object_storage.objects.copy_all.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Filter == "help" {
				common.PrintFilterHelp()
				return nil
			}

			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runCopyAll(ctx, objectService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Src, "src", "", manager.T("cli.auth.object_storage.objects.copy_all.src"))
	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.objects.copy_all.dst"))
	cmd.Flags().StringVar(&opts.StorageClass, "storage-class", "", manager.T("cli.auth.object_storage.objects.copy_all.storage_class"))
	cmd.Flags().StringVar(&opts.Filter, "filter", "", manager.T("cli.auth.object_storage.objects.copy_all.filter"))

	return cmd
}

func runCopyAll(ctx context.Context, objectService objSdk.ObjectService, args []string, opts copyAllOptions, rawMode bool) error {
	if objectService == nil {
		return nil
	}

	src := opts.Src

	if len(args) > 0 {
		src = args[0]
	}

	if src == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o caminho dos arquivos que deseja copiar como argumento ou usar a flag --src")

		return nil
	}

	dst := opts.Dst

	if len(args) > 1 {
		dst = args[1]
	}

	if dst == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o caminho de destino no bucket como argumento ou usar a flag --dst")

		return nil
	}

	copyOpts := objSdk.CopyAllOptions{
		StorageClass: opts.StorageClass,
	}

	if opts.Filter != "" {
		var filter *[]objSdk.FilterOptions
		if err := json.Unmarshal([]byte(opts.Filter), &filter); err != nil {
			return fmt.Errorf("--filter JSON inválido: %w", err)
		}

		copyOpts.Filter = filter
	}

	fmt.Println("Copiando...")

	bucketNameSrc, objectKeySrc := common.ParseBucketNameAndObjectKey(src)
	bucketNameDst, objectKeyDst := common.ParseBucketNameAndObjectKey(dst)

	copyAllResult, err := objectService.CopyAll(
		ctx,
		objSdk.CopyPath{
			BucketName: bucketNameSrc,
			ObjectKey:  objectKeySrc,
		}, objSdk.CopyPath{
			BucketName: bucketNameDst,
			ObjectKey:  objectKeyDst,
		},
		&copyOpts,
	)
	if err != nil {
		return err
	}

	if copyAllResult.ErrorCount > 0 {
		beautiful.NewOutput(rawMode).PrintError("não foi possível fazer a cópia de alguns objetos")
	} else {
		fmt.Fprintln(os.Stderr, "✓ Cópia realizado com sucesso!")
	}

	return nil
}
