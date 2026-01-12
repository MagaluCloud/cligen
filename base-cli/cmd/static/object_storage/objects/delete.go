package objects

import (
	"context"
	"fmt"
	"os"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/charmbracelet/huh"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/objects/common"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type deleteOptions struct {
	Dst        string
	ObjVersion string
}

func DeleteCommand(ctx context.Context, objectService objSdk.ObjectService) *cobra.Command {
	manager := i18n.GetInstance()
	var opts deleteOptions

	cmd := &cobra.Command{
		Use:   "delete [dst]",
		Short: manager.T("cli.auth.object_storage.objects.delete.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runDelete(ctx, objectService, args, opts, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Dst, "dst", "", manager.T("cli.auth.object_storage.objects.dst"))
	cmd.Flags().StringVar(&opts.ObjVersion, "obj-version", "", manager.T("cli.auth.object_storage.objects.obj_version"))

	return cmd
}

func runDelete(ctx context.Context, objectService objSdk.ObjectService, args []string, opts deleteOptions, rawMode bool) error {
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

	var input string
	huh.NewInput().
		Title(fmt.Sprintf("This command will delete the object at %s, and its result is NOT reversible. Please confirm by retyping: %s", path, path)).
		Value(&input).
		Run()

	if input != path {
		fmt.Println("Não foi possível deletar. O texto digitado não corresponde ao caminho do objeto informado!")

		return nil
	}

	bucketName, objectKey := common.ParseBucketNameAndObjectKey(path)

	var deleteOptions *objSdk.DeleteOptions

	if opts.ObjVersion != "" {
		deleteOptions = &objSdk.DeleteOptions{VersionID: opts.ObjVersion}
	}

	err := objectService.Delete(ctx, bucketName, objectKey, deleteOptions)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "✓ Deleção realizada com sucesso!")

	return nil
}
