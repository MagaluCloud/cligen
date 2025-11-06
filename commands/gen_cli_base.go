package commands

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	baseCLIDir     = "base-cli"
	baseCLIGenDir  = "base-cli-gen"
	destCLIDir     = "tmp-cli"
	dirPermissions = 0755
)

func GenCLICmd() *cobra.Command {
	return &cobra.Command{
		Use:   "gen-cli-base",
		Short: "Gerar a base da CLI",
		Run: func(cmd *cobra.Command, args []string) {
			if err := genCliBase(); err != nil {
				log.Fatalf("Erro ao gerar base da CLI: %v", err)
			}
		},
	}
}

func genCliBase() error {
	srcDirs := []string{baseCLIDir, baseCLIGenDir}

	if err := ensureDestDirReady(); err != nil {
		return err
	}

	for _, srcDir := range srcDirs {
		if err := verifySourceDir(srcDir); err != nil {
			return err
		}
		if err := copyDirectory(srcDir, destCLIDir); err != nil {
			return fmt.Errorf("erro ao copiar %s: %w", srcDir, err)
		}
	}

	return nil
}

func ensureDestDirReady() error {
	if _, err := os.Stat(destCLIDir); err == nil {
		if err := os.RemoveAll(destCLIDir); err != nil {
			return fmt.Errorf("erro ao remover diretório de destino: %w", err)
		}
	}
	return nil
}

func verifySourceDir(srcDir string) error {
	info, err := os.Stat(srcDir)
	if err != nil {
		return fmt.Errorf("diretório de origem não encontrado: %s: %w", srcDir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s não é um diretório", srcDir)
	}
	return nil
}

func copyDirectory(srcDir, dstDir string) error {
	return filepath.WalkDir(srcDir, func(srcPath string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, srcPath)
		if err != nil {
			return fmt.Errorf("erro ao calcular caminho relativo: %w", err)
		}

		dstPath := filepath.Join(dstDir, relPath)

		if d.IsDir() {
			if err := os.MkdirAll(dstPath, dirPermissions); err != nil {
				return fmt.Errorf("erro ao criar diretório %s: %w", dstPath, err)
			}
			return nil
		}

		return copyFile(srcPath, dstPath)
	})
}

func copyFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo de origem %s: %w", srcPath, err)
	}
	defer src.Close()

	if err := os.MkdirAll(filepath.Dir(dstPath), dirPermissions); err != nil {
		return fmt.Errorf("erro ao criar diretório pai para %s: %w", dstPath, err)
	}

	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("erro ao criar arquivo de destino %s: %w", dstPath, err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("erro ao copiar conteúdo de %s para %s: %w", srcPath, dstPath, err)
	}

	return nil
}
