# Comparação: mgc kubernetes vs ./tmp-cli/cli kubernetes

**Data:** 24/10/2025  
**Referência:** `mgc` (comandos oficiais que precisam ser replicados em `./tmp-cli/cli`)  
**Nota:** Flags globais não incluídas nesta análise (já documentadas anteriormente)

---

## 1. Comando Principal: `kubernetes`

### MGC (Referência)
```
APIs related to the Kubernetes product.

Aliases:
  kubernetes, k8s, kube, kub

Commands:
  cluster     Endpoints related to the creation, listing, deletion, and retrieval of the kubeconfig for a Kubernetes cluster.
  flavor      Endpoints related to listing available flavors for a Kubernetes cluster.
  info        Endpoints related to listing available flavors (resource configurations) and Kubernetes versions for a Kubernetes cluster.
  nodepool    Endpoints related to the creation, listing of nodepools and nodes, updating, and deletion of nodepools for a Kubernetes cluster.
  version     Endpoints related to listing available Kubernetes versions for a Kubernetes cluster.
```

### TMP-CLI (Atual)
```
Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.

Package kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
This package allows you to manage Kubernetes clusters, node pools, flavors, and versions.

Commands:
  clusters            Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
  flavors             Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
  nodepools           Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
  versions            Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
```

### ❌ Problemas Identificados:
1. **Descrição divergente:**
   - MGC: "APIs related to the Kubernetes product." (concisa)
   - TMP-CLI: Descrição repetitiva e verbosa com "Package kubernetes..." redundante
2. **Aliases FALTANDO:** TMP-CLI não tem os aliases: `k8s`, `kube`, `kub`
3. **Grupo FALTANDO:** `info` não existe no TMP-CLI
4. **Nomenclatura divergente (Singular vs Plural):**
   - MGC: `cluster` (singular) → TMP-CLI: `clusters` (plural)
   - MGC: `flavor` (singular) → TMP-CLI: `flavors` (plural)
   - MGC: `nodepool` (singular) → TMP-CLI: `nodepools` (plural)
   - MGC: `version` (singular) → TMP-CLI: `versions` (plural)
5. **Descrições dos subcomandos:**
   - MGC: Descrições específicas e informativas para cada comando
   - TMP-CLI: Mesma descrição genérica repetida para todos os comandos

---

## 2. Subcomando FALTANDO: `info`

### MGC (Referência)
```
Endpoints related to listing available flavors (resource configurations) and Kubernetes versions for a Kubernetes cluster.

Commands:
  flavors     Lists all available flavors
  versions    Lists all available versions
```

### TMP-CLI (Atual)
```
❌ SUBCOMANDO COMPLETAMENTE AUSENTE

Os comandos equivalentes estão espalhados:
- info/flavors → flavor/list (mas flavor é grupo de nível 2)
- info/versions → version/list (mas version é grupo de nível 2)
```

### ❌ Problema:
**Grupo organizacional `info` ausente**, que agrupa comandos de informação sobre recursos disponíveis.

---

## 3. Comando: `cluster` vs `clusters`

### MGC (Referência)
```
Endpoints related to the creation, listing, deletion, and retrieval of the kubeconfig for a Kubernetes cluster.

Commands:
  create      Create a cluster
  delete      Delete a cluster
  get         List a cluster by uuid
  kubeconfig  Get kubeconfig cluster
  list        List all clusters
  update      Patches a cluster
```

### TMP-CLI (Atual - comando `clusters`)
```
Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.

Dqui1   <-- ⚠️ BUG

Commands:
  create              Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
  delete              Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
  get                 Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
  get-kube-config     Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
  list                Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
  update              Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Nome do grupo divergente:** `cluster` (MGC) vs `clusters` (TMP-CLI)
3. **Descrição divergente:** TMP-CLI usa descrição genérica
4. **Nome de comando divergente:** 
   - MGC: `kubeconfig` → TMP-CLI: `get-kube-config`
5. **Descrições dos subcomandos:** Todas genéricas ao invés de específicas

---

## 4. Comando: `nodepool` vs `nodepools`

### MGC (Referência)
```
Endpoints related to the creation, listing of nodepools and nodes, updating, and deletion of nodepools for a Kubernetes cluster.

Commands:
  create      Create a node pool
  delete      Delete node pool by node_pool_id
  get         Get node pool by node_pool_id
  list        V1alpha0 List node pools by cluster_id
  nodes       List nodes from a node pool by node_pool_id
  update      Patch node pool replicas by node_pool_id
```

### TMP-CLI (Atual - comando `nodepools`)
```
Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.

Dqui1   <-- ⚠️ BUG

Commands:
  create              Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
  delete              Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
  get                 Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
  list                Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
  nodes               Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
  update              Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Nome do grupo divergente:** `nodepool` (MGC) vs `nodepools` (TMP-CLI)
3. **Descrição divergente:** TMP-CLI usa descrição genérica
4. **Descrições dos subcomandos:** Todas genéricas ao invés de específicas

---

## 5. Comando: `flavor` vs `flavors`

### MGC (Referência)
```
Endpoints related to listing available flavors for a Kubernetes cluster.

Commands:
  list        Lists all available flavors
```

### TMP-CLI (Atual - comando `flavors`)
```
Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.

Dqui1   <-- ⚠️ BUG

Commands:
  list                Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Nome do grupo divergente:** `flavor` (MGC) vs `flavors` (TMP-CLI)
3. **Descrição divergente:** TMP-CLI usa descrição genérica
4. **Descrições dos subcomandos:** Genéricas ao invés de específicas

---

## 6. Comando: `version` vs `versions`

### MGC (Referência)
```
Endpoints related to listing available Kubernetes versions for a Kubernetes cluster.

Commands:
  list        Lists all available versions
```

### TMP-CLI (Atual - comando `versions`)
```
Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.

Dqui1   <-- ⚠️ BUG

Commands:
  list                Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Nome do grupo divergente:** `version` (MGC) vs `versions` (TMP-CLI)
3. **Descrição divergente:** TMP-CLI usa descrição genérica
4. **Descrições dos subcomandos:** Genéricas ao invés de específicas

---

## 7. Comando: `cluster create` vs `clusters create`

### MGC (Referência)
```
Creates a Kubernetes cluster in Magalu Cloud.

Examples:
  mgc kubernetes cluster create --allowed-cidrs='["192.168.1.0/24","10.0.0.0/16"]' --name="cluster-example" --node-pools='[...]' --version="v1.32.3"

Flags:
  --allowed-cidrs array(string)   List of allowed CIDR blocks for API server access.
  --cli.list-links enum[=table]   List all available links for this command
  --cluster-ipv4-cidr string      The IP address CIDR used by the Pods in the cluster.
  --description string            A brief description of the Kubernetes cluster.
  --enabled-bastion               [Deprecated] Enables the use of a bastion host
  --enabled-server-group          Enables the use of a server group with anti-affinity policy
  --name string                   Kubernetes cluster name. (required)
  --node-pools array(object)      An array representing a set of nodes
  --services-ipv4-cidr string     The IPv4 subnet CIDR used by Kubernetes Services.
  --version string                The Kubernetes version for the cluster
  --zone string                   [Deprecated] Identifier of the zone
```

### TMP-CLI (Atual)
```
Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.

doto3   <-- ⚠️ BUG

Flags:
  --allowed-cidrs         (sem descrição)
  --cluster-ipv4cidr      (sem descrição)
  --description           (sem descrição)
  --enabled-server-group  (sem descrição)
  --name                  (required) (sem descrição)
  --node-pools            (sem descrição)
  --services-ip-v4cidr    (sem descrição)
  --version               (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** Genérica ao invés de específica
3. **Flags FALTANDO:**
   - `--cli.list-links`
   - `--enabled-bastion`
   - `--zone`
4. **Flags com nome divergente:**
   - MGC: `--cluster-ipv4-cidr` → TMP-CLI: `--cluster-ipv4cidr` (sem hífen)
   - MGC: `--services-ipv4-cidr` → TMP-CLI: `--services-ip-v4cidr` (hífen em lugar diferente)
5. **Flags sem descrição:** Todas as flags sem descrição
6. **Falta de Examples:** MGC tem exemplos completos

---

## 8. Comando: `cluster list` vs `clusters list`

### MGC (Referência)
```
Lists all clusters for a user.

Flags: (nenhuma flag local - só as globais)
```

### TMP-CLI (Atual)
```
Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.

doto3   <-- ⚠️ BUG

Usage: cli kubernetes clusters list [Expand] [Limit] [Offset] [Sort] [flags]

Flags:
  --expand          (sem descrição)
  --limit            (required) (sem descrição)
  --offset           (required) (sem descrição)
  --sort             (required) (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** Genérica ao invés de específica
3. **Argumentos posicionais incorretos:** Usage mostra flags como argumentos posicionais
4. **Flags EXTRAS (não existem no MGC):**
   - `--expand`
   - `--limit`
   - `--offset`
   - `--sort`
5. **Flags sem descrição:** Todas as flags sem descrição
6. **Marcação incorreta de required:** Todas as flags marcadas como required

---

## 9. Comando: `cluster kubeconfig` vs `clusters get-kube-config`

### MGC (Referência)
```
Get kubeconfig cluster

Usage:
  mgc kubernetes cluster kubeconfig [cluster-id] [flags]

Flags:
  --cluster-id uuid   Cluster's UUID. (required)
```

### TMP-CLI (Atual - comando `get-kube-config`)
```
Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.

doto3   <-- ⚠️ BUG

Usage: cli kubernetes clusters get-kube-config [clusterID] [flags]

Flags:
  --cluster-id      (required) (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Nome do comando divergente:** `kubeconfig` (MGC) vs `get-kube-config` (TMP-CLI)
3. **Descrição do comando:** Genérica ao invés de específica
4. **Argumentos posicionais:** MGC usa kebab-case `[cluster-id]`, TMP-CLI usa camelCase `[clusterID]`
5. **Flag sem descrição:** `--cluster-id` sem descrição

---

## 10. Comando: `nodepool create` vs `nodepools create`

### MGC (Referência - primeiras linhas)
```
Creates a node pool in a Kubernetes cluster.

Usage:
  mgc kubernetes nodepool create [cluster-id] [flags]

Examples:
  mgc kubernetes nodepool create --auto-scale.max-replicas=5 --availability-zones='["a","b","c"]' --flavor="cloud-k8s.gp1.small" --name="nodepool-example" --replicas=3

Flags:
  --auto-scale object                 Object specifying properties for updating workload resources
  --auto-scale.max-replicas integer   Maximum number of replicas for autoscaling
  --auto-scale.min-replicas integer   Minimum number of replicas for autoscaling
  --availability-zones array(enum)    List of availability zones
  --cli.list-links enum[=table]       List all available links
  --cluster-id uuid                   Cluster's UUID. (required)
  --flavor string                     Definition of the CPU, RAM, and storage capacity (required)
  --max-pods-per-node integer         Maximum number of Pods allowed per node (range: 8 - 110)
  --name string                       Name of the node pool (required)
  --replicas integer                  Number of replicas (min: 0)
  --tags array(string)                Labels for organizing resources
  --taints array(object)              Taints to apply to nodes
```

### TMP-CLI (Atual - comando `nodepools create`)
```
Kubernetes provides a client for interacting with the Magalu Cloud Kubernetes API.

doto3   <-- ⚠️ BUG

Flags:
  --auto-scale.max-replicas     (sem descrição)
  --auto-scale.min-replicas     (sem descrição)
  --availability-zones          (sem descrição)
  --cluster-id                  (sem descrição)
  --flavor                      (required) (sem descrição)
  --max-pods-per-node           (sem descrição)
  --name                        (required) (sem descrição)
  --replicas                    (sem descrição)
  --tags                        (sem descrição)
  --taints                      (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** Genérica ao invés de específica
3. **Flags FALTANDO:**
   - `--cli.list-links`
   - `--auto-scale` (objeto pai)
4. **Flags sem descrição:** Todas as flags sem descrição (MGC tem descrições detalhadas incluindo tabela de flavors)
5. **Falta de Examples:** MGC tem exemplos

---

## Resumo Geral de Problemas

### 🐛 BUGS CRÍTICOS:
- String **"Dqui1"** aparece em todos os 4 subcomandos de nível 2
- String **"doto3"** aparece em todos os comandos leaf

### 📝 NOMENCLATURA DIVERGENTE (Singular vs Plural):
Padrão **SISTEMÁTICO** de pluralização incorreta em **TODOS** os grupos:
- `cluster` → TMP-CLI: `clusters`
- `flavor` → TMP-CLI: `flavors`
- `nodepool` → TMP-CLI: `nodepools`
- `version` → TMP-CLI: `versions`

**Impacto:** 100% dos grupos têm nome divergente

### 🔄 NOMES DE COMANDOS DIVERGENTES:
- `kubeconfig` → TMP-CLI: `get-kube-config`

### 🔧 NOMES DE FLAGS DIVERGENTES:
- `--cluster-ipv4-cidr` → TMP-CLI: `--cluster-ipv4cidr` (sem hífen interno)
- `--services-ipv4-cidr` → TMP-CLI: `--services-ip-v4cidr` (hífen em lugar diferente)

### ❌ GRUPO FALTANDO:
- `info` (grupo organizacional) não existe no TMP-CLI

### 🚫 ALIASES FALTANDO:
- `k8s`, `kube`, `kub` (no comando principal)

### ➕ FLAGS EXTRAS (que não existem no MGC):
- `--expand`, `--limit`, `--offset`, `--sort` em `cluster list`

### 📝 DESCRIÇÕES:
- Todas as descrições são genéricas e repetitivas
- Faltam descrições de todas as flags
- Faltam Examples em comandos que os têm no MGC

### 🔧 FLAGS FALTANDO:
- `--cli.list-links` (em comandos create)
- `--enabled-bastion` (em cluster create)
- `--zone` (em cluster create)
- Objetos pai: `--auto-scale` (apenas propriedades leaf existem)

---

## 📊 Estatísticas

### Subcomandos (Nível 2)
- **Total MGC:** 5 grupos
- **Total TMP-CLI:** 4 grupos
- **Faltando:** 1 grupo (info)
- **Nome divergente:** 100% (4/4 dos grupos presentes - todos com plural incorreto)
- **Bugs "Dqui1":** 100% (4/4)

### Comandos Leaf
- **Comandos presentes:** Todos os comandos principais existem
- **Nome divergente:** 1 comando (kubeconfig vs get-kube-config)
- **Bugs "doto3":** 100% (todos os comandos leaf)

### Aliases
- **Total MGC (principal):** 4 (kubernetes, k8s, kube, kub)
- **Total TMP-CLI:** 1 (só o nome principal)
- **Compatibilidade:** 25% (1/4)

### Flags
- **Sem descrição:** 100%
- **Nome divergente:** ~10% (cluster-ipv4-cidr, services-ipv4-cidr)
- **Extras não solicitadas:** 4 flags em cluster list
- **Required incorreto:** 100% em cluster list

---

## ✅ Checklist de Ações

### P0 - Crítico (Nomenclatura Sistemática)
- [ ] **Renomear TODOS os grupos para singular:**
  - `clusters` → `cluster`
  - `flavors` → `flavor`
  - `nodepools` → `nodepool`
  - `versions` → `version`
- [ ] Remover string "Dqui1" de todos os subcomandos nível 2 (4 ocorrências)
- [ ] Remover string "doto3" de todos os comandos leaf
- [ ] Adicionar descrições em TODAS as flags

### P1 - Alto (Compatibilidade)
- [ ] **Implementar grupo `info` com subcomandos:**
  - `info/flavors`
  - `info/versions`
- [ ] Renomear `get-kube-config` → `kubeconfig`
- [ ] Corrigir nomes de flags:
  - `--cluster-ipv4cidr` → `--cluster-ipv4-cidr`
  - `--services-ip-v4cidr` → `--services-ipv4-cidr`
- [ ] Adicionar aliases: `k8s`, `kube`, `kub`
- [ ] Corrigir descrições de comandos (específicas, não genéricas)
- [ ] Remover argumentos posicionais incorretos do Usage
- [ ] Padronizar formato argumentos (kebab-case, não camelCase)

### P2 - Médio (Funcionalidades)
- [ ] Adicionar `--cli.list-links` em comandos de criação
- [ ] Adicionar flags faltantes:
  - `--enabled-bastion` (cluster create)
  - `--zone` (cluster create)
  - Objetos pai: `--auto-scale`
- [ ] Remover flags extras de `cluster list`:
  - `--expand`
  - `--limit`
  - `--offset`
  - `--sort`
- [ ] Adicionar Examples nos comandos apropriados
- [ ] Corrigir marcação de required em flags de list

### P3 - Baixo (Polish)
- [ ] Melhorar formatação geral do help
- [ ] Padronizar estrutura de descrições

