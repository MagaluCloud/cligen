# Comparação Detalhada: mgc virtual-machine vs ./tmp-cli/cli virtual-machine

## Sumário Executivo

A comparação entre `mgc virtual-machine` e `./tmp-cli/cli virtual-machine` revela **problemas arquiteturais graves** e perda massiva de funcionalidade:

### 🔴 Problemas Críticos Identificados

1. **BUGS VISUAIS**: Strings "Dqui1" e "doto3" aparecem ubiquamente
2. **NOME DE GRUPO DIVERGENTE**: `machine-types` → `instance-types`
3. **SUBGRUPO AUSENTE**: `images custom` completamente ausente (5 comandos perdidos)
4. **ARQUITETURA DIVERGENTE**: `instances network-interface` (subgrupo) → comandos flat (`attach-network-interface`, `detach-network-interface`)
5. **COMANDOS DIVERGENTES**: `password` → `get-first-windows-password`, `init-logs` → `init-log`
6. **COMANDO AUSENTE**: `reboot` não existe em `./tmp-cli/cli`
7. **COMANDOS EXTRAS**: `list-all` em 3 grupos (não existe em `mgc`)
8. **PERDA MASSIVA DE FUNCIONALIDADE**: Flags complexas de objetos aninhados ausentes em `instances create`
9. **0% DE DESCRIÇÕES NAS FLAGS**: Nenhuma flag em `./tmp-cli/cli` possui descrição
10. **ALIASES AUSENTES**: 6 aliases do produto principal ausentes
11. **FLAGS INCORRETAMENTE MARCADAS COMO REQUIRED**: Em `list`, flags opcionais marcadas como `(required)`
12. **BUGS DE FORMATAÇÃO**: Falta de espaços entre nomes de comandos e descrições

---

## 1. Comando Principal

### mgc virtual-machine

```
Virtual Machine Api Product

Usage:
  mgc virtual-machine [flags]
  mgc virtual-machine [command]

Aliases:
  virtual-machine, vms, vm, virtual-machines, machines, vmachine

Commands:
  images        Operations with images for instances.
  instances     Operations with instances, including create, delete, start, stop, reboot and other actions.
  machine-types Operations with machine types for instances.
  snapshots     Operations with snapshots for instances.
```

### ./tmp-cli/cli virtual-machine

```
Compute provides functionality to interact with the MagaluCloud compute service.

Package compute provides functionality to interact with the MagaluCloud compute service.
This package allows managing virtual machine instances, images, instance types, and snapshots.

Available Commands:

Other commands:
  images              Compute provides functionality to interact with the MagaluCloud compute service.
  instance-types      Compute provides functionality to interact with the MagaluCloud compute service.
  instances           Compute provides functionality to interact with the MagaluCloud compute service.
  snapshots           Compute provides functionality to interact with the MagaluCloud compute service.
```

### ⚠️ Diferenças Identificadas

| Aspecto | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **Nome do grupo `machine-types`** | `machine-types` | `instance-types` | 🔴 Divergente |
| **Aliases** | 6 aliases | 0 aliases | 🔴 Todos ausentes |
| **Descrições dos comandos** | Específicas e úteis | Genéricas e repetitivas | ❌ Divergente |

---

## 2. Grupo: images

### mgc virtual-machine images

```
Operations with images for instances.

Usage:
  mgc virtual-machine images [flags]
  mgc virtual-machine images [command]

Commands:
  custom      custom
  list        Retrieves all images.
```

### ./tmp-cli/cli virtual-machine images

```
Compute provides functionality to interact with the MagaluCloud compute service.

Dqui1

Available Commands:

Other commands:
  list                Compute provides functionality to interact with the MagaluCloud compute service.
  list-all            Compute provides functionality to interact with the MagaluCloud compute service.
```

### ⚠️ Diferenças CRÍTICAS

| Aspecto | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **Bug visual** | Não | "Dqui1" aparece | 🔴 CRÍTICO |
| **Subgrupo `custom`** | ✅ Existe (5 comandos) | 🔴 **COMPLETAMENTE AUSENTE** | 🔴 CRÍTICO |
| **Comando `list-all`** | ❌ Não existe | ⚠️ Existe | ⚠️ Extra |

### 2.1. Subgrupo: images custom (AUSENTE em ./tmp-cli/cli)

**mgc virtual-machine images custom**:

```
Operations with images for instances. | custom

Usage:
  mgc virtual-machine images custom [flags]
  mgc virtual-machine images custom [command]

Commands:
  create      Create Custom Image
  delete      Delete Image Custom
  get         Get Custom Image By Id
  list        Retrieves all active custom images.
  update      Patch Image Custom
```

**./tmp-cli/cli**: ❌ **SUBGRUPO COMPLETAMENTE AUSENTE**

**⚠️ IMPACTO CRÍTICO**: O subgrupo `images custom` e seus **5 comandos** (`create`, `delete`, `get`, `list`, `update`) estão completamente ausentes em `./tmp-cli/cli`. Isso representa uma **perda de funcionalidade massiva** para gerenciamento de imagens customizadas.

---

## 3. Grupo: instances

### mgc virtual-machine instances

```
Operations with instances, including create, delete, start, stop, reboot and other actions.

Usage:
  mgc virtual-machine instances [flags]
  mgc virtual-machine instances [command]

Commands:
  create            Create an instance.
  delete            Delete an instance.
  get               Retrieve the details of a specific instance.
  init-logs         Get Instance Console Logs
  list              List all instances.
  network-interface network-interface
  password          Retrieve the first windows admin password
  reboot            Reboot an instance.
  rename            Renames an instance.
  retype            Changes an instance machine-type.
  start             Starts an instance.
  stop              Stops an instance.
  suspend           Suspends instance.
```

### ./tmp-cli/cli virtual-machine instances

```
Compute provides functionality to interact with the MagaluCloud compute service.

Dqui1

Available Commands:

Other commands:
  attach-network-interfaceCompute provides functionality to interact with the MagaluCloud compute service.
  create              Compute provides functionality to interact with the MagaluCloud compute service.
  delete              Compute provides functionality to interact with the MagaluCloud compute service.
  detach-network-interfaceCompute provides functionality to interact with the MagaluCloud compute service.
  get                 Compute provides functionality to interact with the MagaluCloud compute service.
  get-first-windows-passwordCompute provides functionality to interact with the MagaluCloud compute service.
  init-log            Compute provides functionality to interact with the MagaluCloud compute service.
  list                Compute provides functionality to interact with the MagaluCloud compute service.
  list-all            Compute provides functionality to interact with the MagaluCloud compute service.
  rename              Compute provides functionality to interact with the MagaluCloud compute service.
  retype              Compute provides functionality to interact with the MagaluCloud compute service.
  start               Compute provides functionality to interact with the MagaluCloud compute service.
  stop                Compute provides functionality to interact with the MagaluCloud compute service.
  suspend             Compute provides functionality to interact with the MagaluCloud compute service.
```

### ⚠️ Diferenças ARQUITETURAIS CRÍTICAS

| Aspecto | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **Bug visual** | Não | "Dqui1" aparece | 🔴 CRÍTICO |
| **Bugs de formatação** | Não | Falta espaço em 3 comandos | 🔴 CRÍTICO |
| **Subgrupo `network-interface`** | ✅ Subgrupo (2 comandos) | ❌ Convertido em comandos flat | 🔴 CRÍTICO |
| **Comando `password`** | `password` | `get-first-windows-password` | ❌ Divergente |
| **Comando `init-logs`** | `init-logs` | `init-log` | ❌ Divergente |
| **Comando `reboot`** | ✅ Existe | 🔴 **AUSENTE** | 🔴 CRÍTICO |
| **Comando `list-all`** | ❌ Não existe | ⚠️ Existe | ⚠️ Extra |

**Bugs de formatação**:
- `attach-network-interfaceCompute` (falta espaço antes de "Compute")
- `detach-network-interfaceCompute` (falta espaço antes de "Compute")
- `get-first-windows-passwordCompute` (falta espaço antes de "Compute")

### 3.1. Arquitetura Divergente: network-interface

**mgc**: Subgrupo `network-interface` com 2 comandos:
- `mgc virtual-machine instances network-interface attach`
- `mgc virtual-machine instances network-interface detach`

**./tmp-cli/cli**: Comandos flat:
- `cli virtual-machine instances attach-network-interface`
- `cli virtual-machine instances detach-network-interface`

**Impacto**: Incompatibilidade total na estrutura de comandos.

### 3.2. instances create - PERDA MASSIVA DE FUNCIONALIDADE

**mgc**: Flags complexas com objetos aninhados:
- `--image` (object): `--image.id`, `--image.name`
- `--machine-type` (object): `--machine-type.id`, `--machine-type.name`
- `--network` (object):
  - `--network.associate-public-ip`
  - `--network.interface` (object):
    - `--network.interface.id`
    - `--network.interface.security-groups` (array)
  - `--network.vpc` (object):
    - `--network.vpc.id`
    - `--network.vpc.name`
- `--volumes` (array of objects)
- `--ssh-key-name`
- `--user-data`
- `--availability-zone`
- `--labels` (array)

**./tmp-cli/cli**: Flags simples e incompletas:
- `--image.id`
- `--image.name`
- `--machine-type.id`
- `--machine-type.name`
- `--ssh-key-name`
- `--user-data`
- `--availability-zone`
- `--labels`
- Bug visual: "doto3"

**Flags completamente ausentes**:
- ❌ `--network` e todas as suas propriedades aninhadas
- ❌ `--network.associate-public-ip`
- ❌ `--network.interface.id`
- ❌ `--network.interface.security-groups`
- ❌ `--network.vpc.id`
- ❌ `--network.vpc.name`
- ❌ `--volumes`

**⚠️ IMPACTO CRÍTICO**: Sem essas flags, **NÃO é possível**:
- Associar um IP público à instância durante a criação
- Especificar uma interface de rede
- Definir security groups
- Escolher uma VPC
- Anexar volumes
- **Criar uma instância funcional na maioria dos casos reais**

### 3.3. instances list

**mgc**: `mgc virtual-machine instances list`
- Flags:
  - `--control.limit` (integer): Limit (max: 1000)
  - `--control.offset` (integer): Offset (max: 2147483647)
  - `--control.sort` (string): Sort (pattern)
  - `--expand` (array): Expand details ['image', 'machine-type', 'machine-types', 'network', 'labels']
  - `--name` (string): name of the instance

**./tmp-cli/cli**: `cli virtual-machine instances list [Sort] [Expand] [Name] [Limit] [Offset]`
- Flags:
  - `--expand`
  - `--limit` (required)
  - `--name` (required)
  - `--offset` (required)
  - `--sort` (required)
- Bug visual: "doto3"

**Divergências**:
- ❌ `--control.limit`, `--control.offset`, `--control.sort` vs `--limit`, `--offset`, `--sort`
- ❌ Todas as flags (exceto `expand`) marcadas incorretamente como `(required)`
- 🔴 Bug visual: "doto3"

---

## 4. Grupo: machine-types / instance-types

### mgc virtual-machine machine-types

```
Operations with machine types for instances.

Usage:
  mgc virtual-machine machine-types [flags]
  mgc virtual-machine machine-types [command]

Commands:
  list        Retrieves all machine-types.
```

### ./tmp-cli/cli virtual-machine instance-types

```
Compute provides functionality to interact with the MagaluCloud compute service.

Dqui1

Available Commands:

Other commands:
  list                Compute provides functionality to interact with the MagaluCloud compute service.
  list-all            Compute provides functionality to interact with the MagaluCloud compute service.
```

### ⚠️ Diferenças

| Aspecto | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **Nome do grupo** | `machine-types` | `instance-types` | 🔴 Divergente |
| **Bug visual** | Não | "Dqui1" aparece | 🔴 CRÍTICO |
| **Comando `list-all`** | ❌ Não existe | ⚠️ Existe | ⚠️ Extra |

---

## 5. Grupo: snapshots

### mgc virtual-machine snapshots

```
Operations with snapshots for instances.

Usage:
  mgc virtual-machine snapshots [flags]
  mgc virtual-machine snapshots [command]

Commands:
  copy        Copy a snapshot of a virtual machine asynchronously to another region.
  create      Create a snapshot of an instance.
  delete      Delete a Snapshot.
  get         Retrieve the details of an snapshot.
  list        Lists all snapshots.
  rename      Renames a snapshot.
  restore     Restore a snapshot to an instance.
```

### ./tmp-cli/cli virtual-machine snapshots

```
Compute provides functionality to interact with the MagaluCloud compute service.

Dqui1

Available Commands:

Other commands:
  copy                Compute provides functionality to interact with the MagaluCloud compute service.
  create              Compute provides functionality to interact with the MagaluCloud compute service.
  delete              Compute provides functionality to interact with the MagaluCloud compute service.
  get                 Compute provides functionality to interact with the MagaluCloud compute service.
  list                Compute provides functionality to interact with the MagaluCloud compute service.
  list-all            Compute provides functionality to interact with the MagaluCloud compute service.
  rename              Compute provides functionality to interact with the MagaluCloud compute service.
  restore             Compute provides functionality to interact with the MagaluCloud compute service.
```

### ⚠️ Diferenças

| Aspecto | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **Bug visual** | Não | "Dqui1" aparece | 🔴 CRÍTICO |
| **Comandos** | 7 comandos | 8 comandos | - |
| **Comando `list-all`** | ❌ Não existe | ⚠️ Existe | ⚠️ Extra |

---

## 6. Problemas Sistemáticos Identificados

### 🔴 Bugs Visuais Críticos
1. **"Dqui1"** aparece em todos os 4 grupos
2. **"doto3"** aparece em todos os comandos específicos
3. **Falta de espaços** entre nome do comando e descrição em 3 comandos de `instances`

### 🔴 Problemas de Nomenclatura
1. **Nome de grupo**: `machine-types` → `instance-types`
2. **Nome de comando**: `password` → `get-first-windows-password`
3. **Nome de comando**: `init-logs` → `init-log` (plural → singular)
4. **Prefixo ausente**: `control.limit` → `limit`, `control.offset` → `offset`, `control.sort` → `sort`

### 🔴 Problemas Arquiteturais CRÍTICOS
1. **Subgrupo ausente**: `images custom` completamente ausente (5 comandos perdidos)
2. **Arquitetura divergente**: `instances network-interface` (subgrupo) → comandos flat
3. **Comando ausente**: `reboot` não existe em `./tmp-cli/cli`

### 🔴 Problemas de Funcionalidade CRÍTICOS
1. **Perda massiva de funcionalidade em `instances create`**:
   - Flags de rede (`--network.*`) completamente ausentes
   - Flags de volumes (`--volumes`) completamente ausentes
   - **Impossível criar instâncias com configuração de rede customizada**
   - **Impossível anexar volumes durante a criação**

### 🔴 Problemas de Conteúdo
1. **0% de descrições nas flags** em `./tmp-cli/cli`
2. **Descrições genéricas** para todos os comandos
3. **Aliases ausentes** (6 aliases do produto principal)

### 🔴 Problemas de Flags
1. **Comandos extras**: `list-all` em 3 grupos (não existe em `mgc`)
2. **Flags incorretamente marcadas como required**: `limit`, `offset`, `sort`, `name`, `expand` em `instances list`

---

## 7. Resumo de Incompatibilidades

| Categoria | Qtd. Problemas | Severidade |
|-----------|----------------|------------|
| **Bugs Visuais** | 3 tipos | 🔴 CRÍTICA |
| **Nomes Divergentes** | 5+ instâncias | 🔴 CRÍTICA |
| **Descrições Ausentes** | 100% das flags | 🔴 CRÍTICA |
| **Arquitetura Divergente** | 2 casos | 🔴 CRÍTICA |
| **Subgrupo Ausente** | 1 (5 comandos) | 🔴 CRÍTICA |
| **Comando Ausente** | 1 (`reboot`) | 🔴 CRÍTICA |
| **Perda de Funcionalidade** | `instances create` | 🔴 CRÍTICA |
| **Comandos Extras** | 3 (`list-all`) | ⚠️ MÉDIA |
| **Aliases Ausentes** | 6 aliases | ⚠️ BAIXA |

---

## 8. Recomendações

### Prioridade CRÍTICA 🔴
1. **Eliminar bugs visuais**: Remover "Dqui1", "doto3" e bugs de formatação
2. **Restaurar nome correto do grupo**: `instance-types` → `machine-types`
3. **Restaurar subgrupo `images custom`**: Implementar os 5 comandos ausentes
4. **Restaurar comando `reboot`**: Adicionar comando ausente
5. **Restaurar flags complexas em `instances create`**: Implementar `--network.*` e `--volumes`
6. **Restaurar arquitetura de subgrupos**: Implementar `instances network-interface` com subcomandos
7. **Corrigir nomes de comandos**: `get-first-windows-password` → `password`, `init-log` → `init-logs`
8. **Adicionar descrições em 100% das flags**
9. **Corrigir flags incorretamente marcadas como required**
10. **Restaurar prefixo `control.`** em flags de paginação

### Prioridade ALTA ⚠️
1. **Adicionar aliases faltando** (6 aliases)
2. **Verificar comandos extras**: Confirmar se `list-all` deve existir
3. **Corrigir bugs de formatação**: Adicionar espaços entre nomes e descrições

### Prioridade MÉDIA
1. **Melhorar descrições dos comandos**: Substituir textos genéricos por descrições específicas

---

## Conclusão

O CLI gerado (`./tmp-cli/cli virtual-machine`) apresenta os **PROBLEMAS MAIS GRAVES** encontrados em todos os produtos analisados:

### Problemas Únicos e Críticos:
1. **Perda MASSIVA de funcionalidade em `instances create`**: Sem flags de rede e volumes, o comando é **praticamente inutilizável** para casos reais
2. **Subgrupo inteiro ausente**: `images custom` com 5 comandos completamente perdidos
3. **Comando crítico ausente**: `reboot` não existe
4. **Nome de grupo divergente**: `machine-types` → `instance-types`
5. **Bugs de formatação únicos**: Falta de espaços em 3 comandos

### Problemas Compartilhados:
1. **Bugs visuais**: "Dqui1" e "doto3"
2. **0% de descrições** nas flags
3. **Arquitetura divergente**: Subgrupo convertido em comandos flat
4. **Comandos extras**: `list-all`
5. **Prefixo `control.` removido**
6. **Flags incorretamente marcadas como required**
7. **Aliases ausentes**

### Impacto Geral:
- **Funcionalidade**: **CRÍTICO** - Criação de instâncias praticamente inutilizável
- **Compatibilidade**: 0% dos comandos são 100% compatíveis
- **Arquitetura**: 25% dos grupos têm divergências estruturais críticas
- **Usabilidade**: Severamente comprometida por bugs visuais e falta de documentação
- **Completude**: 1 subgrupo + 1 comando completamente ausentes (perda de ~20% da funcionalidade)

**Recomendação**: ❌ **COMPLETAMENTE INAPTO PARA PRODUÇÃO**. O produto `virtual-machine` tem os problemas **MAIS GRAVES** de todos os produtos analisados. A perda de flags críticas em `instances create` torna o CLI **inutilizável** para a maioria dos casos reais. Requer **refatoração completa** antes de qualquer uso.

