package clients

import (
	"context"
	"fmt"

	"github.com/magaluCloud/mgccli/beautiful"
	"github.com/magaluCloud/mgccli/cmd/common/auth"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
	"github.com/magaluCloud/mgccli/i18n"
	"github.com/spf13/cobra"
)

type ListClientResult struct {
	UUID                             string   `json:"uuid,omitempty"`
	ClientID                         string   `json:"client_id,omitempty"`
	Name                             string   `json:"name,omitempty"`
	Description                      string   `json:"description,omitempty"`
	Status                           string   `json:"client_approval_status,omitempty"`
	Scopes                           []string `json:"scopes,omitempty"`
	ScopesDefault                    []string `json:"scopes_default,omitempty"`
	TermOfUse                        string   `json:"term_of_use,omitempty"`
	ClientPrivacyTermUrl             string   `json:"client_privacy_term_url,omitempty"`
	Audiences                        []string `json:"audiences,omitempty"`
	OidcAudiences                    []string `json:"oidc_audience,omitempty"`
	AlwaysRequireLogin               bool     `json:"always_require_login"`
	BackchannelLogoutSessionEnabled  bool     `json:"backchannel_logout_session_enabled"`
	BackchannelLogoutUri             string   `json:"backchannel_logout_uri,omitempty"`
	RefreshTokenCustomExpiresEnabled bool     `json:"refresh_token_custom_expires_enabled"`
	RefreshTokenExp                  int      `json:"refresh_token_expiration,omitempty"`
	AccessTokenExp                   int      `json:"access_token_expiration,omitempty"`
	RedirectURIs                     []string `json:"redirect_uris,omitempty"`
	Icon                             string   `json:"icon,omitempty"`
}

// ListCommand cria o comando de listar todos os clientes
func ListCommand(ctx context.Context) *cobra.Command {
	manager := i18n.GetInstance()

	cmd := &cobra.Command{
		Use:   "list",
		Short: manager.T("cli.auth.clients.list.short"),
		Long:  manager.T("cli.auth.clients.list.long"),
		RunE: func(cmd *cobra.Command, args []string) error {
			raw, _ := cmd.Root().PersistentFlags().GetBool("raw")

			return runList(ctx, raw)
		},
	}

	return cmd
}

// runList executa o processo de exibir todos clientes
func runList(ctx context.Context, rawMode bool) error {
	auth := ctx.Value(cmdutils.CTX_AUTH_KEY).(auth.Auth)

	clients, err := auth.ListClients(ctx)
	if err != nil {
		return fmt.Errorf("erro ao listar os clientes: %w", err)
	}

	result := []*ListClientResult{}

	for _, client := range clients {
		clientInfo := &ListClientResult{
			UUID:                             client.UUID,
			ClientID:                         client.ClientID,
			Name:                             client.Name,
			Description:                      client.Description,
			Status:                           client.Status,
			TermOfUse:                        client.TermOfUse,
			ClientPrivacyTermUrl:             client.ClientPrivacyTermUrl,
			Audiences:                        client.Audience,
			AlwaysRequireLogin:               client.AlwaysRequireLogin,
			OidcAudiences:                    client.OidcAudience,
			BackchannelLogoutSessionEnabled:  client.BackchannelLogoutSessionRequired,
			BackchannelLogoutUri:             client.BackchannelLogoutUri,
			RefreshTokenCustomExpiresEnabled: client.RefreshTokenCustomExpiresEnabled,
			RefreshTokenExp:                  client.RefreshTokenExp,
			AccessTokenExp:                   client.AccessTokenExp,
			RedirectURIs:                     client.RedirectURIs,
			Icon:                             client.Icon,
		}

		for _, scope := range client.Scopes {
			clientInfo.Scopes = append(clientInfo.Scopes, scope.Name)
		}
		for _, scope := range client.ScopesDefault {
			clientInfo.ScopesDefault = append(clientInfo.ScopesDefault, scope.Name)
		}

		result = append(result, clientInfo)
	}

	if len(result) == 0 {
		fmt.Println("Nenhum cliente foi criado ainda.")
	} else {
		beautiful.NewOutput(rawMode).PrintData(result)
	}

	return nil
}
