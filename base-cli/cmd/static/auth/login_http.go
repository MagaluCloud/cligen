package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const serverShutdownTimeout = 500 * time.Millisecond

// const clientID = "mAOEtX-uU-6IH-w6t3FiLbnaj8z0FTtEmLeiB1IJO04" // CLI NOVA
const clientID = "cw9qpaUl2nBiC8PVjNFN5jZeb2vTd_1S5cYs1FhEXh0" // CLI ANTIGA
const redirectURI = "http://localhost:8095/callback"

func startCallbackServer(ctx context.Context, isHeadless bool, auth *auth) (resultChan chan *authResult, cancel func(), err error) {
	// Host includes the port, then listen to specific address + port, ex: "localhost:8095"
	addr := "127.0.0.1:8095"

	if envListenAddr := os.Getenv("MGC_LISTEN_ADDRESS"); envListenAddr != "" {
		addr = envListenAddr
	}

	// Listen so we can fail early on bad address, before starting goroutine
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, nil, err
	}

	resultChan = make(chan *authResult, 1)
	callbackChan := make(chan *authResult, 1)
	cancelChan := make(chan struct{}, 1)

	handler := &callbackHandler{
		auth,
		callbackChan,
		ctx,
	}
	srv := &http.Server{Addr: addr, Handler: handler}

	// serve HTTP until shutdown happened, then report result via channel
	serveAndReportResult := func() {
		serverErrorChan := make(chan error, 1)
		signalChan := make(chan os.Signal, 1)
		serverDoneChan := make(chan *authResult, 1)

		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		waitChannelsAndShutdownServer := func() {
			var result *authResult

			select {
			case err := <-serverErrorChan:
				result = &authResult{err: fmt.Errorf("Could not serve HTTP: %w", err)}

			case sig := <-signalChan:
				result = &authResult{err: fmt.Errorf("Canceled by signal: %v", sig)}

			case <-cancelChan:
				result = &authResult{err: fmt.Errorf("Canceled by user")}

			case result = <-callbackChan:
			}

			signal.Stop(signalChan)

			ctx, cancelShutdown := context.WithTimeout(context.Background(), serverShutdownTimeout)
			defer cancelShutdown()

			// sync: unblocks serveAndReportResult()/srv.Serve()
			if err := srv.Shutdown(ctx); err != nil {
				srv.Close() // aggressively try to close it
			}

			// sync: notify serveAndReportResult() we're done
			serverDoneChan <- result
		}
		go waitChannelsAndShutdownServer()
		if !isHeadless {
			if err := srv.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
				// sync: unblock waitChannelsAndShutdownServer()
				serverErrorChan <- err
			}
		}

		result := <-serverDoneChan // sync: wait server shutdown by waitChannelsAndShutdownServer()

		close(callbackChan)
		close(cancelChan)

		close(serverErrorChan)
		close(signalChan)
		close(serverDoneChan)

		resultChan <- result
	}

	cancel = func() {
		defer func() {
			// serveAndReportResult() will close channels when done.
			// That means there is nothing to cancel and we should do nothing else, just ignore.
			_ = recover()
		}()

		cancelChan <- struct{}{} // exit waitChannelsAndShutdownServer()
		<-resultChan             // wait serveAndReportResult(), discard as results are not meaningful
	}

	go serveAndReportResult()

	return resultChan, cancel, nil
}

type callbackHandler struct {
	auth *auth
	done chan *authResult
	ctx  context.Context
}

func (h *callbackHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url := req.URL

	// Roteamento para as diferentes URLs
	switch url.Path {
	case "/callback": // /callback
		h.handleCallback(w, req)
	case "/term":
		h.handleTermsRedirect(w, req)
	case "/privacy":
		h.handlePrivacyRedirect(w, req)
	default:
		err := fmt.Errorf("Unknown Path: %s", url)
		if err := h.showErrorPage(w, "Unknown Path", err); err != nil {
			fmt.Errorf("could not show error page: %w", err)
		}
	}
}

func (h *callbackHandler) handleCallback(w http.ResponseWriter, req *http.Request) {
	auth := h.auth

	authCode := req.URL.Query().Get("code")
	body, err := auth.requestAuthTokenWithAuthorizationCode(h.ctx, authCode)
	if err != nil {
		if err := h.showErrorPage(w, body, err); err != nil {
			fmt.Errorf("could not show error page: %w", err)
		}
		h.done <- &authResult{err: fmt.Errorf("Could not request auth token with authorization code: %w", err)}
		return
	}

	if err := h.showSuccessPage(w); err != nil {
		fmt.Errorf("could not show whole Succes Page: %w", err)
	}

	h.done <- &authResult{value: authCode}
}

func (h *callbackHandler) handleTermsRedirect(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "https://magalu.cloud/termos-legais/termos-de-uso-magalu-cloud/", http.StatusPermanentRedirect)
}

func (h *callbackHandler) handlePrivacyRedirect(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "https://magalu.cloud/termos-legais/politica-de-privacidade/", http.StatusPermanentRedirect)
}

func (h *callbackHandler) showErrorPage(w http.ResponseWriter, errorDescription string, err error) error {
	w.Header().Set("Content-Type", "text/html")

	data := map[string]interface{}{
		"Title":            "Error",
		"ErrorDescription": errorDescription,
	}

	buf := bytes.NewBuffer(nil)
	if err := h.auth.htmlTempl.Execute(buf, data); err != nil {
		return fmt.Errorf("could not execute template: %w", err)
	}

	if _, err := io.WriteString(w, buf.String()); err != nil {
		return fmt.Errorf("could not write response: %w", err)
	}

	return nil
}
func (h *callbackHandler) showSuccessPage(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "text/html")

	data := map[string]interface{}{
		"Title": "Success",
		"Lines": []string{"You have successfully logged in.", "The page can be closed now."},
	}

	buf := bytes.NewBuffer(nil)
	if err := h.auth.htmlTempl.Execute(buf, data); err != nil {
		return fmt.Errorf("could not execute template: %w", err)
	}

	if _, err := io.WriteString(w, buf.String()); err != nil {
		return fmt.Errorf("could not write response: %w", err)
	}

	return nil
}

type LoginResult struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type auth struct {
	codeVerifier *codeVerifier
	httpClient   *http.Client
	htmlTempl    *template.Template
}

func (a *auth) requestAuthTokenWithAuthorizationCode(ctx context.Context, authCode string) (string, error) {
	if a.codeVerifier == nil {
		return "", fmt.Errorf("no code verification provided, first execute a code challenge request")
	}

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")
	data.Set("code", authCode)
	data.Set("code_verifier", a.codeVerifier.value)

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://id.magalu.com/oauth/token", strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return "", err
	}

	resp, err := a.httpClient.Do(r)

	if err != nil || resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return string(body), fmt.Errorf("could not read response body: %w", err)
		}
		defer resp.Body.Close()

		if err == nil {
			return string(body), fmt.Errorf("bad response from auth server, status %d", resp.StatusCode)
		}
		return string(body), err
	}

	var result LoginResult
	defer resp.Body.Close()
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	fmt.Println(result)
	return "", nil
}

func (a *auth) codeChallengeToURL(scopes []string) (*url.URL, error) {
	loginUrl, err := url.Parse("https://id.magalu.com/login")
	if err != nil {
		return nil, err
	}
	query := loginUrl.Query()
	query.Add("response_type", "code")
	query.Add("client_id", clientID)
	query.Add("redirect_uri", redirectURI)
	query.Add("code_challenge", a.codeVerifier.CodeChallengeS256())
	query.Add("code_challenge_method", "S256")
	query.Add("scope", strings.Join(scopes, " "))
	query.Add("choose_tenants", "true")

	loginUrl.RawQuery = query.Encode()

	return loginUrl, nil
}
