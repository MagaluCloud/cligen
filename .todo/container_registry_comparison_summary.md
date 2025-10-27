# Resumo Executivo - Comparação MGC vs TMP-CLI (container-registry)

## 📊 Tabela Geral de Subcomandos

| Subcomando | MGC | TMP-CLI | Status |
|------------|-----|---------|--------|
| **Nível 1: container-registry** |
| Comando principal | ✅ | ✅ | ⚠️ Descrição divergente + sem aliases |
| **Nível 2: Subcomandos** |
| credentials | ✅ | ✅ | ⚠️ Bug "Dqui1" + nomes divergentes |
| images | ✅ | ✅ | ⚠️ Bug "Dqui1" + list-all extra |
| proxy-caches | ✅ | ❌ | 🔴 **SUBCOMANDO INTEIRO FALTANDO** |
| registries | ✅ | ✅ | ⚠️ Bug "Dqui1" + list-all extra |
| repositories | ✅ | ✅ | ⚠️ Bug "Dqui1" + list-all extra |

---

## 📊 Comandos de `credentials`

| Comando | MGC | TMP-CLI | Status |
|---------|-----|---------|--------|
| `list` | ✅ | ❌ | 🔴 TMP-CLI usa `get` |
| `get` | ❌ | ✅ | 🔴 Nome divergente (deveria ser `list`) |
| `password` | ✅ | ❌ | 🔴 TMP-CLI usa `reset-password` |
| `reset-password` | ❌ | ✅ | 🔴 Nome divergente (deveria ser `password`) |

---

## 📊 Comandos de `images`

| Comando | MGC | TMP-CLI | Status |
|---------|-----|---------|--------|
| `delete` | ✅ | ✅ | ⚠️ Bug "doto3" + repository-name incorreto |
| `get` | ✅ | ✅ | ⚠️ Bug "doto3" |
| `list` | ✅ | ✅ | ⚠️ Bug "doto3" + repository-name incorreto |
| `list-all` | ❌ | ✅ | 🔴 **COMANDO EXTRA** |

---

## 📊 Comandos de `proxy-caches` (FALTANDO)

| Comando | MGC | TMP-CLI | Status |
|---------|-----|---------|--------|
| `create` | ✅ | ❌ | 🔴 **COMANDO FALTANDO** |
| `delete` | ✅ | ❌ | 🔴 **COMANDO FALTANDO** |
| `get` | ✅ | ❌ | 🔴 **COMANDO FALTANDO** |
| `list` | ✅ | ❌ | 🔴 **COMANDO FALTANDO** |
| `status` | ✅ | ❌ | 🔴 **COMANDO FALTANDO** |
| `update` | ✅ | ❌ | 🔴 **COMANDO FALTANDO** |

**Total:** 6 comandos completamente ausentes

---

## 📊 Comandos de `registries`

| Comando | MGC | TMP-CLI | Status |
|---------|-----|---------|--------|
| `create` | ✅ | ✅ | ⚠️ Bug "doto3" + flags faltando |
| `delete` | ✅ | ✅ | ⚠️ Bug "doto3" |
| `get` | ✅ | ✅ | ⚠️ Bug "doto3" |
| `list` | ✅ | ✅ | ⚠️ Bug "doto3" + flags sem prefixo |
| `list-all` | ❌ | ✅ | 🔴 **COMANDO EXTRA** |

---

## 📊 Comandos de `repositories`

| Comando | MGC | TMP-CLI | Status |
|---------|-----|---------|--------|
| `delete` | ✅ | ✅ | ⚠️ Bug "doto3" |
| `get` | ✅ | ✅ | ⚠️ Bug "doto3" |
| `list` | ✅ | ✅ | ⚠️ Bug "doto3" + flags sem prefixo |
| `list-all` | ❌ | ✅ | 🔴 **COMANDO EXTRA** |

---

## 📊 Aliases do Comando Principal

| Alias | MGC | TMP-CLI | Status |
|-------|-----|---------|--------|
| `container-registry` | ✅ | ✅ | OK |
| `cr` | ✅ | ❌ | 🔴 **FALTANDO** |
| `registry` | ✅ | ❌ | 🔴 **FALTANDO** |

---

## 🔴 Problema CRÍTICO: Mudança de Tipo de Parâmetro

### `repository-id` vs `repository-name`

| Comando | MGC | TMP-CLI | Impacto |
|---------|-----|---------|---------|
| `images list` | `--repository-id` (UUID) | `--repository-name` (Nome) | 🔴 **TIPO DIFERENTE** |
| `images delete` | `--repository-id` (UUID) | `--repository-name` (Nome) | 🔴 **TIPO DIFERENTE** |
| `images get` | `--repository-id` (UUID) | `--repository-name` (Nome) | 🔴 **TIPO DIFERENTE** |

**Severidade:** CRÍTICA  
**Impacto:** Funcionalidade diferente! MGC aceita UUID, TMP-CLI aceita Nome.  
**Consequência:** Scripts e automações não funcionam, comportamento da API pode ser diferente.

---

## 📊 Padrões de Flags em Comandos `list`

### registries list

| Flag | MGC | TMP-CLI | Problema |
|------|-----|---------|----------|
| `--control.limit` | ✅ | ❌ | TMP-CLI: `--limit` (sem prefixo) |
| `--control.offset` | ✅ | ❌ | TMP-CLI: `--offset` (sem prefixo) |
| `--control.sort` | ✅ | ❌ | TMP-CLI: `--sort` (sem prefixo) |
| `--name` | ✅ | ❌ | 🔴 **FALTANDO** |
| **Required incorreto** | ❌ | ✅ | TMP-CLI marca todas como required |
| **Argumentos posicionais** | ❌ | ✅ | TMP-CLI mostra incorretamente no Usage |

### images list

| Flag | MGC | TMP-CLI | Problema |
|------|-----|---------|----------|
| `--control.limit` | ✅ | ❌ | TMP-CLI: `--limit` (sem prefixo) |
| `--control.offset` | ✅ | ❌ | TMP-CLI: `--offset` (sem prefixo) |
| `--control.sort` | ✅ | ❌ | TMP-CLI: `--sort` (sem prefixo) |
| `--expand` | ✅ | ✅ | OK (mas sem descrição) |
| `--name` | ✅ | ❌ | 🔴 **FALTANDO** |
| `--registry-id` | ✅ | ✅ | OK (mas sem descrição) |
| `--repository-id` | ✅ | ❌ | TMP-CLI: `--repository-name` 🔴 **TIPO DIFERENTE** |
| `--repository-name` | ❌ | ✅ | Deveria ser `--repository-id` |
| **Required incorreto** | ❌ | ✅ | TMP-CLI marca todas como required |
| **Argumentos posicionais** | `[registry-id] [repository-id]` | `[registryID] [repositoryName] [Offset] [Limit] [Sort] [Expand]` | Completamente divergente |

### repositories list

| Flag | MGC | TMP-CLI | Problema |
|------|-----|---------|----------|
| `--control.limit` | ✅ | ❌ | TMP-CLI: `--limit` (sem prefixo) |
| `--control.offset` | ✅ | ❌ | TMP-CLI: `--offset` (sem prefixo) |
| `--control.sort` | ✅ | ❌ | TMP-CLI: `--sort` (sem prefixo) |
| `--name` | ✅ | ❌ | 🔴 **FALTANDO** |
| `--registry-id` | ✅ | ✅ | OK (mas sem descrição) |
| **Required incorreto** | ❌ | ✅ | TMP-CLI marca todas como required |
| **Argumentos posicionais** | `[registry-id]` | `[registryID] [Offset] [Limit] [Sort]` | Divergente |

---

## 📊 Flags em Comandos `create`

### registries create

| Flag | MGC | TMP-CLI | Problema |
|------|-----|---------|----------|
| `--cli.list-links` | ✅ | ❌ | 🔴 **FALTANDO** |
| `--name` | ✅ | ✅ | OK (mas sem descrição) |
| `--proxy-cache-id` | ✅ | ❌ | 🔴 **FALTANDO** |

### proxy-caches create (FALTANDO)

| Flag | MGC | TMP-CLI | Problema |
|------|-----|---------|----------|
| `--access-key` | ✅ | ❌ | Comando não existe |
| `--access-secret` | ✅ | ❌ | Comando não existe |
| `--description` | ✅ | ❌ | Comando não existe |
| `--name` | ✅ | ❌ | Comando não existe |
| `--provider` | ✅ | ❌ | Comando não existe |
| `--url` | ✅ | ❌ | Comando não existe |

---

## 📊 Flags em Comandos `delete`

### images delete

| Flag | MGC | TMP-CLI | Problema |
|------|-----|---------|----------|
| `--digest-or-tag` | ✅ | ✅ | OK (mas sem descrição) |
| `--registry-id` | ✅ | ✅ | OK (mas sem descrição) |
| `--repository-id` | ✅ | ❌ | TMP-CLI: `--repository-name` 🔴 **TIPO DIFERENTE** |
| `--repository-name` | ❌ | ✅ | Deveria ser `--repository-id` |

---

## 🐛 Bugs Encontrados

| Bug | Onde Aparece | Quantidade |
|-----|--------------|------------|
| String "Dqui1" | Todos os subcomandos nível 2 | 4 ocorrências |
| String "doto3" | Todos os comandos leaf (nível 3) | ~13 ocorrências |
| Descrições ausentes | Todas as flags | 100% das flags |
| Argumentos posicionais incorretos | Comandos list e delete | Múltiplos comandos |
| Required incorreto | Múltiplos comandos | ~80% dos comandos list |
| Tipo de parâmetro incorreto | `repository-name` vs `repository-id` | 🔴 **CRÍTICO** |

---

## 📝 Estatísticas Gerais

### Subcomandos (Nível 2)
- **Total MGC:** 5 grupos
- **Total TMP-CLI:** 4 grupos
- **Faltando:** 1 grupo completo (proxy-caches)
- **Compatibilidade:** 80% (4/5)
- **Bugs "Dqui1":** 100% (4/4 dos grupos presentes)

### Comandos Leaf (Nível 3)
- **Total MGC:** 19 comandos
- **Total TMP-CLI:** 13 comandos presentes + 3 extras = 16 total
- **Faltando:** 6 comandos (todo proxy-caches)
- **Nome divergente:** 2 comandos (credentials list/get, password/reset-password)
- **Comandos extras:** 3 comandos (list-all)
- **Compatibilidade:** 68% (13/19, mas 2 com nomes divergentes)
- **Bugs "doto3":** 100% (todos os comandos leaf)

### Aliases
- **Total MGC:** 3 aliases
- **Total TMP-CLI:** 1 (só o nome principal)
- **Compatibilidade:** 33% (1/3)

### Descrições
- **Descrições de comandos corretas:** 0% (todas genéricas)
- **Descrições de flags:** 0% (todas ausentes)

### Flags
- **Flags com prefixo correto:** ~20% (maioria sem `control.`)
- **Flags com descrição:** 0%
- **Flags marcadas corretamente como required:** ~20%
- **Flags com tipo correto:** ~80% (problema grave com repository-id)
- **Comandos com `--cli.list-links`:** 0% (deveria estar em create)

---

## 🎯 Prioridades de Correção

### P0 - CRÍTICO (Bloqueia funcionalidade)
| Item | Impacto | Ocorrências |
|------|---------|-------------|
| **Corrigir `--repository-name` → `--repository-id`** | Funcionalidade quebrada | 3 comandos |
| **Implementar grupo `proxy-caches`** | Funcionalidade ausente | 6 comandos |
| Remover "Dqui1" | Visual/profissional | 4 |
| Remover "doto3" | Visual/profissional | ~13 |
| Adicionar descrições nas flags | Usabilidade | 100% das flags |

### P1 - ALTO (Incompatibilidade com referência)
| Item | Impacto | Ocorrências |
|------|---------|-------------|
| Renomear `get` → `list` em credentials | Compatibilidade | 1 |
| Renomear `reset-password` → `password` | Compatibilidade | 1 |
| Adicionar aliases | UX/Compatibilidade | 2 aliases |
| Adicionar prefixo `control.` | Compatibilidade | ~12 flags |
| Corrigir descrições de comandos | Documentação | Todos os comandos |
| Remover argumentos posicionais incorretos | Clareza/docs | ~4 comandos |
| Padronizar formato argumentos (kebab-case) | Consistência | Todos os comandos |
| Corrigir marcação required | Funcionalidade | ~80% dos list |

### P2 - MÉDIO (Melhorias importantes)
| Item | Impacto | Ocorrências |
|------|---------|-------------|
| Adicionar `--cli.list-links` | Funcionalidade | ~1 comando |
| Adicionar `--proxy-cache-id` | Funcionalidade | 1 flag |
| Adicionar flag `--name` em lists | Funcionalidade | 3 comandos |
| Remover comandos `list-all` | Compatibilidade | 3 comandos |

### P3 - BAIXO (Polish)
- Melhorar formatação geral
- Padronizar estrutura de descrições

---

## ⚠️ Observações Especiais

### 1. Problema de Tipo CRÍTICO

O mais grave problema encontrado é a divergência `--repository-id` vs `--repository-name`:

**MGC (correto):**
```bash
mgc container-registry images list --registry-id <UUID> --repository-id <UUID>
```

**TMP-CLI (incorreto):**
```bash
cli container-registry images list --registry-id <UUID> --repository-name <NOME>
```

**Impacto:**
- ❌ Scripts existentes não funcionam
- ❌ Tipo de parâmetro diferente (UUID vs String)
- ❌ Possível comportamento diferente da API
- ❌ Documentação não se aplica

### 2. Grupo Inteiro Ausente

O grupo `proxy-caches` está completamente ausente no TMP-CLI:
- 6 comandos faltando
- Funcionalidade de proxy cache não disponível
- Possível problema no processo de geração (grupo foi esquecido?)

### 3. Nomenclatura de Comandos

Divergências nos comandos `credentials`:
- `list` → `get`: Semanticamente diferente (list sugere múltiplos, get sugere único)
- `password` → `reset-password`: Mais verboso mas mais claro

**Decisão necessária:** Manter compatibilidade (usar nomes MGC) ou melhorar clareza (manter TMP-CLI)?

---

## ✅ Checklist Rápido (Top 10 Ações)

1. [ ] 🔴 **URGENTE: Corrigir `repository-name` → `repository-id`** - Tipo incorreto
2. [ ] 🔴 **URGENTE: Implementar grupo `proxy-caches` completo** - 6 comandos
3. [ ] **Remover bugs de strings** ("Dqui1" e "doto3") - ~17 ocorrências
4. [ ] **Adicionar descrições em TODAS as flags** - Impacto massivo na usabilidade
5. [ ] **Corrigir marcação de flags required** - Evita erros do usuário
6. [ ] **Adicionar prefixo `control.`** em limit/offset/sort - ~12 flags
7. [ ] **Renomear comandos** get→list, reset-password→password
8. [ ] **Adicionar aliases** (cr, registry)
9. [ ] **Adicionar flags faltantes** (name, proxy-cache-id, cli.list-links)
10. [ ] **Corrigir descrições de comandos** - Todas específicas ao invés de genéricas

---

## 💡 Insights

### Padrões Confirmados (vs audit e block-storage)
✅ Mesmos bugs visuais ("Dqui1", "doto3")  
✅ Mesmas descrições genéricas  
✅ Mesmos problemas de prefixo `control.`  
✅ Mesmos comandos `list-all` extras  
✅ Mesma falta de descrições de flags

### Problemas NOVOS (específicos de container-registry)
🆕 **Grupo inteiro faltando** (proxy-caches) - Não visto antes  
🆕 **Mudança de tipo de parâmetro** (repository-id vs repository-name) - Gravíssimo  
🆕 **Nomenclatura de comandos** (list/get, password/reset-password)  
🆕 **Formato de argumentos posicionais** (kebab-case vs camelCase)

### Conclusão
Os problemas sistemáticos se confirmam em container-registry, MAS aparecem novos problemas graves:
1. Grupos inteiros podem ser esquecidos na geração
2. Tipos de parâmetros podem ser mudados incorretamente
3. Nomenclatura pode divergir semanticamente

Isso sugere que além dos problemas sistemáticos já identificados, há **problemas adicionais de mapeamento SDK → CLI** que precisam ser investigados.

