# Makefile para CLIgen - Gerador de CLI a partir de SDK
# Versão: 1.0.0

# ==================================================================================== #
# VARIÁVEIS
# ==================================================================================== #

# Configuração do binário
BINARY_NAME := cligen
BUILD_DIR := .
TMP_CLI_DIR := tmp-cli
TMP_SDK_DIR := tmp-sdk
BASE_CLI_GEN_DIR := base-cli-gen
BASE_CLI_CUSTOM_CMD_DIR := base-cli-custom/cmd

# Configuração do Go
GO := go
GOFLAGS := -v
GOOS ?= $(shell $(GO) env GOOS)
GOARCH ?= $(shell $(GO) env GOARCH)

# Versionamento
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Flags de build
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT) -X main.BuildDate=$(BUILD_DATE)"

# Cores para output (se terminal suporta)
NO_COLOR := \033[0m
OK_COLOR := \033[32;01m
ERROR_COLOR := \033[31;01m
WARN_COLOR := \033[33;01m
INFO_COLOR := \033[36;01m

# ==================================================================================== #
# TARGETS DE AJUDA
# ==================================================================================== #

.PHONY: help
help: ## Exibe esta mensagem de ajuda
	@echo "$(INFO_COLOR)CLIgen - Gerador de CLI a partir de SDK$(NO_COLOR)"
	@echo ""
	@echo "$(OK_COLOR)Uso:$(NO_COLOR)"
	@echo "  make [target]"
	@echo ""
	@echo "$(OK_COLOR)Targets disponíveis:$(NO_COLOR)"
	@awk 'BEGIN {FS = ":.*##"; printf ""} /^[a-zA-Z_-]+:.*?##/ { printf "  $(INFO_COLOR)%-18s$(NO_COLOR) %s\n", $$1, $$2 } /^##@/ { printf "\n$(WARN_COLOR)%s$(NO_COLOR)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

# ==================================================================================== #
##@ Desenvolvimento
# ==================================================================================== #

.PHONY: build
build: ## Compila o binário do cligen
	@echo "$(INFO_COLOR)Compilando $(BINARY_NAME)...$(NO_COLOR)"
	@$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "$(OK_COLOR)✓ Compilação concluída: $(BINARY_NAME)$(NO_COLOR)"

.PHONY: build-all
build-all: ## Compila para múltiplas plataformas
	@echo "$(INFO_COLOR)Compilando para múltiplas plataformas...$(NO_COLOR)"
	@GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	@GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .
	@GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	@GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .
	@GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .
	@echo "$(OK_COLOR)✓ Compilação multi-plataforma concluída$(NO_COLOR)"

.PHONY: install
install: build ## Instala o binário no GOPATH/bin
	@echo "$(INFO_COLOR)Instalando $(BINARY_NAME)...$(NO_COLOR)"
	@$(GO) install $(LDFLAGS) .
	@echo "$(OK_COLOR)✓ Instalação concluída$(NO_COLOR)"

.PHONY: deps
deps: ## Baixa e verifica dependências
	@echo "$(INFO_COLOR)Baixando dependências...$(NO_COLOR)"
	@$(GO) mod download
	@$(GO) mod verify
	@echo "$(OK_COLOR)✓ Dependências verificadas$(NO_COLOR)"

.PHONY: tidy
tidy: ## Organiza go.mod e go.sum
	@echo "$(INFO_COLOR)Organizando dependências...$(NO_COLOR)"
	@$(GO) mod tidy -v
	@echo "$(OK_COLOR)✓ Dependências organizadas$(NO_COLOR)"

# ==================================================================================== #
##@ Execução
# ==================================================================================== #

.PHONY: new-run
new-run: build ## Executa o fluxo completo (clone SDK + geração de código + geração base)
	@echo "$(INFO_COLOR)Executando fluxo completo...$(NO_COLOR)"
	@./$(BINARY_NAME) clone-sdk
	@./$(BINARY_NAME) gen-cli-code
	@./$(BINARY_NAME) gen-cli-base
	@echo "$(OK_COLOR)✓ Fluxo completo executado$(NO_COLOR)"

.PHONY: run
run: build ## Executa o fluxo de geração (código + base)
	@echo "$(INFO_COLOR)Executando geração...$(NO_COLOR)"
	@./$(BINARY_NAME) gen-cli-code
	@./$(BINARY_NAME) gen-cli-base
	@echo "$(OK_COLOR)✓ Geração concluída$(NO_COLOR)"

.PHONY: run-cli
run-cli: ## Compila e executa o CLI gerado
	@echo "$(INFO_COLOR)Compilando CLI gerado...$(NO_COLOR)"
	@cd $(TMP_CLI_DIR) && $(GO) mod tidy
	@bash build_cli.sh
	@echo "$(OK_COLOR)✓ CLI gerado compilado$(NO_COLOR)"

# ==================================================================================== #
##@ Limpeza
# ==================================================================================== #

.PHONY: clean
clean: ## Remove arquivos temporários e binários
	@echo "$(INFO_COLOR)Limpando arquivos temporários...$(NO_COLOR)"
	@rm -rf $(TMP_CLI_DIR)
	@rm -rf $(TMP_SDK_DIR)
	@rm -rf $(BASE_CLI_GEN_DIR)
	@rm -rf $(BASE_CLI_CUSTOM_CMD_DIR)
	@rm -f $(BINARY_NAME)
	@echo "$(OK_COLOR)✓ Limpeza concluída$(NO_COLOR)"

.PHONY: clean-all
clean-all: clean ## Remove todos os arquivos gerados (incluindo builds multi-plataforma)
	@echo "$(INFO_COLOR)Limpando todos os builds...$(NO_COLOR)"
	@rm -f $(BINARY_NAME)-*
	@echo "$(OK_COLOR)✓ Limpeza completa concluída$(NO_COLOR)"

# ==================================================================================== #
##@ Qualidade de Código
# ==================================================================================== #

.PHONY: fmt
fmt: ## Formata o código com go fmt
	@echo "$(INFO_COLOR)Formatando código...$(NO_COLOR)"
	@$(GO) fmt ./commands/... ./config/... ./str_utils/...
	@echo "$(OK_COLOR)✓ Código formatado$(NO_COLOR)"

.PHONY: vet
vet: ## Executa go vet
	@echo "$(INFO_COLOR)Executando go vet...$(NO_COLOR)"
	@$(GO) vet ./commands/... ./config/... ./str_utils/...
	@echo "$(OK_COLOR)✓ Verificação concluída$(NO_COLOR)"

.PHONY: lint
lint: ## Executa golangci-lint (requer instalação)
	@echo "$(INFO_COLOR)Executando linter...$(NO_COLOR)"
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./commands/... ./config/... ./str_utils/...; \
		echo "$(OK_COLOR)✓ Lint concluído$(NO_COLOR)"; \
	else \
		echo "$(WARN_COLOR)⚠ golangci-lint não encontrado. Instale com: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NO_COLOR)"; \
	fi

.PHONY: test
test: ## Executa os testes
	@echo "$(INFO_COLOR)Executando testes...$(NO_COLOR)"
	@$(GO) test -v -race -coverprofile=coverage.out ./commands/... ./config/... ./str_utils/...
	@echo "$(OK_COLOR)✓ Testes concluídos$(NO_COLOR)"

.PHONY: test-coverage
test-coverage: test ## Executa testes e gera relatório de cobertura
	@echo "$(INFO_COLOR)Gerando relatório de cobertura...$(NO_COLOR)"
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(OK_COLOR)✓ Relatório gerado: coverage.html$(NO_COLOR)"

.PHONY: check
check: fmt vet test ## Executa todas as verificações de qualidade (fmt, vet, test)

# ==================================================================================== #
##@ CI/CD
# ==================================================================================== #

.PHONY: ci-deps
ci-deps: ## Instala dependências para CI (sem cache)
	@echo "$(INFO_COLOR)Instalando dependências para CI...$(NO_COLOR)"
	@$(GO) mod download
	@$(GO) mod verify
	@echo "$(OK_COLOR)✓ Dependências instaladas$(NO_COLOR)"

.PHONY: ci-build
ci-build: ci-deps ## Build otimizado para CI
	@echo "$(INFO_COLOR)Build CI...$(NO_COLOR)"
	@$(GO) build $(LDFLAGS) -o $(BINARY_NAME) .
	@echo "$(OK_COLOR)✓ Build CI concluído$(NO_COLOR)"

.PHONY: ci-test
ci-test: ## Testes otimizados para CI
	@echo "$(INFO_COLOR)Executando testes CI...$(NO_COLOR)"
	@$(GO) test -v -race -coverprofile=coverage.out -covermode=atomic ./commands/... ./config/... ./str_utils/...
	@echo "$(OK_COLOR)✓ Testes CI concluídos$(NO_COLOR)"

.PHONY: ci-lint
ci-lint: ## Lint otimizado para CI
	@echo "$(INFO_COLOR)Executando lint CI...$(NO_COLOR)"
	@$(GO) fmt ./commands/... ./config/... ./str_utils/... 
	@$(GO) vet ./commands/... ./config/... ./str_utils/...
	@echo "$(OK_COLOR)✓ Lint CI concluído$(NO_COLOR)"

.PHONY: ci-full
ci-full: ci-lint ci-test ci-build ## Pipeline completo de CI (lint, test, build)

.PHONY: ci-release
ci-release: clean build-all ## Prepara release com binários multi-plataforma
	@echo "$(INFO_COLOR)Preparando release...$(NO_COLOR)"
	@mkdir -p dist
	@mv $(BINARY_NAME)-* dist/ 2>/dev/null || true
	@echo "$(OK_COLOR)✓ Release preparado em dist/$(NO_COLOR)"

# ==================================================================================== #
##@ Informações
# ==================================================================================== #

.PHONY: version
version: ## Exibe informações de versão
	@echo "$(INFO_COLOR)Informações de Versão:$(NO_COLOR)"
	@echo "  Versão:    $(VERSION)"
	@echo "  Commit:    $(GIT_COMMIT)"
	@echo "  Data:      $(BUILD_DATE)"
	@echo "  Go:        $(shell $(GO) version)"
	@echo "  GOOS:      $(GOOS)"
	@echo "  GOARCH:    $(GOARCH)"

.PHONY: info
info: version ## Alias para version

# ==================================================================================== #
##@ Outros
# ==================================================================================== #

.PHONY: verify
verify: ## Verifica se o ambiente está configurado corretamente
	@echo "$(INFO_COLOR)Verificando ambiente...$(NO_COLOR)"
	@echo -n "  Go: "
	@command -v go > /dev/null && echo "$(OK_COLOR)✓$(NO_COLOR)" || echo "$(ERROR_COLOR)✗$(NO_COLOR)"
	@echo -n "  Git: "
	@command -v git > /dev/null && echo "$(OK_COLOR)✓$(NO_COLOR)" || echo "$(ERROR_COLOR)✗$(NO_COLOR)"
	@echo "$(OK_COLOR)✓ Verificação concluída$(NO_COLOR)"

.PHONY: all
all: clean check build ## Executa limpeza, verificações e build
	@echo "$(OK_COLOR)✓ Pipeline completo executado$(NO_COLOR)"
