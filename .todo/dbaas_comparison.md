# Comparação: mgc dbaas vs ./tmp-cli/cli dbaas

**Data:** 24/10/2025  
**Referência:** `mgc` (comandos oficiais que precisam ser replicados em `./tmp-cli/cli`)  
**Nota:** Flags globais não incluídas nesta análise (já documentadas anteriormente)

---

## 1. Comando Principal: `dbaas`

### MGC (Referência)
```
DBaaS API Product.

Aliases:
  dbaas, database, db

Commands:
  clusters         Clusters management.
  engines          Engines available for database instances.
  instance-types   Instance Types available for database instances.
  instances        Database instances management.
  parameter-groups Parameter groups management.
  replicas         Database replicas management.
  snapshots        Snapshots management.
```

### TMP-CLI (Atual)
```
Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.

Package dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
This package allows you to manage database instances, clusters, replicas, engines, instance types, and parameters.

Commands:
  clusters            Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  engines             Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  instance-types      Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  instances           Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  parameters          Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  parameters-group    Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  replicas            Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
```

### ❌ Problemas Identificados:
1. **Descrição divergente:**
   - MGC: "DBaaS API Product." (concisa)
   - TMP-CLI: Descrição repetitiva e verbosa com "Package dbaas..." redundante
2. **Aliases FALTANDO:** TMP-CLI não tem os aliases: `database`, `db`
3. **Subcomando FALTANDO:** `snapshots` não existe no nível principal do TMP-CLI
4. **Subcomando com NOME DIVERGENTE:** 
   - MGC: `parameter-groups` → TMP-CLI: `parameters-group` (ordem invertida)
5. **Subcomando EXTRA:** `parameters` não existe no MGC (é subcomando de parameter-groups e engines)
6. **Descrições dos subcomandos:**
   - MGC: Descrições específicas e informativas para cada comando
   - TMP-CLI: Mesma descrição genérica repetida para todos os comandos

---

## 2. Problema CRÍTICO: Estrutura de Snapshots Completamente Diferente

### MGC (Referência) - Estrutura Hierárquica
```
dbaas/
  ├── snapshots/
  │   ├── clusters-snapshots/
  │   │   ├── create
  │   │   ├── delete
  │   │   ├── get
  │   │   ├── list
  │   │   ├── restore
  │   │   └── update
  │   └── instances-snapshots/
  │       ├── create
  │       ├── delete
  │       ├── get
  │       ├── list
  │       ├── restore
  │       └── update
```

### TMP-CLI (Atual) - Comandos Misturados
```
dbaas/
  ├── instances/
  │   ├── create
  │   ├── create-snapshot       <-- Comando de snapshot MISTURADO
  │   ├── delete
  │   ├── delete-snapshot       <-- Comando de snapshot MISTURADO
  │   ├── get
  │   ├── get-snapshot          <-- Comando de snapshot MISTURADO
  │   ├── list
  │   ├── list-snapshots        <-- Comando de snapshot MISTURADO
  │   ├── list-all-snapshots    <-- Comando de snapshot MISTURADO
  │   ├── restore-snapshot      <-- Comando de snapshot MISTURADO
  │   ├── update-snapshot       <-- Comando de snapshot MISTURADO
  │   ├── resize
  │   ├── start
  │   ├── stop
  │   └── update
  └── (snapshots não existe como grupo separado)
```

### ❌ Problemas Gravíssimos:
1. **Grupo `snapshots` FALTANDO completamente** no nível principal
2. **Comandos de snapshots misturados com comandos de instances** no TMP-CLI
3. **Estrutura hierárquica perdida:** MGC separa clusters-snapshots de instances-snapshots
4. **Nomenclatura divergente:** MGC usa caminho hierárquico, TMP-CLI usa prefixo no nome
5. **12 comandos de snapshots** fora do lugar correto

---

## 3. Comando: `dbaas clusters`

### MGC (Referência)
```
Clusters management.

Commands:
  create      Creates a new database high availability cluster.
  delete      Deletes a database cluster.
  get         Database cluster details.
  list        List all database clusters.
  resize      Resizes a cluster database.
  start       Starts a database cluster.
  stop        Stops a database cluster.
  update      Database cluster update.
```

### TMP-CLI (Atual)
```
Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.

Dqui1   <-- ⚠️ BUG

Commands:
  create              Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  delete              Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  get                 Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  list                Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  list-all            Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  resize              Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  start               Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  stop                Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  update              Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Descrição divergente:** TMP-CLI usa descrição genérica
3. **Comando EXTRA:** `list-all` não existe no MGC
4. **Descrições dos subcomandos:** Todas genéricas ao invés de específicas

---

## 4. Comando: `dbaas instances` (com comandos de snapshot misturados)

### MGC (Referência)
```
Database instances management.

Commands:
  create      Creates a new database instance.
  delete      Deletes a database instance.
  get         Database instance details.
  list        List all database instances.
  resize      Resizes a database instance.
  start       Starts a database instance.
  stop        Stops a database instance.
  update      Database instance update.
```

### TMP-CLI (Atual)
```
Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.

Dqui1   <-- ⚠️ BUG

Commands:
  create              Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  create-snapshot     Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.  <-- EXTRA
  delete              Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  delete-snapshot     Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.  <-- EXTRA
  get                 Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  get-snapshot        Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.  <-- EXTRA
  list                Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  list-all            Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.  <-- EXTRA
  list-all-snapshots  Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.  <-- EXTRA
  list-snapshots      Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.  <-- EXTRA
  resize              Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  restore-snapshot    Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.  <-- EXTRA
  start               Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  stop                Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  update              Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  update-snapshot     Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.  <-- EXTRA
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Descrição divergente:** TMP-CLI usa descrição genérica
3. **8 Comandos de snapshot MISTURADOS** aqui (deveriam estar em `snapshots/instances-snapshots/`)
4. **Comando EXTRA:** `list-all` (padrão sistemático)
5. **Descrições dos subcomandos:** Todas genéricas ao invés de específicas

---

## 5. Comando: `dbaas replicas`

### MGC (Referência)
```
Database replicas management.

Commands:
  create      Replica Create.
  delete      Deletes a replica instance.
  get         Replica Detail.
  list        Replicas List.
  resize      Replica Resize.
  start       Replica Start.
  stop        Replica Stop.
```

### TMP-CLI (Atual)
```
Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.

Dqui1   <-- ⚠️ BUG

Commands:
  create              Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  delete              Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  get                 Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  list                Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  list-all            Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  resize              Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  start               Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  stop                Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Descrição divergente:** TMP-CLI usa descrição genérica
3. **Comando EXTRA:** `list-all`
4. **Descrições dos subcomandos:** Todas genéricas

---

## 6. Comando: `dbaas parameter-groups` vs `parameters-group`

### MGC (Referência)
```
Parameter groups management.

Commands:
  create      Creates a new parameter group.
  delete      Deletes a parameter group.
  get         Parameter group details.
  list        List all Parameter Groups
  parameters  parameters
  update      Parameter group update.
```

### TMP-CLI (Atual - comando `parameters-group`)
```
Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.

Dqui1   <-- ⚠️ BUG

Commands:
  create              Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  delete              Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  get                 Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  list                Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  list-all            Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  update              Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Nome do grupo divergente:** `parameter-groups` (MGC) vs `parameters-group` (TMP-CLI)
3. **Descrição divergente:** TMP-CLI usa descrição genérica
4. **Subcomando FALTANDO:** `parameters` não existe em parameters-group
5. **Comando EXTRA:** `list-all`
6. **Descrições dos subcomandos:** Todas genéricas

---

## 7. Grupo EXTRA: `dbaas parameters`

### MGC (Referência)
```
❌ NÃO EXISTE como grupo de nível 2

(parameters é subcomando de parameter-groups e engines)
```

### TMP-CLI (Atual)
```
Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.

Dqui1   <-- ⚠️ BUG

Commands:
  create              Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  delete              Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  list                Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  list-all            Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  update              Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
```

### ❌ Problema:
**Grupo extra que não deveria existir no nível 2** - parameters deveria ser subcomando, não grupo

---

## 8. Comando: `dbaas engines`

### MGC (Referência)
```
Engines available for database instances.

Commands:
  get         Engine detail.
  list        List available engines.
  parameters  List available engine parameters.
```

### TMP-CLI (Atual)
```
Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.

Dqui1   <-- ⚠️ BUG

Commands:
  get                 Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  list                Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  list-all            Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  list-engine-parametersDbaas provides...   <-- ⚠️ BUG: SEM ESPAÇO entre nome e descrição
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Bug gravíssimo:** `list-engine-parameters` SEM ESPAÇO antes da descrição
3. **Descrição divergente:** TMP-CLI usa descrição genérica
4. **Nome divergente:** MGC: `parameters` → TMP-CLI: `list-engine-parameters`
5. **Comando EXTRA:** `list-all`
6. **Descrições dos subcomandos:** Todas genéricas

---

## 9. Comando: `dbaas instance-types`

### MGC (Referência)
```
Instance Types available for database instances.

Aliases:
  instance-types, instance_types

Commands:
  get         Instance Type detail.
  list        List available instance types.
```

### TMP-CLI (Atual)
```
Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.

Dqui1   <-- ⚠️ BUG

Commands:
  get                 Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  list                Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
  list-all            Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Descrição divergente:** TMP-CLI usa descrição genérica
3. **Comando EXTRA:** `list-all`
4. **Aliases:** Não é possível verificar se foram implementados
5. **Descrições dos subcomandos:** Todas genéricas

---

## 10. Comando: `instances create`

### MGC (Referência)
```
Creates a new database instance asynchronously for a tenant.

Examples:
  mgc dbaas instances create --availability-zone="br-ne1-a" --volume.size=30

Flags:
  --availability-zone enum          Availability Zone (one of "br-ne1-a", "br-ne1-b", "br-se1-a", "br-se1-b" or "br-se1-c")
  --backup-retention-days integer   Backup Retention Days
  --backup-start-at time            Backup Start At
  --cli.list-links enum[=table]     List all available links
  --deletion-protected              Deletion Protected
  --engine-id uuid                  Engine Id (required)
  --instance-type-id uuid           Instance Type Id (required)
  --name string                     Name (max: 100) (required)
  --parameter-group-id uuid         Parameter group Id
  --password string                 Password (max: 50) (required)
  --security-groups array(uuid)     Security Group IDs
  --user string                     User (max: 25) (required)
  --volume object                   Instance Volume Request (properties: size and type) (required)
  --volume.size integer             The size of the volume (in GiB). (range: 10 - 50000)
  --volume.type enum                The type of the volume. (one of "CLOUD_NVME15K", "CLOUD_NVME20K", "CLOUD_NVME30K", "CLOUD_NVME40K" or "CLOUD_NVME50K")
```

### TMP-CLI (Atual)
```
Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.

doto3   <-- ⚠️ BUG

Flags:
  --availability-zone         (sem descrição)
  --backup-retention-days     (sem descrição)
  --backup-start-at           (sem descrição)
  --engine-id                 (sem descrição)
  --instance-type-id          (sem descrição)
  --name                      (required) (sem descrição)
  --parameter-group-id        (sem descrição)
  --password                  (required) (sem descrição)
  --user                      (required) (sem descrição)
  --volume.size               (required) (sem descrição)
  --volume.type               (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** Genérica ao invés de específica
3. **Flags FALTANDO:**
   - `--cli.list-links`
   - `--deletion-protected`
   - `--security-groups`
   - `--volume` (objeto pai)
4. **Flags sem descrição:** Todas as flags sem descrição
5. **Falta de Examples:** MGC tem exemplos, TMP-CLI não
6. **Marcação incorreta:** Algumas flags não são required mas estão marcadas

---

## 11. Comando: `instances list`

### MGC (Referência)
```
Returns a list of database instances for a x-tenant-id.

Examples:
  mgc dbaas instances list --status="ACTIVE"

Flags:
  --control.expand enum       Instance extra attributes or relations to show
  --control.limit integer     The maximum number of items per page. (range: 1 - 25)
  --control.offset integer    The number of items to skip (min: 0)
  --deletion-protected        Deletion Protected
  --engine-id uuid            Value referring to engine Id.
  --parameter-group-id uuid   Value referring to parameter group Id.
  --status enum               Value referring to instance status.
  --volume.size integer       Value referring to volume size.
  --volume.size-gt integer    Value referring to volume size greater than.
  --volume.size-gte integer   Value referring to volume size greater than or equal to.
  --volume.size-lt integer    Value referring to volume size less than.
  --volume.size-lte integer   Value referring to volume size less than or equal to.
```

### TMP-CLI (Atual)
```
Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.

doto3   <-- ⚠️ BUG

Usage: cli dbaas instances list [Offset] [Limit] [Status] [VolumeSizeLte] [ExpandedFields] [EngineID] [VolumeSize] [VolumeSizeGt] [VolumeSizeGte] [VolumeSizeLt] [flags]

Flags:
  --engine-id        (required) (sem descrição)
  --expanded-fields             (sem descrição)
  --limit            (required) (sem descrição)
  --offset           (required) (sem descrição)
  --status           (required) (sem descrição)
  --volume-size      (required) (sem descrição)
  --volume-size-gt   (required) (sem descrição)
  --volume-size-gte  (required) (sem descrição)
  --volume-size-lt   (required) (sem descrição)
  --volume-size-lte  (required) (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** Genérica ao invés de específica
3. **Argumentos posicionais incorretos:** Usage mostra flags como argumentos posicionais em camelCase
4. **Nomes de flags divergentes:**
   - MGC: `--control.expand` → TMP-CLI: `--expanded-fields`
   - MGC: `--control.limit` → TMP-CLI: `--limit` (sem prefixo)
   - MGC: `--control.offset` → TMP-CLI: `--offset` (sem prefixo)
5. **Flags FALTANDO:**
   - `--deletion-protected`
   - `--parameter-group-id`
6. **Flags sem descrição:** Todas as flags sem descrição
7. **Marcação incorreta de required:** TODAS as flags marcadas como required (incorreto)
8. **Falta de Examples:** MGC tem exemplos

---

## 12. Comando: `clusters list`

### MGC (Referência)
```
Returns a list of database clusters for a x-tenant-id.

Examples:
  mgc dbaas clusters list --status="ACTIVE"

Flags:
  --control.limit integer     The maximum number of items per page. (range: 1 - 25)
  --control.offset integer    The number of items to skip (min: 0)
  --deletion-protected        Deletion Protected
  --engine-id uuid            Value referring to engine Id.
  --parameter-group-id uuid   Value referring to parameter group Id.
  --status enum               Value referring to cluster status.
  --volume.size integer       Value referring to volume size.
  --volume.size-gt integer    Greater than.
  --volume.size-gte integer   Greater than or equal to.
  --volume.size-lt integer    Less than.
  --volume.size-lte integer   Less than or equal to.
```

### TMP-CLI (Atual)
```
Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.

doto3   <-- ⚠️ BUG

Usage: cli dbaas clusters list [Offset] [Limit] [VolumeSize] [VolumeSizeLte] [ParameterGroupID] [Status] [EngineID] [VolumeSizeGt] [VolumeSizeGte] [VolumeSizeLt] [flags]

Flags:
  --engine-id        (required) (sem descrição)
  --limit            (required) (sem descrição)
  --offset           (required) (sem descrição)
  --parameter-group-id  (required) (sem descrição)
  --status           (required) (sem descrição)
  --volume-size      (required) (sem descrição)
  --volume-size-gt   (required) (sem descrição)
  --volume-size-gte  (required) (sem descrição)
  --volume-size-lt   (required) (sem descrição)
  --volume-size-lte  (required) (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** Genérica ao invés de específica
3. **Argumentos posicionais incorretos:** Flags mostradas como args posicionais em camelCase
4. **Nomes de flags divergentes:**
   - MGC: `--control.limit` → TMP-CLI: `--limit`
   - MGC: `--control.offset` → TMP-CLI: `--offset`
5. **Flag FALTANDO:** `--deletion-protected`
6. **Flags sem descrição:** Todas as flags sem descrição
7. **Marcação incorreta de required:** TODAS marcadas como required
8. **Falta de Examples**

---

## 13. Comparação de Snapshots: `snapshots instances-snapshots list` vs `instances list-snapshots`

### MGC (Referência)
```
Caminho: mgc dbaas snapshots instances-snapshots list

List all snapshots.

Usage:
  mgc dbaas snapshots instances-snapshots list [instance-id] [flags]

Flags:
  --control.limit integer    The maximum number of items per page. (range: 1 - 50)
  --control.offset integer   The number of items to skip (min: 0)
  --instance-id uuid         Value referring to instance Id. (required)
  --status enum              Value referring to snapshot status.
  --type enum                Snapshot Type Filter
```

### TMP-CLI (Atual)
```
Caminho: cli dbaas instances list-snapshots

Dbaas provides a client for interacting with the Magalu Cloud Database as a Service (DBaaS) API.

doto3   <-- ⚠️ BUG

Usage: cli dbaas instances list-snapshots [instanceID] [flags]

Flags:
  --instance-id      (required) (sem descrição)
  --limit                       (sem descrição)
  --offset                      (sem descrição)
  --status                      (sem descrição)
  --type                        (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Estrutura completamente diferente:** 
   - MGC: `snapshots/instances-snapshots/list`
   - TMP-CLI: `instances/list-snapshots`
3. **Descrição do comando:** Genérica ao invés de específica
4. **Nomes de flags divergentes:**
   - MGC: `--control.limit` → TMP-CLI: `--limit`
   - MGC: `--control.offset` → TMP-CLI: `--offset`
5. **Flags sem descrição:** Todas as flags sem descrição
6. **Argumentos posicionais:** MGC usa kebab-case `[instance-id]`, TMP-CLI usa camelCase `[instanceID]`

---

## Resumo Geral de Problemas

### 🐛 BUGS CRÍTICOS:
- String **"Dqui1"** aparece em todos os 7 subcomandos de nível 2
- String **"doto3"** aparece em todos os comandos leaf
- String **"list-engine-parametersDbaas"** SEM ESPAÇO (bug de formatação)

### 🏗️ ESTRUTURA COMPLETAMENTE DIVERGENTE:
- **Snapshots:** Grupo inteiro faltando no nível principal, comandos misturados em instances
- **Hierarquia perdida:** MGC separa clusters-snapshots de instances-snapshots
- **12 comandos de snapshots** no lugar errado

### ❌ NOMENCLATURA DIVERGENTE:
- `parameter-groups` → TMP-CLI: `parameters-group` (ordem invertida)
- `parameters` (subcomando) → TMP-CLI: grupo de nível 2 (EXTRA)
- `engines/parameters` → TMP-CLI: `engines/list-engine-parameters`
- `control.expand` → TMP-CLI: `expanded-fields`

### 🚫 ALIASES FALTANDO:
- `database`, `db` (no comando principal)
- `instance_types` (em instance-types)

### ➕ COMANDOS/GRUPOS EXTRAS:
- `list-all` em: clusters, engines, instance-types, instances, parameters, parameters-group, replicas (7 grupos)
- `parameters` como grupo de nível 2 (não deveria existir)

### 📝 DESCRIÇÕES:
- Todas as descrições são genéricas e repetitivas
- Faltam descrições de todas as flags
- Faltam Examples em comandos que os têm no MGC

### 🔧 FLAGS:
#### Padrão de problemas em TODOS os comandos list:
- Falta prefixo `control.` em: `limit`, `offset`, `expand`
- Flags marcadas incorretamente como `required` (100% dos list)
- Flags sem descrição (100%)
- Argumentos posicionais incorretos (camelCase) no Usage

#### Flags FALTANDO (padrão):
- `--cli.list-links` (em comandos create)
- `--deletion-protected` (em vários list)
- `--parameter-group-id` (em instances list)
- `--security-groups` (em instances create)
- Objetos pai: `--volume` (apenas propriedades leaf existem)

---

## 📊 Estatísticas

### Subcomandos (Nível 2)
- **Total MGC:** 7 grupos
- **Total TMP-CLI:** 7 grupos (mas 1 faltando, 1 extra, 1 nome divergente)
- **Faltando:** 1 grupo (snapshots)
- **Nome divergente:** 1 grupo (parameter-groups vs parameters-group)
- **Extra:** 1 grupo (parameters)
- **Bugs "Dqui1":** 100% (7/7)

### Comandos Leaf - Situação Complexa
- **MGC snapshots:** 12 comandos (6 instances + 6 clusters)
- **TMP-CLI:** Comandos de snapshots misturados em instances (8 comandos)
- **Comandos extras `list-all`:** 7+ comandos
- **Bugs "doto3":** 100% (todos os comandos leaf)

### Aliases
- **Total MGC (principal):** 3 (dbaas, database, db)
- **Total TMP-CLI:** 1 (só o nome principal)
- **Compatibilidade:** 33% (1/3)

### Flags
- **Sem prefixo `control.`:** 100% dos list
- **Sem descrição:** 100%
- **Required incorreto:** 100% dos list (TODAS marcadas incorretamente)
- **Nomenclatura divergente:** ~20% (expand, limit, offset)

---

## ✅ Checklist de Ações

### P0 - Crítico (Arquitetura)
- [ ] **REESTRUTURAR SNAPSHOTS:** Criar grupo `snapshots` no nível principal
- [ ] **Mover comandos:** Tirar 8 comandos de snapshot de `instances` para `snapshots/instances-snapshots/`
- [ ] **Adicionar cluster snapshots:** Implementar `snapshots/clusters-snapshots/` com 6 comandos
- [ ] Remover string "Dqui1" de todos os subcomandos nível 2 (7 ocorrências)
- [ ] Remover string "doto3" de todos os comandos leaf
- [ ] Corrigir bug "list-engine-parametersDbaas" (falta espaço)
- [ ] Adicionar descrições em TODAS as flags

### P1 - Alto (Nomenclatura e Compatibilidade)
- [ ] Renomear `parameters-group` → `parameter-groups`
- [ ] Remover grupo `parameters` do nível 2 (deve ser subcomando)
- [ ] Renomear `list-engine-parameters` → `parameters` (em engines)
- [ ] Renomear `expanded-fields` → `control.expand`
- [ ] Adicionar prefixo `control.` em: limit, offset, expand
- [ ] Adicionar aliases: `database`, `db`
- [ ] Corrigir marcação de flags required (remover de TODAS as flags de list)
- [ ] Corrigir descrições de comandos (específicas, não genéricas)
- [ ] Remover argumentos posicionais incorretos do Usage
- [ ] Padronizar formato argumentos (kebab-case, não camelCase)

### P2 - Médio (Funcionalidades)
- [ ] Adicionar `--cli.list-links` em comandos de criação
- [ ] Adicionar flag `--deletion-protected` onde aplicável
- [ ] Adicionar flag `--parameter-group-id` em instances list
- [ ] Adicionar flag `--security-groups` em instances create
- [ ] Adicionar objetos pai completos (volume, etc.)
- [ ] Adicionar Examples nos comandos apropriados
- [ ] Remover comandos `list-all` (7+ comandos)

### P3 - Baixo (Polish)
- [ ] Melhorar formatação geral do help
- [ ] Padronizar estrutura de descrições

