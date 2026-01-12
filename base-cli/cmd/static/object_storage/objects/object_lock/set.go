package objectlock

import (
	"context"
	"fmt"
	"os"
	"time"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/objects/common"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type setOptions struct {
	Dst             string
	RetainUntilDate string
}

func SetCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts setOptions

	cmd := &cobra.Command{
		Use:   "set [dst]",
		Short: manager.T("cli.auth.object_storage.objects.object_lock.set.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runSet(ctx, objectService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.objects.dst"))
	cmd.Flags().StringVar(&opts.RetainUntilDate, "retain-until-date", "", manager.T("cli.auth.object_storage.objects.object_lock.retain_until_date"))

	cmd.MarkFlagRequired("retain-until-date")

	return cmd
}

func runSet(ctx context.Context, objectService objSdk.ObjectService, args []string, opts setOptions, rawMode bool) error {
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

	retainUntilDate, err := time.Parse(time.RFC3339, opts.RetainUntilDate)
	if err != nil {
		return fmt.Errorf("invalid ISO 8601 format: %w", err)
	}

	err = objectService.LockObject(ctx, bucketName, objectKey, retainUntilDate)
	if err != nil {
		return fmt.Errorf("failed to set object lock: %w", err)
	}

	fmt.Fprintln(os.Stderr, "✓ Object lock set successfully")

	return nil
}
