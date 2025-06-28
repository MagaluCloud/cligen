package commands

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func GenCLICmd() *cobra.Command {
	return &cobra.Command{
		Use:   "gen-cli-base",
		Short: "Gerar a base da CLI",
		Run: func(cmd *cobra.Command, args []string) {
			genCliBase()
		},
	}
}

func genCliBase() {
	// Copiar o código de base-cli/* para o diretório tmp-cli/
	srcDir := "base-cli"
	srcGenDir := "base-cli-gen"
	dstDir := "tmp-cli"

	// Remover diretório de destino se existir
	if _, err := os.Stat(dstDir); err == nil {
		os.RemoveAll(dstDir)
	}

	// Criar diretório de destino se não existir
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		os.MkdirAll(dstDir, 0755)
	}

	runCopyDir(srcDir, dstDir)
	runCopyDir(srcGenDir, dstDir)

}

func runCopyDir(srcDir string, dstDir string) {
	// Copiar arquivos do diretório de origem para o diretório de destino
	files, err := os.ReadDir(srcDir)
	if err != nil {
		log.Fatalf("Erro ao ler diretório de origem: %v", err)
	}

	for _, file := range files {
		srcPath := filepath.Join(srcDir, file.Name())
		dstPath := filepath.Join(dstDir, file.Name())

		if file.IsDir() {
			// Se for um diretório, copiar recursivamente
			copyDir(srcPath, dstPath)
		} else {
			// Se for um arquivo, copiar normalmente
			copyFile(srcPath, dstPath)
		}
	}
}

func copyDir(srcDir, dstDir string) {
	// Criar diretório de destino
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		log.Fatalf("Erro ao criar diretório de destino: %v", err)
	}

	// Ler conteúdo do diretório de origem
	files, err := os.ReadDir(srcDir)
	if err != nil {
		log.Fatalf("Erro ao ler diretório de origem: %v", err)
	}

	// Copiar cada item do diretório
	for _, file := range files {
		srcPath := filepath.Join(srcDir, file.Name())
		dstPath := filepath.Join(dstDir, file.Name())

		if file.IsDir() {
			// Se for um diretório, copiar recursivamente
			copyDir(srcPath, dstPath)
		} else {
			// Se for um arquivo, copiar normalmente
			copyFile(srcPath, dstPath)
		}
	}
}

func copyFile(srcPath, dstPath string) {
	src, err := os.Open(srcPath)
	if err != nil {
		log.Fatalf("Erro ao abrir arquivo de origem: %v", err)
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		log.Fatalf("Erro ao criar arquivo de destino: %v", err)
	}
	defer dst.Close()

	io.Copy(dst, src)
}
