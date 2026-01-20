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

type downloadAllOptions struct {
	Dst    string
	Src    string
	Filter string
}

func DownloadAllCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts downloadAllOptions

	cmd := &cobra.Command{
		Use:   "download-all [src]",
		Short: manager.T("cli.auth.object_storage.objects.download_all.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Filter == "help" {
				common.PrintFilterHelp()
				return nil
			}

			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runDownloadAll(ctx, objectService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Src, "src", "", manager.T("cli.auth.object_storage.objects.download_all.src"))
	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.objects.download_all.dst"))
	cmd.Flags().StringVar(&opts.Filter, "filter", "", manager.T("cli.auth.object_storage.objects.download_all.filter"))

	return cmd
}

func runDownloadAll(ctx context.Context, objectService objSdk.ObjectService, args []string, opts downloadAllOptions, rawMode bool) error {
	if objectService == nil {
		return nil
	}

	path := opts.Src

	if len(args) > 0 {
		path = args[0]
	}

	if path == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o nome do bucket como argumento ou usar a flag --src")

		return nil
	}

	bucketName, objectKey := common.ParseBucketNameAndObjectKey(path)

	downloadOpts := objSdk.DownloadAllOptions{
		Prefix: "",
	}

	if objectKey != "" {
		downloadOpts.Prefix = objectKey
	}

	if opts.Filter != "" {
		var filter *[]objSdk.FilterOptions
		if err := json.Unmarshal([]byte(opts.Filter), &filter); err != nil {
			return fmt.Errorf("--filter JSON inválido: %w", err)
		}

		downloadOpts.Filter = filter
	}

	dst := opts.Dst

	if dst == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed to get current working directory: %w", err)
		}

		dst = cwd
	}

	fmt.Println("Downloading...")

	downloadAllResult, err := objectService.DownloadAll(ctx, bucketName, dst, &downloadOpts)
	if err != nil {
		return err
	}

	if downloadAllResult.ErrorCount > 0 {
		beautiful.NewOutput(rawMode).PrintError("não foi possível fazer o download de alguns objetos")
	} else {
		fmt.Fprintln(os.Stderr, "✓ Download realizado com sucesso!")
	}

	return nil
}
