package versioning

import (
	"context"
	"fmt"

	"github.com/magaluCloud/mgccli/beautiful"
	objectstorage "github.com/magaluCloud/mgccli/cmd/common/object_storage"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type getOptions struct {
	Bucket string
}

// GetCommand cria o comando de exibir as informações de versionamento do bucket
func GetCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()
	var opts getOptions

	cmd := &cobra.Command{
		Use:   "get [bucket]",
		Short: manager.T("cli.auth.object_storage.buckets.versioning.get.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runGet(ctx, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Bucket, "bucket", "", manager.T("cli.auth.object_storage.buckets.dst"))

	return cmd
}

// runGet executa o processo de exibir as informações de versionamento do bucket
func runGet(ctx context.Context, args []string, opts getOptions, rawMode bool) error {
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

	info, err := bucketService.GetVersioningStatus(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("erro ao pegar as informações de versionamento: %w", err)
	}

	beautiful.NewOutput(rawMode).PrintData(info)

	return nil
}
