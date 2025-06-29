package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

func CloneSDKCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "clone-sdk",
		Short: "Clona o SDK do MagaluCloud",
		Long:  "Clona o repositório github.com/MagaluCloud/mgc-sdk-go no diretório tmp-sdk/",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cloneSDK()
		},
	}

	return cmd
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

	fmt.Printf("Clonando SDK de %s para %s...\n", sdkURL, sdkDir)

	// Executar comando git clone
	cloneCmd := exec.Command("git", "clone", sdkURL, sdkDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr

	if err := cloneCmd.Run(); err != nil {
		return fmt.Errorf("erro ao clonar SDK: %v", err)
	}

	// remove .github inside sdkDir
	os.RemoveAll(filepath.Join(sdkDir, ".github")) // evitar cazzo
	os.RemoveAll(filepath.Join(sdkDir, ".git"))    // evitar cazzo

	fmt.Printf("SDK clonado com sucesso em %s\n", sdkDir)
	return nil
}
