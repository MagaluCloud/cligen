package common

import (
	"fmt"
	"os"

	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/beautiful"
)

func PrintFilterHelp() {
	fmt.Fprintln(os.Stderr, "Filter format:")

	beautiful.NewOutput(false).PrintData([]objSdk.FilterOptions{
		{Include: "string", Exclude: "string"},
	})
}
