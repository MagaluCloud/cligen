package objects

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/objects/common"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type listOptions struct {
	Dst       string
	Filter    string
	Recursive bool
}

type objectInfo struct {
	Key          string     `json:"key"`
	Size         int64      `json:"size"`
	LastModified *time.Time `json:"last_modified,omitempty"`
	ETag         string     `json:"etag,omitempty"`
	StorageClass string     `json:"storage_class,omitempty"`
}

func ListCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts listOptions

	cmd := &cobra.Command{
		Use:   "list [dst]",
		Short: manager.T("cli.auth.object_storage.objects.list.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.Filter == "help" {
				common.PrintFilterHelp()
				return nil
			}

			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runList(ctx, objectService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.objects.list.dst"))
	cmd.Flags().StringVar(&opts.Filter, "filter", "", manager.T("cli.auth.object_storage.objects.list.filter"))
	cmd.Flags().BoolVar(&opts.Recursive, "recursive", false, manager.T("cli.auth.object_storage.objects.list.recursive"))

	return cmd
}

func runList(ctx context.Context, objectService objSdk.ObjectService, args []string, opts listOptions, rawMode bool) error {
	if objectService == nil {
		return nil
	}

	path := opts.Dst

	if len(args) > 0 {
		path = args[0]
	}

	if path == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o caminho do bucket como argumento ou usar a flag --dst")

		return nil
	}

	delimiter := ""

	if !opts.Recursive {
		delimiter = "/"
	}

	limit := 99999999

	listOpts := objSdk.ObjectListOptions{
		Delimiter: delimiter,
		Limit:     &limit,
	}

	bucketName, objectKey := common.ParseBucketNameAndObjectKey(path)

	if objectKey != "" {
		listOpts.Prefix = objectKey
	}

	if opts.Filter != "" {
		var filter *[]objSdk.FilterOptions
		if err := json.Unmarshal([]byte(opts.Filter), &filter); err != nil {
			return fmt.Errorf("--filter JSON inválido: %w", err)
		}

		listOpts.Filter = filter
	}

	objects, err := objectService.List(ctx, bucketName, listOpts)
	if err != nil {
		return err
	}

	results := []objectInfo{}

	for _, object := range objects {
		var lastModified *time.Time

		lastModified = &object.LastModified

		if lastModified.IsZero() {
			lastModified = nil
		}

		results = append(results, objectInfo{
			Key:          object.Key,
			Size:         object.Size,
			LastModified: lastModified,
			ETag:         object.ETag,
			StorageClass: object.StorageClass,
		})
	}

	beautiful.NewOutput(rawMode).PrintData(results)

	return nil
}
