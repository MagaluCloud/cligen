# CLI Generator (cligen)

[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/magaluCloud/cligen)
[![CI](https://github.com/magaluCloud/cligen/workflows/CI/badge.svg)](https://github.com/magaluCloud/cligen/actions/workflows/ci.yml)
[![Release](https://github.com/magaluCloud/cligen/workflows/Release/badge.svg)](https://github.com/magaluCloud/cligen/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/magaluCloud/cligen)](https://goreportcard.com/report/github.com/magaluCloud/cligen)

> ⚠️ **PRODUCTION USE WARNING / AVISO DE USO EM PRODUÇÃO** ⚠️
> 
> Production use is NOT recommended yet.
> 
> O uso em produção ainda NÃO é recomendado.

Um gerador de código que cria automaticamente o código fonte da CLI baseado no SDK do MagaluCloud.

## 📋 Descrição

O **cligen** é uma ferramenta de linha de comando desenvolvida em Go que automatiza a criação de CLIs (Command Line Interfaces) baseadas no SDK do MagaluCloud. Ele gera código estruturado e funcional a partir de configurações YAML, facilitando o desenvolvimento de interfaces de linha de comando para serviços cloud.

## 🚀 Funcionalidades

- **Clone do SDK**: Baixa automaticamente a versão mais recente do SDK do MagaluCloud
- **Geração de Código**: Cria código Go estruturado baseado em templates
- **Base CLI**: Gera uma estrutura base de CLI usando Cobra
- **Configuração Flexível**: Suporte a configurações via YAML para personalização

## 📦 Pré-requisitos

- Go 1.25.3 ou superior
- Git
- Make
- Acesso à internet (para clonar o SDK)

## 🛠️ Instalação

### Instalação via Release (Recomendado)

Baixe o binário pré-compilado para sua plataforma na [página de releases](https://github.com/magaluCloud/cligen/releases):

```bash
# Linux/macOS
curl -L -o cligen https://github.com/magaluCloud/cligen/releases/latest/download/cligen-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m)
chmod +x cligen
sudo mv cligen /usr/local/bin/
```

### Instalação via Source

1. Clone o repositório:
```bash
git clone https://github.com/magaluCloud/cligen.git
cd cligen
```

2. Instale as dependências e compile:
```bash
make install
```

Ou manualmente:
```bash
go mod tidy
make build
```

## 🎯 Uso

### Comandos Disponíveis

#### 1. Clone do SDK
```bash
./cligen clone-sdk
```
Baixa o SDK do MagaluCloud para o diretório `tmp-sdk/` com a versão configurada em `config/config.yaml`.

#### 2. Geração de Código da CLI
```bash
./cligen gen-cli-code
```
Gera o código da CLI baseado nas configurações do arquivo `config/config.yaml`.

#### 3. Geração da Base da CLI
```bash
./cligen gen-cli-base
```
Copia os arquivos base da CLI para o diretório `tmp-cli/`.

#### 4. Execução Completa
```bash
make new-run
```
Executa todo o fluxo: clone do SDK, geração de código e base da CLI.

### Fluxo de Trabalho Recomendado

1. **Primeira execução** (com clone do SDK):
```bash
make new-run
```

2. **Execuções subsequentes** (apenas geração):
```bash
make run
```

3. **Testar a CLI gerada**:
```bash
make run-cli
```

4. **Limpar arquivos temporários**:
```bash
make clean
```

## ⚙️ Configuração

O arquivo `config/config.yaml` define a estrutura da CLI gerada:

```yaml
version: v0.3.42
menus:
  - name: "profile"
    menus:
      - name: "availability-zones"
        alias: ["azs", "az"]
        sdk_package: "availabilityzones"
      
      - name: "ssh-keys"
        alias: ["ssh-keys", "ssh"]
        sdk_package: "sshkeys"

  - name: "virtual-machine"
    alias: ["vm", "virtual-machines", "vms", "compute"]
    sdk_package: "compute"
```

### Estrutura da Configuração

- **version**: Versão do SDK a ser utilizada
- **menus**: Lista de menus/comandos da CLI
  - **name**: Nome do comando
  - **alias**: Aliases alternativos para o comando
  - **sdk_package**: Pacote do SDK correspondente
  - **menus**: Subcomandos aninhados

## 📁 Estrutura do Projeto

```
cligen/
├── base-cli/           # Código base da CLI
├── base-cli-gen/       # Código gerado
├── commands/           # Comandos do cligen
│   ├── gen_cli_code/   # Lógica de geração de código
│   └── sdk_structure/  # Estrutura do SDK
├── config/             # Configurações
├── cobra_utils/        # Utilitários do Cobra
├── str_utils/          # Utilitários de string
├── main.go             # Ponto de entrada
├── Makefile            # Comandos de build
└── README.md           # Este arquivo
```

## 🔧 Desenvolvimento

### Estrutura de Comandos

O projeto utiliza o framework Cobra para gerenciar os comandos:

- `clone-sdk`: Clona o SDK do MagaluCloud
- `gen-cli-code`: Gera o código da CLI
- `gen-cli-base`: Gera a base da CLI
- `gen-cli-sdk-structure`: Analisa a estrutura do SDK

### Templates

Os templates de código estão localizados em `commands/gen_cli_code/` e incluem:

- Templates para comandos raiz
- Templates para grupos de serviços
- Templates para pacotes
- Templates para produtos

## 🧪 Testando

Para testar a CLI gerada:

```bash
make run-cli
```

Isso irá:
1. Compilar a CLI gerada
2. Executar a CLI com os comandos disponíveis

## 🧹 Limpeza

Para limpar arquivos temporários:

```bash
make clean
```

Remove os diretórios:
- `tmp-cli/`
- `tmp-sdk/`
- `base-cli-gen/`
- Executável `cligen`

## 📝 Dependências

- **github.com/spf13/cobra**: Framework para criação de CLIs
- **github.com/MagaluCloud/mgc-sdk-go**: SDK do MagaluCloud
- **gopkg.in/yaml.v3**: Parser YAML

## 🔨 Makefile

O projeto inclui um Makefile profissional com diversos comandos úteis. Para ver todos os comandos disponíveis:

```bash
make help
```

### Comandos Principais

| Comando | Descrição |
|---------|-----------|
| `make build` | Compila o binário |
| `make test` | Executa os testes |
| `make check` | Executa fmt, vet e test |
| `make clean` | Remove arquivos temporários |
| `make ci-full` | Pipeline completo de CI |
| `make ci-release` | Prepara release multi-plataforma |


## 🔄 CI/CD

O projeto utiliza GitHub Actions para CI/CD:

### Workflows Disponíveis

1. **CI** (`.github/workflows/ci.yml`): Executado em cada push/PR
   - Lint e formatação
   - Testes com cobertura
   - Build para verificação

2. **Release** (`.github/workflows/release.yml`): Executado em tags
   - Build multi-plataforma
   - Geração de checksums
   - Criação automática de release

### Criando uma Release

Para criar uma nova release:

```bash
# Criar e push da tag
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0

# O GitHub Actions irá automaticamente:
# 1. Compilar para todas as plataformas
# 2. Gerar checksums
# 3. Criar a release no GitHub
```

## 🤝 Contribuição

Ao contribuir para o projeto, siga estas diretrizes:

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Garanta que o código passa nas verificações:
   ```bash
   make check
   ```
4. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
5. Push para a branch (`git push origin feature/AmazingFeature`)
6. Abra um Pull Request

### Padrões de Código

- Use `go fmt` para formatar o código
- Execute `go vet` antes de fazer commit
- Mantenha a cobertura de testes acima de 70%
- Documente funções e pacotes públicos

## 📄 Licença

Este projeto está sob a licença especificada no arquivo [LICENSE](LICENSE).
