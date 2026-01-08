package objectstorage

import (
	"fmt"

	sdk "github.com/MagaluCloud/mgc-sdk-go/client"
	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/common/auth"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/buckets"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/objects"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

// ObjectStorageCmd cria e configura o comando de object storage
func ObjectStorageCmd(parent *cobra.Command) {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:     "object-storage",
		Short:   manager.T("cli.object_storage.short"),
		Long:    manager.T("cli.object_storage.long"),
		Aliases: []string{"object", "objects", "objs", "os"},
		GroupID: "products",
	}

	ctx := parent.Context()

	sdkCoreConfig := ctx.Value(cmdutils.CTX_SDK_KEY).(sdk.CoreClient)
	authCtx := ctx.Value(cmdutils.CTX_AUTH_KEY).(auth.Auth)

	accessKeyID := authCtx.GetAccessKeyID()
	secretAccessKey := authCtx.GetSecretAccessKey()

	var bucketService objSdk.BucketService = nil
	var objectService objSdk.ObjectService = nil

	objectStorageService, err := objSdk.New(&sdkCoreConfig, accessKeyID, secretAccessKey)
	if err != nil {
		beautiful.NewOutput(false).PrintError(fmt.Errorf("erro ao acessar o service: %w", err).Error())
	} else {
		bucketService = objectStorageService.Buckets()
		objectService = objectStorageService.Objects()
	}

	// Adicionar subcomandos
	cmd.AddCommand(buckets.BucketsCommand(ctx, bucketService))
	cmd.AddCommand(objects.ObjectsCommand(ctx, objectService))

	parent.AddCommand(cmd)
}
