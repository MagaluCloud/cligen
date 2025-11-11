package gen

import (
	sdk "github.com/MagaluCloud/mgc-sdk-go/client"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/spf13/cobra"
)

func RootGen(parent *cobra.Command) {
	sdkCoreConfig := parent.Context().Value(cmdutils.CTX_SDK_KEY).(*sdk.CoreClient)
	// JUST KEEP THIS FILE AS IT IS
}
