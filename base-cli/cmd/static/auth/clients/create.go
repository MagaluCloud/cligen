package clients

import (
	"context"
	"fmt"
	"os"

	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/common/auth"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type CreateOptions struct {
	Name                             string
	Description                      string
	RedirectURIs                     string
	BackchannelLogoutSessionEnabled  bool
	ClientTermsURL                   string
	ClientPrivacyTermURL             string
	Audiences                        string
	Email                            string
	Reason                           string
	Icon                             string
	AccessTokenExp                   int
	AlwaysRequireLogin               bool
	BackchannelLogoutURI             string
	OidcAudience                     string
	RefreshTokenCustomExpiresEnabled bool
	RefreshTokenExp                  int
	SupportURL                       string
	GrantTypes                       string
}

// CreateCommand cria o comando de criar um novo cliente
func CreateCommand(ctx context.Context) *cobra.Command {
	var opts CreateOptions
	var params auth.CreateClientParams

	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "create",
		Short: manager.T("cli.auth.clients.create.short"),
		Long:  manager.T("cli.auth.clients.create.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			params = auth.CreateClientParams{
				Name:                 opts.Name,
				Description:          opts.Description,
				ClientTermsURL:       opts.ClientTermsURL,
				ClientPrivacyTermURL: opts.ClientPrivacyTermURL,
				RedirectURIs:         opts.RedirectURIs,
				AccessTokenExp:       &opts.AccessTokenExp,
				RefreshTokenExp:      &opts.RefreshTokenExp,
			}

			NilIfNotChanged(cmd, BackchannelLogoutSessionEnabled, &params.BackchannelLogoutSessionEnabled, opts.BackchannelLogoutSessionEnabled)
			NilIfNotChanged(cmd, Audiences, &params.Audiences, opts.Audiences)
			NilIfNotChanged(cmd, Reason, &params.Reason, opts.Reason)
			NilIfNotChanged(cmd, Icon, &params.Icon, opts.Icon)
			NilIfNotChanged(cmd, AlwaysRequireLogin, &params.AlwaysRequireLogin, opts.AlwaysRequireLogin)
			NilIfNotChanged(cmd, BackchannelLogoutURI, &params.BackchannelLogoutURI, opts.BackchannelLogoutURI)
			NilIfNotChanged(cmd, OidcAudience, &params.OidcAudience, opts.OidcAudience)
			NilIfNotChanged(cmd, RefreshTokenCustomExpiresEnabled, &params.RefreshTokenCustomExpiresEnabled, opts.RefreshTokenCustomExpiresEnabled)
			NilIfNotChanged(cmd, SupportURL, &params.SupportURL, opts.SupportURL)
			NilIfNotChanged(cmd, Email, &params.Email, opts.Email)
			NilIfNotChanged(cmd, GrantTypes, &params.GrantTypes, opts.GrantTypes)

			return runCreate(ctx, params, raw)
		},
	}

	cmd.Flags().StringVar(&opts.Name, Name, "", manager.T("cli.auth.clients.create.name"))
	cmd.Flags().StringVar(&opts.Description, Description, "", manager.T("cli.auth.clients.create.description"))
	cmd.Flags().StringVar(&opts.RedirectURIs, RedirectURIs, "", manager.T("cli.auth.clients.create.redirect_uris"))
	cmd.Flags().BoolVar(&opts.BackchannelLogoutSessionEnabled, BackchannelLogoutSessionEnabled, false, manager.T("cli.auth.clients.create.backchannel_logout_session"))
	cmd.Flags().StringVar(&opts.ClientTermsURL, ClientTermsURL, "", manager.T("cli.auth.clients.create.client_term_url"))
	cmd.Flags().StringVar(&opts.ClientPrivacyTermURL, ClientPrivacyTermURL, "", manager.T("cli.auth.clients.create.client_privacy_term_url"))
	cmd.Flags().StringVar(&opts.Audiences, Audiences, "", manager.T("cli.auth.clients.create.audiences"))
	cmd.Flags().StringVar(&opts.Email, Email, "", manager.T("cli.auth.clients.create.email"))
	cmd.Flags().StringVar(&opts.Reason, Reason, "", manager.T("cli.auth.clients.create.request_reason"))
	cmd.Flags().StringVar(&opts.Icon, Icon, "", manager.T("cli.auth.clients.create.icon"))
	cmd.Flags().IntVar(&opts.AccessTokenExp, AccessTokenExp, 7200, manager.T("cli.auth.clients.create.access_token_expiration"))
	cmd.Flags().BoolVar(&opts.AlwaysRequireLogin, AlwaysRequireLogin, false, manager.T("cli.auth.clients.create.always_require_login"))
	cmd.Flags().StringVar(&opts.BackchannelLogoutURI, BackchannelLogoutURI, "", manager.T("cli.auth.clients.create.backchannel_logout_uri"))
	cmd.Flags().StringVar(&opts.OidcAudience, OidcAudience, "", manager.T("cli.auth.clients.create.oidc_audience"))
	cmd.Flags().BoolVar(&opts.RefreshTokenCustomExpiresEnabled, RefreshTokenCustomExpiresEnabled, false, manager.T("cli.auth.clients.create.refresh_token_custom_expires_enabled"))
	cmd.Flags().IntVar(&opts.RefreshTokenExp, RefreshTokenExp, 15552000, manager.T("cli.auth.clients.create.refresh_token_exp"))
	cmd.Flags().StringVar(&opts.SupportURL, SupportURL, "", manager.T("cli.auth.clients.create.support_url"))
	cmd.Flags().StringVar(&opts.GrantTypes, GrantTypes, "", manager.T("cli.auth.clients.create.grant_types"))

	cmd.MarkFlagRequired(Name)
	cmd.MarkFlagRequired(ClientTermsURL)
	cmd.MarkFlagRequired(ClientPrivacyTermURL)
	cmd.MarkFlagRequired(Description)
	cmd.MarkFlagRequired(RedirectURIs)

	return cmd
}

// runCreate executa o processo de criar um novo cliente
func runCreate(ctx context.Context, opts auth.CreateClientParams, rawMode bool) error {
	auth := ctx.Value(cmdutils.CTX_AUTH_KEY).(auth.Auth)

	if opts.BackchannelLogoutSessionEnabled != nil && *opts.BackchannelLogoutSessionEnabled && opts.BackchannelLogoutURI == nil {
		return fmt.Errorf("backchannel-logout-uri is required when backchannel-logout-session is true")
	}

	result, err := auth.CreateClient(ctx, opts)
	if err != nil {
		return fmt.Errorf("erro ao criar o cliente: %w", err)
	}

	fmt.Fprintln(os.Stderr, "Client created successfully! We'll analise your requisition and approve your client. You can check the approval status using client list command.")
	beautiful.NewOutput(rawMode).PrintData(result)

	return nil
}
