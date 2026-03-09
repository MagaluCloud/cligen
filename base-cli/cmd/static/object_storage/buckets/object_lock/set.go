package objectlock

import (
	"context"
	"fmt"
	"os"

	"github.com/magaluCloud/mgccli/beautiful"
	objectstorage "github.com/magaluCloud/mgccli/cmd/common/object_storage"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	cobrautils "github.com/magaluCloud/mgccli/cobra_utils/flags"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type SetFlags struct {
	Dst   string
	Days  uint
	Years uint
}

type SetParams struct {
	Dst   string
	Days  *uint
	Years *uint
}

// SetCommand cria o comando de setar o bloqueio de objetos
func SetCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()
	var flags SetFlags
	var params SetParams

	cmd := &cobra.Command{
		Use:   "set [dst]",
		Short: manager.T("cli.auth.object_storage.buckets.object_lock.set.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			params.Dst = flags.Dst

			cobrautils.NilIfNotChanged(cmd, "days", &params.Days, flags.Days)
			cobrautils.NilIfNotChanged(cmd, "years", &params.Years, flags.Years)

			return runSet(ctx, args, params, raw)
		},
	}

	cmd.Flags().StringVar(&flags.Dst, "dst", "", manager.T("cli.auth.object_storage.buckets.dst"))
	cmd.Flags().UintVar(&flags.Days, "days", 0, manager.T("cli.auth.object_storage.buckets.object_lock.set.days"))
	cmd.Flags().UintVar(&flags.Years, "years", 0, manager.T("cli.auth.object_storage.buckets.object_lock.set.years"))

	return cmd
}

// runSet executa o processo de setar o bloqueio de objetos
func runSet(ctx context.Context, args []string, opts SetParams, rawMode bool) error {
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

	if opts.Days == nil && opts.Years == nil {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o parâmetro 'days' ou 'years'")
		return nil
	}

	validity := opts.Days

	if validity == nil {
		validity = opts.Years
	}

	unit := "days"
	if opts.Days == nil {
		unit = "years"
	}

	err = bucketService.LockBucket(ctx, bucketName, uint(*validity), unit)
	if err != nil {
		return fmt.Errorf("erro ao bloquear: %w", err)
	}

	fmt.Fprintln(os.Stderr, "✓ Bloqueado com sucesso!")

	return nil
}
