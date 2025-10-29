package auth

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	_ "embed"

	"github.com/pkg/browser"
	"github.com/spf13/cobra"
)

var (
	//go:embed html.template
	htmlTemplate string

	htmlTempl *template.Template
)

func Login(ctx context.Context) *cobra.Command {
	var headless bool
	var qrcode bool
	var show bool

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to the Magalu Cloud",
		Long:  `Login to the Magalu Cloud`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Login to the Magalu Cloud")
			err := defaultLogin(ctx)
			if err != nil {
				fmt.Println("Error: %v", err)
			}
		},
	}

	cmd.Flags().BoolVar(&headless, "headless", false, "Login to the Magalu Cloud without a browser")
	cmd.Flags().BoolVar(&qrcode, "qrcode", false, "Login to the Magalu Cloud with a QR code")
	cmd.Flags().BoolVar(&show, "show", false, "Show the access token after the login process")
	return cmd
}

// default login with browser
type authResult struct {
	value string
	err   error
}

func defaultLogin(ctx context.Context) error {

	htmlTempl, err := template.New("html").Parse(htmlTemplate)
	if err != nil {
		return err
	}
	fmt.Println("defaultLogin")
	auth := &auth{
		httpClient:   &http.Client{},
		codeVerifier: newVerifier(),
		htmlTempl:    htmlTempl,
	}
	resultChan, cancel, err := startCallbackServer(ctx, false, auth)
	if err != nil {
		return err
	}
	defer cancel()

	authUrl, err := auth.codeChallengeToURL([]string{"mke.write", "api-consulta.read", "openid", "mcr.read", "dbaas.write",
		"cpo:read", "cpo:write", "evt:event-tr", "network.read", "network.write", "object-storage.write", "object-storage.read",
		"block-storage.read", "block-storage.write", "mke.read", "virtual-machine.read", "virtual-machine.write", "dbaas.read",
		"mcr.write", "gdb:ssh-pkey-r", "gdb:ssh-pkey-w", "pa:sa:manage", "lba.loadbalancer.read", "lba.loadbalancer.write", "gdb:azs-r",
		"lbaas.read", "lbaas.write", "iam:read", "iam:write"})
	if err != nil {
		return err
	}
	fmt.Println(authUrl)

	browser.OpenURL(authUrl.String())
	result := <-resultChan
	if result.err != nil {
		return result.err
	}

	fmt.Println(result.value)
	return nil
}

// qrcode login
func qrcodeLogin(ctx context.Context) {
	fmt.Println("qrcodeLogin")
}

// headless login
func headlessLogin(ctx context.Context) {
	fmt.Println("headlessLogin")
}
