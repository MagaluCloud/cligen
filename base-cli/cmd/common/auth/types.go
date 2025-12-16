package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	TenantIDWithType string `json:"tenant"`
	ScopesStr        string `json:"scope"`
	Email            string `json:"email"`
}

// TokenResponse representa a resposta do servidor de autenticação OAuth
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

// AuthResult representa o resultado de uma tentativa de autenticação
type AuthResult struct {
	Token *TokenResponse
	Error error
}

// LoginOptions configura opções para o processo de login
type LoginOptions struct {
	Headless bool // Login sem abrir navegador
	QRCode   bool // Exibir QR code para login
	Show     bool // Mostrar token de acesso após login
}

// AuthService define a interface para serviços de autenticação
type AuthService interface {
	// Login inicia o fluxo de autenticação OAuth
	Login(ctx context.Context, opts LoginOptions) (*TokenResponse, error)
}

// TemplateData representa os dados passados para o template HTML
type TemplateData struct {
	Title            string
	Lines            []string
	ErrorDescription string
}

// Tenant representa as informações de um tenant
type Tenant struct {
	UUID        string `json:"uuid"`
	Name        string `json:"legal_name"`
	Email       string `json:"email"`
	IsManaged   bool   `json:"is_managed"`
	IsDelegated bool   `json:"is_delegated"`
}

// TenantResult representa o retorno da alteração do tenant
type TenantResult struct {
	AccessToken  string `json:"access_token"`
	CreatedAt    int    `json:"created_at"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"scope_type"`
}

// TokenExchangeResult representa o resultado da troca de token
type TokenExchangeResult struct {
	TenantID     string    `json:"uuid"`
	CreatedAt    time.Time `json:"created_at"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Scope        []string  `json:"scope"`
}

// ApiKeys representa o retorno da requisição de API Keys
type ApiKeys struct {
	UUID          string  `json:"uuid"`
	Name          string  `json:"name"`
	ApiKey        string  `json:"api_key"`
	Description   string  `json:"description"`
	KeyPairID     string  `json:"key_pair_id"`
	KeyPairSecret string  `json:"key_pair_secret"`
	StartValidity string  `json:"start_validity"`
	EndValidity   *string `json:"end_validity,omitempty"`
	RevokedAt     *string `json:"revoked_at,omitempty"`
	TenantName    *string `json:"tenant_name,omitempty"`
	Tenant        struct {
		UUID      string `json:"uuid"`
		LegalName string `json:"legal_name"`
	} `json:"tenant"`
	Scopes []struct {
		UUID        string `json:"uuid"`
		Name        string `json:"name"`
		Title       string `json:"title"`
		ConsentText string `json:"consent_text"`
		Icon        string `json:"icon"`
		APIProduct  struct {
			UUID string `json:"uuid"`
			Name string `json:"name"`
		} `json:"api_product"`
	} `json:"scopes"`
}

type Scope struct {
	Name  string `json:"name"`
	Title string `json:"title"`
	UUID  string `json:"uuid"`
}

type APIProduct struct {
	Name   string  `json:"name"`
	Scopes []Scope `json:"scopes"`
	UUID   string  `json:"uuid"`
}

type Platform struct {
	APIProducts []APIProduct `json:"api_products"`
	Name        string       `json:"name"`
	UUID        string       `json:"uuid"`
}

type PlatformsResponse []Platform

type ScopesCreate struct {
	ID            string `json:"id"`
	RequestReason string `json:"request_reason"`
}

type CreateApiKey struct {
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	TenantID      string         `json:"tenant_id"`
	ScopesList    []ScopesCreate `json:"scopes"`
	StartValidity string         `json:"start_validity"`
	EndValidity   string         `json:"end_validity"`
}

type ApiKeyResult struct {
	UUID string `json:"uuid,omitempty"`
	Used bool   `json:"used,omitempty"`
}

type Clients struct {
	UUID        string `json:"uuid,omitempty"`
	ClientID    string `json:"client_id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
	Scopes      []struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
	} `json:"scopes,omitempty"`
	ScopesDefault []struct {
		UUID string `json:"uuid"`
		Name string `json:"name"`
	} `json:"scopes_default,omitempty"`
	TermOfUse                        string   `json:"client_term_url,omitempty"`
	ClientPrivacyTermUrl             string   `json:"client_privacy_term_url,omitempty"`
	Audience                         []string `json:"audience,omitempty"`
	OidcAudience                     []string `json:"oidc_audience,omitempty"`
	AlwaysRequireLogin               bool     `json:"always_require_login,omitempty"`
	BackchannelLogoutSessionRequired bool     `json:"backchannel_logout_session_required,omitempty"`
	BackchannelLogoutUri             string   `json:"backchannel_logout_uri,omitempty"`
	RefreshTokenCustomExpiresEnabled bool     `json:"refresh_token_custom_expires_enabled,omitempty"`
	RefreshTokenExp                  int      `json:"refresh_token_exp,omitempty"`
	AccessTokenExp                   int      `json:"access_token_exp,omitempty"`
	RedirectURIs                     []string `json:"redirect_uris,omitempty"`
	Icon                             string   `json:"icon,omitempty"`
}

type CreateClientParams struct {
	Name                             string
	Description                      string
	RedirectURIs                     string
	BackchannelLogoutSessionEnabled  *bool
	ClientTermsURL                   string
	ClientPrivacyTermURL             string
	Audiences                        *string
	Email                            *string
	Reason                           *string
	Icon                             *string
	AccessTokenExp                   *int
	AlwaysRequireLogin               *bool
	BackchannelLogoutURI             *string
	OidcAudience                     *string
	RefreshTokenCustomExpiresEnabled *bool
	RefreshTokenExp                  *int
	SupportURL                       *string
	GrantTypes                       *string
}

type createClientScopes struct {
	UUID     string `json:"id"`
	Reason   string `json:"request_reason"`
	Optional bool   `json:"optional"`
}

type CreateClient struct {
	Name                             string               `json:"name" jsonschema:"description=Name of new client,example=Client Name" mgc:"positional"`
	Description                      string               `json:"description" jsonschema:"description=Description of new client,example=Client description" mgc:"positional"`
	Scopes                           []createClientScopes `json:"scopes" jsonschema:"description=List of scopes (separated by space),example=openid profile" mgc:"positional"`
	RedirectURIs                     []string             `json:"redirect_uris" jsonschema:"description=Redirect URIs (separated by space)" mgc:"positional"`
	Icon                             *string              `json:"icon,omitempty" jsonschema:"description=URL for client icon" mgc:"positional"`
	AccessTokenExp                   int                  `json:"access_token_exp,omitempty" jsonschema:"description=Access token expiration (in seconds),example=7200" mgc:"positional"`
	AlwaysRequireLogin               *bool                `json:"always_require_login,omitempty" jsonschema:"description=Must ignore active Magalu ID session and always require login,example=false" mgc:"positional"`
	ClientPrivacyTermUrl             string               `json:"client_privacy_term_url" jsonschema:"description=URL to privacy term" mgc:"positional"`
	ClientTermUrl                    string               `json:"client_term_url" jsonschema:"description=URL to terms of use" mgc:"positional"`
	Audience                         []string             `json:"audience,omitempty" jsonschema:"description=Client audiences (separated by space),example=public" mgc:"positional"`
	BackchannelLogoutSessionEnabled  *bool                `json:"backchannel_logout_session_required,omitempty" jsonschema:"description=Client requires backchannel logout session,example=false" mgc:"positional"`
	BackchannelLogoutUri             *string              `json:"backchannel_logout_uri,omitempty" jsonschema:"description=Backchannel logout URI" mgc:"positional"`
	OidcAudience                     []string             `json:"oidc_audience,omitempty" jsonschema:"description=Audiences for ID token, should be the Client ID values" mgc:"positional"`
	RefreshTokenCustomExpiresEnabled *bool                `json:"refresh_token_custom_expires_enabled,omitempty" jsonschema:"description=Use custom value for refresh token expiration,example=false" mgc:"positional"`
	RefreshTokenExp                  int                  `json:"refresh_token_exp,omitempty" jsonschema:"description=Custom refresh token expiration value (in seconds),example=15778476" mgc:"positional"`
	Reason                           string               `json:"request_reason,omitempty" jsonschema:"description=Note to inform the reason for creating the client. Will help with the application approval process" mgc:"positional"`
	SupportUrl                       *string              `json:"support_url,omitempty" jsonschema:"description=URL for client support" mgc:"positional"`
	GrantTypes                       []string             `json:"grant_types,omitempty" jsonschema:"description=Grant types the client can request for token generation (separated by space)" mgc:"positional"`
	Email                            *string              `json:"email,omitempty" jsonschema:"description=Email for client support" mgc:"positional"`
}

type UpdateClient struct {
	Name                             *string  `json:"name" jsonschema:"description=Name of new client,example=Client Name" mgc:"positional"`
	Description                      *string  `json:"description" jsonschema:"description=Description of new client,example=Client description" mgc:"positional"`
	RedirectURIs                     []string `json:"redirect_uris" jsonschema:"description=Redirect URIs (separated by space)" mgc:"positional"`
	Icon                             *string  `json:"icon,omitempty" jsonschema:"description=URL for client icon" mgc:"positional"`
	AccessTokenExp                   *int     `json:"access_token_exp,omitempty" jsonschema:"description=Access token expiration (in seconds),example=7200" mgc:"positional"`
	AlwaysRequireLogin               *bool    `json:"always_require_login,omitempty" jsonschema:"description=Must ignore active Magalu ID session and always require login,example=false" mgc:"positional"`
	ClientPrivacyTermUrl             *string  `json:"client_privacy_term_url" jsonschema:"description=URL to privacy term" mgc:"positional"`
	ClientTermUrl                    *string  `json:"client_term_url" jsonschema:"description=URL to terms of use" mgc:"positional"`
	Audience                         []string `json:"audience,omitempty" jsonschema:"description=Client audiences (separated by space),example=public" mgc:"positional"`
	OidcAudience                     []string `json:"oidc_audience,omitempty" jsonschema:"description=Audiences for ID token, should be the Client ID values" mgc:"positional"`
	BackchannelLogoutSessionEnabled  *bool    `json:"backchannel_logout_session_required,omitempty" jsonschema:"description=Client requires backchannel logout session,example=false" mgc:"positional"`
	BackchannelLogoutUri             *string  `json:"backchannel_logout_uri,omitempty" jsonschema:"description=Backchannel logout URI" mgc:"positional"`
	RefreshTokenCustomExpiresEnabled *bool    `json:"refresh_token_custom_expires_enabled,omitempty" jsonschema:"description=Use custom value for refresh token expiration,example=false" mgc:"positional"`
	RefreshTokenExp                  *int     `json:"refresh_token_exp,omitempty" jsonschema:"description=Custom refresh token expiration value (in seconds),example=15778476" mgc:"positional"`
	Reason                           *string  `json:"request_reason,omitempty" jsonschema:"description=Note to inform the reason for creating the client. Will help with the application approval process" mgc:"positional"`
	SupportUrl                       *string  `json:"support_url,omitempty" jsonschema:"description=URL for client support" mgc:"positional"`
}

type CreateClientResult struct {
	UUID         string `json:"uuid,omitempty"`
	ClientID     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret,omitempty"`
}

type UpdateClientParams struct {
	ID                               string
	Name                             *string
	Description                      *string
	RedirectURIs                     *string
	BackchannelLogoutSessionEnabled  *bool
	ClientTermsURL                   *string
	ClientPrivacyTermURL             *string
	Audiences                        *string
	Reason                           *string
	Icon                             *string
	AccessTokenExp                   *int
	AlwaysRequireLogin               *bool
	BackchannelLogoutURI             *string
	OidcAudience                     *string
	RefreshTokenCustomExpiresEnabled *bool
	RefreshTokenExp                  *int
	SupportURL                       *string
}

type UpdateClientResult struct {
	UUID     string `json:"uuid,omitempty"`
	ClientID string `json:"client_id,omitempty"`
}
