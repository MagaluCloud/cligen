package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/magaluCloud/cligen/config"

	"github.com/spf13/cobra"
)

const (
	sdkURL        = "https://github.com/MagaluCloud/mgc-sdk-go.git"
	sdkDir        = "tmp-sdk"
	githubAPIBase = "https://api.github.com/repos/MagaluCloud/mgc-sdk-go/releases/tags"
	httpTimeout   = 30 * time.Second
	gitCloneDepth = "1"
)

// GitHubRelease representa a estrutura de resposta da API do GitHub
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	HTMLURL string `json:"html_url"`
}

func CloneSDKCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "clone-sdk",
		Short: "Clona o SDK do MagaluCloud",
		Long:  "Clona o repositório github.com/MagaluCloud/mgc-sdk-go no diretório tmp-sdk/ usando a última versão publicada",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cloneSDK()
		},
	}
}

func getLatestRelease() (*GitHubRelease, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar configuração: %w", err)
	}

	url := fmt.Sprintf("%s/%s", githubAPIBase, cfg.Version)

	client := &http.Client{Timeout: httpTimeout}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao acessar API do GitHub: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro na API do GitHub: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta da API: %w", err)
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta da API: %w", err)
	}

	return &release, nil
}

func cloneSDK() error {
	if err := ensureSDKDirRemoved(); err != nil {
		return err
	}

	fmt.Println("Obtendo a última versão publicada via API do GitHub...")
	release, err := getLatestRelease()
	if err != nil {
		fmt.Printf("Aviso: Não foi possível obter a última versão via API (%v). Clonando branch padrão...\n", err)
		return cloneWithFallback(sdkURL, sdkDir, "")
	}

	fmt.Printf("Última versão encontrada: %s (%s)\n", release.TagName, release.Name)
	fmt.Printf("Clonando SDK versão %s de %s para %s...\n", release.TagName, sdkURL, sdkDir)

	if err := executeGitClone(release.TagName); err != nil {
		fmt.Printf("Aviso: Falha ao clonar tag %s (%v). Tentando branch padrão...\n", release.TagName, err)
		return cloneWithFallback(sdkURL, sdkDir, "")
	}

	if err := cleanupSDKDir(); err != nil {
		return fmt.Errorf("erro ao limpar diretório do SDK: %w", err)
	}

	fmt.Printf("SDK versão %s clonado com sucesso em %s\n", release.TagName, sdkDir)
	return nil
}

func ensureSDKDirRemoved() error {
	if _, err := os.Stat(sdkDir); err == nil {
		fmt.Printf("Diretório %s já existe. Removendo...\n", sdkDir)
		if err := os.RemoveAll(sdkDir); err != nil {
			return fmt.Errorf("erro ao remover diretório existente: %w", err)
		}
	}
	return nil
}

func executeGitClone(branchOrTag string) error {
	cloneCmd := exec.Command("git", "clone", "--depth", gitCloneDepth, "--branch", branchOrTag, sdkURL, sdkDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	return cloneCmd.Run()
}

func cleanupSDKDir() error {
	dirsToRemove := []string{".github", ".git"}
	for _, dir := range dirsToRemove {
		path := filepath.Join(sdkDir, dir)
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("erro ao remover %s: %w", dir, err)
		}
	}
	return nil
}

func cloneWithFallback(sdkURL, sdkDir, preferredBranch string) error {
	branches := []string{"main", "master"}
	if preferredBranch != "" {
		branches = append([]string{preferredBranch}, branches...)
	}

	for _, branch := range branches {
		if branch == "" {
			continue
		}
		fmt.Printf("Tentando clonar branch %s...\n", branch)
		if err := executeGitClone(branch); err == nil {
			if err := cleanupSDKDir(); err != nil {
				return fmt.Errorf("erro ao limpar diretório do SDK: %w", err)
			}
			fmt.Printf("SDK clonado com sucesso da branch %s em %s\n", branch, sdkDir)
			return nil
		}
	}

	return fmt.Errorf("erro ao clonar SDK: todas as branches falharam")
}
