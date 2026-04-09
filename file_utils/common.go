package file_utils

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func WriteTemplateToFile[T any](tmpl *template.Template, data T, filePath string) error {
	buf := bytes.NewBuffer(nil)
	if err := tmpl.Execute(buf, data); err != nil {
		return fmt.Errorf("erro ao executar template: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório: %w", err)
	}

	if err := os.WriteFile(filePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("erro ao escrever arquivo: %w", err)
	}

	return nil
}

func ToRelativePath(absPath string) string {
	wd, err := os.Getwd()
	if err != nil {
		return absPath
	}

	rel, err := filepath.Rel(filepath.Dir(wd), absPath)
	if err != nil {
		return absPath
	}

	return rel
}
