package static

import (
	"mgccli/cmd/static/config"

	sdk "github.com/MagaluCloud/mgc-sdk-go/client"
	"github.com/spf13/cobra"
)

func RootStatic(sdkCoreConfig sdk.CoreClient, root *cobra.Command) {

	config.ConfigCmd(sdkCoreConfig, root)

}
