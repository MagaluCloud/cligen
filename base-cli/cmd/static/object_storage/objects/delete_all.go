package objects

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/charmbracelet/huh"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/objects/common"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type deleteAllOptions struct {
	Dst       string
	Filter    string
	BatchSize int
}

func DeleteAllCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts deleteAllOptions

	cmd := &cobra.Command{
		Use:   "delete-all [dst]",
		Short: manager.T("cli.auth.object_storage.objects.delete_all.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Filter == "help" {
				common.PrintFilterHelp()
				return nil
			}

			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runDeleteAll(ctx, objectService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.objects.delete_all.dst"))
	cmd.Flags().StringVar(&opts.Filter, "filter", "", manager.T("cli.auth.object_storage.objects.delete_all.filter"))
	cmd.Flags().IntVar(&opts.BatchSize, "batch-size", 1000, manager.T("cli.auth.object_storage.objects.delete_all.batch_size"))

	return cmd
}

func runDeleteAll(ctx context.Context, objectService objSdk.ObjectService, args []string, opts deleteAllOptions, rawMode bool) error {
	if objectService == nil {
		return nil
	}

	bucketName := opts.Dst

	if len(args) > 0 {
		bucketName = args[0]
	}

	if bucketName == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o nome do bucket como argumento ou usar a flag --dst")

		return nil
	}

	if opts.BatchSize <= 0 || opts.BatchSize > 1000 {
		return fmt.Errorf("invalid --batch-size. Range: 1 - 1000")
	}

	deleteOpts := objSdk.DeleteAllOptions{
		BatchSize: &opts.BatchSize,
	}

	if opts.Filter != "" {
		var filter *[]objSdk.FilterOptions
		if err := json.Unmarshal([]byte(opts.Filter), &filter); err != nil {
			return fmt.Errorf("--filter JSON inválido: %w", err)
		}

		deleteOpts.Filter = filter
	}

	var input string
	huh.NewInput().
		Title(fmt.Sprintf("This command will delete all objects at %s, and its result is NOT reversible. Please confirm by retyping: %s", bucketName, bucketName)).
		Value(&input).
		Run()

	if input != bucketName {
		fmt.Println("Não foi possível deletar. O texto digitado não corresponde ao nome do bucket informado!")

		return nil
	}

	progress := beautiful.NewPTermProgress(nil, nil)
	defer progress.Finish()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	go func() {
		<-sigCh
		cancel()
	}()

	ctx = objSdk.WithProgress(ctx, progress)

	deletionResult, err := objectService.DeleteAll(ctx, bucketName, &deleteOpts)
	if err != nil {
		progress.Finish()
		cancel()
		return err
	}

	if deletionResult.ErrorCount > 0 {
		beautiful.NewOutput(rawMode).PrintError("não foi possível deletar alguns objetos")
	} else {
		fmt.Fprintln(os.Stderr, "✓ Deleção realizada com sucesso!")
	}

	return nil
}
