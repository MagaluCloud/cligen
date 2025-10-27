# Comparação Detalhada: mgc profile vs ./tmp-cli/cli profile

## Sumário Executivo

A comparação entre `mgc profile` e `./tmp-cli/cli profile` revela divergências de nomenclatura e problemas sistemáticos no CLI gerado:

### 🔴 Problemas Críticos Identificados

1. **BUG VISUAL**: String "defaultLongDesc 1" aparece no output principal
2. **BUGS VISUAIS**: Strings "Dqui1" e "doto3" aparecem nos subgrupos
3. **NOME DE GRUPO DIVERGENTE**: `ssh-keys` → `keys`
4. **0% DE DESCRIÇÕES NAS FLAGS**: Nenhuma flag em `./tmp-cli/cli` possui descrição
5. **FLAGS GLOBAIS AUSENTES**: `--cli.retry-until`, `-t/--cli.timeout`, `-o/--output`, `--env`, `--server-url`
6. **FLAGS INCORRETAMENTE MARCADAS COMO REQUIRED**: Em `list`, flags opcionais marcadas como `(required)`
7. **FLAG EXTRA**: `show-blocked` em `availability-zones list` não existe em `mgc`
8. **ALIASES AUSENTES**: Todos os aliases ausentes em `./tmp-cli/cli`

---

## 1. Comando Principal

### mgc profile

```
The profile group provides commands to view and modify user account settings. 
It allows users to manage their SSH keys, update personal information, and configure other 
account-related preferences. This group is essential for maintaining secure access and 
personalizing the user experience within the system.

Usage:
  mgc profile [flags]
  mgc profile [command]

Commands:
  availability-zones Manage Availability Zones
  ssh-keys           Manage SSH Keys
```

### ./tmp-cli/cli profile

```
Manage account settings, including SSH keys and related configurations.

defaultLongDesc 1

Available Commands:

Other commands:
  availability-zones  
  keys                
```

### ⚠️ Diferenças Identificadas

| Aspecto | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **Bug visual** | Não | "defaultLongDesc 1" aparece | 🔴 CRÍTICO |
| **Nome do grupo `ssh-keys`** | `ssh-keys` | `keys` | 🔴 CRÍTICO |
| **Descrição do produto** | Completa e detalhada | Simples | ⚠️ Divergente |
| **Descrições dos comandos** | Específicas | Vazias | 🔴 CRÍTICO |
| **Flags globais** | 7 flags | 5 flags (faltam 3) | ❌ Incompleto |

---

## 2. Grupo: ssh-keys / keys

### mgc profile ssh-keys

```
Manage SSH Keys

Usage:
  mgc profile ssh-keys [flags]
  mgc profile ssh-keys [command]

Aliases:
  ssh-keys, ssh_keys

Commands:
  create      Register new SSH key
  delete      Delete SSH Key
  get         Retrieve a SSH key
  list        List SSH keys
```

### ./tmp-cli/cli profile keys

```


Dqui1

Available Commands:

Other commands:
  create              
  delete              
  get                 
  list                
```

### ⚠️ Diferenças

| Aspecto | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **Nome do grupo** | `ssh-keys` | `keys` | 🔴 CRÍTICO |
| **Bug visual** | Não | "Dqui1" aparece | 🔴 CRÍTICO |
| **Descrição do grupo** | "Manage SSH Keys" | Vazio (apenas linha em branco) | 🔴 CRÍTICO |
| **Descrições dos comandos** | Específicas | Vazias | 🔴 CRÍTICO |
| **Aliases** | `ssh-keys`, `ssh_keys` | Nenhum | ❌ Faltando |

### 2.1. ssh-keys create

**mgc**: `mgc profile ssh-keys create`
- Flags:
  - `--key` (required): The SSH public key (max 16384 chars). Supported types: ssh-rsa, ssh-dss, ecdsa-sha, ssh-ed25519, sk-ecdsa-sha, sk-ssh-ed25519
  - `--name` (required): The SSH Key name (max 45 chars)
- Descrição completa e detalhada

**./tmp-cli/cli**: `cli profile keys create`
- Flags:
  - `--key` (required)
  - `--name` (required)
- Bug visual: "doto3"
- 0% das flags possuem descrições

**Divergências**:
- 🔴 Nome do grupo: `ssh-keys` vs `keys`
- 🔴 Bug visual: "doto3"
- ❌ Faltam descrições completas nas flags (tipo de chave, limites de caracteres, etc.)
- ❌ Falta descrição do comando

### 2.2. ssh-keys delete

**mgc**: `mgc profile ssh-keys delete [key-id]`
- Flags:
  - `--key-id` (uuid, required): Key Id

**./tmp-cli/cli**: `cli profile keys delete [keyID]`
- Flags:
  - `--key-id` (required)
- Bug visual: "doto3"

**Divergências**:
- 🔴 Nome do grupo: `ssh-keys` vs `keys`
- 🔴 Bug visual: "doto3"
- ❌ Tipo da flag ausente (uuid)
- ❌ Falta descrição da flag

### 2.3. ssh-keys list

**mgc**: `mgc profile ssh-keys list`
- Flags:
  - `--control.limit` (integer): Limit
  - `--control.offset` (integer): Offset
  - `--control.sort` (string): Sort
- Descrição: "List the SSH keys. It is possible sort this list with parameters id, name, key_type"

**./tmp-cli/cli**: `cli profile keys list [Limit] [Offset] [Sort]`
- Flags:
  - `--limit` (required)
  - `--offset` (required)
  - `--sort` (required)
- Bug visual: "doto3"

**Divergências**:
- 🔴 Nome do grupo: `ssh-keys` vs `keys`
- 🔴 Bug visual: "doto3"
- ❌ `--control.limit`, `--control.offset`, `--control.sort` vs `--limit`, `--offset`, `--sort` (prefixo `control.` ausente)
- ❌ Flags incorretamente marcadas como `(required)` - devem ser opcionais
- ❌ Falta descrição do comando

---

## 3. Grupo: availability-zones

### mgc profile availability-zones

```
Manage Availability Zones

Usage:
  mgc profile availability-zones [flags]
  mgc profile availability-zones [command]

Aliases:
  availability-zones, availability_zones

Commands:
  list        List all availability zones.
```

### ./tmp-cli/cli profile availability-zones

```


Dqui1

Available Commands:

Other commands:
  list                
```

### ⚠️ Diferenças

| Aspecto | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **Bug visual** | Não | "Dqui1" aparece | 🔴 CRÍTICO |
| **Descrição do grupo** | "Manage Availability Zones" | Vazio | 🔴 CRÍTICO |
| **Descrição do comando** | "List all availability zones." | Vazio | 🔴 CRÍTICO |
| **Aliases** | `availability-zones`, `availability_zones` | Nenhum | ❌ Faltando |

### 3.1. availability-zones list

**mgc**: `mgc profile availability-zones list`
- Flags: Nenhuma flag local (apenas global flags)
- Descrição: "List all the availability zones"

**./tmp-cli/cli**: `cli profile availability-zones list [ShowBlocked]`
- Flags:
  - `--show-blocked` (required)
- Bug visual: "doto3"

**Divergências**:
- 🔴 Bug visual: "doto3"
- ⚠️ Flag extra: `--show-blocked` não existe em `mgc`
- ❌ Flag marcada incorretamente como `(required)`
- ❌ Falta descrição da flag

---

## 4. Problemas Sistemáticos Identificados

### 🔴 Bugs Visuais Críticos
1. **"defaultLongDesc 1"** aparece no comando principal
2. **"Dqui1"** aparece em ambos os subgrupos
3. **"doto3"** aparece em todos os comandos específicos

### 🔴 Problemas de Nomenclatura
1. **Nome de grupo**: `ssh-keys` → `keys` (perda de clareza semântica)
2. **Prefixo ausente**: `control.limit` → `limit`, `control.offset` → `offset`, `control.sort` → `sort`

### 🔴 Problemas de Conteúdo
1. **0% de descrições nas flags** em `./tmp-cli/cli`
2. **0% de descrições dos comandos** em `./tmp-cli/cli`
3. **Descrições dos grupos vazias** (apenas linhas em branco)
4. **Aliases ausentes** em todos os grupos (2 grupos afetados)

### 🔴 Problemas de Flags
1. **Flags globais ausentes**: `--cli.retry-until`, `-t/--cli.timeout`, `-o/--output`, `--env`, `--server-url`
2. **Flags incorretamente marcadas como required**:
   - `ssh-keys list`: `limit`, `offset`, `sort`
   - `availability-zones list`: `show-blocked`
3. **Flag extra**: `show-blocked` em `availability-zones list` não existe em `mgc`
4. **Tipos de flags ausentes**: `uuid` em `--key-id` não especificado

---

## 5. Resumo de Incompatibilidades

| Categoria | Qtd. Problemas | Severidade |
|-----------|----------------|------------|
| **Bugs Visuais** | 3 tipos | 🔴 CRÍTICA |
| **Nomes Divergentes** | 1 grupo + 3 flags | 🔴 CRÍTICA |
| **Descrições Ausentes** | 100% das flags e comandos | 🔴 CRÍTICA |
| **Flags Incorretas** | 4 flags marcadas como required | 🔴 CRÍTICA |
| **Flags Extras** | 1 flag (`show-blocked`) | ⚠️ MÉDIA |
| **Flags Globais Ausentes** | 5 flags | 🔴 CRÍTICA |
| **Aliases Ausentes** | 2 grupos | ⚠️ BAIXA |

---

## 6. Recomendações

### Prioridade CRÍTICA 🔴
1. **Eliminar bugs visuais**: Remover "defaultLongDesc 1", "Dqui1" e "doto3" completamente
2. **Restaurar nome correto do grupo**: `keys` → `ssh-keys`
3. **Adicionar descrições em 100% das flags e comandos**
4. **Restaurar prefixo `control.`**: `limit` → `control.limit`, etc.
5. **Corrigir flags incorretamente marcadas como required**
6. **Restaurar flags globais ausentes**

### Prioridade ALTA ⚠️
1. **Adicionar aliases faltando**
2. **Verificar flag extra**: Confirmar se `show-blocked` deve existir em `availability-zones list`
3. **Adicionar tipos de flags**: Especificar `uuid` onde apropriado

### Prioridade MÉDIA
1. **Melhorar descrições dos grupos**: Adicionar descrições claras e úteis

---

## Conclusão

O CLI gerado (`./tmp-cli/cli profile`) apresenta **problemas sistemáticos graves**:

### Problemas Únicos de `profile`:
1. **Novo bug visual**: "defaultLongDesc 1" (não visto em outros produtos)
2. **Nome de grupo simplificado incorretamente**: `ssh-keys` → `keys`
3. **Flag extra não documentada**: `show-blocked`

### Problemas Compartilhados com Outros Produtos:
1. **Bugs visuais**: "Dqui1" e "doto3"
2. **0% de descrições** nas flags e comandos
3. **Flags globais ausentes**
4. **Prefixo `control.` removido** (padrão visto em block-storage, audit, etc.)
5. **Flags incorretamente marcadas como required**
6. **Aliases ausentes**

### Impacto Geral:
- **Compatibilidade**: 0% dos comandos são 100% compatíveis devido à divergência de nomes
- **Usabilidade**: Severamente comprometida por bugs visuais e falta de documentação
- **Clareza**: Nome do grupo `keys` é ambíguo comparado a `ssh-keys`

**Recomendação**: ❌ **NÃO APTO PARA PRODUÇÃO**. Requer correção de bugs visuais, restauração do nome correto do grupo e adição de descrições completas.

