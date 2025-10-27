# Resumo Executivo: Comparação profile

## 📊 Visão Geral

| Métrica | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **Grupos Principais** | 2 | 2 | ✅ |
| **Bugs Visuais** | 0 | 3 tipos ("defaultLongDesc 1", "Dqui1", "doto3") | 🔴 |
| **Nomes de Grupos Divergentes** | - | 1 (`keys` vs `ssh-keys`) | 🔴 |
| **Flags com Descrição** | 100% | 0% | 🔴 |
| **Comandos com Descrição** | 100% | 0% | 🔴 |
| **Aliases Ausentes** | - | 2 grupos | ⚠️ |
| **Flags Globais Ausentes** | - | 5 flags | 🔴 |
| **Flags Extras** | - | 1 flag (`show-blocked`) | ⚠️ |

---

## 📋 Tabela de Comandos Completa

### Grupo: ssh-keys / keys

| Comando | mgc | ./tmp-cli/cli | Divergências | Severidade |
|---------|-----|---------------|--------------|------------|
| **Grupo** | ✅ `ssh-keys` | 🔴 `keys` | **Nome do grupo divergente**; Bug "Dqui1"; 0% descrições | 🔴 |
| **create** | ✅ | ✅ | Bug "doto3"; 0% descrições; Falta tipo UUID; Falta limites de caracteres | 🔴 |
| **delete** | ✅ | ✅ | Bug "doto3"; 0% descrições; Falta tipo UUID | 🔴 |
| **get** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **list** | ✅ | ✅ | `control.limit/offset/sort` vs `limit/offset/sort`; Flags incorretas como required; Bug "doto3" | 🔴 |

**Aliases**: mgc tem (`ssh-keys`, `ssh_keys`), ./tmp-cli/cli não tem ❌

---

### Grupo: availability-zones

| Comando | mgc | ./tmp-cli/cli | Divergências | Severidade |
|---------|-----|---------------|--------------|------------|
| **Grupo** | ✅ `availability-zones` | ✅ `availability-zones` | Bug "Dqui1"; 0% descrições | 🔴 |
| **list** | ✅ | ✅ | Flag extra `--show-blocked` (required); Bug "doto3"; 0% descrições | 🔴 |

**Aliases**: mgc tem (`availability-zones`, `availability_zones`), ./tmp-cli/cli não tem ❌

---

## 🔴 Problemas Críticos Sistemáticos

### 1. Bugs Visuais

| Bug | Localização | Frequência |
|-----|-------------|------------|
| **"defaultLongDesc 1"** | Comando principal `profile` | 1 ocorrência |
| **"Dqui1"** | Grupos `keys` e `availability-zones` | 2 ocorrências |
| **"doto3"** | Todos os comandos específicos | 5 ocorrências |

### 2. Nomes Divergentes

| mgc | ./tmp-cli/cli | Impacto | Severidade |
|-----|---------------|---------|------------|
| `profile ssh-keys` | `profile keys` | **Perda de clareza semântica**. "keys" é ambíguo (API keys? SSH keys? Chaves de criptografia?) | 🔴 CRÍTICO |
| `--control.limit` | `--limit` | Perda de contexto. Prefixo `control.` indica que é um parâmetro de controle de paginação | 🔴 |
| `--control.offset` | `--offset` | Perda de contexto | 🔴 |
| `--control.sort` | `--sort` | Perda de contexto | 🔴 |

### 3. Descrições Ausentes

| Elemento | mgc | ./tmp-cli/cli | Taxa de Perda |
|----------|-----|---------------|---------------|
| **Descrições de Grupos** | 2 descrições completas | 0 (vazias ou ausentes) | 100% |
| **Descrições de Comandos** | 5 descrições específicas | 0 (vazias) | 100% |
| **Descrições de Flags** | 100% descritas | 0% descritas | 100% |

**Exemplo crítico**:
- **mgc** `--key`: "The SSH public key. The supported key types are: ssh-rsa, ssh-dss, ecdsa-sha, ssh-ed25519, sk-ecdsa-sha, sk-ssh-ed25519 (max character count: 16384)"
- **./tmp-cli/cli** `--key`: (sem descrição)

### 4. Flags Incorretas

| Flag | Comando | mgc | ./tmp-cli/cli | Problema |
|------|---------|-----|---------------|----------|
| `--limit` | `ssh-keys list` | Opcional | `(required)` | ❌ Incorretamente obrigatória |
| `--offset` | `ssh-keys list` | Opcional | `(required)` | ❌ Incorretamente obrigatória |
| `--sort` | `ssh-keys list` | Opcional | `(required)` | ❌ Incorretamente obrigatória |
| `--show-blocked` | `availability-zones list` | ❌ Não existe | `(required)` | ⚠️ Flag extra |

### 5. Flags Globais Ausentes

| Flag | Presente em mgc | Presente em ./tmp-cli/cli |
|------|-----------------|---------------------------|
| `--cli.retry-until` | ✅ | ❌ |
| `-t/--cli.timeout` | ✅ | ❌ |
| `-o/--output` | ✅ | ❌ |
| `--env` | ✅ | ❌ |
| `--server-url` | ✅ | ❌ |

---

## 📈 Estatísticas de Divergências

| Categoria | Total de Ocorrências | % de Impacto |
|-----------|---------------------|--------------|
| **Bugs Visuais** | 3 tipos (8 ocorrências) | 100% dos comandos/grupos |
| **Nomes Divergentes (Grupos)** | 1 | 50% dos grupos |
| **Nomes Divergentes (Flags)** | 3 | 60% das flags de paginação |
| **Descrições Ausentes (Grupos)** | 2 | 100% |
| **Descrições Ausentes (Comandos)** | 5 | 100% |
| **Descrições Ausentes (Flags)** | Todas | 100% |
| **Aliases Ausentes** | 2 | 100% dos grupos com aliases |
| **Flags Incorretas** | 4 | 80% das flags de lista |
| **Flags Globais Ausentes** | 5 | 71% das flags globais (5/7) |

---

## 🎯 Análise de Impacto

### 🔴 Impacto CRÍTICO

#### 1. Nome do Grupo `keys` vs `ssh-keys`
**Problema**: Ambiguidade semântica.

**mgc**: `mgc profile ssh-keys create` - Clara indicação de que são chaves SSH.

**./tmp-cli/cli**: `cli profile keys create` - "keys" pode ser:
- SSH keys?
- API keys?
- Chaves de criptografia?
- Chaves de acesso?

**Impacto**: Usuários podem confundir com outros tipos de chaves, especialmente em um contexto onde `--api-key` é uma flag global.

#### 2. Perda de Documentação
- **0% das flags** possuem descrições
- Usuários não sabem:
  - Quais tipos de chaves SSH são suportados
  - Limites de caracteres
  - Tipo de dado esperado (UUID, string, integer)
  - Formato esperado (base64, plain text, etc.)

#### 3. Flags Incorretamente Obrigatórias
- Flags de paginação (`limit`, `offset`, `sort`) marcadas como `(required)`
- Usuários não conseguem listar sem especificar paginação
- Comportamento divergente da referência (`mgc`)

### ⚠️ Impacto MÉDIO

#### 1. Prefixo `control.` Removido
- `control.limit` → `limit`
- Perda de contexto: o prefixo indica que são parâmetros de controle da consulta
- Inconsistência com boas práticas de API design

#### 2. Flag Extra `show-blocked`
- Não existe em `mgc`
- Não documentada
- Pode estar expondo funcionalidade não pronta ou não aprovada

---

## 🔍 Problemas Únicos de `profile`

### 1. **Novo Bug Visual**: "defaultLongDesc 1"
- Aparece apenas em `profile` (não visto em outros produtos)
- Sugere problema no template de descrição padrão

### 2. **Simplificação Incorreta de Nome**
- `ssh-keys` → `keys`
- Outros produtos mantêm nomes compostos (ex: `nat-gateways`, `security-groups`)
- Inconsistência no padrão de simplificação

### 3. **Flag Extra Misteriosa**
- `show-blocked` em `availability-zones list`
- Não existe em `mgc`
- Não documentada
- Origem desconhecida

---

## 📊 Comparação com Outros Produtos

| Problema | audit | block-storage | network | profile |
|----------|-------|---------------|---------|---------|
| Bug "Dqui1" | ✅ | ✅ | ✅ | ✅ |
| Bug "doto3" | ✅ | ✅ | ✅ | ✅ |
| Bug visual único | ❌ | ❌ | ❌ | ✅ "defaultLongDesc 1" |
| Prefixo `control.` removido | ✅ | ✅ | ❌ | ✅ |
| 0% descrições | ✅ | ✅ | ✅ | ✅ |
| Nome de grupo divergente | ❌ | ❌ | ✅ (`public-i-ps`, `v-p-cs`) | ✅ (`keys`) |
| Flags extras | ❌ | ❌ | ⚠️ | ⚠️ (`show-blocked`) |

---

## 🎯 Prioridades de Correção

### 🔴 CRÍTICO (Bloqueante para Produção)
1. **Eliminar bugs visuais** ("defaultLongDesc 1", "Dqui1", "doto3")
2. **Restaurar nome correto do grupo**: `keys` → `ssh-keys`
3. **Adicionar descrições em 100% das flags e comandos**
4. **Corrigir flags incorretamente obrigatórias**
5. **Restaurar prefixo `control.`** em flags de paginação
6. **Restaurar flags globais ausentes** (5 flags)

### ⚠️ ALTO (Afeta Usabilidade)
1. **Adicionar aliases faltando** (2 grupos)
2. **Investigar e documentar flag extra** (`show-blocked`)
3. **Adicionar tipos de flags** (UUID, integer, string)

### 📝 MÉDIO (Aprimoramentos)
1. **Melhorar descrições dos grupos**
2. **Adicionar exemplos de uso** (seguindo padrão de `mgc`)

---

## 💡 Conclusão

O produto `profile` apresenta:

### Problemas Críticos:
1. **3 tipos de bugs visuais** (incluindo um novo: "defaultLongDesc 1")
2. **Nome de grupo ambíguo**: `keys` não é claro (deveria ser `ssh-keys`)
3. **0% de documentação**: Flags e comandos sem descrição
4. **Flags incorretas**: 4 flags opcionais marcadas como obrigatórias
5. **Flag extra não documentada**: `show-blocked`

### Impacto Geral:
- **Compatibilidade**: 0% devido à divergência de nome do grupo
- **Usabilidade**: Severamente comprometida por falta de documentação
- **Clareza**: Nome `keys` é ambíguo e confuso
- **Confiabilidade**: Flags incorretas podem quebrar scripts existentes

### Severidade:
- **Profile** tem um dos **piores impactos** entre todos os produtos analisados devido à:
  - Ambiguidade do nome do grupo
  - Novo tipo de bug visual
  - Flag extra misteriosa

**Recomendação**: ❌ **NÃO APTO PARA PRODUÇÃO**. Requer correção imediata do nome do grupo, eliminação de bugs visuais e adição de documentação completa.

