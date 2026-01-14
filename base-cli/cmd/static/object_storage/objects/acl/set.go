package acl

import (
	"context"
	"fmt"
	"os"

	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/static/object_storage/common"
	cobrautils "github.com/magaluCloud/mgccli/cobra_utils/flags"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type setFlags struct {
	Dst        string
	GrantWrite string
	Private    bool
	PublicRead bool
}

type setParams struct {
	Dst        string
	GrantWrite *string
	Private    *bool
	PublicRead *bool
}

func SetCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()
	var flags setFlags
	var params setParams

	cmd := &cobra.Command{
		Use:   "set [dst]",
		Short: manager.T("cli.auth.object_storage.objects.acl.set.short"),
		RunE: func(cmd *cobra.Command, args []string) error {
			if flags.GrantWrite == "help" {
				common.PrintGrantWriteHelp()
				return nil
			}

			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			params.Dst = flags.Dst

			cobrautils.NilIfNotChanged(cmd, "grant-write", &params.GrantWrite, flags.GrantWrite)
			cobrautils.NilIfNotChanged(cmd, "private", &params.Private, flags.Private)
			cobrautils.NilIfNotChanged(cmd, "public-read", &params.PublicRead, flags.PublicRead)

			return runSet(ctx, args, params, raw)
		},
	}

	cmd.Flags().StringVar(&flags.Dst, "dst", "", manager.T("cli.auth.object_storage.objects.dst"))
	cmd.Flags().StringVar(&flags.GrantWrite, "grant-write", "", manager.T("cli.auth.object_storage.buckets.acl.grant_write"))
	cmd.Flags().BoolVar(&flags.Private, "private", false, manager.T("cli.auth.object_storage.buckets.acl.private"))
	cmd.Flags().BoolVar(&flags.PublicRead, "public-read", false, manager.T("cli.auth.object_storage.buckets.acl.public_read"))

	return cmd
}

func runSet(ctx context.Context, args []string, opts setParams, rawMode bool) error {
	path := opts.Dst

	if len(args) > 0 {
		path = args[0]
	}

	if path == "" {
		beautiful.NewOutput(rawMode).PrintError("é necessário fornecer o caminho do objeto como argumento ou usar a flag --dst")

		return nil
	}

	aclOptions := common.Options{
		GrantWrite: opts.GrantWrite,
		Private:    opts.Private,
		PublicRead: opts.PublicRead,
	}

	if common.PermissionsIsEmpty(aclOptions) {
		return fmt.Errorf("needs to pass either grant permissions or canned info")
	}

	err := common.SetACL(ctx, path, aclOptions)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stderr, "✓ ACL definido com sucesso!")

	return nil
}
