package auth

import (
	"context"
	"fmt"
	"net/http"

	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
)

type authRoundTripper struct {
	parent http.RoundTripper
	auth   Auth
}

func (rt *authRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	token := rt.auth.GetAccessToken(req.Context())

	if token == "" {
		return nil, fmt.Errorf("não foi possível obter o token de acesso. Você esqueceu de fazer login?")
	}

	req.Header.Set("Authorization", "Bearer "+token)

	return rt.parent.RoundTrip(req)
}

func (c *OAuthClient) AuthenticatedHttpClientFromContext(ctx context.Context) *http.Client {
	auth := ctx.Value(cmdutils.CTX_AUTH_KEY).(Auth)

	transport := c.httpClient.Transport

	if transport == nil {
		transport = http.DefaultTransport
	}

	transport = &authRoundTripper{parent: transport, auth: auth}

	return &http.Client{
		Transport: transport,
		Timeout:   c.httpClient.Timeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects, return an error to use the last response.
			return http.ErrUseLastResponse
		}}
}
