package clients

import (
	"context"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/common/auth"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type UpdateOptions struct {
	ID                               string
	Name                             string
	Description                      string
	RedirectURIs                     string
	BackchannelLogoutSessionEnabled  bool
	ClientTermsURL                   string
	ClientPrivacyTermURL             string
	Audiences                        string
	Reason                           string
	Icon                             string
	AccessTokenExp                   int
	AlwaysRequireLogin               bool
	BackchannelLogoutURI             string
	OidcAudience                     string
	RefreshTokenCustomExpiresEnabled bool
	RefreshTokenExp                  int
	SupportURL                       string
}

// UpdateCommand cria o comando de atualizar as informações de um cliente
func UpdateCommand(ctx context.Context) *cobra.Command {
	var opts UpdateOptions
	var params auth.UpdateClientParams

	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:     "update",
		Short:   manager.T("cli.auth.clients.update.short"),
		Long:    manager.T("cli.auth.clients.update.long"),
		Example: `mgc auth clients update --access-token-expiration=7200 --audiences="public" --description="Client description" --name="Client Name" --refresh-token-exp=15552000`,
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			params = auth.UpdateClientParams{
				ID: opts.ID,
			}

			NilIfNotChanged(cmd, Name, &params.Name, opts.Name)
			NilIfNotChanged(cmd, Description, &params.Description, opts.Description)
			NilIfNotChanged(cmd, RedirectURIs, &params.RedirectURIs, opts.RedirectURIs)
			NilIfNotChanged(cmd, BackchannelLogoutSessionEnabled, &params.BackchannelLogoutSessionEnabled, opts.BackchannelLogoutSessionEnabled)
			NilIfNotChanged(cmd, ClientTermsURL, &params.ClientTermsURL, opts.ClientTermsURL)
			NilIfNotChanged(cmd, ClientPrivacyTermURL, &params.ClientPrivacyTermURL, opts.ClientPrivacyTermURL)
			NilIfNotChanged(cmd, Audiences, &params.Audiences, opts.Audiences)
			NilIfNotChanged(cmd, Reason, &params.Reason, opts.Reason)
			NilIfNotChanged(cmd, Icon, &params.Icon, opts.Icon)
			NilIfNotChanged(cmd, AccessTokenExp, &params.AccessTokenExp, opts.AccessTokenExp)
			NilIfNotChanged(cmd, AlwaysRequireLogin, &params.AlwaysRequireLogin, opts.AlwaysRequireLogin)
			NilIfNotChanged(cmd, BackchannelLogoutURI, &params.BackchannelLogoutURI, opts.BackchannelLogoutURI)
			NilIfNotChanged(cmd, OidcAudience, &params.OidcAudience, opts.OidcAudience)
			NilIfNotChanged(cmd, RefreshTokenCustomExpiresEnabled, &params.RefreshTokenCustomExpiresEnabled, opts.RefreshTokenCustomExpiresEnabled)
			NilIfNotChanged(cmd, RefreshTokenExp, &params.RefreshTokenExp, opts.RefreshTokenExp)
			NilIfNotChanged(cmd, SupportURL, &params.SupportURL, opts.SupportURL)

			return runUpdate(ctx, params, raw)
		},
	}

	cmd.Flags().StringVar(&opts.ID, ID, "", manager.T("cli.auth.clients.update.uuid"))
	cmd.Flags().StringVar(&opts.Name, Name, "", manager.T("cli.auth.clients.create.name"))
	cmd.Flags().StringVar(&opts.Description, Description, "", manager.T("cli.auth.clients.create.description"))
	cmd.Flags().StringVar(&opts.RedirectURIs, RedirectURIs, "", manager.T("cli.auth.clients.create.redirect_uris"))
	cmd.Flags().BoolVar(&opts.BackchannelLogoutSessionEnabled, BackchannelLogoutSessionEnabled, false, manager.T("cli.auth.clients.create.backchannel_logout_session"))
	cmd.Flags().StringVar(&opts.ClientTermsURL, ClientTermsURL, "", manager.T("cli.auth.clients.create.client_term_url"))
	cmd.Flags().StringVar(&opts.ClientPrivacyTermURL, ClientPrivacyTermURL, "", manager.T("cli.auth.clients.create.client_privacy_term_url"))
	cmd.Flags().StringVar(&opts.Audiences, Audiences, "", manager.T("cli.auth.clients.create.audiences"))
	cmd.Flags().StringVar(&opts.Reason, Reason, "", manager.T("cli.auth.clients.create.request_reason"))
	cmd.Flags().StringVar(&opts.Icon, Icon, "", manager.T("cli.auth.clients.create.icon"))
	cmd.Flags().IntVar(&opts.AccessTokenExp, AccessTokenExp, 7200, manager.T("cli.auth.clients.create.access_token_expiration"))
	cmd.Flags().BoolVar(&opts.AlwaysRequireLogin, AlwaysRequireLogin, false, manager.T("cli.auth.clients.create.always_require_login"))
	cmd.Flags().StringVar(&opts.BackchannelLogoutURI, BackchannelLogoutURI, "", manager.T("cli.auth.clients.create.backchannel_logout_uri"))
	cmd.Flags().StringVar(&opts.OidcAudience, OidcAudience, "", manager.T("cli.auth.clients.create.oidc_audience"))
	cmd.Flags().BoolVar(&opts.RefreshTokenCustomExpiresEnabled, RefreshTokenCustomExpiresEnabled, false, manager.T("cli.auth.clients.create.refresh_token_custom_expires_enabled"))
	cmd.Flags().IntVar(&opts.RefreshTokenExp, RefreshTokenExp, 15552000, manager.T("cli.auth.clients.create.refresh_token_exp"))
	cmd.Flags().StringVar(&opts.SupportURL, SupportURL, "", manager.T("cli.auth.clients.create.support_url"))

	cmd.MarkFlagRequired(ID)

	return cmd
}

// runUpdate executa o processo de atualizar as informações de um cliente
func runUpdate(ctx context.Context, opts auth.UpdateClientParams, rawMode bool) error {
	var confirm bool
	huh.NewConfirm().Title(fmt.Sprintf("This operation may disable your client %s until updates are approved by the ID Magalu. Do you wish to continue?", opts.ID)).
		Affirmative("Yes").
		Negative("No").Value(&confirm).Run()
	if !confirm {
		return nil
	}

	auth := ctx.Value(cmdutils.CTX_AUTH_KEY).(auth.Auth)

	result, err := auth.UpdateClient(ctx, opts)
	if err != nil {
		return fmt.Errorf("erro ao atualizar o cliente: %w", err)
	}

	beautiful.NewOutput(rawMode).PrintData(result)

	return nil
}
