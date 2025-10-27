# Comparação: mgc load-balancer/lbaas vs ./tmp-cli/cli lbaas

**Data:** 24/10/2025  
**Referência:** `mgc` (comandos oficiais que precisam ser replicados em `./tmp-cli/cli`)  
**Nota:** Flags globais não incluídas nesta análise (já documentadas anteriormente)

---

## NOTA IMPORTANTE: Nome do Comando Principal

**MGC usa `load-balancer`** como comando principal, com aliases `lb` e `lbaas`:
```bash
mgc load-balancer [command]
mgc lb [command]           # alias
mgc lbaas [command]        # alias
```

**TMP-CLI usa `lbaas`** como comando principal (sem `load-balancer`).

Para esta comparação, usarei:
- MGC: `mgc load-balancer` (comando principal)
- TMP-CLI: `cli lbaas`

---

## 1. Comando Principal: `load-balancer` vs `lbaas`

### MGC (Referência)
```
Lbaas API: create and manage Load Balancers

Aliases:
  load-balancer, lb, lbaas

Commands:
  network-acls          Network Load Balancer ACLs
  network-backends      Network Load Balancer Backends (Target Pools)
  network-certificates  Network Load Balancer TLS Certificates
  network-healthchecks  Network Load Balancer Health Checks
  network-listeners     Network Load Balancer Listeners
  network-loadbalancers Network Load Balancer
```

### TMP-CLI (Atual)
```
Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.

Package lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
This package allows you to manage network load balancers, listeners, backends, health checks, certificates, and ACLs.

Commands:
  network-a-c-ls      Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  network-backend-targetsLbaas provides...   <-- ⚠️ BUG SEM ESPAÇO
  network-backends    Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  network-certificatesLbaas provides...      <-- ⚠️ BUG SEM ESPAÇO
  network-health-checksLbaas provides...     <-- ⚠️ BUG SEM ESPAÇO
  network-listeners   Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  network-load-balancersLbaas provides...    <-- ⚠️ BUG SEM ESPAÇO
```

### ❌ Problemas Identificados:
1. **Descrição divergente:**
   - MGC: "Lbaas API: create and manage Load Balancers" (concisa)
   - TMP-CLI: Descrição repetitiva e verbosa com "Package lbaas..." redundante
2. **Comando principal divergente:**
   - MGC: `load-balancer` (com aliases)
   - TMP-CLI: `lbaas` apenas (sem suporte para `load-balancer`)
3. **Bugs gravíssimos de formatação - NOMES com espaços ou SEM ESPAÇO:**
   - `network-acls` → TMP-CLI: **`network-a-c-ls`** (com espaços entre letras!)
   - 4 comandos **SEM ESPAÇO** entre nome e descrição:
     - `network-backend-targetsLbaas`
     - `network-certificatesLbaas`
     - `network-health-checksLbaas`
     - `network-load-balancersLbaas`
4. **Grupo EXTRA:** `network-backend-targets` não existe no MGC (é subcomando de network-backends)
5. **Nomenclatura divergente:**
   - MGC: `network-loadbalancers` → TMP-CLI: `network-load-balancers` (com hífen extra)
   - MGC: `network-healthchecks` → TMP-CLI: `network-health-checks` (com hífen extra)
6. **Descrições dos subcomandos:** Todas genéricas ao invés de específicas

---

## 2. Bug CRÍTICO: `network-acls` vs `network-a-c-ls`

### MGC (Referência)
```
network-acls      Network Load Balancer ACLs
```

### TMP-CLI (Atual)
```
network-a-c-ls    Lbaas provides...
```

### ❌ Problema GRAVÍSSIMO:
**Nome do grupo com ESPAÇOS entre as letras!**
- `acls` → `a-c-ls` (separou cada letra com hífen)
- Provavelmente bug no parser que trata "ACLs" como acrônimo

---

## 3. Bug CRÍTICO: Comandos SEM ESPAÇO antes da Descrição

Múltiplos comandos têm bug de formatação onde não há espaço entre o nome e a descrição:

```
network-backend-targetsLbaas provides...    (sem espaço)
network-certificatesLbaas provides...       (sem espaço)
network-health-checksLbaas provides...      (sem espaço)
network-load-balancersLbaas provides...     (sem espaço)
```

**Impacto:** Output ilegível, comandos parecem ter nomes incorretos

---

## 4. Comando: `network-loadbalancers` vs `network-load-balancers`

### MGC (Referência)
```
Network Load Balancer

Commands:
  create      Create Load Balancer
  delete      Delete Load Balancer by ID
  get         Get Load Balancer by ID
  list        List Load Balancers
  replace     Update Load Balancer by ID
```

### TMP-CLI (Atual - comando `network-load-balancers`)
```
Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.

Dqui1   <-- ⚠️ BUG

Commands:
  create              Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  delete              Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  get                 Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  list                Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  list-all            Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  update              Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Nome do grupo divergente:**
   - MGC: `network-loadbalancers` (sem hífen em loadbalancers)
   - TMP-CLI: `network-load-balancers` (com hífen)
3. **Descrição divergente:** TMP-CLI usa descrição genérica
4. **Comando divergente:**
   - MGC: `replace` → TMP-CLI: `update`
5. **Comando EXTRA:** `list-all` não existe no MGC
6. **Descrições dos subcomandos:** Todas genéricas

---

## 5. Comando: `network-backends`

### MGC (Referência)
```
Network Load Balancer Backends (Target Pools)

Commands:
  create      Create Backend
  delete      Delete Backend by ID
  get         Get Backend by ID
  list        List Backends
  replace     Update Backend by ID
  targets     targets
```

### TMP-CLI (Atual)
```
Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.

Dqui1   <-- ⚠️ BUG

Commands:
  create              Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  delete              Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  get                 Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  list                Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  list-all            Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  update              Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Descrição divergente:** TMP-CLI usa descrição genérica
3. **Subcomando FALTANDO:** `targets` não existe no TMP-CLI
4. **Comando divergente:** `replace` → TMP-CLI: `update`
5. **Comando EXTRA:** `list-all`
6. **Descrições dos subcomandos:** Todas genéricas

### NOTA: Grupo EXTRA no nível principal
O TMP-CLI tem `network-backend-targets` como grupo de nível 2, quando deveria ser `network-backends/targets` (subcomando).

---

## 6. Comando: `network-listeners`

### MGC (Referência)
```
Network Load Balancer Listeners

Commands:
  create      Create Listener
  delete      Delete Listener by ID
  get         Get Listener by ID
  list        List Listeners
  replace     Update Listener by ID
```

### TMP-CLI (Atual)
```
Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.

Dqui1   <-- ⚠️ BUG

Commands:
  create              Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  delete              Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  get                 Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  list                Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  list-all            Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  update              Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Descrição divergente:** TMP-CLI usa descrição genérica
3. **Comando divergente:** `replace` → TMP-CLI: `update`
4. **Comando EXTRA:** `list-all`
5. **Descrições dos subcomandos:** Todas genéricas

---

## 7. Comando: `network-acls` vs `network-a-c-ls`

### MGC (Referência)
```
Network Load Balancer ACLs

Commands:
  create      Create Load Balancer ACL
  delete      Delete Load Balancer ACL
  replace     Replace Load Balancer ACLs
```

### TMP-CLI (Atual - comando `network-a-c-ls`)
```
Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.

Dqui1   <-- ⚠️ BUG

Commands:
  create              Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  delete              Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
  replace             Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Bug gravíssimo:** Nome do grupo `network-a-c-ls` com hífens entre cada letra
3. **Descrição divergente:** TMP-CLI usa descrição genérica
4. **Descrições dos subcomandos:** Todas genéricas

**Nota:** Este grupo NÃO tem o comando `list`, então não tem `list-all` extra

---

## 8. Comando: `network-loadbalancers list`

### MGC (Referência)
```
List Load Balancers

Flags:
  --control.limit integer    Items Per Page (min: 1)
  --control.offset integer   Page Number (min: 0)
  --control.sort string      Sort: Name of the field...
```

### TMP-CLI (Atual - comando `network-load-balancers list`)
```
Lbaas provides a client for interacting with the Magalu Cloud Load Balancer as a Service (LBaaS) API.

doto3   <-- ⚠️ BUG

Usage: cli lbaas network-load-balancers list [Offset] [Limit] [Sort] [flags]

Flags:
  --limit            (required) (sem descrição)
  --offset           (required) (sem descrição)
  --sort             (required) (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** Genérica ao invés de específica
3. **Argumentos posicionais incorretos:** Usage mostra flags como argumentos posicionais
4. **Nomes de flags divergentes:**
   - MGC: `--control.limit` → TMP-CLI: `--limit` (sem prefixo)
   - MGC: `--control.offset` → TMP-CLI: `--offset` (sem prefixo)
   - MGC: `--control.sort` → TMP-CLI: `--sort` (sem prefixo)
5. **Flags sem descrição:** Todas as flags sem descrição
6. **Marcação incorreta de required:** Todas as flags marcadas como required

---

## Resumo Geral de Problemas

### 🐛 BUGS CRÍTICOS DE FORMATAÇÃO:
- **Nome com hífens incorretos:** `network-acls` → `network-a-c-ls` (ACLs separado em letras)
- **4 comandos SEM ESPAÇO** entre nome e descrição:
  - `network-backend-targetsLbaas`
  - `network-certificatesLbaas`
  - `network-health-checksLbaas`
  - `network-load-balancersLbaas`
- String **"Dqui1"** em todos os subcomandos de nível 2
- String **"doto3"** em todos os comandos leaf

### 📝 NOMENCLATURA DIVERGENTE:
**Comando principal:**
- MGC: `load-balancer` (com aliases `lb`, `lbaas`)
- TMP-CLI: `lbaas` apenas

**Grupos de nível 2:**
- `network-loadbalancers` → TMP-CLI: `network-load-balancers` (hífen extra)
- `network-healthchecks` → TMP-CLI: `network-health-checks` (hífen extra)
- `network-acls` → TMP-CLI: `network-a-c-ls` (bug de parsing de acrônimo)

**Comandos:**
- `replace` → TMP-CLI: `update` (em múltiplos grupos)

### ❌ ESTRUTURA DIVERGENTE:
- `network-backends/targets` → TMP-CLI: promoveu para `network-backend-targets` (grupo de nível 2)
- Subcomando `targets` faltando em `network-backends`

### ➕ COMANDOS EXTRAS:
- `list-all` em múltiplos grupos
- `network-backend-targets` (deveria ser subcomando, não grupo)

### 📝 DESCRIÇÕES:
- Todas as descrições são genéricas e repetitivas
- Faltam descrições de todas as flags

### 🔧 FLAGS:
- Falta prefixo `control.` em: `limit`, `offset`, `sort`
- Flags marcadas incorretamente como `required`
- Flags sem descrição

---

## 📊 Estatísticas

### Comando Principal
- **Nome divergente:** MGC usa `load-balancer`, TMP-CLI usa `lbaas`
- **Aliases implementados:** Desconhecido (não é possível verificar se `load-balancer` e `lb` funcionam)

### Subcomandos (Nível 2)
- **Total MGC:** 6 grupos
- **Total TMP-CLI:** 7 grupos (1 extra: network-backend-targets)
- **Bugs de formatação gravíssimos:** 5 grupos afetados
  - 1 com nome incorreto (`network-a-c-ls`)
  - 4 sem espaço antes da descrição
- **Bugs "Dqui1":** 100% dos grupos

### Nomenclatura
- **Grupos com hífen extra:** 2 (loadbalancers, healthchecks)
- **Grupos com nome completamente diferente:** 1 (acls → a-c-ls)

### Comandos
- **Comando divergente:** `replace` vs `update` (em 4 grupos)
- **Comandos extras `list-all`:** 4+ comandos
- **Bugs "doto3":** 100% dos comandos leaf

### Flags
- **Sem prefixo `control.`:** 100% dos list
- **Sem descrição:** 100%
- **Required incorreto:** 100% dos list

---

## ✅ Checklist de Ações

### P0 - Crítico (Bugs de Formatação)
- [ ] **Corrigir nome `network-a-c-ls` → `network-acls`** (bug de parsing de acrônimo)
- [ ] **Corrigir 4 comandos SEM ESPAÇO** antes da descrição:
  - `network-backend-targetsLbaas`
  - `network-certificatesLbaas`
  - `network-health-checksLbaas`
  - `network-load-balancersLbaas`
- [ ] Remover string "Dqui1" de todos os subcomandos nível 2
- [ ] Remover string "doto3" de todos os comandos leaf
- [ ] Adicionar descrições em TODAS as flags

### P1 - Alto (Nomenclatura e Estrutura)
- [ ] **Decidir comando principal:** Implementar suporte para `load-balancer` e `lb` como aliases?
- [ ] **Corrigir nomes de grupos com hífen extra:**
  - `network-load-balancers` → `network-loadbalancers`
  - `network-health-checks` → `network-healthchecks`
- [ ] **Reestruturar `network-backend-targets`:**
  - Remover do nível 2
  - Adicionar como subcomando `network-backends/targets`
- [ ] **Renomear comando:** `update` → `replace` (em 4 grupos)
- [ ] Adicionar prefixo `control.` em: limit, offset, sort
- [ ] Corrigir marcação de flags required
- [ ] Corrigir descrições de comandos (específicas, não genéricas)
- [ ] Remover argumentos posicionais incorretos do Usage

### P2 - Médio (Funcionalidades)
- [ ] Remover comandos `list-all` (4+ comandos)
- [ ] Adicionar subcomando `targets` em `network-backends`

### P3 - Baixo (Polish)
- [ ] Melhorar formatação geral do help
- [ ] Padronizar estrutura de descrições

