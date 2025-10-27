# Comparação Detalhada: mgc network vs ./tmp-cli/cli network

## Sumário Executivo

A comparação entre `mgc network` e `./tmp-cli/cli network` revela diferenças **arquiteturais críticas** e problemas sistemáticos graves no CLI gerado:

### 🔴 Problemas Críticos Identificados

1. **BUGS VISUAIS**: Strings "Dqui1" e "doto3" aparecem aleatoriamente no output
2. **ARQUITETURA DIVERGENTE**: Subgrupos em `mgc` foram convertidos em comandos flat em `./tmp-cli/cli`
3. **NOMES DE GRUPOS INCONSISTENTES**: `public-ips` → `public-i-ps`, `vpcs` → `v-p-cs`, `subnetpools` → `subnet-pools`
4. **NOMES DE COMANDOS DIVERGENTES**: `attach` → `attach-security-group`, `attach` → `attach-to-port`
5. **NOMES DE FLAGS INCONSISTENTES**: Remoção sistemática de hífens (ex: `vpc-id` → `vpcid`, `cidr-block` → `cidrblock`)
6. **0% DE DESCRIÇÕES NAS FLAGS**: Nenhuma flag em `./tmp-cli/cli` possui descrição
7. **FLAGS GLOBAIS AUSENTES**: `--cli.retry-until`, `-t/--cli.timeout`, `--env`, `--region`, `--server-url`, `-o/--output`, `--x-zone` completamente ausentes
8. **SUBGRUPO AUSENTE**: `security-groups rules` existe em `mgc` mas não em `./tmp-cli/cli`
9. **COMANDOS EXTRAS**: `create` e `list` existem em `./tmp-cli/cli network rules` mas não em `mgc`
10. **FLAGS INCORRETAMENTE MARCADAS COMO REQUIRED**: Em `list`, flags opcionais como `limit`, `offset`, `sort` marcadas como `(required)`

---

## 1. Comando Principal

### mgc network

```
VPC Api Product

Usage:
  mgc network [flags]
  mgc network [command]

Commands:
  nat-gateways    Operations related to Nat Gateway
  ports           Operations related to Ports
  public-ips      Operations related to Public IPs
  rules           Operations related to Rules
  security-groups Operations related to Security Groups
  subnetpools     Operations related to Subnet Pools
  subnets         Operations related to Subnets
  vpcs            Operations related to VPCs
```

### ./tmp-cli/cli network

```
Network provides a client for interacting with the Magalu Cloud Network API.

Package network provides a client for interacting with the Magalu Cloud Network API.
This package allows you to manage VPCs, subnets, ports, security groups, rules, public IPs, subnet pools, and NAT gateways.

Available Commands:

Other commands:
  nat-gateways        Network provides a client for interacting with the Magalu Cloud Network API.
  ports               Network provides a client for interacting with the Magalu Cloud Network API.
  public-i-ps         Network provides a client for interacting with the Magalu Cloud Network API.
  rules               Network provides a client for interacting with the Magalu Cloud Network API.
  security-groups     Network provides a client for interacting with the Magalu Cloud Network API.
  subnet-pools        Network provides a client for interacting with the Magalu Cloud Network API.
  subnets             Network provides a client for interacting with the Magalu Cloud Network API.
  v-p-cs              Network provides a client for interacting with the Magalu Cloud Network API.
```

### ⚠️ Diferenças Identificadas

| Aspecto | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **Nome do grupo `public-ips`** | `public-ips` | `public-i-ps` | ❌ Divergente |
| **Nome do grupo `vpcs`** | `vpcs` | `v-p-cs` | ❌ Divergente |
| **Nome do grupo `subnetpools`** | `subnetpools` | `subnet-pools` | ❌ Divergente |
| **Descrições dos comandos** | Específicas e úteis | Genéricas e repetitivas | ❌ Divergente |
| **Flags globais** | 7 flags (incluindo `--cli.retry-until`, `-t/--cli.timeout`, `-o/--output`) | 5 flags (faltam 3) | ❌ Incompleto |

---

## 2. Grupo: nat-gateways

### mgc network nat-gateways

```
Operations related to Nat Gateway

Usage:
  mgc network nat-gateways [flags]
  mgc network nat-gateways [command]

Aliases:
  nat-gateways, nat_gateways

Commands:
  create      Create a new NAT Gateway resource
  delete      Delete a NAT Gateway from a VPC
  get         Detail a NAT Gateway from a VPC
  list        List the NAT Gateways from a VPC
```

### ./tmp-cli/cli network nat-gateways

```
Network provides a client for interacting with the Magalu Cloud Network API.

Dqui1

Available Commands:

Other commands:
  create              Network provides a client for interacting with the Magalu Cloud Network API.
  delete              Network provides a client for interacting with the Magalu Cloud Network API.
  get                 Network provides a client for interacting with the Magalu Cloud Network API.
  list                Network provides a client for interacting with the Magalu Cloud Network API.
```

### ⚠️ Diferenças

| Aspecto | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **Bug visual** | Não | "Dqui1" aparece | 🔴 CRÍTICO |
| **Aliases** | `nat-gateways`, `nat_gateways` | Nenhum | ❌ Faltando |
| **Descrições dos comandos** | Específicas | Genéricas | ❌ Divergente |

### 2.1. nat-gateways create

**mgc**:
- Flags: `--description`, `--name` (required), `--vpc-id` (required), `--zone` (required)
- Todas as flags têm descrições completas

**./tmp-cli/cli**:
- Flags: `--description`, `--name` (required), `--vpcid` (required), `--zone` (required)
- Bug visual: "doto3" aparece
- 0% das flags possuem descrições

**Divergências**:
- ❌ `--vpc-id` vs `--vpcid` (falta hífen)
- 🔴 Bug visual: "doto3"
- ❌ Faltam descrições nas flags

### 2.2. nat-gateways list

**mgc**:
- Flags: `--items-per-page` (range: 1-100), `--page` (min: 1), `--sort`, `--vpc-id` (required)

**./tmp-cli/cli**:
- Flags: `--limit` (required), `--offset` (required), `--sort` (required), `--vpc-id` (required)
- Bug visual: "doto3"

**Divergências**:
- ❌ `--items-per-page` e `--page` vs `--limit` e `--offset` (nomenclatura diferente)
- ❌ `--limit`, `--offset`, `--sort` marcadas incorretamente como `(required)`
- 🔴 Bug visual: "doto3"

---

## 3. Grupo: ports

### mgc network ports

```
Operations related to Ports

Usage:
  mgc network ports [flags]
  mgc network ports [command]

Commands:
  attach      Attach Security Group
  delete      Delete Port
  detach      Detach Security Group
  get         Port Details
  list        Details of a Port list
  update      Update a Port
```

### ./tmp-cli/cli network ports

```
Network provides a client for interacting with the Magalu Cloud Network API.

Dqui1

Available Commands:

Other commands:
  attach-security-groupNetwork provides a client for interacting with the Magalu Cloud Network API.
  delete              Network provides a client for interacting with the Magalu Cloud Network API.
  detach-security-groupNetwork provides a client for interacting with the Magalu Cloud Network API.
  get                 Network provides a client for interacting with the Magalu Cloud Network API.
  list                Network provides a client for interacting with the Magalu Cloud Network API.
  update              Network provides a client for interacting with the Magalu Cloud Network API.
```

### ⚠️ Diferenças

| Aspecto | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **Bug visual** | Não | "Dqui1" aparece | 🔴 CRÍTICO |
| **Nome comando `attach`** | `attach` | `attach-security-group` | ❌ Divergente |
| **Nome comando `detach`** | `detach` | `detach-security-group` | ❌ Divergente |
| **Bug de formatação** | Não | Falta espaço entre nome e descrição | 🔴 CRÍTICO |

### 3.1. ports attach

**mgc**: `mgc network ports attach [port-id] [security-group-id]`
- Flags: `--port-id` (required), `--security-group-id` (required)

**./tmp-cli/cli**: `cli network ports attach-security-group [portID] [securityGroupID]`
- Flags: `--port-id` (required), `--security-group-id` (required)
- Bug visual: "doto3"

**Divergências**:
- ❌ Nome do comando: `attach` vs `attach-security-group`
- 🔴 Bug visual: "doto3"
- ❌ Faltam descrições nas flags

---

## 4. Grupo: public-ips

### mgc network public-ips

```
Operations related to Public IPs

Usage:
  mgc network public-ips [flags]
  mgc network public-ips [command]

Aliases:
  public-ips, public_ips

Commands:
  attach      Attach Public IP
  delete      Delete Public IP
  detach      Detach Public IP
  get         Public IP Details
  list        Tenant's public IP list
```

### ./tmp-cli/cli network public-i-ps

```
Network provides a client for interacting with the Magalu Cloud Network API.

Dqui1

Available Commands:

Other commands:
  attach-to-port      Network provides a client for interacting with the Magalu Cloud Network API.
  delete              Network provides a client for interacting with the Magalu Cloud Network API.
  detach-from-port    Network provides a client for interacting with the Magalu Cloud Network API.
  get                 Network provides a client for interacting with the Magalu Cloud Network API.
  list                Network provides a client for interacting with the Magalu Cloud Network API.
```

### ⚠️ Diferenças

| Aspecto | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **Nome do grupo** | `public-ips` | `public-i-ps` | 🔴 CRÍTICO |
| **Bug visual** | Não | "Dqui1" aparece | 🔴 CRÍTICO |
| **Nome comando `attach`** | `attach` | `attach-to-port` | ❌ Divergente |
| **Nome comando `detach`** | `detach` | `detach-from-port` | ❌ Divergente |
| **Aliases** | `public-ips`, `public_ips` | Nenhum | ❌ Faltando |

### 4.1. public-ips attach

**mgc**: `mgc network public-ips attach [public-ip-id] [port-id]`
- Flags: `--port-id` (required), `--public-ip-id` (required)

**./tmp-cli/cli**: `cli network public-i-ps attach-to-port [publicIPID] [portID]`
- Flags: `--port-id` (required), `--public-ipid` (required)
- Bug visual: "doto3"

**Divergências**:
- 🔴 Nome do grupo: `public-ips` vs `public-i-ps`
- ❌ Nome do comando: `attach` vs `attach-to-port`
- ❌ Nome da flag: `--public-ip-id` vs `--public-ipid`
- 🔴 Bug visual: "doto3"

---

## 5. Grupo: rules

### mgc network rules

```
Operations related to Rules

Usage:
  mgc network rules [flags]
  mgc network rules [command]

Commands:
  delete      Delete a Rule
  get         Rule Details
```

### ./tmp-cli/cli network rules

```
Network provides a client for interacting with the Magalu Cloud Network API.

Dqui1

Available Commands:

Other commands:
  create              Network provides a client for interacting with the Magalu Cloud Network API.
  delete              Network provides a client for interacting with the Magalu Cloud Network API.
  get                 Network provides a client for interacting with the Magalu Cloud Network API.
  list                Network provides a client for interacting with the Magalu Cloud Network API.
```

### ⚠️ Diferenças

| Aspecto | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **Bug visual** | Não | "Dqui1" aparece | 🔴 CRÍTICO |
| **Comando `create`** | ❌ Não existe | ✅ Existe | ⚠️ Extra |
| **Comando `list`** | ❌ Não existe | ✅ Existe | ⚠️ Extra |

**⚠️ ALERTA**: `./tmp-cli/cli` possui comandos `create` e `list` que não existem em `mgc network rules`. Estes comandos podem existir em `mgc network security-groups rules` (ver próxima seção).

---

## 6. Grupo: security-groups

### mgc network security-groups

```
Operations related to Security Groups

Usage:
  mgc network security-groups [flags]
  mgc network security-groups [command]

Aliases:
  security-groups, security_groups

Commands:
  create      Create Security Group
  delete      Delete a security group
  get         Security Group Details
  list        List Security Groups by Tenant
  rules       rules
```

### ./tmp-cli/cli network security-groups

```
Network provides a client for interacting with the Magalu Cloud Network API.

Dqui1

Available Commands:

Other commands:
  create              Network provides a client for interacting with the Magalu Cloud Network API.
  delete              Network provides a client for interacting with the Magalu Cloud Network API.
  get                 Network provides a client for interacting with the Magalu Cloud Network API.
  list                Network provides a client for interacting with the Magalu Cloud Network API.
```

### ⚠️ Diferenças

| Aspecto | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **Bug visual** | Não | "Dqui1" aparece | 🔴 CRÍTICO |
| **Subgrupo `rules`** | ✅ Existe | ❌ Não existe | 🔴 CRÍTICO |
| **Aliases** | `security-groups`, `security_groups` | Nenhum | ❌ Faltando |

### 6.1. Subgrupo: security-groups rules

**mgc network security-groups rules**:

```
Operations related to Security Groups | rules

Usage:
  mgc network security-groups rules [flags]
  mgc network security-groups rules [command]

Commands:
  create      Create Rule
  list        List Rules
```

**./tmp-cli/cli**: ❌ **SUBGRUPO COMPLETAMENTE AUSENTE**

**⚠️ IMPACTO CRÍTICO**: O subgrupo `security-groups rules` e seus 2 comandos (`create`, `list`) estão completamente ausentes em `./tmp-cli/cli`. Isso representa uma perda de funcionalidade significativa.

---

## 7. Grupo: subnetpools

### mgc network subnetpools

```
Operations related to Subnet Pools

Usage:
  mgc network subnetpools [flags]
  mgc network subnetpools [command]

Commands:
  create             Create a Subnet Pool in a tenant
  create-book-cidr   Book Subnetpool
  create-unbook-cidr Unbook Subnetpool
  delete             Delete Subnet Pool by ID
  get                Get Subnet Pool by ID
  list               List Subnet Pools by Tenant
```

### ./tmp-cli/cli network subnet-pools

```
Network provides a client for interacting with the Magalu Cloud Network API.

Dqui1

Available Commands:

Other commands:
  book-c-i-d-r        Network provides a client for interacting with the Magalu Cloud Network API.
  create              Network provides a client for interacting with the Magalu Cloud Network API.
  delete              Network provides a client for interacting with the Magalu Cloud Network API.
  get                 Network provides a client for interacting with the Magalu Cloud Network API.
  list                Network provides a client for interacting with the Magalu Cloud Network API.
  unbook-c-i-d-r      Network provides a client for interacting with the Magalu Cloud Network API.
```

### ⚠️ Diferenças

| Aspecto | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **Nome do grupo** | `subnetpools` | `subnet-pools` | ❌ Divergente |
| **Bug visual** | Não | "Dqui1" aparece | 🔴 CRÍTICO |
| **Comando `create-book-cidr`** | `create-book-cidr` | `book-c-i-d-r` | ❌ Divergente |
| **Comando `create-unbook-cidr`** | `create-unbook-cidr` | `unbook-c-i-d-r` | ❌ Divergente |

---

## 8. Grupo: subnets

### mgc network subnets

```
Operations related to Subnets

Usage:
  mgc network subnets [flags]
  mgc network subnets [command]

Commands:
  delete      Delete a Subnet
  get         Subnet Details
  update      Update Subnet
```

### ./tmp-cli/cli network subnets

```
Network provides a client for interacting with the Magalu Cloud Network API.

Dqui1

Available Commands:

Other commands:
  delete              Network provides a client for interacting with the Magalu Cloud Network API.
  get                 Network provides a client for interacting with the Magalu Cloud Network API.
  update              Network provides a client for interacting with the Magalu Cloud Network API.
```

### ⚠️ Diferenças

| Aspecto | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **Bug visual** | Não | "Dqui1" aparece | 🔴 CRÍTICO |
| **Comandos** | Idênticos | Idênticos | ✅ OK |

---

## 9. Grupo: vpcs

### mgc network vpcs

```
Operations related to VPCs

Usage:
  mgc network vpcs [flags]
  mgc network vpcs [command]

Commands:
  create      Create a new Virtual Private Cloud (VPC)
  delete      Delete VPC
  get         VPC Details
  list        List VPC
  ports       ports
  public-ips  public-ips
  subnets     subnets
```

### ./tmp-cli/cli network v-p-cs

```
Network provides a client for interacting with the Magalu Cloud Network API.

Dqui1

Available Commands:

Other commands:
  create              Network provides a client for interacting with the Magalu Cloud Network API.
  create-port         Network provides a client for interacting with the Magalu Cloud Network API.
  create-public-i-p   Network provides a client for interacting with the Magalu Cloud Network API.
  create-subnet       Network provides a client for interacting with the Magalu Cloud Network API.
  delete              Network provides a client for interacting with the Magalu Cloud Network API.
  get                 Network provides a client for interacting with the Magalu Cloud Network API.
  list                Network provides a client for interacting with the Magalu Cloud Network API.
  list-ports          Network provides a client for interacting with the Magalu Cloud Network API.
  list-public-i-ps    Network provides a client for interacting with the Magalu Cloud Network API.
  list-subnets        Network provides a client for interacting with the Magalu Cloud Network API.
  rename              Network provides a client for interacting with the Magalu Cloud Network API.
```

### ⚠️ Diferenças ARQUITETURAIS CRÍTICAS

| Aspecto | mgc | ./tmp-cli/cli | Status |
|---------|-----|---------------|--------|
| **Nome do grupo** | `vpcs` | `v-p-cs` | 🔴 CRÍTICO |
| **Bug visual** | Não | "Dqui1" aparece | 🔴 CRÍTICO |
| **Arquitetura** | Subgrupos (`ports`, `public-ips`, `subnets`) | Comandos flat (`create-port`, `list-ports`, etc.) | 🔴 CRÍTICO |
| **Comando `rename`** | ❌ Não existe | ✅ Existe | ⚠️ Extra |

### 🔴 PROBLEMA ARQUITETURAL CRÍTICO

Em `mgc`, o grupo `vpcs` possui 3 **subgrupos**:
- `mgc network vpcs ports` (com comandos `create`, `list`)
- `mgc network vpcs public-ips` (com comandos `create`, `list`)
- `mgc network vpcs subnets` (com comandos `create`, `list`)

Em `./tmp-cli/cli`, esses subgrupos foram **eliminados** e convertidos em comandos "flat":
- `create-port`, `list-ports`
- `create-public-i-p`, `list-public-i-ps`
- `create-subnet`, `list-subnets`

**Comandos divergentes**:

| mgc | ./tmp-cli/cli | Status |
|-----|---------------|--------|
| `mgc network vpcs ports create` | `cli network v-p-cs create-port` | 🔴 Incompatível |
| `mgc network vpcs ports list` | `cli network v-p-cs list-ports` | 🔴 Incompatível |
| `mgc network vpcs public-ips create` | `cli network v-p-cs create-public-i-p` | 🔴 Incompatível |
| `mgc network vpcs public-ips list` | `cli network v-p-cs list-public-i-ps` | 🔴 Incompatível |
| `mgc network vpcs subnets create` | `cli network v-p-cs create-subnet` | 🔴 Incompatível |
| `mgc network vpcs subnets list` | `cli network v-p-cs list-subnets` | 🔴 Incompatível |

### 9.1. vpcs subnets create

**mgc**: `mgc network vpcs subnets create [vpc-id]`
- Flags: `--cidr-block` (required), `--description`, `--dns-nameservers`, `--ip-version` (required), `--name` (required), `--subnetpool-id`, `--vpc-id` (required)

**./tmp-cli/cli**: `cli network v-p-cs create-subnet [vpcID] [IPVersion] [Name] [CIDRBlock]`
- Flags: `--cidrblock` (required), `--description`, `--dnsnameservers`, `--ipversion` (required), `--name` (required), `--subnet-pool-id`, `--vpc-id` (required), `--zone`
- Bug visual: "doto3"

**Divergências de flags**:
- ❌ `--cidr-block` vs `--cidrblock`
- ❌ `--dns-nameservers` vs `--dnsnameservers`
- ❌ `--ip-version` vs `--ipversion`
- ❌ `--subnetpool-id` vs `--subnet-pool-id`
- ⚠️ Flag extra: `--zone` (não existe em mgc)

---

## 10. Problemas Sistemáticos Identificados

### 🔴 Bugs Visuais Críticos
1. **"Dqui1"** aparece em todos os grupos
2. **"doto3"** aparece em todos os comandos específicos
3. **Falta de espaços** entre nome do comando e descrição (ex: `attach-security-groupNetwork`)

### 🔴 Problemas de Nomenclatura
1. **Nomes de grupos**: `public-ips` → `public-i-ps`, `vpcs` → `v-p-cs`, `subnetpools` → `subnet-pools`
2. **Nomes de comandos**: `attach` → `attach-security-group`, `attach` → `attach-to-port`, `create-book-cidr` → `book-c-i-d-r`
3. **Nomes de flags**: Remoção sistemática de hífens (ex: `vpc-id` → `vpcid`, `cidr-block` → `cidrblock`, `ip-version` → `ipversion`, `dns-nameservers` → `dnsnameservers`)

### 🔴 Problemas de Conteúdo
1. **0% de descrições nas flags** em `./tmp-cli/cli`
2. **Descrições genéricas** para todos os comandos
3. **Aliases ausentes** em todos os grupos

### 🔴 Problemas Arquiteturais
1. **Subgrupo ausente**: `security-groups rules` completamente ausente
2. **Arquitetura divergente**: `vpcs` com subgrupos em `mgc` vs comandos flat em `./tmp-cli/cli`
3. **Comandos extras**: `network rules create` e `list` não existem em `mgc`
4. **Comando extra**: `vpcs rename` não existe em `mgc`

### 🔴 Problemas de Flags
1. **Flags globais ausentes**: `--cli.retry-until`, `-t/--cli.timeout`, `--env`, `--region`, `--server-url`, `-o/--output`, `--x-zone`
2. **Flags incorretamente marcadas como required**: `limit`, `offset`, `sort` em comandos `list`
3. **Nomenclatura divergente**: `items-per-page` e `page` vs `limit` e `offset`

---

## 11. Resumo de Incompatibilidades

| Categoria | Qtd. Problemas | Severidade |
|-----------|----------------|------------|
| **Bugs Visuais** | 3 tipos | 🔴 CRÍTICA |
| **Nomes Divergentes** | 15+ instâncias | 🔴 CRÍTICA |
| **Descrições Ausentes** | 100% das flags | 🔴 CRÍTICA |
| **Arquitetura Divergente** | 2 casos | 🔴 CRÍTICA |
| **Comandos Faltando** | 1 subgrupo (2 comandos) | 🔴 CRÍTICA |
| **Comandos Extras** | 3 comandos | ⚠️ MÉDIA |
| **Flags Globais Ausentes** | 7 flags | 🔴 CRÍTICA |
| **Aliases Ausentes** | 4 grupos | ⚠️ BAIXA |

---

## 12. Recomendações

### Prioridade CRÍTICA 🔴
1. **Eliminar bugs visuais**: Remover "Dqui1" e "doto3" completamente
2. **Corrigir nomes de grupos**: `public-i-ps` → `public-ips`, `v-p-cs` → `vpcs`, `subnet-pools` → `subnetpools`
3. **Restaurar arquitetura de subgrupos**: Implementar `security-groups rules` e `vpcs` com subgrupos
4. **Normalizar nomes de flags**: Manter hífens consistentemente (ex: `vpcid` → `vpc-id`, `cidrblock` → `cidr-block`)
5. **Adicionar descrições em 100% das flags**
6. **Restaurar flags globais ausentes**

### Prioridade ALTA ⚠️
1. **Corrigir nomes de comandos**: `attach-security-group` → `attach`, `attach-to-port` → `attach`
2. **Remover flags incorretamente marcadas como required**
3. **Adicionar aliases faltando**
4. **Verificar comandos extras**: Confirmar se `network rules create/list` e `vpcs rename` devem existir

### Prioridade MÉDIA
1. **Melhorar descrições dos comandos**: Substituir textos genéricos por descrições específicas
2. **Padronizar nomenclatura de paginação**: Decidir entre `items-per-page/page` ou `limit/offset`

---

## Conclusão

O CLI gerado (`./tmp-cli/cli network`) apresenta **problemas sistemáticos graves** em múltiplas dimensões:
- **Bugs visuais críticos** que comprometem a experiência do usuário
- **Divergências arquiteturais fundamentais** (subgrupos vs comandos flat)
- **Inconsistências massivas de nomenclatura** em grupos, comandos e flags
- **Perda de funcionalidade** (subgrupo `security-groups rules` ausente)
- **0% de documentação** nas flags

Estes problemas sugerem **falhas sistemáticas no processo de geração de código**, não apenas erros pontuais. É necessária uma **revisão completa do gerador de CLI** para corrigir estes padrões antes de gerar novos produtos.

