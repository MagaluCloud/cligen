package versioning

import (
	"context"
	"fmt"
	"os"

	"github.com/magaluCloud/mgccli/beautiful"
	objectstorage "github.com/magaluCloud/mgccli/cmd/common/object_storage"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type enableOptions struct {
	Bucket string
}

// EnableCommand cria o comando de habilitar o versionamento do bucket
func EnableCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()
	var opts enableOptions

	cmd := &cobra.Command{
		Use:   "enable [bucket]",
		Short: manager.T("cli.auth.object_storage.buckets.versioning.enable.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runEnable(ctx, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Bucket, "bucket", "", manager.T("cli.auth.object_storage.buckets.dst"))

	return cmd
}

// runEnable executa o processo de habilitar o versionamento do bucket
func runEnable(ctx context.Context, args []string, opts enableOptions, rawMode bool) error {
	objectStorageService, err := objectstorage.NewObjectStorage(ctx)
	if err != nil {
		return cmdutils.NewCliError(err.Error())
	}
	bucketService := objectStorageService.GetBucketService()

	bucketName := opts.Bucket

	if len(args) > 0 {
		bucketName = args[0]
	}

	if bucketName == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o nome do bucket como argumento ou usar a flag --bucket")

		return nil
	}

	err = bucketService.EnableVersioning(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("erro ao habilitar o versionamento: %w", err)
	}

	fmt.Fprintln(os.Stderr, "✓ Habilitado o versionamento do bucket!")

	return nil
}
