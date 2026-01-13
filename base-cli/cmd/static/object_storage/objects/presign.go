package objects

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/objects/common"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type presignOptions struct {
	Dst       string
	Method    string
	ExpiresIn string
}

func PresignCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts presignOptions

	cmd := &cobra.Command{
		Use:   "presign [dst]",
		Short: manager.T("cli.auth.object_storage.objects.presign.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runPresign(ctx, objectService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.objects.dst"))
	cmd.Flags().StringVar(&opts.Method, "method", http.MethodGet, manager.T("cli.auth.object_storage.objects.presign.method"))
	cmd.Flags().StringVar(&opts.ExpiresIn, "expires-in", "5m", manager.T("cli.auth.object_storage.objects.presign.expires_in"))

	return cmd
}

func runPresign(ctx context.Context, objectService objSdk.ObjectService, args []string, opts presignOptions, rawMode bool) error {
	if objectService == nil {
		return nil
	}

	path := opts.Dst

	if len(args) > 0 {
		path = args[0]
	}

	if path == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o caminho do objeto como argumento ou usar a flag --dst")

		return nil
	}

	bucketName, objectKey := common.ParseBucketNameAndObjectKey(path)

	var expiresIn *time.Duration

	if opts.ExpiresIn != "" {
		duration, err := time.ParseDuration(opts.ExpiresIn)
		if err != nil {
			return fmt.Errorf("invalid format to --expires-in: %w", err)
		}

		seconds := time.Duration(duration.Seconds()) * time.Second
		expiresIn = &seconds
	}

	presignedOpts := objSdk.GetPresignedURLOptions{
		Method:          opts.Method,
		ExpiryInSeconds: expiresIn,
	}

	presignedURL, err := objectService.GetPresignedURL(ctx, bucketName, objectKey, presignedOpts)
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "%s\n", presignedURL.URL)

	return nil
}
