# Comparação: mgc audit vs ./tmp-cli/cli audit

**Data:** 24/10/2025  
**Referência:** `mgc` (comandos oficiais que precisam ser replicados em `./tmp-cli/cli`)

---

## 1. Comando Principal: `audit`

### MGC (Referência)
```
Cloud Events API Product.

Usage:
  mgc audit [flags]
  mgc audit [command]

Commands:
  event-types Endpoints related to listing types of events emitted by other products.
  events      Endpoints related to listing events emitted by other products.
```

### TMP-CLI (Atual)
```
Audit provides functionality to interact with the MagaluCloud audit service.

Package audit provides functionality to interact with the MagaluCloud audit service.
This package allows listing audit events and event types.

Available Commands:

Other commands:
  event-types         Audit provides functionality to interact with the MagaluCloud audit service.
  events              Audit provides functionality to interact with the MagaluCloud audit service.
```

### ❌ Problemas Identificados:
1. **Descrição divergente:** 
   - MGC: "Cloud Events API Product." (concisa e clara)
   - TMP-CLI: Descrição muito repetitiva e verbosa
2. **Descrições dos subcomandos:**
   - MGC: Descrições específicas e informativas para cada comando
   - TMP-CLI: Mesma descrição genérica repetida para ambos os comandos

---

## 2. Global Flags

### ⚠️ Flags FALTANDO no TMP-CLI:
- `--cli.retry-until` (retry com condições)
- `--cli.timeout` (timeout para execução)
- `--output` / `-o` (formatação de output)
- `--env` (ambiente: pre-prod ou prod)
- `--region` (região: br-mgl1, br-ne1, br-se1, global)
- `--server-url` (URL manual do servidor)

### ✅ Flags presentes em ambos:
- `--api-key`
- `--debug`
- `--no-confirm`
- `--raw` / `-r`

### ➕ Flag EXTRA no TMP-CLI:
- `--lang` (definir idioma do CLI - não existe no mgc)

---

## 3. Comando: `audit event-types`

### MGC (Referência)
```
Endpoints related to listing types of events emitted by other products.

Commands:
  list        Lists all event types.
```

### TMP-CLI (Atual)
```
Audit provides functionality to interact with the MagaluCloud audit service.

Dqui1   <-- ⚠️ STRING ESTRANHA/BUG

Commands:
  list                Audit provides functionality to interact with the MagaluCloud audit service.
  list-all            Audit provides functionality to interact with the MagaluCloud audit service.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output (provavelmente lixo de código)
2. **Descrição divergente:** TMP-CLI usa descrição genérica ao invés da específica
3. **Comando EXTRA:** `list-all` não existe no mgc

---

## 4. Comando: `audit events`

### MGC (Referência)
```
Endpoints related to listing events emitted by other products.

Commands:
  list        Lists all events.
```

### TMP-CLI (Atual)
```
Audit provides functionality to interact with the MagaluCloud audit service.

Dqui1   <-- ⚠️ STRING ESTRANHA/BUG

Commands:
  list                Audit provides functionality to interact with the MagaluCloud audit service.
  list-all            Audit provides functionality to interact with the MagaluCloud audit service.
```

### ❌ Problemas Identificados:
1. **Bug:** String "Dqui1" aparece no output
2. **Descrição divergente:** TMP-CLI usa descrição genérica
3. **Comando EXTRA:** `list-all` não existe no mgc

---

## 5. Comando: `audit event-types list`

### MGC (Referência)
```
Lists all types of events emitted by other products.

Flags:
  --control.limit integer    Limit: Number of items per page
  --control.offset integer   Offset for pagination
```

### TMP-CLI (Atual)
```
Audit provides functionality to interact with the MagaluCloud audit service.

doto3   <-- ⚠️ STRING ESTRANHA/BUG

Flags:
  --limit           (sem descrição)
  --offset          (sem descrição)
  --tenant-id       (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** TMP-CLI usa descrição genérica ao invés de específica
3. **Nomes de flags divergentes:**
   - MGC: `--control.limit` e `--control.offset`
   - TMP-CLI: `--limit` e `--offset` (sem prefixo "control.")
4. **Flags sem descrição:** TMP-CLI não mostra descrições das flags
5. **Flag EXTRA:** `--tenant-id` não existe no mgc

---

## 6. Comando: `audit events list`

### MGC (Referência)
```
Lists all events emitted by other products.

Flags:
  --authid string            Auth ID: Identification of the actor of the action
  --control.limit integer    Limit: Number of items per page
  --control.offset integer   Offset for pagination
  --correlationid string     Correlation ID: Correlation between event chain
  --id string                Identification of the event
  --product-like string      In which producer product an event occurred ('like' operation)
  --source-like string       Source: Context in which the event occurred ('like' operation)
  --time date-time           Timestamp of when the occurrence happened
  --type-like string         Type of event related to the originating occurrence ('like' operation)
```

### TMP-CLI (Atual)
```
Audit provides functionality to interact with the MagaluCloud audit service.

doto3   <-- ⚠️ STRING ESTRANHA/BUG

Flags:
  --auth-id         (sem descrição)
  --data            (sem descrição) <-- ⚠️ FLAG EXTRA
  --id              (sem descrição)
  --limit           (sem descrição)
  --offset          (sem descrição)
  --product-like    (sem descrição)
  --source-like     (sem descrição)
  --tenant-id       (sem descrição) <-- ⚠️ FLAG EXTRA
  --time            (sem descrição)
  --type-like       (sem descrição)
```

### ❌ Problemas Identificados:
1. **Bug:** String "doto3" aparece no output
2. **Descrição do comando:** Genérica ao invés de específica
3. **Nomes de flags divergentes:**
   - MGC: `--authid` → TMP-CLI: `--auth-id` (com hífen)
   - MGC: `--correlationid` → TMP-CLI: **FALTANDO**
   - MGC: `--control.limit` → TMP-CLI: `--limit` (sem prefixo)
   - MGC: `--control.offset` → TMP-CLI: `--offset` (sem prefixo)
4. **Flags sem descrição:** TMP-CLI não mostra descrições
5. **Flags EXTRAS no TMP-CLI:**
   - `--data` (não existe no mgc)
   - `--tenant-id` (não existe no mgc)
6. **Flags FALTANDO no TMP-CLI:**
   - `--correlationid` (existe no mgc)

---

## 7. Comando EXTRA: `audit event-types list-all`

⚠️ **Este comando NÃO existe no mgc**

```
Flags:
  --tenant-id       (sem descrição)
```

**Decisão necessária:** Este comando deve ser removido ou existe uma razão específica para sua existência?

---

## 8. Comando EXTRA: `audit events list-all`

⚠️ **Este comando NÃO existe no mgc**

```
Flags:
  --auth-id         
  --data            
  --id              
  --product-like    
  --source-like     
  --tenant-id       
  --time            
  --type-like       
```

**Decisão necessária:** Este comando deve ser removido ou existe uma razão específica para sua existência?

---

## Resumo de Ações Necessárias

### 🔴 BUGS CRÍTICOS:
1. Remover strings estranhas: "Dqui1" e "doto3" do output
2. Corrigir descrições genéricas e repetitivas

### 🟡 DIVERGÊNCIAS DE NOMENCLATURA:
1. `--authid` → precisa ser `--authid` (não `--auth-id`)
2. Adicionar prefixo `--control.` para `limit` e `offset`
3. Adicionar `--correlationid` em `events list`

### 🟠 DESCRIÇÕES FALTANDO:
1. Todas as flags precisam de descrições detalhadas
2. Comandos precisam de descrições específicas (não genéricas)
3. Comando principal precisa da descrição correta: "Cloud Events API Product."

### 🔵 GLOBAL FLAGS FALTANDO:
1. `--cli.retry-until`
2. `--cli.timeout`
3. `--output` / `-o`
4. `--env`
5. `--region`
6. `--server-url`

### ⚪ COMANDOS/FLAGS EXTRAS (decisão necessária):
1. Comandos `list-all` em event-types e events
2. Flag `--lang` (global)
3. Flag `--tenant-id` (em vários comandos)
4. Flag `--data` (em events list)

### ✅ REMOVER/CORRIGIR:
1. Descrições duplicadas no comando principal
2. Seção "Package audit provides..." desnecessária
3. Formato de apresentação dos comandos (MGC é mais limpo)

