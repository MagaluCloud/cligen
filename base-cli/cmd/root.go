package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stoewer/go-strcase"
)

const (
	loggerConfigKey = "logging"
	defaultRegion   = "br-se1"
	apiKeyEnvVar    = "MGC_API_KEY"
)

func normalizeFlagName(f *pflag.FlagSet, name string) pflag.NormalizedName {
	name = strcase.KebabCase(name)
	return pflag.NormalizedName(name)
}

func Execute(version string) (err error) {
	vv := fmt.Sprintf("%s (%s/%s)",
		version,
		runtime.GOOS,
		runtime.GOARCH)

	rootCmd := &cobra.Command{
		Use:     "mgc",
		Version: vv,
		Short:   "Magalu Cloud CLI",
		Long: `
	███╗   ███╗ ██████╗  ██████╗     ██████╗██╗     ██╗
	████╗ ████║██╔════╝ ██╔════╝    ██╔════╝██║     ██║
	██╔████╔██║██║  ███╗██║         ██║     ██║     ██║
	██║╚██╔╝██║██║   ██║██║         ██║     ██║     ██║
	██║ ╚═╝ ██║╚██████╔╝╚██████╗    ╚██████╗███████╗██║
	╚═╝     ╚═╝ ╚═════╝  ╚═════╝     ╚═════╝╚══════╝╚═╝

Magalu Cloud CLI is a command-line interface for the Magalu Cloud.
It allows you to interact with the Magalu Cloud to manage your resources.
`,
		SilenceErrors: true, // ####    Hack: true to avoid panic on error / false to debug error
		SilenceUsage:  true, // ####
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	rootCmd.SetGlobalNormalizationFunc(normalizeFlagName)

	rootCmd.AddGroup(&cobra.Group{
		ID:    "catalog",
		Title: "Products:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "settings",
		Title: "Settings:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "other",
		Title: "Other commands:",
	})
	rootCmd.SetHelpCommandGroupID("other")
	rootCmd.SetCompletionCommandGroupID("other")

	rootCmd.InitDefaultHelpFlag()

	err = rootCmd.Execute()
	if err != nil {
		return err
	}
	return nil
}
