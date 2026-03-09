package policy

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

type deleteOptions struct {
	Dst string
}

// DeleteCommand cria o comando de deletar a política do bucket
func DeleteCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()
	var opts deleteOptions

	cmd := &cobra.Command{
		Use:   "delete [dst]",
		Short: manager.T("cli.auth.object_storage.buckets.policy.delete.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runDelete(ctx, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.buckets.dst"))

	return cmd
}

// runDelete executa o processo de deletar a política do bucket
func runDelete(ctx context.Context, args []string, opts deleteOptions, rawMode bool) error {
	objectStorageService, err := objectstorage.NewObjectStorage(ctx)
	if err != nil {
		return cmdutils.NewCliError(err.Error())
	}
	bucketService := objectStorageService.GetBucketService()

	bucketName := opts.Dst

	if len(args) > 0 {
		bucketName = args[0]
	}

	if bucketName == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o nome do bucket como argumento ou usar a flag --dst")

		return nil
	}

	err = bucketService.DeletePolicy(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("erro ao deletar a política: %w", err)
	}

	fmt.Fprintln(os.Stderr, "✓ Deletado com sucesso!")

	return nil
}
