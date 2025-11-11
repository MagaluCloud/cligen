package static

import (
	sdk "github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/magaluCloud/mgccli/cmd/static/auth"
	"github.com/magaluCloud/mgccli/cmd/static/config"
	"github.com/magaluCloud/mgccli/cmd/static/i18n"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/spf13/cobra"
)

func RootStatic(parent *cobra.Command) {
	i18n.I18nCmd(parent)

	sdkCoreConfig := parent.Context().Value(cmdutils.CTX_SDK_KEY).(sdk.CoreClient)
	config.ConfigCmd(parent, sdkCoreConfig)
	auth.AuthCmd(parent.Context(), parent, sdkCoreConfig)
}
