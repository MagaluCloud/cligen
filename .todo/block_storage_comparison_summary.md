# Resumo Executivo - Comparação MGC vs TMP-CLI (block-storage)

## 📊 Tabela Geral de Comandos

| Comando | MGC | TMP-CLI | Status |
|---------|-----|---------|--------|
| **Nível 1: block-storage** |
| Comando principal | ✅ | ✅ | ⚠️ Descrição divergente + sem aliases |
| **Nível 2: Subcomandos** |
| schedulers | ✅ | ✅ | ⚠️ Bug "Dqui1" + descrição divergente |
| snapshots | ✅ | ✅ | ⚠️ Bug "Dqui1" + descrição divergente |
| volume-types | ✅ | ✅ | ⚠️ Bug "Dqui1" + descrição divergente |
| volumes | ✅ | ✅ | ⚠️ Bug "Dqui1" + descrição divergente |

---

## 📊 Comandos de `schedulers`

| Comando | MGC | TMP-CLI | Status |
|---------|-----|---------|--------|
| `attach` | ✅ | ❌ | 🔴 TMP-CLI usa `attach-volume` |
| `attach-volume` | ❌ | ✅ | 🔴 Nome divergente (deveria ser `attach`) |
| `create` | ✅ | ✅ | ⚠️ Bug "doto3" + flags faltando |
| `delete` | ✅ | ✅ | ⚠️ Bug "doto3" |
| `detach` | ✅ | ❌ | 🔴 TMP-CLI usa `detach-volume` |
| `detach-volume` | ❌ | ✅ | 🔴 Nome divergente (deveria ser `detach`) |
| `get` | ✅ | ✅ | ⚠️ Bug "doto3" |
| `list` | ✅ | ✅ | ⚠️ Bug "doto3" + flags sem prefixo |
| `list-all` | ❌ | ✅ | 🔴 **COMANDO EXTRA** |

---

## 📊 Comandos de `snapshots`

| Comando | MGC | TMP-CLI | Status |
|---------|-----|---------|--------|
| `copy` | ✅ | ❌ | 🔴 **COMANDO FALTANDO** |
| `create` | ✅ | ✅ | ⚠️ Bug "doto3" + flags faltando |
| `delete` | ✅ | ✅ | ⚠️ Bug "doto3" |
| `get` | ✅ | ✅ | ⚠️ Bug "doto3" |
| `list` | ✅ | ✅ | ⚠️ Bug "doto3" + flags sem prefixo |
| `list-all` | ❌ | ✅ | 🔴 **COMANDO EXTRA** |
| `rename` | ✅ | ✅ | ⚠️ Bug "doto3" |

---

## 📊 Comandos de `volume-types`

| Comando | MGC | TMP-CLI | Status |
|---------|-----|---------|--------|
| `list` | ✅ | ✅ | ⚠️ Bug "doto3" + flags extras |
| `list-all` | ❌ | ✅ | 🔴 **COMANDO EXTRA** |

---

## 📊 Comandos de `volumes`

| Comando | MGC | TMP-CLI | Status |
|---------|-----|---------|--------|
| `attach` | ✅ | ✅ | ⚠️ Bug "doto3" |
| `create` | ✅ | ✅ | ⚠️ Bug "doto3" + flags faltando |
| `delete` | ✅ | ✅ | ⚠️ Bug "doto3" |
| `detach` | ✅ | ✅ | ⚠️ Bug "doto3" |
| `extend` | ✅ | ✅ | ⚠️ Bug "doto3" |
| `get` | ✅ | ✅ | ⚠️ Bug "doto3" |
| `list` | ✅ | ✅ | ⚠️ Bug "doto3" + flags sem prefixo |
| `list-all` | ❌ | ✅ | 🔴 **COMANDO EXTRA** |
| `rename` | ✅ | ✅ | ⚠️ Bug "doto3" |
| `retype` | ✅ | ✅ | ⚠️ Bug "doto3" |

---

## 📊 Aliases do Comando Principal

| Alias | MGC | TMP-CLI | Status |
|-------|-----|---------|--------|
| `bs` | ✅ | ❌ | 🔴 **FALTANDO** |
| `blocks` | ✅ | ❌ | 🔴 **FALTANDO** |
| `blst` | ✅ | ❌ | 🔴 **FALTANDO** |
| `block` | ✅ | ❌ | 🔴 **FALTANDO** |
| `volumes` | ✅ | ❌ | 🔴 **FALTANDO** |

---

## 📊 Padrões de Flags em Comandos `list`

### schedulers list

| Flag | MGC | TMP-CLI | Problema |
|------|-----|---------|----------|
| `--control.limit` | ✅ | ❌ | TMP-CLI: `--limit` (sem prefixo) |
| `--control.offset` | ✅ | ❌ | TMP-CLI: `--offset` (sem prefixo) |
| `--control.sort` | ✅ | ❌ | TMP-CLI: `--sort` (sem prefixo) |
| `--expand` | ✅ | ✅ | OK (mas sem descrição) |
| **Required incorreto** | ❌ | ✅ | TMP-CLI marca todas como required |
| **Argumentos posicionais** | ❌ | ✅ | TMP-CLI mostra incorretamente no Usage |

### volumes list

| Flag | MGC | TMP-CLI | Problema |
|------|-----|---------|----------|
| `--control.limit` | ✅ | ❌ | TMP-CLI: `--limit` (sem prefixo) |
| `--control.offset` | ✅ | ❌ | TMP-CLI: `--offset` (sem prefixo) |
| `--control.sort` | ✅ | ❌ | TMP-CLI: `--sort` (sem prefixo) |
| `--expand` | ✅ | ✅ | OK (mas sem descrição) |
| `--name` | ✅ | ❌ | 🔴 **FALTANDO** |
| **Required incorreto** | ❌ | ✅ | TMP-CLI marca todas como required |

### volume-types list

| Flag | MGC | TMP-CLI | Problema |
|------|-----|---------|----------|
| `--allows-encryption` | ✅ | ✅ | OK (mas sem descrição + required incorreto) |
| `--availability-zone` | ✅ | ✅ | OK (mas sem descrição + required incorreto) |
| `--name` | ✅ | ✅ | OK (mas sem descrição + required incorreto) |
| `--limit` | ❌ | ✅ | 🔴 **FLAG EXTRA** (não existe no MGC) |
| `--offset` | ❌ | ✅ | 🔴 **FLAG EXTRA** (não existe no MGC) |
| `--sort` | ❌ | ✅ | 🔴 **FLAG EXTRA** (não existe no MGC) |
| **Argumentos posicionais** | ❌ | ✅ | TMP-CLI mostra incorretamente no Usage |

---

## 📊 Flags em Comandos `create`

### schedulers create

| Flag/Estrutura | MGC | TMP-CLI | Problema |
|----------------|-----|---------|----------|
| `--cli.list-links` | ✅ | ❌ | 🔴 **FALTANDO** |
| `--description` | ✅ | ✅ | OK (mas sem descrição) |
| `--name` | ✅ | ✅ | OK (mas sem descrição) |
| `--policy` (objeto) | ✅ | ❌ | Falta objeto pai |
| `--policy.frequency` (objeto) | ✅ | ❌ | Falta objeto intermediário |
| `--policy.frequency.daily` (objeto) | ✅ | ❌ | Falta objeto intermediário |
| `--policy.frequency.daily.start-time` | ✅ | ✅ | OK (mas sem descrição + required incorreto) |
| `--policy.retention-in-days` | ✅ | ✅ | OK (mas sem descrição + required incorreto) |
| `--snapshot` (objeto) | ✅ | ❌ | Falta objeto pai |
| `--snapshot.type` | ✅ | ✅ | OK (mas sem descrição) |

### snapshots create

| Flag/Estrutura | MGC | TMP-CLI | Problema |
|----------------|-----|---------|----------|
| `--cli.list-links` | ✅ | ❌ | 🔴 **FALTANDO** |
| `--description` | ✅ | ✅ | OK (mas sem descrição) |
| `--name` | ✅ | ✅ | OK (mas sem descrição) |
| `--type` | ✅ | ✅ | OK (mas sem descrição) |
| `--source-snapshot` (objeto) | ✅ | ❌ | 🔴 **FALTANDO** |
| `--source-snapshot.id` | ✅ | ❌ | 🔴 **FALTANDO** |
| `--source-snapshot.name` | ✅ | ❌ | 🔴 **FALTANDO** |
| `--volume` (objeto) | ✅ | ❌ | 🔴 **FALTANDO** |
| `--volume.id` | ✅ | ❌ | 🔴 **FALTANDO** |
| `--volume.name` | ✅ | ❌ | 🔴 **FALTANDO** |
| **Examples** | ✅ | ❌ | Falta seção de exemplos |

### volumes create

| Flag/Estrutura | MGC | TMP-CLI | Problema |
|----------------|-----|---------|----------|
| `--cli.list-links` | ✅ | ❌ | 🔴 **FALTANDO** |
| `--availability-zone` | ✅ | ✅ | OK (mas sem descrição) |
| `--encrypted` | ✅ | ✅ | TMP-CLI marca como required (deveria ser opcional) |
| `--name` | ✅ | ✅ | OK (mas sem descrição) |
| `--size` | ✅ | ✅ | OK (mas sem descrição) |
| `--type` (objeto) | ✅ | ❌ | Falta objeto pai |
| `--type.id` | ✅ | ✅ | OK (mas sem descrição) |
| `--type.name` | ✅ | ✅ | OK (mas sem descrição) |
| `--snapshot` (objeto) | ✅ | ❌ | 🔴 **FALTANDO** |
| `--snapshot.id` | ✅ | ❌ | 🔴 **FALTANDO** |
| `--snapshot.name` | ✅ | ❌ | 🔴 **FALTANDO** |
| **Examples** | ✅ | ❌ | Falta seção de exemplos |

---

## 📊 Flags em `schedulers attach`

| Flag/Estrutura | MGC: `attach` | TMP-CLI: `attach-volume` | Problema |
|----------------|---------------|--------------------------|----------|
| `--cli.list-links` | ✅ | ❌ | 🔴 **FALTANDO** |
| `--id` | ✅ | ✅ | OK (mas sem descrição) |
| `--volume` (objeto) | ✅ | ❌ | 🔴 **FALTANDO** |
| `--volume.id` | ✅ | ❌ | 🔴 **FALTANDO** |
| `--volume.name` | ✅ | ❌ | 🔴 **FALTANDO** |
| **Examples** | ✅ | ❌ | Falta seção de exemplos |

---

## 🐛 Bugs Encontrados

| Bug | Onde Aparece | Quantidade |
|-----|--------------|------------|
| String "Dqui1" | Todos os subcomandos nível 2 | 4 ocorrências |
| String "doto3" | Todos os comandos leaf (nível 3+) | ~25 ocorrências |
| Descrições ausentes | Todas as flags | 100% das flags |
| Argumentos posicionais incorretos | Comandos list | Múltiplos comandos |
| Required incorreto | Múltiplos comandos | ~70% dos comandos |

---

## 📝 Estatísticas Gerais

### Comandos
- **Total MGC:** 25 comandos leaf
- **Total TMP-CLI:** 28 comandos leaf
- **Faltando:** 1 comando (`snapshots copy`)
- **Extras:** 4 comandos (`list-all` em 4 grupos)
- **Nomes divergentes:** 2 comandos (`attach` e `detach` em schedulers)
- **Compatibilidade:** 88% (22/25 comandos presentes e com nome correto)

### Aliases
- **Total MGC:** 6 aliases (incluindo nome principal)
- **Total TMP-CLI:** 1 (só o nome principal)
- **Compatibilidade:** 17% (1/6)

### Descrições
- **Descrições de comandos corretas:** 0% (todas genéricas)
- **Descrições de flags:** 0% (todas ausentes)
- **Seções Rules/Notes:** 0% (todas ausentes)
- **Seções Examples:** 0% (todas ausentes)

### Flags
- **Flags com prefixo correto:** ~30% (maioria sem `control.`)
- **Flags com descrição:** 0%
- **Flags marcadas corretamente como required:** ~30%
- **Comandos com `--cli.list-links`:** 0% (deveria estar em ~40%)

---

## 🎯 Prioridades de Correção

### P0 - CRÍTICO (Bloqueia uso correto)
| Item | Impacto | Ocorrências |
|------|---------|-------------|
| Remover "Dqui1" | Visual/profissional | 4 |
| Remover "doto3" | Visual/profissional | ~25 |
| Adicionar descrições nas flags | Usabilidade | 100% das flags |
| Corrigir marcação required | Funcionalidade | ~70% dos comandos |
| Adicionar comando `snapshots copy` | Funcionalidade ausente | 1 |

### P1 - ALTO (Incompatibilidade com referência)
| Item | Impacto | Ocorrências |
|------|---------|-------------|
| Renomear `attach-volume` → `attach` | Compatibilidade | 1 |
| Renomear `detach-volume` → `detach` | Compatibilidade | 1 |
| Adicionar aliases | UX/Compatibilidade | 5 aliases |
| Adicionar prefixo `control.` | Compatibilidade | ~15 flags |
| Corrigir descrições de comandos | Documentação | Todos os comandos |
| Adicionar objetos pai faltantes | Funcionalidade | ~10 estruturas |
| Adicionar flags `.id`/`.name` | Funcionalidade | ~8 pares |
| Remover argumentos posicionais | Clareza/docs | ~6 comandos |

### P2 - MÉDIO (Melhorias importantes)
| Item | Impacto | Ocorrências |
|------|---------|-------------|
| Adicionar `--cli.list-links` | Funcionalidade | ~10 comandos |
| Adicionar seções Rules/Notes | Documentação | ~5 comandos |
| Adicionar Examples | Documentação | ~8 comandos |
| Remover flags extras | Compatibilidade | 3 flags em volume-types list |
| Adicionar flag `--name` em volumes list | Funcionalidade | 1 |
| Decidir sobre `list-all` | Decisão de design | 4 comandos |

### P3 - BAIXO (Polish)
- Melhorar formatação geral
- Padronizar estrutura de descrições
- Revisar documentação inline

---

## ✅ Checklist Rápido (Top 10 Ações)

1. [ ] **Remover bugs de strings** ("Dqui1" e "doto3") - ~29 ocorrências
2. [ ] **Adicionar descrições em TODAS as flags** - Impacto massivo na usabilidade
3. [ ] **Corrigir marcação de flags required** - Evita erros do usuário
4. [ ] **Adicionar prefixo `control.`** em limit/offset/sort - ~15 flags
5. [ ] **Implementar comando `snapshots copy`** - Funcionalidade ausente
6. [ ] **Renomear comandos** attach-volume→attach, detach-volume→detach
7. [ ] **Adicionar aliases** (bs, blocks, blst, block, volumes)
8. [ ] **Adicionar objetos e flags faltantes** (volume, snapshot, type, etc.)
9. [ ] **Corrigir descrições de comandos** - Todas específicas ao invés de genéricas
10. [ ] **Adicionar Examples, Rules e Notes** - Melhora documentação

---

## 💡 Observação Importante

Os mesmos padrões de problemas identificados em `audit` se repetem em `block-storage`:
- Bugs de strings estranhas ("Dqui1", "doto3")
- Descrições genéricas e repetitivas
- Flags sem descrição
- Comandos `list-all` extras
- Falta de prefixo `control.`
- Marcação incorreta de `required`

Isso sugere que estes são **problemas sistemáticos do gerador de CLI**, não problemas isolados. A correção deve ser feita no gerador para resolver todos os produtos de uma vez.

