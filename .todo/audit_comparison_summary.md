# Resumo Executivo - Comparação MGC vs TMP-CLI (audit)

## 📊 Tabela de Comandos

| Comando | MGC (Referência) | TMP-CLI | Status |
|---------|------------------|---------|--------|
| `audit` | ✅ Existe | ✅ Existe | ⚠️ Descrição divergente |
| `audit event-types` | ✅ Existe | ✅ Existe | ⚠️ Descrição divergente + Bug "Dqui1" |
| `audit event-types list` | ✅ Existe | ✅ Existe | ⚠️ Flags divergentes + Bug "doto3" |
| `audit event-types list-all` | ❌ NÃO existe | ✅ Existe | 🔴 **COMANDO EXTRA** |
| `audit events` | ✅ Existe | ✅ Existe | ⚠️ Descrição divergente + Bug "Dqui1" |
| `audit events list` | ✅ Existe | ✅ Existe | ⚠️ Flags divergentes + Bug "doto3" |
| `audit events list-all` | ❌ NÃO existe | ✅ Existe | 🔴 **COMANDO EXTRA** |

---

## 📋 Tabela de Flags Globais

| Flag | MGC | TMP-CLI | Prioridade |
|------|-----|---------|------------|
| `--api-key` | ✅ | ✅ | OK |
| `--debug` | ✅ | ✅ | OK |
| `--no-confirm` | ✅ | ✅ | OK |
| `--raw` / `-r` | ✅ | ✅ | OK |
| `--cli.retry-until` | ✅ | ❌ | 🔴 **FALTANDO** |
| `--cli.timeout` / `-t` | ✅ | ❌ | 🔴 **FALTANDO** |
| `--output` / `-o` | ✅ | ❌ | 🔴 **FALTANDO** |
| `--env` | ✅ | ❌ | 🔴 **FALTANDO** |
| `--region` | ✅ | ❌ | 🔴 **FALTANDO** |
| `--server-url` | ✅ | ❌ | 🔴 **FALTANDO** |
| `--lang` | ❌ | ✅ | 🟡 **EXTRA** |

---

## 📋 Tabela de Flags: `event-types list`

| Flag | MGC | TMP-CLI | Observação |
|------|-----|---------|------------|
| `--control.limit` | ✅ | ❌ | TMP-CLI usa `--limit` (sem prefixo) |
| `--control.offset` | ✅ | ❌ | TMP-CLI usa `--offset` (sem prefixo) |
| `--limit` | ❌ | ✅ | Deveria ser `--control.limit` |
| `--offset` | ❌ | ✅ | Deveria ser `--control.offset` |
| `--tenant-id` | ❌ | ✅ | 🟡 **EXTRA** |

---

## 📋 Tabela de Flags: `events list`

| Flag | MGC | TMP-CLI | Observação |
|------|-----|---------|------------|
| `--authid` | ✅ | ❌ | TMP-CLI usa `--auth-id` (com hífen) |
| `--control.limit` | ✅ | ❌ | TMP-CLI usa `--limit` |
| `--control.offset` | ✅ | ❌ | TMP-CLI usa `--offset` |
| `--correlationid` | ✅ | ❌ | 🔴 **FALTANDO** |
| `--id` | ✅ | ✅ | OK |
| `--product-like` | ✅ | ✅ | OK |
| `--source-like` | ✅ | ✅ | OK |
| `--time` | ✅ | ✅ | OK |
| `--type-like` | ✅ | ✅ | OK |
| `--auth-id` | ❌ | ✅ | Deveria ser `--authid` |
| `--limit` | ❌ | ✅ | Deveria ser `--control.limit` |
| `--offset` | ❌ | ✅ | Deveria ser `--control.offset` |
| `--data` | ❌ | ✅ | 🟡 **EXTRA** |
| `--tenant-id` | ❌ | ✅ | 🟡 **EXTRA** |

---

## 🐛 Bugs Encontrados

1. **String "Dqui1"** aparece em:
   - `audit event-types --help`
   - `audit events --help`

2. **String "doto3"** aparece em:
   - `audit event-types list --help`
   - `audit events list --help`

3. **Descrições ausentes:** TMP-CLI não mostra descrições das flags

---

## 📝 Estatísticas

### Comandos
- **Total no MGC:** 5 comandos
- **Total no TMP-CLI:** 7 comandos
- **Comandos extras:** 2 (`list-all` em event-types e events)
- **Compatibilidade:** 71% (5/7)

### Global Flags
- **Total no MGC:** 10 flags
- **Total no TMP-CLI:** 5 flags
- **Flags faltando:** 6 flags importantes
- **Flags extras:** 1 flag (`--lang`)
- **Compatibilidade:** 40% (4/10)

### Qualidade do Output
- **Descrições consistentes:** ❌ Não (MGC: concisas e específicas | TMP-CLI: repetitivas e genéricas)
- **Formato limpo:** ❌ Não (bugs de strings estranhas)
- **Documentação completa:** ❌ Não (falta descrição nas flags)

---

## 🎯 Prioridades de Correção

### P0 - CRÍTICO (Quebra experiência do usuário)
1. ❌ Remover strings "Dqui1" e "doto3"
2. ❌ Adicionar descrições em todas as flags
3. ❌ Corrigir descrição do comando principal: "Cloud Events API Product."

### P1 - ALTO (Incompatibilidade com referência)
1. ❌ Adicionar flags globais faltantes:
   - `--cli.retry-until`
   - `--cli.timeout`
   - `--output`
   - `--env`
   - `--region`
   - `--server-url`
2. ❌ Corrigir nomenclatura de flags:
   - `--auth-id` → `--authid`
   - `--limit` → `--control.limit`
   - `--offset` → `--control.offset`
3. ❌ Adicionar flag faltante: `--correlationid` em `events list`

### P2 - MÉDIO (Decisão necessária)
1. ⚠️ Decidir sobre comandos `list-all`:
   - Remover ou justificar existência
2. ⚠️ Decidir sobre flags extras:
   - `--lang` (global)
   - `--tenant-id` (vários comandos)
   - `--data` (events list)

### P3 - BAIXO (Melhorias estéticas)
1. 📝 Melhorar formatação geral do help
2. 📝 Padronizar estrutura de descrições
3. 📝 Revisar documentação inline

---

## ✅ Checklist de Ações

- [ ] Corrigir bug das strings "Dqui1" e "doto3"
- [ ] Adicionar descrições em todas as flags
- [ ] Atualizar descrição do comando `audit`
- [ ] Atualizar descrições dos subcomandos `event-types` e `events`
- [ ] Renomear `--auth-id` para `--authid`
- [ ] Adicionar prefixo `--control.` em `limit` e `offset`
- [ ] Adicionar flag `--correlationid` em `events list`
- [ ] Implementar flags globais faltantes (6 flags)
- [ ] Decidir sobre comandos `list-all`
- [ ] Decidir sobre flags extras (`--lang`, `--tenant-id`, `--data`)
- [ ] Testar compatibilidade após correções

