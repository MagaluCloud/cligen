# GitHub Actions Workflows

Este documento descreve os workflows disponíveis e como utilizá-los.

## Workflows Disponíveis

### 1. CI Workflow (`ci.yml`)

Executado automaticamente em:
- Push para branches `main` e `develop`
- Pull Requests para `main` e `develop`
- Manualmente via `workflow_dispatch`

#### Jobs

1. **Lint**: Verifica formatação e qualidade do código
2. **Test**: Executa testes com cobertura
3. **Build**: Compila o binário
4. **Build Multi-Platform**: Compila para múltiplas plataformas (apenas no `main`)

#### Como Usar

O workflow é automático, mas você pode executá-lo manualmente:

1. Vá para Actions no GitHub
2. Selecione "CI"
3. Clique em "Run workflow"

### 2. Release Workflow (`release.yml`)

Executado quando:
- Uma tag começando com `v` é criada (ex: `v1.0.0`)
- Manualmente via `workflow_dispatch`

#### Processo

1. Compila binários para todas as plataformas
2. Gera checksums SHA256
3. Cria release notes automáticas
4. Publica release no GitHub

#### Como Criar uma Release

```bash
# Certifique-se de estar na branch main
git checkout main
git pull origin main

# Crie a tag
git tag -a v1.0.0 -m "Release v1.0.0"

# Push da tag
git push origin v1.0.0

# O workflow será acionado automaticamente
```

#### Release Manual

1. Vá para Actions no GitHub
2. Selecione "Release"
3. Clique em "Run workflow"
4. Insira a versão (ex: `v1.0.0`)
5. Execute

## Personalizando Workflows

### Usando o Makefile em Seus Workflows

O Makefile foi projetado para ser facilmente usado em workflows customizados:

```yaml
name: Custom Workflow

on:
  push:
    branches: [ feature/* ]

jobs:
  custom-build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
          cache: true
      
      # Usar comandos do Makefile
      - name: Verificar qualidade
        run: make check
      
      - name: Build
        run: make build
      
      - name: Teste personalizado
        run: |
          make run
          ./cligen --version
```

### Exemplo: Deploy Automático

```yaml
name: Deploy

on:
  release:
    types: [published]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Build release
        run: make ci-release
      
      - name: Deploy to server
        run: |
          # Seus comandos de deploy aqui
          scp dist/cligen-linux-amd64 user@server:/path/
```

### Exemplo: Build Noturno

```yaml
name: Nightly Build

on:
  schedule:
    - cron: '0 0 * * *'  # Todo dia à meia-noite
  workflow_dispatch:

jobs:
  nightly:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: '1.21'
      
      - name: Build nightly
        env:
          VERSION: nightly-${{ github.run_number }}
        run: make ci-release
      
      - name: Upload artifacts
        uses: actions/upload-artifact@v4
        with:
          name: nightly-builds
          path: dist/*
          retention-days: 7
```

### Exemplo: Teste em Múltiplas Versões do Go

```yaml
name: Multi-Go Test

on:
  pull_request:

jobs:
  test:
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: ['1.21', '1.22', '1.23']
    
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}
      
      - name: Test with Go ${{ matrix.go-version }}
        run: make ci-test
```

## Variáveis de Ambiente

### Disponíveis nos Workflows

| Variável | Descrição | Exemplo |
|----------|-----------|---------|
| `VERSION` | Versão do build | `v1.0.0` |
| `GOOS` | Sistema operacional | `linux`, `darwin`, `windows` |
| `GOARCH` | Arquitetura | `amd64`, `arm64` |
| `GITHUB_REF` | Referência Git | `refs/tags/v1.0.0` |
| `GITHUB_SHA` | Hash do commit | `abc123...` |

### Usando Variáveis

```yaml
- name: Build com versão customizada
  env:
    VERSION: custom-${{ github.run_number }}
    GOOS: linux
    GOARCH: arm64
  run: make build
```

## Secrets Necessários

Os workflows atuais usam apenas secrets padrão do GitHub:

- `GITHUB_TOKEN`: Automaticamente fornecido pelo GitHub Actions

### Adicionando Secrets Customizados

Se você precisar de secrets adicionais (ex: para deploy):

1. Vá para Settings > Secrets and variables > Actions
2. Clique em "New repository secret"
3. Adicione seu secret

Use no workflow:

```yaml
- name: Deploy
  env:
    API_KEY: ${{ secrets.MY_API_KEY }}
  run: ./deploy.sh
```

## Troubleshooting

### Build Falha

**Problema**: O build falha com erro de dependências

**Solução**:
```yaml
- name: Limpar cache
  run: go clean -cache -modcache -testcache
  
- name: Reinstalar dependências
  run: make deps
```

### Timeout

**Problema**: Job excede o timeout

**Solução**:
```yaml
jobs:
  build:
    timeout-minutes: 30  # Aumentar timeout
```

### Artefatos Não Encontrados

**Problema**: Artefatos não são gerados

**Solução**: Verifique se o diretório existe antes do upload:
```yaml
- name: Upload artifacts
  if: always()  # Upload mesmo se houver falha
  run: |
    ls -la dist/
    # Upload do artefato
```

## Boas Práticas

### 1. Cache de Dependências

Sempre use cache para Go:

```yaml
- uses: actions/setup-go@v5
  with:
    go-version: '1.21'
    cache: true
    cache-dependency-path: |
      **/go.mod
      **/go.sum
```

### 2. Versionamento

Use sempre versionamento semântico:
- `v1.0.0` - Release estável
- `v1.0.0-rc1` - Release candidate
- `v1.0.0-beta1` - Beta
- `v1.0.0-alpha1` - Alpha

### 3. Build Condicional

Execute builds pesados apenas quando necessário:

```yaml
- name: Build multi-platform
  if: github.ref == 'refs/heads/main'
  run: make build-all
```

### 4. Paralelização

Use jobs paralelos quando possível:

```yaml
jobs:
  lint:
    runs-on: ubuntu-latest
    # ...
  
  test:
    runs-on: ubuntu-latest
    # ...
  
  build:
    needs: [lint, test]  # Espera lint e test
    runs-on: ubuntu-latest
    # ...
```

## Monitoramento

### Badges

Adicione badges ao README:

```markdown
[![CI](https://github.com/magaluCloud/cligen/workflows/CI/badge.svg)](https://github.com/magaluCloud/cligen/actions/workflows/ci.yml)
[![Release](https://github.com/magaluCloud/cligen/workflows/Release/badge.svg)](https://github.com/magaluCloud/cligen/actions/workflows/release.yml)
```

### Notificações

Configure notificações no GitHub:
1. Vá para Settings > Notifications
2. Configure para receber notificações de workflow failures

## Recursos Adicionais

- [GitHub Actions Documentation](https://docs.github.com/actions)
- [Go Setup Action](https://github.com/actions/setup-go)
- [Makefile Documentation](../MAKEFILE.md)

## Suporte

Para problemas com workflows:
1. Verifique os logs no GitHub Actions
2. Execute localmente: `make ci-full`
3. Abra uma [issue](https://github.com/magaluCloud/cligen/issues)

