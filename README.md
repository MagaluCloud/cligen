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

## Descrição

O **cligen** é uma ferramenta de linha de comando desenvolvida em Go que automatiza a criação de CLIs (Command Line Interfaces) baseadas no SDK do MagaluCloud. Ele gera código estruturado e funcional a partir de configurações JSON, facilitando o desenvolvimento de interfaces de linha de comando para serviços cloud.

## Funcionalidades

- **Clone do SDK**: Baixa automaticamente a versão mais recente do SDK do MagaluCloud
- **Geração Automática de Configuração**: Analisa o SDK e gera automaticamente a estrutura de configuração baseada nos pacotes disponíveis
- **Manipulação Visual de Configuração**: Interface web para editar, reordenar e personalizar a configuração de forma intuitiva
- **Geração de Código**: Cria código Go estruturado baseado em templates modulares
- **Base CLI**: Gera uma estrutura base de CLI usando Cobra
- **Configuração Flexível**: Suporte a configurações via JSON com campos avançados para descrições, aliases, grupos e métodos

## Pré-requisitos

- Go 1.25.3 ou superior
- Git
- Make
- Acesso à internet (para clonar o SDK)

## Instalação

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

## Uso

### Comandos Disponíveis

#### 1. Clone do SDK
```bash
./cligen clone-sdk
```
Baixa o SDK do MagaluCloud para o diretório `tmp-sdk/` com a versão configurada em `config/config.json`.

#### 2. Geração Automática de Configuração
```bash
./cligen generate-config
```
Analisa o SDK clonado e gera automaticamente o arquivo `config/config.json` com todos os pacotes disponíveis. Este comando detecta os métodos, parâmetros e estruturas do SDK, criando uma configuração inicial completa que pode ser personalizada posteriormente.

#### 3. Manipulação Visual da Configuração
```bash
./cligen manipulate-config
```
Inicia um servidor web local (padrão na porta 9080) que permite editar a configuração de forma visual e intuitiva. Através da interface web você pode:
- Visualizar toda a estrutura de menus e métodos
- Editar descrições curtas e longas
- Habilitar ou desabilitar menus e métodos
- Reordenar menus e submenus
- Criar novos menus e grupos
- Mover elementos entre diferentes níveis
- Atualizar parâmetros e configurações de métodos

Acesse `http://localhost:9080` após iniciar o servidor.

#### 4. Geração de Código da CLI
```bash
./cligen gen-cli-code
```
Gera o código da CLI baseado nas configurações do arquivo `config/config.json`. Este comando utiliza a estrutura modular de templates para criar comandos, subcomandos e handlers baseados na configuração.

#### 5. Geração da Base da CLI
```bash
./cligen gen-cli-base
```
Copia os arquivos base da CLI para o diretório `tmp-cli/`.

#### 6. Execução Completa
```bash
make new-run
```
Executa todo o fluxo: clone do SDK, geração de código e base da CLI.

### Fluxo de Trabalho Recomendado

1. **Primeira execução** (com clone do SDK e geração inicial):
```bash
make new-run
```

2. **Gerar ou atualizar configuração automaticamente**:
```bash
./cligen generate-config
```
Este comando analisa o SDK e atualiza o `config/config.json` com novos pacotes e métodos encontrados, preservando suas personalizações existentes.

3. **Personalizar a configuração** (opcional):
```bash
./cligen manipulate-config
```
Use a interface web para ajustar descrições, reordenar menus, habilitar ou desabilitar funcionalidades, e fazer outras personalizações.

4. **Execuções subsequentes** (apenas geração):
```bash
make run
```

5. **Testar a CLI gerada**:
```bash
make run-cli
```

6. **Limpar arquivos temporários**:
```bash
make clean
```

## Configuração

O arquivo `config/config.json` define a estrutura completa da CLI gerada. Este arquivo pode ser criado manualmente ou gerado automaticamente através do comando `generate-config`.

### Exemplo de Estrutura

```json
{
  "cli_version": "v1.0.0",
  "sdk_branch": "",
  "sdk_tag": "v0.3.42",
  "tag_or_branch_or_latest": "tag",
  "show_git_error": false,
  "show_logs": false,
  "menus": [
    {
      "id": "uuid-gerado",
      "name": "profile",
      "enabled": true,
      "description": "Comandos de perfil",
      "long_description": "Gerencia configurações e perfis do usuário",
      "is_group": true,
      "menus": [
        {
          "id": "uuid-gerado",
          "name": "availability-zones",
          "enabled": true,
          "description": "Zonas de disponibilidade",
          "long_description": "Lista e gerencia zonas de disponibilidade",
          "alias": ["azs", "az"],
          "sdk_package": "availabilityzones",
          "methods": [
            {
              "name": "List",
              "description": "Lista zonas de disponibilidade",
              "long_description": "Retorna todas as zonas de disponibilidade disponíveis na região",
              "parameters": [...],
              "returns": [...]
            }
          ]
        }
      ]
    }
  ]
}
```

### Estrutura da Configuração

A configuração suporta uma estrutura hierárquica rica com os seguintes campos principais:

**Config (nível raiz):**
- **cli_version**: Versão da CLI gerada
- **sdk_tag**: Tag do SDK a ser utilizada
- **sdk_branch**: Branch do SDK (alternativa à tag)
- **tag_or_branch_or_latest**: Define se usa tag, branch ou latest
- **show_git_error**: Exibir erros do Git durante o clone
- **show_logs**: Exibir logs detalhados durante a geração
- **menus**: Lista de menus/comandos da CLI

**Menu:**
- **id**: Identificador único do menu (gerado automaticamente)
- **sdk_name**: Nome original da entidade (pacote ou método) no SDK
- **cli_name**: Nome do comando na CLI
- **enabled**: Se o menu está habilitado (padrão: true)
- **description**: Descrição curta do comando
- **long_description**: Descrição detalhada do comando
- **alias**: Lista de aliases alternativos para o comando
- **sdk_package**: Pacote do SDK correspondente
- **is_group**: Indica se é um grupo de comandos (sem métodos próprios)
- **cli_group**: Grupo CLI para organização
- **service_interface**: Interface do serviço no SDK
- **sdk_file**: Arquivo específico do SDK a ser usado
- **custom_file**: Arquivo customizado para sobrescrever comportamento
- **parent_menu_id**: ID do menu pai (para submenus)
- **menus**: Subcomandos aninhados
- **methods**: Lista de métodos expostos como comandos

**Method:**
- **name**: Nome do método
- **description**: Descrição curta do método
- **long_description**: Descrição detalhada do método
- **parameters**: Lista de parâmetros do método
- **returns**: Lista de valores de retorno
- **comments**: Comentários adicionais
- **confirmation**: Configuração de confirmação antes da execução
- **is_service**: Se é um método de serviço
- **service_import**: Import do serviço
- **sdk_file**: Arquivo específico do SDK
- **custom_file**: Arquivo customizado

**Parameter:**
- **name**: Nome do parâmetro
- **type**: Tipo do parâmetro
- **description**: Descrição do parâmetro
- **is_primitive**: Se é um tipo primitivo
- **is_pointer**: Se é um ponteiro
- **is_optional**: Se é opcional
- **is_array**: Se é um array
- **is_positional**: Se é um argumento posicional
- **positional_index**: Índice para argumentos posicionais
- **struct**: Estrutura aninhada (para tipos complexos)
- **alias_type**: Tipo alias

A configuração é gerada automaticamente pelo comando `generate-config`, que analisa o SDK e detecta todos os pacotes, métodos e estruturas disponíveis. Você pode então personalizar usando a interface web do `manipulate-config` ou editando o JSON manualmente.

## Estrutura do Projeto

```
cligen/
├── base-cli/                    # Código base da CLI
├── base-cli-gen/                # Código gerado
├── commands/                    # Comandos do cligen
│   └── gen_cli_code/            # Lógica de geração de código
│       ├── code/                # Módulos de geração de código
│       │   ├── gomod/           # Geração do go.mod
│       │   ├── menu/            # Geração de menus/grupos
│       │   ├── menu_item/       # Geração de itens de menu
│       │   ├── module/          # Geração de módulos
│       │   └── root_gen/        # Geração do comando raiz
│       ├── generate_config/     # Geração automática de configuração
│       └── manipulate_config/   # Interface web para manipulação
│           ├── static/          # Arquivos estáticos (CSS, JS)
│           └── templates/       # Templates HTML
├── config/                      # Configurações
│   └── config.json              # Arquivo de configuração principal
├── cobra_utils/                 # Utilitários do Cobra
├── file_utils/                  # Utilitários de arquivo
├── str_utils/                   # Utilitários de string
├── main.go                      # Ponto de entrada
├── Makefile                     # Comandos de build
└── README.md                    # Este arquivo
```

## Desenvolvimento

### Estrutura de Comandos

O projeto utiliza o framework Cobra para gerenciar os comandos:

- `clone-sdk`: Clona o SDK do MagaluCloud para o diretório `tmp-sdk/`
- `generate-config`: Analisa o SDK e gera automaticamente o arquivo `config/config.json`
- `manipulate-config`: Inicia servidor web para edição visual da configuração
- `gen-cli-code`: Gera o código da CLI baseado na configuração
- `gen-cli-base`: Copia os arquivos base da CLI para `tmp-cli/`

### Arquitetura de Geração de Código

A geração de código foi refatorada para uma arquitetura modular, onde cada componente tem sua responsabilidade específica:

- **root_gen**: Gera o comando raiz da CLI com todos os subcomandos
- **menu**: Gera grupos de comandos (menus) que agrupam módulos relacionados
- **module**: Gera módulos individuais que expõem métodos do SDK como comandos
- **menu_item**: Gera itens de menu específicos dentro de módulos
- **gomod**: Gera o arquivo `go.mod` com todas as dependências necessárias

Cada módulo possui seu próprio template e lógica de geração, facilitando a manutenção e extensão do sistema.

### Geração Automática de Configuração

O sistema de geração automática de configuração (`generate_config`) analisa o SDK clonado e:

1. Detecta todos os pacotes disponíveis no SDK
2. Analisa as interfaces e métodos de cada pacote
3. Extrai informações sobre parâmetros, tipos de retorno e estruturas
4. Gera uma configuração JSON completa com todos os menus e métodos encontrados
5. Preserva configurações existentes quando possível

Isso permite que a configuração seja sempre sincronizada com o SDK, facilitando a adição de novos serviços e métodos.

## Testando

Para testar a CLI gerada:

```bash
make run-cli
```

Isso irá:
1. Compilar a CLI gerada
2. Executar a CLI com os comandos disponíveis

## Limpeza

Para limpar arquivos temporários:

```bash
make clean
```

Remove os diretórios:
- `tmp-cli/`
- `tmp-sdk/`
- `base-cli-gen/`
- Executável `cligen`

## Dependências

- **github.com/spf13/cobra**: Framework para criação de CLIs
- **github.com/MagaluCloud/mgc-sdk-go**: SDK do MagaluCloud
- **github.com/gin-gonic/gin**: Framework web para interface de manipulação de configuração
- **golang.org/x/tools/go/packages**: Análise de pacotes Go para geração automática de configuração

## Makefile

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


## CI/CD

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

## Contribuição

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

## Licença

Este projeto está sob a licença especificada no arquivo [LICENSE](LICENSE).
