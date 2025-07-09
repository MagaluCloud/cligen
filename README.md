# CLI Generator (cligen)

> ⚠️ **WARNING / AVISO** ⚠️
> 
> This is NOT an official Magalu Cloud project. This is a community-driven tool that generates CLI code based on the Magalu Cloud SDK.
> 
> Este NÃO é um projeto oficial da Magalu Cloud. Esta é uma ferramenta desenvolvida pela comunidade que gera código de CLI baseado no SDK da Magalu Cloud.

> 🚫 **PRODUCTION USE WARNING / AVISO DE USO EM PRODUÇÃO** 🚫
> 
> This project is intended for development and learning purposes only. Production use is NOT recommended yet.
> 
> Este projeto é destinado apenas para fins de desenvolvimento e aprendizado. O uso em produção ainda NÃO é recomendado.

Um gerador de código que cria automaticamente o código fonte da CLI baseado no SDK do MagaluCloud.

## 📋 Descrição

O **cligen** é uma ferramenta de linha de comando desenvolvida em Go que automatiza a criação de CLIs (Command Line Interfaces) baseadas no SDK do MagaluCloud. Ele gera código estruturado e funcional a partir de configurações YAML, facilitando o desenvolvimento de interfaces de linha de comando para serviços cloud.

## 🚀 Funcionalidades

- **Clone do SDK**: Baixa automaticamente a versão mais recente do SDK do MagaluCloud
- **Geração de Código**: Cria código Go estruturado baseado em templates
- **Base CLI**: Gera uma estrutura base de CLI usando Cobra
- **Configuração Flexível**: Suporte a configurações via YAML para personalização

## 📦 Pré-requisitos

- Go 1.24.2 ou superior
- Git
- Acesso à internet (para clonar o SDK)

## 🛠️ Instalação

1. Clone o repositório:
```bash
git clone https://github.com/geffersonFerraz/cligen.git
cd cligen
```

2. Instale as dependências:
```bash
go mod tidy
```

3. Compile o projeto:
```bash
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

## 🤝 Contribuição

1. Faça um fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request
