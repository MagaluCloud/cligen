# Comparação: mgc container-registry vs ./tmp-cli/cli container-registry

**Data:** 24/10/2025  
**Referência:** `mgc` (comandos oficiais que precisam ser replicados em `./tmp-cli/cli`)  
**Nota:** Flags globais não incluídas nesta análise (já documentadas anteriormente)

---

## 1. Comando Principal: `container-registry`

### MGC (Referência)
```
Magalu Container Registry product API.

Aliases:
  container-registry, cr, registry

Commands:
  credentials  Routes related to credentials to login to Docker.
  images       Routes related to listing and deletion of images.
  proxy-caches Routes related to creating, listing and deletion of proxy-caches.
  registries   Routes related to creation, listing and deletion of registries.
  repositories Routes related to listing and deletion of repositories.
```

### TMP-CLI (Atual)
```
Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.

Package containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
This package allows you to manage container registries, repositories, images, and credentials.

Commands:
  credentials         Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
  images              Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
  registries          Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
  repositories        Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
```

### ❌ Problemas Identificados:
1. **Descrição divergente:**
   - MGC: "Magalu Container Registry product API." (concisa)
   - TMP-CLI: Descrição repetitiva e verbosa com "Package containerregistry..." redundante
2. **Aliases FALTANDO:** TMP-CLI não tem os aliases: `cr`, `registry`
3. **Subcomando FALTANDO:** `proxy-caches` não existe no TMP-CLI (6 comandos ausentes)
4. **Descrições dos subcomandos:**
   - MGC: Descrições específicas e informativas para cada comando
   - TMP-CLI: Mesma descrição genérica repetida para todos os comandos

---

## 2. Subcomando FALTANDO: `proxy-caches`

### MGC (Referência)
```
Routes related to creating, listing and deletion of proxy-caches.

Commands:
  create      Create a proxy cache
  delete      Delete a proxy cache by proxy_cache_id
  get         Get a proxy cache by proxy_cache_id
  list        List all proxy-caches
  status      status
  update      Update a proxy cache by proxy_cache_id
```

### TMP-CLI (Atual)
```
❌ SUBCOMANDO COMPLETAMENTE AUSENTE
```

### ❌ Problema:
**Grupo inteiro de comandos não existe no TMP-CLI** - 6 comandos faltando:
- `proxy-caches create`
- `proxy-caches delete`
- `proxy-caches get`
- `proxy-caches list`
- `proxy-caches status`
- `proxy-caches update`

---

## 3. Comando: `container-registry credentials`

### MGC (Referência)
```
Routes related to credentials to login to Docker.

Commands:
  list        Get credentials for container registry
  password    Reset password
```

### TMP-CLI (Atual)
```
Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.

Dqui1   <-- ⚠️ BUG

Commands:
  get                 Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
  reset-password      Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Descrição divergente:** TMP-CLI usa descrição genérica
3. **Nomes de comandos divergentes:**
   - MGC: `list` → TMP-CLI: `get`
   - MGC: `password` → TMP-CLI: `reset-password`
4. **Descrições dos subcomandos:** Todas genéricas ao invés de específicas

---

## 4. Comando: `container-registry images`

### MGC (Referência)
```
Routes related to listing and deletion of images.

Commands:
  delete      Delete image by digest or tag
  get         Get image details
  list        List images in container registry repository
```

### TMP-CLI (Atual)
```
Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.

Dqui1   <-- ⚠️ BUG

Commands:
  delete              Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
  get                 Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
  list                Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
  list-all            Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Descrição divergente:** TMP-CLI usa descrição genérica
3. **Comando EXTRA:** `list-all` não existe no MGC
4. **Descrições dos subcomandos:** Todas genéricas ao invés de específicas

---

## 5. Comando: `container-registry registries`

### MGC (Referência)
```
Routes related to creation, listing and deletion of registries.

Commands:
  create      Create a container registry
  delete      Delete a container registry by registry_id
  get         Get registry information
  list        List all container registries
```

### TMP-CLI (Atual)
```
Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.

Dqui1   <-- ⚠️ BUG

Commands:
  create              Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
  delete              Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
  get                 Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
  list                Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
  list-all            Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Descrição divergente:** TMP-CLI usa descrição genérica
3. **Comando EXTRA:** `list-all` não existe no MGC
4. **Descrições dos subcomandos:** Todas genéricas ao invés de específicas

---

## 6. Comando: `container-registry repositories`

### MGC (Referência)
```
Routes related to listing and deletion of repositories.

Commands:
  delete      Delete a container registry repository by repository_id.
  get         Get a container registry repository by repository_id
  list        List all container registry repositories
```

### TMP-CLI (Atual)
```
Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.

Dqui1   <-- ⚠️ BUG

Commands:
  delete              Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
  get                 Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
  list                Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
  list-all            Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Descrição divergente:** TMP-CLI usa descrição genérica
3. **Comando EXTRA:** `list-all` não existe no MGC
4. **Descrições dos subcomandos:** Todas genéricas ao invés de específicas

---

## 7. Comando: `credentials list` vs `credentials get`

### MGC (Referência)
```
Return container registry user's authentication credentials.

Usage:
  mgc container-registry credentials list [flags]

Flags: (nenhuma flag local)
```

### TMP-CLI (Atual - comando `get`)
```
Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.

doto3   <-- ⚠️ BUG

Usage: cli container-registry credentials get [flags]

Flags: (nenhuma flag local)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Nome do comando divergente:** MGC: `list` vs TMP-CLI: `get`
3. **Descrição do comando:** Genérica ao invés de específica

---

## 8. Comando: `credentials password` vs `credentials reset-password`

### MGC (Referência)
```
Reset container registry user's password.

Usage:
  mgc container-registry credentials password [flags]

Flags: (nenhuma flag local)
```

### TMP-CLI (Atual - comando `reset-password`)
```
Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.

doto3   <-- ⚠️ BUG

Usage: cli container-registry credentials reset-password [flags]

Flags: (nenhuma flag local)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Nome do comando divergente:** MGC: `password` vs TMP-CLI: `reset-password`
3. **Descrição do comando:** Genérica ao invés de específica

---

## 9. Comando: `registries create`

### MGC (Referência)
```
Creates a container registry in Magalu Cloud.

Flags:
  --cli.list-links enum[=table]   List all available links for this command
  --name string                   A unique, global name for the container registry. It must be written in lowercase 
                                  letters and consists only of numbers and letters, up to a limit of 63 characters. (required)
  --proxy-cache-id string         Proxy Cache UUID.
```

### TMP-CLI (Atual)
```
Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.

doto3   <-- ⚠️ BUG

Flags:
  --name             (required) (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** Genérica ao invés de específica
3. **Flags FALTANDO:**
   - `--cli.list-links`
   - `--proxy-cache-id`
4. **Flag sem descrição:** `--name` aparece sem descrição

---

## 10. Comando: `registries list`

### MGC (Referência)
```
List user's container registries.

Flags:
  --control.limit integer    Limit (min: 1)
  --control.offset integer   Offset (min: 0)
  --control.sort string      Fields to use as reference to sort.
  --name string              Name
```

### TMP-CLI (Atual)
```
Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.

doto3   <-- ⚠️ BUG

Usage: cli container-registry registries list [Offset] [Limit] [Sort] [flags]

Flags:
  --limit            (required) (sem descrição)
  --offset           (required) (sem descrição)
  --sort             (required) (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** Genérica ao invés de específica
3. **Argumentos posicionais incorretos:** Usage mostra `[Offset] [Limit] [Sort]` como se fossem argumentos posicionais
4. **Nomes de flags divergentes:**
   - MGC: `--control.limit` → TMP-CLI: `--limit` (sem prefixo)
   - MGC: `--control.offset` → TMP-CLI: `--offset` (sem prefixo)
   - MGC: `--control.sort` → TMP-CLI: `--sort` (sem prefixo)
5. **Flag FALTANDO:** `--name` (filtro por nome)
6. **Flags sem descrição:** Todas as flags sem descrição
7. **Marcação incorreta de required:** Todas as flags marcadas como required quando não deveriam ser

---

## 11. Comando: `images list`

### MGC (Referência)
```
List all images in container registry repository

Usage:
  mgc container-registry images list [registry-id] [repository-id] [flags]

Flags:
  --control.limit integer    Limit (min: 1)
  --control.offset integer   Offset (min: 0)
  --control.sort string      Fields to use as reference to sort.
  --expand array(string)     You can get more detailed info about: ['tags_details', 'extra_attr', 'manifest_media_type', 'media_type']
  --name string              Used to filter images in response
  --registry-id uuid         Container Registry's UUID. (required)
  --repository-id uuid       Repository's UUID. (required)
```

### TMP-CLI (Atual)
```
Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.

doto3   <-- ⚠️ BUG

Usage: cli container-registry images list [registryID] [repositoryName] [Offset] [Limit] [Sort] [Expand] [flags]

Flags:
  --expand          (sem descrição)
  --limit            (required) (sem descrição)
  --offset           (required) (sem descrição)
  --registry-id      (required) (sem descrição)
  --repository-name  (required) (sem descrição)
  --sort             (required) (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** Genérica ao invés de específica
3. **Argumentos posicionais divergentes:**
   - MGC: `[registry-id] [repository-id]`
   - TMP-CLI: `[registryID] [repositoryName] [Offset] [Limit] [Sort] [Expand]`
4. **Nomenclatura divergente:**
   - MGC: `--repository-id` (UUID) → TMP-CLI: `--repository-name` (Nome)
   - Isso muda o tipo de parâmetro esperado!
5. **Nomes de flags divergentes:**
   - MGC: `--control.limit` → TMP-CLI: `--limit` (sem prefixo)
   - MGC: `--control.offset` → TMP-CLI: `--offset` (sem prefixo)
   - MGC: `--control.sort` → TMP-CLI: `--sort` (sem prefixo)
6. **Flag FALTANDO:** `--name` (filtro por nome de imagem)
7. **Flags sem descrição:** Todas as flags sem descrição
8. **Marcação incorreta de required:** Todas as flags marcadas como required

---

## 12. Comando: `images delete`

### MGC (Referência)
```
Delete repository image by digest or tag

Usage:
  mgc container-registry images delete [registry-id] [repository-id] [digest-or-tag] [flags]

Flags:
  --digest-or-tag string   Digest or tag of an image (required)
  --registry-id uuid       Container Registry's UUID. (required)
  --repository-id uuid     Repository's UUID. (required)
```

### TMP-CLI (Atual)
```
Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.

doto3   <-- ⚠️ BUG

Usage: cli container-registry images delete [registryID] [repositoryName] [digestOrTag] [flags]

Flags:
  --digest-or-tag    (required) (sem descrição)
  --registry-id      (required) (sem descrição)
  --repository-name  (required) (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** Genérica ao invés de específica
3. **Argumentos posicionais divergentes:**
   - MGC: `[registry-id] [repository-id] [digest-or-tag]`
   - TMP-CLI: `[registryID] [repositoryName] [digestOrTag]` (formato camelCase)
4. **Nomenclatura divergente:**
   - MGC: `--repository-id` (UUID) → TMP-CLI: `--repository-name` (Nome)
   - Isso muda o tipo de parâmetro esperado!
5. **Flags sem descrição:** Todas as flags sem descrição

---

## 13. Comando: `repositories list`

### MGC (Referência)
```
List all user's repositories in the container registry.

Usage:
  mgc container-registry repositories list [registry-id] [flags]

Flags:
  --control.limit integer    Limit (min: 1)
  --control.offset integer   Offset (min: 0)
  --control.sort string      Fields to use as reference to sort.
  --name string              Used to filter repositories in response
  --registry-id uuid         Container Registry's UUID. (required)
```

### TMP-CLI (Atual)
```
Containerregistry provides a client for interacting with the Magalu Cloud Container Registry API.

doto3   <-- ⚠️ BUG

Usage: cli container-registry repositories list [registryID] [Offset] [Limit] [Sort] [flags]

Flags:
  --limit            (required) (sem descrição)
  --offset           (required) (sem descrição)
  --registry-id      (required) (sem descrição)
  --sort             (required) (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** Genérica ao invés de específica
3. **Argumentos posicionais divergentes:**
   - MGC: `[registry-id]`
   - TMP-CLI: `[registryID] [Offset] [Limit] [Sort]`
4. **Nomes de flags divergentes:**
   - MGC: `--control.limit` → TMP-CLI: `--limit` (sem prefixo)
   - MGC: `--control.offset` → TMP-CLI: `--offset` (sem prefixo)
   - MGC: `--control.sort` → TMP-CLI: `--sort` (sem prefixo)
5. **Flag FALTANDO:** `--name` (filtro por nome)
6. **Flags sem descrição:** Todas as flags sem descrição
7. **Marcação incorreta de required:** Todas as flags marcadas como required

---

## 14. Comandos `proxy-caches` (FALTANDO)

### MGC - proxy-caches create
```
Creates a proxy cache in Magalu Cloud.

Flags:
  --access-key string      A string consistent with provider access_id.
  --access-secret string   A string consistent with provider access_secret.
  --description string     A string.
  --name string            A unique name for each tenant, used for the proxy-cache. (required)
  --provider string        A provider identifier string. (required)
  --url string             An Endpoint URL for the proxied registry. (required)
```

### MGC - proxy-caches list
```
List user's proxy caches.

Flags:
  --control.limit integer    Limit (min: 1)
  --control.offset integer   Offset (min: 0)
  --control.sort string      Fields to use as reference to sort.
```

### MGC - Outros comandos proxy-caches
- `delete` - Delete a proxy cache by proxy_cache_id
- `get` - Get a proxy cache by proxy_cache_id
- `status` - status
- `update` - Update a proxy cache by proxy_cache_id

### TMP-CLI
```
❌ NENHUM DESTES COMANDOS EXISTE
```

---

## Resumo Geral de Problemas

### 🐛 BUGS CRÍTICOS:
- String **"Dqui1"** aparece em todos os 4 subcomandos de nível 2
- String **"doto3"** aparece em todos os comandos leaf

### ❌ COMANDOS/SUBCOMANDOS FALTANDO:
- Grupo inteiro `proxy-caches` (6 comandos faltando)

### ➕ COMANDOS EXTRAS:
- `list-all` em: images, registries, repositories (3 comandos)

### 🔄 NOMENCLATURA DIVERGENTE:
**Comandos:**
- `credentials list` → TMP-CLI usa `credentials get`
- `credentials password` → TMP-CLI usa `credentials reset-password`

**Flags:**
- `--repository-id` (UUID) → TMP-CLI usa `--repository-name` (Nome) ⚠️ **TIPO DIFERENTE**
- Falta prefixo `control.` em: limit, offset, sort

**Argumentos posicionais:**
- MGC: kebab-case (`[registry-id]`)
- TMP-CLI: camelCase (`[registryID]`)

### 🚫 ALIASES FALTANDO:
- `cr`, `registry` (no comando principal)

### 📝 DESCRIÇÕES:
- Todas as descrições são genéricas e repetitivas
- Faltam descrições de todas as flags
- Faltam descrições específicas dos comandos

### 🔧 FLAGS:
#### Padrão de problemas em TODOS os comandos list:
- Falta prefixo `control.` em: `limit`, `offset`, `sort`
- Flags marcadas incorretamente como `required`
- Flags sem descrição
- Argumentos posicionais incorretos no Usage
- Flag `--name` (filtro) faltando em vários list

#### Flags FALTANDO:
- `--cli.list-links` (em comandos create)
- `--proxy-cache-id` (em registries create)
- `--name` (em registries list, images list, repositories list)

#### Problema GRAVE de Tipo:
- `--repository-id` vs `--repository-name`: Muda de UUID para Nome, afetando funcionalidade!

---

## 📊 Estatísticas

### Subcomandos (Nível 2)
- **Total MGC:** 5 grupos
- **Total TMP-CLI:** 4 grupos
- **Faltando:** 1 grupo inteiro (proxy-caches)
- **Bugs "Dqui1":** 100% (4/4 dos grupos presentes)

### Comandos Leaf (Nível 3)
- **Total MGC:** 19 comandos
- **Total TMP-CLI:** 13 comandos (mais 3 list-all = 16 total)
- **Faltando:** 6 comandos (todo o proxy-caches)
- **Nome divergente:** 2 comandos (credentials)
- **Bugs "doto3":** 100% (todos os comandos leaf)
- **Comandos extras:** 3 (list-all)

### Aliases
- **Total MGC:** 3 (container-registry, cr, registry)
- **Total TMP-CLI:** 1 (só o nome principal)
- **Compatibilidade:** 33% (1/3)

### Flags
- **Sem prefixo `control.`:** 100% dos list
- **Sem descrição:** 100%
- **Required incorreto:** ~80% dos comandos list
- **Tipo divergente:** `repository-id` vs `repository-name` (crítico)

---

## ✅ Checklist de Ações

### P0 - Crítico
- [ ] Remover string "Dqui1" de todos os subcomandos nível 2
- [ ] Remover string "doto3" de todos os comandos leaf
- [ ] Adicionar descrições em TODAS as flags
- [ ] Corrigir descrição do comando principal
- [ ] **Implementar grupo `proxy-caches` completo (6 comandos)**
- [ ] **Corrigir `--repository-name` → `--repository-id`** (problema de tipo crítico)

### P1 - Alto
- [ ] Adicionar aliases: `cr`, `registry`
- [ ] Renomear `credentials get` → `credentials list`
- [ ] Renomear `credentials reset-password` → `credentials password`
- [ ] Adicionar prefixo `control.` em: limit, offset, sort
- [ ] Corrigir marcação de flags required
- [ ] Remover argumentos posicionais incorretos do Usage
- [ ] Adicionar descrições específicas para cada comando
- [ ] Padronizar argumentos posicionais (kebab-case, não camelCase)

### P2 - Médio
- [ ] Adicionar flag `--cli.list-links` em comandos de criação
- [ ] Adicionar flag `--proxy-cache-id` em `registries create`
- [ ] Adicionar flag `--name` em: registries list, images list, repositories list
- [ ] Remover comandos `list-all` (3 comandos)

### P3 - Baixo
- [ ] Melhorar formatação geral do help
- [ ] Padronizar estrutura de descrições

