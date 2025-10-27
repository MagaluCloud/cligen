# Comparação: mgc block-storage vs ./tmp-cli/cli block-storage

**Data:** 24/10/2025  
**Referência:** `mgc` (comandos oficiais que precisam ser replicados em `./tmp-cli/cli`)  
**Nota:** Flags globais não incluídas nesta análise (já documentadas em audit_comparison.md)

---

## 1. Comando Principal: `block-storage`

### MGC (Referência)
```
Block Storage API Product

Aliases:
  block-storage, bs, blocks, blst, block, volumes

Commands:
  schedulers   Operations with schedulers for snapshot creation and retention.
  snapshots    Operations with snapshots for volumes.
  volume-types Operations with volume types for volumes.
  volumes      Operations with volumes, including create, delete, extend, retype, list and other actions.
```

### TMP-CLI (Atual)
```
Blockstorage provides functionality to interact with the MagaluCloud block storage service.

Package blockstorage provides functionality to interact with the MagaluCloud block storage service.
This package allows managing volumes, volume types, and snapshots.

Commands:
  schedulers          Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  snapshots           Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  volume-types        Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  volumes             Blockstorage provides functionality to interact with the MagaluCloud block storage service.
```

### ❌ Problemas Identificados:
1. **Descrição divergente:**
   - MGC: "Block Storage API Product" (concisa)
   - TMP-CLI: Descrição repetitiva e verbosa com "Package blockstorage..." redundante
2. **Aliases FALTANDO:** TMP-CLI não tem os aliases: `bs`, `blocks`, `blst`, `block`, `volumes`
3. **Descrições dos subcomandos:**
   - MGC: Descrições específicas e informativas para cada comando
   - TMP-CLI: Mesma descrição genérica repetida para todos os comandos

---

## 2. Comando: `block-storage schedulers`

### MGC (Referência)
```
Operations with schedulers for snapshot creation and retention.

Commands:
  attach      Attach volume on scheduler
  create      Create a scheduler.
  delete      Delete a scheduler.
  detach      Detach volume on scheduler
  get         Retrieve the details of a specific scheduler.
  list        List all schedulers.
```

### TMP-CLI (Atual)
```
Blockstorage provides functionality to interact with the MagaluCloud block storage service.

Dqui1   <-- ⚠️ BUG

Commands:
  attach-volume       Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  create              Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  delete              Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  detach-volume       Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  get                 Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  list                Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  list-all            Blockstorage provides functionality to interact with the MagaluCloud block storage service.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Descrição divergente:** TMP-CLI usa descrição genérica
3. **Nomes de comandos divergentes:**
   - MGC: `attach` → TMP-CLI: `attach-volume`
   - MGC: `detach` → TMP-CLI: `detach-volume`
4. **Comando EXTRA:** `list-all` não existe no MGC
5. **Descrições dos subcomandos:** Todas genéricas ao invés de específicas

---

## 3. Comando: `block-storage snapshots`

### MGC (Referência)
```
Operations with snapshots for volumes.

Commands:
  copy        Copy a object snapshot to another region.
  create      Create a snapshot.
  delete      Delete a snapshot.
  get         Retrieve the details of a specific snapshot.
  list        List all snapshots.
  rename      Rename a snapshot.
```

### TMP-CLI (Atual)
```
Blockstorage provides functionality to interact with the MagaluCloud block storage service.

Dqui1   <-- ⚠️ BUG

Commands:
  create              Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  delete              Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  get                 Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  list                Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  list-all            Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  rename              Blockstorage provides functionality to interact with the MagaluCloud block storage service.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Descrição divergente:** TMP-CLI usa descrição genérica
3. **Comando FALTANDO:** `copy` não existe no TMP-CLI
4. **Comando EXTRA:** `list-all` não existe no MGC
5. **Descrições dos subcomandos:** Todas genéricas ao invés de específicas

---

## 4. Comando: `block-storage volume-types`

### MGC (Referência)
```
Operations with volume types for volumes.

Commands:
  list        List all volume types.
```

### TMP-CLI (Atual)
```
Blockstorage provides functionality to interact with the MagaluCloud block storage service.

Dqui1   <-- ⚠️ BUG

Commands:
  list                Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  list-all            Blockstorage provides functionality to interact with the MagaluCloud block storage service.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Descrição divergente:** TMP-CLI usa descrição genérica
3. **Comando EXTRA:** `list-all` não existe no MGC
4. **Descrições dos subcomandos:** Genéricas ao invés de específicas

---

## 5. Comando: `block-storage volumes`

### MGC (Referência)
```
Operations with volumes, including create, delete, extend, retype, list and other actions.

Commands:
  attach      Attach the volume to an instance.
  create      Create a new volume.
  delete      Delete a volume.
  detach      Detach the volume from an instance.
  extend      Extend the size of a volume.
  get         Retrieve the details of a specific volume.
  list        List all volumes.
  rename      Rename a volume.
  retype      Change the type of a volume.
```

### TMP-CLI (Atual)
```
Blockstorage provides functionality to interact with the MagaluCloud block storage service.

Dqui1   <-- ⚠️ BUG

Commands:
  attach              Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  create              Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  delete              Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  detach              Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  extend              Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  get                 Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  list                Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  list-all            Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  rename              Blockstorage provides functionality to interact with the MagaluCloud block storage service.
  retype              Blockstorage provides functionality to interact with the MagaluCloud block storage service.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Descrição divergente:** TMP-CLI usa descrição genérica
3. **Comando EXTRA:** `list-all` não existe no MGC
4. **Descrições dos subcomandos:** Todas genéricas ao invés de específicas

---

## 6. Comando: `block-storage schedulers create`

### MGC (Referência)
```
Creates a schedule for snapshot creation.

Flags:
  --cli.list-links enum[=table]              List all available links for this command
  --description string                       Description
  --name string                              Name (required)
  --policy object                            Policy (properties: frequency and retention_in_days) (required)
  --policy.frequency object                  Policy: Frequency (single property: daily)
  --policy.frequency.daily object            Frequency: DailyFrequency (single property: start_time)
  --policy.frequency.daily.start-time time   DailyFrequency: Start Time
  --policy.retention-in-days integer         Policy: Retention In Days (min: 1)
  --snapshot object                          Snapshot (single property: type) (required)
  --snapshot.type enum                       Snapshot: SnapshotType (one of "instant" or "object")
```

### TMP-CLI (Atual)
```
Blockstorage provides functionality to interact with the MagaluCloud block storage service.

doto3   <-- ⚠️ BUG

Flags:
  --description     (sem descrição)
  --name            (required) (sem descrição)
  --policy.frequency.daily.start-time  (required) (sem descrição)
  --policy.retention-in-days           (required) (sem descrição)
  --snapshot.type                      (required) (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** Genérica ao invés de específica
3. **Flags FALTANDO:**
   - `--cli.list-links`
   - `--policy` (objeto pai)
   - `--policy.frequency` (objeto intermediário)
   - `--policy.frequency.daily` (objeto intermediário)
   - `--snapshot` (objeto pai)
4. **Flags sem descrição:** Todas as flags aparecem sem descrição
5. **Marcação incorreta de required:** `--policy.retention-in-days` e `--policy.frequency.daily.start-time` marcados como required, mas no MGC fazem parte de objetos opcionais

---

## 7. Comando: `block-storage schedulers list`

### MGC (Referência)
```
Retrieve a list of Schedulers for the currently authenticated tenant.

#### Notes
- Use the expand argument to obtain additional details about the Volume.

Flags:
  --control.limit integer     Limit
  --control.offset integer    Offset
  --control.sort string       Sort (pattern: ^(^[\w-]+:(asc|desc)(,[\w-]+:(asc|desc))*)?$)
  --expand array(enum)       Expand
```

### TMP-CLI (Atual)
```
Blockstorage provides functionality to interact with the MagaluCloud block storage service.

doto3   <-- ⚠️ BUG

Usage: cli block-storage schedulers list [Limit] [Offset] [Sort] [Expand] [flags]

Flags:
  --expand          (sem descrição)
  --limit            (required) (sem descrição)
  --offset           (required) (sem descrição)
  --sort             (required) (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** Genérica ao invés de específica (falta a seção Notes)
3. **Argumentos posicionais incorretos:** Usage mostra `[Limit] [Offset] [Sort] [Expand]` como se fossem argumentos posicionais
4. **Nomes de flags divergentes:**
   - MGC: `--control.limit` → TMP-CLI: `--limit` (sem prefixo)
   - MGC: `--control.offset` → TMP-CLI: `--offset` (sem prefixo)
   - MGC: `--control.sort` → TMP-CLI: `--sort` (sem prefixo)
5. **Flags sem descrição:** Todas as flags sem descrição
6. **Marcação incorreta de required:** Todas as flags marcadas como required quando não deveriam ser

---

## 8. Comando: `block-storage schedulers attach`

### MGC (Referência)
```
Attach volume on scheduler.

Usage:
  mgc block-storage schedulers attach [id] [flags]

Examples:
  mgc block-storage schedulers attach --volume.id="..." --volume.name="..."

Flags:
  --cli.list-links enum[=table]   List all available links for this command
  --id uuid                       Id (required)
  --volume object                 Volume (at least one of: id or name) (required)
  --volume.id string              Volume: Id (min character count: 1)
  --volume.name string            Volume: Name (between 1 and 255 characters)
```

### TMP-CLI (Atual - comando `attach-volume`)
```
Blockstorage provides functionality to interact with the MagaluCloud block storage service.

doto3   <-- ⚠️ BUG

Usage: cli block-storage schedulers attach-volume [id] [flags]

Flags:
  --id               (required) (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Nome do comando divergente:** MGC: `attach` vs TMP-CLI: `attach-volume`
3. **Descrição do comando:** Genérica ao invés de específica
4. **Flags FALTANDO:**
   - `--cli.list-links`
   - `--volume` (objeto)
   - `--volume.id`
   - `--volume.name`
5. **Flags sem descrição:** Flag `--id` sem descrição
6. **Falta de examples:** MGC tem exemplos, TMP-CLI não

---

## 9. Comando: `block-storage snapshots create`

### MGC (Referência)
```
Create a Snapshot for the currently authenticated tenant.

The Snapshot can be used when it reaches the "available" state and the "completed" status.

#### Rules
- The Snapshot name must be unique; otherwise, the creation will be disallowed.
- Creating Snapshots from restored Volumes may lead to future conflicts...

#### Notes
- Use the **block-storage volume list** command to retrieve a list of all Volumes...

Examples:
  mgc block-storage snapshots create --source-snapshot.id="..." --volume.id="..."

Flags:
  --cli.list-links enum[=table]   List all available links
  --description string            Description (required)
  --name string                   Name (between 3 and 50 characters) (required)
  --source-snapshot object        Source Snapshot (at least one of: id or name)
  --source-snapshot.id string     Source Snapshot: Id (min character count: 1)
  --source-snapshot.name string   Source Snapshot: Name (between 1 and 255 characters)
  --type enum                     SnapshotType (one of "instant" or "object")
  --volume object                 Volume (at least one of: id or name)
  --volume.id string              Volume: Id (min character count: 1)
  --volume.name string            Volume: Name (between 1 and 255 characters)
```

### TMP-CLI (Atual)
```
Blockstorage provides functionality to interact with the MagaluCloud block storage service.

doto3   <-- ⚠️ BUG

Flags:
  --description      (required) (sem descrição)
  --name             (required) (sem descrição)
  --type             (required) (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** Genérica ao invés de descrição detalhada com Rules e Notes
3. **Flags FALTANDO:**
   - `--cli.list-links`
   - `--source-snapshot` (objeto)
   - `--source-snapshot.id`
   - `--source-snapshot.name`
   - `--volume` (objeto)
   - `--volume.id`
   - `--volume.name`
4. **Flags sem descrição:** Todas as flags sem descrição
5. **Falta de examples:** MGC tem exemplos, TMP-CLI não

---

## 10. Comando FALTANDO: `block-storage snapshots copy`

### MGC (Referência)
```
Copy a object snapshot cross region for the currently authenticated tenant.

#### Rules
- The copy only be accepted when the destiny region is different from origin region.
- The copy only be accepted if the snapshot name in destiny region is different from input name.
- The copy only be accepted if the user has access to destiny region.

#### Notes
- Utilize the **block-storage snapshots list** command to retrieve a list...

Usage:
  mgc block-storage snapshots copy [id] [flags]

Flags:
  --cli.list-links enum[=table]   List all available links
  --destination-region enum       Regions (one of "br-mgl1", "br-ne1" or "br-se1") (required)
  --id uuid                       Id (required)
```

### TMP-CLI (Atual)
```
❌ COMANDO NÃO EXISTE
```

### ❌ Problema:
**Comando completamente ausente no TMP-CLI**

---

## 11. Comando: `block-storage volumes create`

### MGC (Referência)
```
Create a Volume for the currently authenticated tenant.

The Volume can be used when it reaches the "available" state and "completed" status.

#### Rules
- The Volume name must be unique; otherwise, the creation will be disallowed.
- The Volume type must be available to use.

#### Notes
- Utilize the **block-storage volume-types list** command...
- Verify the state and status of your Volume using the **block-storage volume get --id [uuid]** command.

Examples:
  mgc block-storage volumes create --snapshot.id="..." --type.id="..."

Flags:
  --availability-zone string      Availability Zone
  --cli.list-links enum[=table]   List all available links
  --encrypted                     Indicates if the volume is encrypted. Default is False.
  --name string                   Name (between 3 and 50 characters) (required)
  --size integer                  Size: Gibibytes (GiB) (range: 10 - 2147483648) (required)
  --snapshot object               Snapshot (at least one of: id or name)
  --snapshot.id string            Snapshot: Id (min character count: 1)
  --snapshot.name string          Snapshot: Name (between 1 and 255 characters)
  --type object                   Type (at least one of: id or name) (required)
  --type.id string                Type: Id (min character count: 1)
  --type.name string              Type: Name (between 1 and 255 characters)
```

### TMP-CLI (Atual)
```
Blockstorage provides functionality to interact with the MagaluCloud block storage service.

doto3   <-- ⚠️ BUG

Flags:
  --availability-zone  (sem descrição)
  --encrypted        (required) (sem descrição)
  --name             (required) (sem descrição)
  --size             (required) (sem descrição)
  --type.id          (sem descrição)
  --type.name        (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** Genérica ao invés de descrição detalhada com Rules e Notes
3. **Flags FALTANDO:**
   - `--cli.list-links`
   - `--snapshot` (objeto)
   - `--snapshot.id`
   - `--snapshot.name`
   - `--type` (objeto pai)
4. **Flags sem descrição:** Todas as flags sem descrição
5. **Marcação incorreta de required:** 
   - `--encrypted` marcado como required (no MGC é opcional, default False)
6. **Falta de examples:** MGC tem exemplos, TMP-CLI não

---

## 12. Comando: `block-storage volumes list`

### MGC (Referência)
```
Retrieve a list of Volumes for the currently authenticated tenant.

#### Notes
- Use the expand argument to obtain additional details about the Volume Type.

Flags:
  --control.limit integer     Limit
  --control.offset integer    Offset
  --control.sort string       Sort (pattern: ^(^[\w-]+:(asc|desc)(,[\w-]+:(asc|desc))*)?$)
  --expand array(string)     Expand: You can get more detailed info about: ['volume_type', 'attachment']
  --name string              Name
```

### TMP-CLI (Atual)
```
Blockstorage provides functionality to interact with the MagaluCloud block storage service.

doto3   <-- ⚠️ BUG

Usage: cli block-storage volumes list [Sort] [Expand] [Limit] [Offset] [flags]

Flags:
  --expand          (sem descrição)
  --limit            (required) (sem descrição)
  --offset           (required) (sem descrição)
  --sort             (required) (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** Genérica ao invés de específica (falta Notes)
3. **Argumentos posicionais incorretos:** Usage mostra flags como argumentos posicionais
4. **Nomes de flags divergentes:**
   - MGC: `--control.limit` → TMP-CLI: `--limit`
   - MGC: `--control.offset` → TMP-CLI: `--offset`
   - MGC: `--control.sort` → TMP-CLI: `--sort`
5. **Flags FALTANDO:** `--name` (filtro por nome)
6. **Flags sem descrição:** Todas as flags sem descrição
7. **Marcação incorreta de required:** Todas as flags marcadas como required quando não deveriam ser

---

## 13. Comando: `block-storage volume-types list`

### MGC (Referência)
```
List Volume Types allowed in the current region.

#### Notes
- Volume types are managed internally. If you wish to use a Volume Type that is not yet available,
  please contact our support team for assistance.

Flags:
  --allows-encryption          Allows-Encryption
  --availability-zone string   Availability-Zone
  --name string                Name
```

### TMP-CLI (Atual)
```
Blockstorage provides functionality to interact with the MagaluCloud block storage service.

doto3   <-- ⚠️ BUG

Usage: cli block-storage volume-types list [Offset] [Limit] [Sort] [AvailabilityZone] [Name] [AllowsEncryption] [flags]

Flags:
  --allows-encryption  (required) (sem descrição)
  --availability-zone  (required) (sem descrição)
  --limit            (required) (sem descrição)
  --name             (required) (sem descrição)
  --offset           (required) (sem descrição)
  --sort             (required) (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** Genérica ao invés de específica (falta Notes)
3. **Argumentos posicionais incorretos:** Usage mostra flags como argumentos posicionais
4. **Flags EXTRAS (não existem no MGC):**
   - `--limit`
   - `--offset`
   - `--sort`
5. **Flags sem descrição:** Todas as flags sem descrição
6. **Marcação incorreta de required:** Todas as flags marcadas como required quando não deveriam ser

---

## Resumo Geral de Problemas

### 🐛 BUGS CRÍTICOS:
- String **"Dqui1"** aparece em todos os subcomandos de nível 2
- String **"doto3"** aparece em todos os comandos leaf (nível 3+)

### ❌ COMANDOS FALTANDO:
- `block-storage snapshots copy`

### ➕ COMANDOS EXTRAS:
- `list-all` em: schedulers, snapshots, volume-types, volumes

### 🔄 NOMENCLATURA DIVERGENTE:
- `schedulers attach` → TMP-CLI usa `schedulers attach-volume`
- `schedulers detach` → TMP-CLI usa `schedulers detach-volume`

### 🚫 ALIASES FALTANDO:
- `bs`, `blocks`, `blst`, `block`, `volumes` (no comando principal)

### 📝 DESCRIÇÕES:
- Todas as descrições são genéricas e repetitivas
- Faltam seções **Rules** e **Notes** dos comandos
- Faltam descrições de todas as flags
- Faltam **Examples** em comandos que os têm no MGC

### 🔧 FLAGS:
#### Padrão de problemas em TODOS os comandos list:
- Falta prefixo `control.` em: `limit`, `offset`, `sort`
- Flags marcadas incorretamente como `required`
- Flags sem descrição
- Argumentos posicionais incorretos no Usage

#### Flags FALTANDO (padrão em vários comandos):
- `--cli.list-links` (em comandos create/modify)
- Objetos pai e intermediários em estruturas aninhadas
- Propriedades `.id` e `.name` para referências a recursos

#### Flags EXTRAS:
- `volume-types list` tem `--limit`, `--offset`, `--sort` que não existem no MGC

---

## 📊 Estatísticas

### Comandos Principais (Nível 1)
- **Aliases:** 0% (0/6 aliases implementados)
- **Descrições:** 0% corretas (genéricas e repetitivas)

### Subcomandos (Nível 2)
- **Total MGC:** 4 grupos (schedulers, snapshots, volume-types, volumes)
- **Total TMP-CLI:** 4 grupos (todos presentes)
- **Bugs:** 100% (4/4 têm bug "Dqui1")
- **Descrições:** 0% corretas (todas genéricas)

### Comandos Leaf (Nível 3+)
- **Total MGC:** 25 comandos
- **Total TMP-CLI:** 28 comandos (3 extras: list-all)
- **Faltando:** 1 comando (snapshots copy)
- **Bugs:** 100% (todos com bug "doto3")
- **Descrições:** 0% corretas
- **Flags com descrição:** 0%
- **Flags com prefixo correto:** ~30% (maioria sem `control.`)

---

## ✅ Checklist de Ações

### P0 - Crítico
- [ ] Remover string "Dqui1" de todos os subcomandos nível 2
- [ ] Remover string "doto3" de todos os comandos leaf
- [ ] Adicionar descrições em TODAS as flags
- [ ] Corrigir descrição do comando principal: "Block Storage API Product"
- [ ] Adicionar comando faltante: `snapshots copy`

### P1 - Alto
- [ ] Adicionar aliases: `bs`, `blocks`, `blst`, `block`, `volumes`
- [ ] Renomear `schedulers attach-volume` → `schedulers attach`
- [ ] Renomear `schedulers detach-volume` → `schedulers detach`
- [ ] Adicionar prefixo `control.` em: limit, offset, sort
- [ ] Corrigir marcação de flags required (remover onde não se aplica)
- [ ] Remover argumentos posicionais incorretos do Usage
- [ ] Adicionar descrições específicas para cada comando
- [ ] Adicionar seções Rules e Notes nos comandos apropriados
- [ ] Adicionar Examples nos comandos apropriados

### P2 - Médio
- [ ] Adicionar flag `--cli.list-links` em comandos de criação/modificação
- [ ] Adicionar objetos pai completos (policy, snapshot, volume, type, etc.)
- [ ] Adicionar propriedades `.id` e `.name` para referências
- [ ] Adicionar flag `--name` em `volumes list`
- [ ] Remover flags extras de `volume-types list`: limit, offset, sort
- [ ] Decidir sobre comandos `list-all` (remover ou justificar)

### P3 - Baixo
- [ ] Melhorar formatação geral do help
- [ ] Padronizar estrutura de descrições

