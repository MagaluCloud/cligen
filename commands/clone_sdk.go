package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/magaluCloud/cligen/config"

	"github.com/spf13/cobra"
)

// GitHubRelease representa a estrutura de resposta da API do GitHub
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	HTMLURL string `json:"html_url"`
}

func CloneSDKCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "clone-sdk",
		Short: "Clona o SDK do MagaluCloud",
		Long:  "Clona o repositório github.com/MagaluCloud/mgc-sdk-go no diretório tmp-sdk/ usando a última versão publicada",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cloneSDK()
		},
	}

	return cmd
}

func getLatestRelease() (*GitHubRelease, error) {

	config, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("erro ao carregar o arquivo de configuração: %v", err)
	}

	url := fmt.Sprintf("https://api.github.com/repos/MagaluCloud/mgc-sdk-go/releases/tags/%s", config.Version)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erro ao acessar API do GitHub: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("erro na API do GitHub: status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta da API: %v", err)
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta da API: %v", err)
	}

	return &release, nil
}

func cloneSDK() error {
	sdkURL := "https://github.com/MagaluCloud/mgc-sdk-go.git"
	sdkDir := "tmp-sdk"

	// Verificar se o diretório já existe
	if _, err := os.Stat(sdkDir); err == nil {
		fmt.Printf("Diretório %s já existe. Removendo...\n", sdkDir)
		if err := os.RemoveAll(sdkDir); err != nil {
			return fmt.Errorf("erro ao remover diretório existente: %v", err)
		}
	}

	// Obter a última versão publicada via API do GitHub
	fmt.Println("Obtendo a última versão publicada via API do GitHub...")
	release, err := getLatestRelease()
	if err != nil {
		fmt.Printf("Aviso: Não foi possível obter a última versão via API (%v). Clonando branch main...\n", err)
		return cloneMainBranch(sdkURL, sdkDir)
	}

	fmt.Printf("Última versão encontrada: %s (%s)\n", release.TagName, release.Name)
	fmt.Printf("Clonando SDK versão %s de %s para %s...\n", release.TagName, sdkURL, sdkDir)

	// Executar comando git clone com a tag específica
	cloneCmd := exec.Command("git", "clone", "--depth", "1", "--branch", release.TagName, sdkURL, sdkDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr

	if err := cloneCmd.Run(); err != nil {
		fmt.Printf("Aviso: Falha ao clonar tag %s (%v). Tentando branch main...\n", release.TagName, err)
		return cloneMainBranch(sdkURL, sdkDir)
	}

	// remove .github inside sdkDir
	os.RemoveAll(filepath.Join(sdkDir, ".github")) // evitar cazzo
	os.RemoveAll(filepath.Join(sdkDir, ".git"))    // evitar cazzo

	fmt.Printf("SDK versão %s clonado com sucesso em %s\n", release.TagName, sdkDir)
	return nil
}

func cloneMainBranch(sdkURL, sdkDir string) error {
	fmt.Printf("Clonando branch main de %s para %s...\n", sdkURL, sdkDir)

	// Executar comando git clone da branch main
	cloneCmd := exec.Command("git", "clone", "--depth", "1", "--branch", "main", sdkURL, sdkDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr

	if err := cloneCmd.Run(); err != nil {
		// Se main também falhar, tenta master
		fmt.Printf("Falha ao clonar branch main, tentando master...\n")
		cloneCmd = exec.Command("git", "clone", "--depth", "1", "--branch", "master", sdkURL, sdkDir)
		cloneCmd.Stdout = os.Stdout
		cloneCmd.Stderr = os.Stderr

		if err := cloneCmd.Run(); err != nil {
			return fmt.Errorf("erro ao clonar SDK: %v", err)
		}
	}

	// remove .github inside sdkDir
	os.RemoveAll(filepath.Join(sdkDir, ".github")) // evitar cazzo
	os.RemoveAll(filepath.Join(sdkDir, ".git"))    // evitar cazzo

	fmt.Printf("SDK clonado com sucesso em %s\n", sdkDir)
	return nil
}
