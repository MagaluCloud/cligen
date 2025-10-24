# Variáveis de Ambiente

Este documento descreve as variáveis de ambiente disponíveis e utilizadas no projeto cligen.

## Variáveis de Desenvolvimento

### `CLI_PANIC_OFF`

Desabilita o tratamento de panic na aplicação. Útil durante a fase de desenvolvimento para facilitar a depuração e identificação de problemas.

**Valores aceitos:** qualquer valor não vazio ativa a opção

**Exemplo:**
```bash
export CLI_PANIC_OFF=1
```

---

### `GEN_CUSTOM_FILE`

Habilita a geração de arquivos que permitem a customização da CLI gerada. Quando ativada, o gerador criará arquivos adicionais que podem ser modificados para personalizar o comportamento da CLI.

**Valores aceitos:** qualquer valor não vazio ativa a opção

**Exemplo:**
```bash
export GEN_CUSTOM_FILE=1
```

---

## Variáveis de Execução da CLI

### `CLI_API_KEY`

Informa uma API-KEY da Magalu para utilizar a CLI. Esta chave é utilizada para autenticação nas requisições à API.

**Valores aceitos:** string UUID contendo a API key válida

**Exemplo:**
```bash
export CLI_API_KEY="sua-api-key-aqui"
```

---

### `EXPLORE_JSON`

Habilita o output experimental de JSON. Quando ativada, a CLI utilizará um formato experimental de visualização para saídas JSON.

**Valores aceitos:** qualquer valor não vazio ativa a opção

**Exemplo:**
```bash
export EXPLORE_JSON=1
```

---

## Variáveis de Internacionalização

As seguintes variáveis de ambiente controlam o idioma da interface da CLI. Elas são verificadas em ordem de precedência.

### `CLI_LANG`

Define o idioma específico para a CLI gerada.

**Valores aceitos:** código do idioma (ex: `pt-BR`, `en-US`, `es-ES`)

**Exemplo:**
```bash
export CLI_LANG="pt-BR"
```

---

### `LANG`

Define o idioma do sistema. Esta é uma variável de ambiente padrão do sistema operacional, mas também é utilizada pela CLI quando `CLI_LANG` não está definida.

**Valores aceitos:** código do idioma com encoding (ex: `pt_BR.UTF-8`, `en_US.UTF-8`)

**Exemplo:**
```bash
export LANG="pt_BR.UTF-8"
```

---

### `LC_ALL`

Define o locale completo do sistema, sobrescrevendo todas as outras configurações de locale. Também é utilizada pela CLI para determinar o idioma.

**Valores aceitos:** código do idioma com encoding (ex: `pt_BR.UTF-8`, `en_US.UTF-8`)

**Exemplo:**
```bash
export LC_ALL="pt_BR.UTF-8"
```

---

## Ordem de Precedência

Para as variáveis de idioma, a ordem de verificação geralmente é:

1. `CLI_LANG` (específica da CLI, maior prioridade)
2. `LANG` (padrão do sistema)
3. `LC_ALL` (sobrescreve configurações de locale)

---

## Exemplo de Uso Completo

```bash
# Configuração para desenvolvimento
export CLI_PANIC_OFF=1
export GEN_CUSTOM_FILE=1

# Configuração de idioma
export CLI_LANG="pt-BR"

# Configuração de API
export CLI_API_KEY="sua-api-key-magalu"

# Habilitar recursos experimentais
export EXPLORE_JSON=1

# Executar a CLI
./cli command
```

