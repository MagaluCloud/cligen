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
	githubAPIBase = "https://api.github.com/repos/MagaluCloud/mgc-sdk-go/releases"
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
			cfg, err := config.LoadConfig()
			if err != nil {
				return fmt.Errorf("erro ao carregar configuração: %w", err)
			}

			g := github{config: cfg}
			return g.cloneSDK()
		},
	}
}

type github struct {
	config *config.Config
}

func (g *github) getLatestRelease() (*GitHubRelease, error) {

	url := fmt.Sprintf("%s/latest", githubAPIBase)

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

func (g *github) cloneSDK() error {
	if err := g.ensureSDKDirRemoved(); err != nil {
		return err
	}

	version := ""
	if g.config.TagOrBranchOrLatest == "branch" {
		if g.config.SDKBranch == "" {
			return fmt.Errorf("branch não configurado")
		}
		version = g.config.SDKBranch
		err := g.executeGitClone(g.config.SDKBranch)
		if err != nil {
			return fmt.Errorf("erro ao clonar branch %s (%v): %w", g.config.SDKBranch, err, err)
		}
	}

	if g.config.TagOrBranchOrLatest == "tag" {
		if g.config.SDKTag == "" {
			return fmt.Errorf("tag não configurado")
		}
		version = g.config.SDKTag
		err := g.executeGitClone(g.config.SDKTag)
		if err != nil {
			return fmt.Errorf("erro ao clonar tag %s (%v): %w", g.config.SDKTag, err, err)
		}
	}

	if g.config.TagOrBranchOrLatest == "latest" {
		fmt.Println("Obtendo a última versão publicada via API do GitHub...")
		release, err := g.getLatestRelease()
		if err != nil {
			return fmt.Errorf("erro ao obter a última versão via API (%v): %w", err, err)
		}

		version = release.TagName
		err = g.executeGitClone(release.TagName)
		if err != nil {
			return fmt.Errorf("erro ao clonar tag %s (%v): %w", release.TagName, err, err)
		}
	}

	if err := g.cleanupSDKDir(); err != nil {
		return fmt.Errorf("erro ao limpar diretório do SDK: %w", err)
	}

	fmt.Printf("SDK versão %s clonado com sucesso em %s\n", version, sdkDir)
	return nil
}

func (g *github) ensureSDKDirRemoved() error {
	if _, err := os.Stat(sdkDir); err == nil {
		fmt.Printf("Diretório %s já existe. Removendo...\n", sdkDir)
		if err := os.RemoveAll(sdkDir); err != nil {
			return fmt.Errorf("erro ao remover diretório existente: %w", err)
		}
	}
	return nil
}

func (g *github) executeGitClone(branchOrTag string) error {
	cloneCmd := exec.Command("git", "clone", "--depth", gitCloneDepth, "--branch", branchOrTag, sdkURL, sdkDir)
	if g.config.ShowGitError {
		cloneCmd.Stdout = os.Stdout
		cloneCmd.Stderr = os.Stderr
	}

	return cloneCmd.Run()
}

func (g *github) cleanupSDKDir() error {
	dirsToRemove := []string{".github", ".git"}
	for _, dir := range dirsToRemove {
		path := filepath.Join(sdkDir, dir)
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("erro ao remover %s: %w", dir, err)
		}
	}
	return nil
}
