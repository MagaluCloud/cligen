# Resumo Executivo: Comparação network

## 📊 Visão Geral

| Métrica | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **Grupos Principais** | 8 | 8 | ✅ |
| **Bugs Visuais** | 0 | 2 tipos ("Dqui1", "doto3") | 🔴 |
| **Nomes de Grupos Divergentes** | - | 3 (`public-i-ps`, `v-p-cs`, `subnet-pools`) | 🔴 |
| **Subgrupos Ausentes** | - | 1 (`security-groups rules`) | 🔴 |
| **Arquitetura Divergente** | - | 1 grupo (`vpcs`) | 🔴 |
| **Comandos Extras** | - | 3 (`rules create/list`, `vpcs rename`) | ⚠️ |
| **Flags com Descrição** | 100% | 0% | 🔴 |
| **Aliases Ausentes** | - | 4 grupos | ⚠️ |
| **Flags Globais Ausentes** | - | 7 flags | 🔴 |

---

## 📋 Tabela de Comandos Completa

### Grupo: nat-gateways

| Comando | mgc | ./tmp-cli/cli | Divergências | Severidade |
|---------|-----|---------------|--------------|------------|
| **Grupo** | ✅ `nat-gateways` | ✅ `nat-gateways` | - | ✅ |
| **create** | ✅ | ✅ | `--vpc-id` vs `--vpcid`; Bug "doto3"; 0% descrições | 🔴 |
| **delete** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **get** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **list** | ✅ | ✅ | `items-per-page/page` vs `limit/offset`; Flags incorretas como required; Bug "doto3" | 🔴 |

**Aliases**: mgc tem (`nat-gateways`, `nat_gateways`), ./tmp-cli/cli não tem ❌

---

### Grupo: ports

| Comando | mgc | ./tmp-cli/cli | Divergências | Severidade |
|---------|-----|---------------|--------------|------------|
| **Grupo** | ✅ `ports` | ✅ `ports` | - | ✅ |
| **attach** | ✅ `attach` | ❌ `attach-security-group` | Nome divergente; Bug "doto3"; Bug formatação (falta espaço) | 🔴 |
| **delete** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **detach** | ✅ `detach` | ❌ `detach-security-group` | Nome divergente; Bug "doto3" | 🔴 |
| **get** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **list** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **update** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |

---

### Grupo: public-ips

| Comando | mgc | ./tmp-cli/cli | Divergências | Severidade |
|---------|-----|---------------|--------------|------------|
| **Grupo** | ✅ `public-ips` | 🔴 `public-i-ps` | **Nome do grupo divergente** | 🔴 |
| **attach** | ✅ `attach` | ❌ `attach-to-port` | Nome divergente; `--public-ip-id` vs `--public-ipid`; Bug "doto3" | 🔴 |
| **delete** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **detach** | ✅ `detach` | ❌ `detach-from-port` | Nome divergente; Bug "doto3" | 🔴 |
| **get** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **list** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |

**Aliases**: mgc tem (`public-ips`, `public_ips`), ./tmp-cli/cli não tem ❌

---

### Grupo: rules

| Comando | mgc | ./tmp-cli/cli | Divergências | Severidade |
|---------|-----|---------------|--------------|------------|
| **Grupo** | ✅ `rules` | ✅ `rules` | - | ✅ |
| **create** | ❌ Não existe | ⚠️ Existe | **Comando extra** (pode estar em `security-groups rules`) | ⚠️ |
| **delete** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **get** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **list** | ❌ Não existe | ⚠️ Existe | **Comando extra** (pode estar em `security-groups rules`) | ⚠️ |

---

### Grupo: security-groups

| Comando | mgc | ./tmp-cli/cli | Divergências | Severidade |
|---------|-----|---------------|--------------|------------|
| **Grupo** | ✅ `security-groups` | ✅ `security-groups` | - | ✅ |
| **create** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **delete** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **get** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **list** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **rules** (subgrupo) | ✅ | 🔴 **AUSENTE** | **Subgrupo completamente ausente** | 🔴 |

**Aliases**: mgc tem (`security-groups`, `security_groups`), ./tmp-cli/cli não tem ❌

#### Subgrupo: security-groups rules (AUSENTE em ./tmp-cli/cli)

| Comando | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **rules** (grupo) | ✅ | 🔴 **AUSENTE** | 🔴 |
| **rules create** | ✅ | 🔴 **AUSENTE** | 🔴 |
| **rules list** | ✅ | 🔴 **AUSENTE** | 🔴 |

---

### Grupo: subnetpools

| Comando | mgc | ./tmp-cli/cli | Divergências | Severidade |
|---------|-----|---------------|--------------|------------|
| **Grupo** | ✅ `subnetpools` | ❌ `subnet-pools` | **Nome do grupo divergente** | 🔴 |
| **create** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **create-book-cidr** | ✅ | ❌ `book-c-i-d-r` | Nome divergente; Bug "doto3" | 🔴 |
| **create-unbook-cidr** | ✅ | ❌ `unbook-c-i-d-r` | Nome divergente; Bug "doto3" | 🔴 |
| **delete** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **get** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **list** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |

---

### Grupo: subnets

| Comando | mgc | ./tmp-cli/cli | Divergências | Severidade |
|---------|-----|---------------|--------------|------------|
| **Grupo** | ✅ `subnets` | ✅ `subnets` | - | ✅ |
| **delete** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **get** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **update** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |

---

### Grupo: vpcs

| Comando | mgc | ./tmp-cli/cli | Divergências | Severidade |
|---------|-----|---------------|--------------|------------|
| **Grupo** | ✅ `vpcs` | 🔴 `v-p-cs` | **Nome do grupo divergente** | 🔴 |
| **create** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **delete** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **get** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **list** | ✅ | ✅ | Bug "doto3"; 0% descrições | ⚠️ |
| **rename** | ❌ Não existe | ⚠️ Existe | **Comando extra** | ⚠️ |
| **ports** (subgrupo) | ✅ | 🔴 Convertido em comandos flat | **Arquitetura divergente** | 🔴 |
| **public-ips** (subgrupo) | ✅ | 🔴 Convertido em comandos flat | **Arquitetura divergente** | 🔴 |
| **subnets** (subgrupo) | ✅ | 🔴 Convertido em comandos flat | **Arquitetura divergente** | 🔴 |

#### Subgrupos convertidos em comandos flat (PROBLEMA ARQUITETURAL)

| mgc (subgrupos) | ./tmp-cli/cli (flat) | Status |
|-----------------|----------------------|--------|
| `vpcs ports create` | `v-p-cs create-port` | 🔴 Incompatível |
| `vpcs ports list` | `v-p-cs list-ports` | 🔴 Incompatível |
| `vpcs public-ips create` | `v-p-cs create-public-i-p` | 🔴 Incompatível |
| `vpcs public-ips list` | `v-p-cs list-public-i-ps` | 🔴 Incompatível |
| `vpcs subnets create` | `v-p-cs create-subnet` | 🔴 Incompatível |
| `vpcs subnets list` | `v-p-cs list-subnets` | 🔴 Incompatível |

**Flags divergentes em `vpcs subnets create`**:
- `--cidr-block` vs `--cidrblock`
- `--dns-nameservers` vs `--dnsnameservers`
- `--ip-version` vs `--ipversion`
- `--subnetpool-id` vs `--subnet-pool-id`
- Flag extra: `--zone` (não existe em mgc)

---

## 🔴 Problemas Críticos Sistemáticos

### 1. Bugs Visuais
- ✅ **"Dqui1"** aparece em todos os 8 grupos
- ✅ **"doto3"** aparece em todos os comandos específicos
- ✅ **Falta de espaço** entre nome e descrição em `ports attach-security-group` e `detach-security-group`

### 2. Nomes Divergentes

#### Grupos
| mgc | ./tmp-cli/cli | Status |
|-----|---------------|--------|
| `public-ips` | `public-i-ps` | 🔴 |
| `vpcs` | `v-p-cs` | 🔴 |
| `subnetpools` | `subnet-pools` | 🔴 |

#### Comandos
| mgc | ./tmp-cli/cli | Status |
|-----|---------------|--------|
| `ports attach` | `ports attach-security-group` | 🔴 |
| `ports detach` | `ports detach-security-group` | 🔴 |
| `public-ips attach` | `public-i-ps attach-to-port` | 🔴 |
| `public-ips detach` | `public-i-ps detach-from-port` | 🔴 |
| `subnetpools create-book-cidr` | `subnet-pools book-c-i-d-r` | 🔴 |
| `subnetpools create-unbook-cidr` | `subnet-pools unbook-c-i-d-r` | 🔴 |

#### Flags (padrão sistemático: remoção de hífens)
| mgc | ./tmp-cli/cli | Ocorrências |
|-----|---------------|-------------|
| `--vpc-id` | `--vpcid` | nat-gateways create |
| `--public-ip-id` | `--public-ipid` | public-ips attach |
| `--cidr-block` | `--cidrblock` | vpcs subnets create |
| `--dns-nameservers` | `--dnsnameservers` | vpcs subnets create |
| `--ip-version` | `--ipversion` | vpcs subnets create |
| `--subnetpool-id` | `--subnet-pool-id` | vpcs subnets create |
| `--security-groups-id` | `--security-groups` | vpcs ports create |

### 3. Arquitetura Divergente

| Problema | Impacto | Severidade |
|----------|---------|------------|
| Subgrupo `security-groups rules` ausente | Perda de 2 comandos (`create`, `list`) | 🔴 CRÍTICO |
| Subgrupos `vpcs` convertidos em comandos flat | Perda de 3 subgrupos, 6 comandos incompatíveis | 🔴 CRÍTICO |

### 4. Flags Globais Ausentes

| Flag | mgc | ./tmp-cli/cli |
|------|-----|---------------|
| `--cli.retry-until` | ✅ | ❌ |
| `-t/--cli.timeout` | ✅ | ❌ |
| `--env` | ✅ | ❌ |
| `--region` | ✅ | ❌ |
| `--server-url` | ✅ | ❌ |
| `-o/--output` | ✅ | ❌ |
| `--x-zone` | ✅ (em alguns comandos) | ❌ |

### 5. Outros Problemas

| Problema | Qtd. | Impacto |
|----------|------|---------|
| **Aliases ausentes** | 4 grupos | Perda de flexibilidade |
| **Descrições de flags ausentes** | 100% | Perda de usabilidade |
| **Descrições genéricas** | 100% dos comandos | Perda de clareza |
| **Comandos extras** | 3 (`rules create/list`, `vpcs rename`) | Inconsistência com referência |
| **Flags incorretas como required** | `list` em nat-gateways | Usabilidade comprometida |
| **Nomenclatura de paginação divergente** | `items-per-page/page` vs `limit/offset` | Inconsistência |

---

## 📈 Estatísticas de Divergências

| Categoria | Total de Ocorrências | % de Impacto |
|-----------|---------------------|--------------|
| **Bugs Visuais** | 2 tipos (ubíquos) | 100% dos comandos |
| **Nomes Divergentes (Grupos)** | 3 | 37.5% dos grupos |
| **Nomes Divergentes (Comandos)** | 8 | ~20% dos comandos |
| **Nomes Divergentes (Flags)** | 15+ | ~30% das flags |
| **Descrições Ausentes (Flags)** | Todas | 100% |
| **Descrições Genéricas (Comandos)** | Todas | 100% |
| **Aliases Ausentes** | 4 | 50% dos grupos com aliases |
| **Subgrupos Ausentes** | 1 | 12.5% dos grupos |
| **Arquitetura Divergente** | 1 grupo (3 subgrupos) | 12.5% dos grupos |
| **Comandos Extras** | 3 | ~7% dos comandos |
| **Flags Globais Ausentes** | 7 | 100% ausentes |

---

## 🎯 Prioridades de Correção

### 🔴 CRÍTICO (Bloqueante para Produção)
1. **Eliminar bugs visuais** ("Dqui1", "doto3", falta de espaços)
2. **Corrigir nomes de grupos** (`public-i-ps` → `public-ips`, `v-p-cs` → `vpcs`, `subnet-pools` → `subnetpools`)
3. **Restaurar arquitetura**:
   - Implementar subgrupo `security-groups rules`
   - Restaurar subgrupos em `vpcs` (`ports`, `public-ips`, `subnets`)
4. **Adicionar descrições em todas as flags** (0% → 100%)
5. **Restaurar flags globais ausentes** (7 flags)
6. **Normalizar nomes de flags** (adicionar hífens consistentemente)

### ⚠️ ALTO (Afeta Usabilidade)
1. **Corrigir nomes de comandos** (`attach-security-group` → `attach`, etc.)
2. **Remover flags incorretas como required** (em `list`)
3. **Adicionar aliases faltando** (4 grupos)
4. **Melhorar descrições dos comandos** (genéricas → específicas)

### 📝 MÉDIO (Aprimoramentos)
1. **Verificar comandos extras** (`rules create/list`, `vpcs rename`)
2. **Padronizar nomenclatura de paginação** (`items-per-page/page` vs `limit/offset`)

---

## 💡 Conclusão

O produto `network` apresenta os **problemas mais graves** identificados até agora:

### Problemas Únicos de `network`:
1. **Nomes de grupos corrompidos** (`public-i-ps`, `v-p-cs`)
2. **Arquitetura divergente em escala maior** (3 subgrupos de `vpcs` convertidos em flat)
3. **Comandos extras** em `rules` (pode ser confusão com `security-groups rules`)

### Problemas Compartilhados com Outros Produtos:
1. **Bugs visuais** ("Dqui1", "doto3")
2. **0% de descrições nas flags**
3. **Flags globais ausentes**
4. **Nomenclatura inconsistente** (remoção de hífens)

### Impacto Geral:
- **Compatibilidade**: 0% dos comandos são 100% compatíveis
- **Arquitetura**: 25% dos grupos têm divergências estruturais críticas
- **Usabilidade**: Severamente comprometida por bugs visuais e falta de documentação
- **Funcionalidade**: 1 subgrupo completamente ausente (perda de funcionalidade)

**Recomendação**: ❌ **NÃO APTO PARA PRODUÇÃO**. Requer refatoração completa do gerador de CLI antes de gerar novamente.

