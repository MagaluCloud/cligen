# Módulo de Autenticação OAuth

Este módulo implementa o fluxo de autenticação OAuth 2.0 com PKCE (Proof Key for Code Exchange) para a Magalu Cloud CLI.

## Arquitetura

O módulo foi estruturado seguindo princípios de separação de responsabilidades e facilidade de manutenção:

```
auth/
├── auth.go          # Entry point do módulo
├── login.go         # Comandos CLI para login
├── config.go        # Configurações centralizadas
├── types.go         # Tipos e interfaces
├── pkce.go          # Implementação PKCE (RFC 7636)
├── client.go        # Cliente OAuth HTTP
├── server.go        # Servidor de callback HTTP
├── service.go       # Orquestração do fluxo de autenticação
└── html.template    # Template HTML para páginas de callback
```

## Componentes

### Config (`config.go`)
Gerencia todas as configurações necessárias para autenticação:
- URLs do OAuth (autorização e token)
- Client ID e scopes
- Configuração do servidor de callback
- Links externos (termos e privacidade)

### Types (`types.go`)
Define todas as estruturas de dados e interfaces:
- `TokenResponse`: Resposta do servidor OAuth
- `AuthResult`: Resultado de tentativa de autenticação
- `LoginOptions`: Opções configuráveis para login
- `AuthService`: Interface para serviços de autenticação

### PKCE (`pkce.go`)
Implementação do PKCE conforme RFC 7636:
- Geração segura de code verifier
- Criação de code challenge usando SHA256
- Encoding Base64 URL-safe

### OAuth Client (`client.go`)
Responsável por comunicação HTTP com o servidor OAuth:
- Construção de URLs de autenticação
- Troca de código de autorização por tokens
- Tratamento de erros HTTP

### Callback Server (`server.go`)
Servidor HTTP temporário para receber callbacks OAuth:
- Gerenciamento de lifecycle do servidor
- Tratamento de rotas (`/callback`, `/term`, `/privacy`)
- Renderização de templates HTML
- Shutdown gracioso

### Service (`service.go`)
Orquestra todo o fluxo de autenticação:
- Integra todos os componentes
- Implementa diferentes métodos de login (browser, headless, QR code)
- Gerencia estado e ciclo de vida

### Login Commands (`login.go`)
Interface CLI usando Cobra:
- Comando `login` com flags configuráveis
- Validação de flags mutuamente exclusivas
- Integração com serviço de autenticação

## Fluxo de Autenticação

1. **Inicialização**
   - Usuário executa `mgc auth login`
   - Sistema cria configuração padrão
   - Gera code verifier PKCE

2. **Autorização**
   - Constrói URL de autorização com PKCE
   - Inicia servidor de callback local
   - Abre navegador com URL de autenticação

3. **Callback**
   - Usuário autentica no navegador
   - Sistema recebe código de autorização
   - Exibe página de sucesso

4. **Token Exchange**
   - Troca código por tokens de acesso
   - Valida resposta do servidor
   - Retorna tokens ao usuário

5. **Finalização**
   - Encerra servidor de callback
   - Exibe mensagem de sucesso
   - Retorna controle ao CLI

## Uso

### Login Padrão (Browser)
```bash
mgc auth login
```

### Login Headless (Device Flow)
```bash
mgc auth login --headless
```

### Login com QR Code
```bash
mgc auth login --qrcode
```

### Exibir Token
```bash
mgc auth login --show
```

## Configuração

### Variáveis de Ambiente

- `MGC_LISTEN_ADDRESS`: Define o endereço do servidor de callback (padrão: `127.0.0.1:8095`)

### Personalização

Para personalizar a configuração, modifique a função `DefaultConfig()` em `config.go`:

```go
config := DefaultConfig()
config.ListenAddr = "localhost:9000"
config.Scopes = append(config.Scopes, "custom.scope")
```

## Segurança

### PKCE
O módulo implementa PKCE (RFC 7636) para proteger contra ataques de interceptação de código:
- Code verifier gerado usando `crypto/rand`
- Code challenge usando SHA256
- Base64 URL-safe encoding

### HTTPS
- Comunicação com servidor OAuth sempre via HTTPS
- Certificados validados pelo cliente HTTP padrão

### Callback Server
- Servidor local escuta apenas em localhost
- Shutdown automático após receber callback
- Timeout configurável para prevenir servidores órfãos

## Tratamento de Erros

O módulo usa wrapping de erros para contexto adequado:

```go
if err != nil {
    return fmt.Errorf("failed to build auth URL: %w", err)
}
```

Todos os erros são propagados até o comando CLI, onde são exibidos ao usuário.

## Extensibilidade

### Adicionar Novo Método de Login

1. Adicionar opção em `LoginOptions` (`types.go`)
2. Implementar método em `Service` (`service.go`)
3. Adicionar flag em `NewLoginCommand` (`login.go`)

### Personalizar Template HTML

Edite `html.template` para modificar páginas de callback.

### Adicionar Armazenamento de Tokens

Implemente lógica de salvamento em `runLogin()` (`login.go`):

```go
// TODO: Salvar token em local seguro
if err := saveToken(token); err != nil {
    return fmt.Errorf("failed to save token: %w", err)
}
```

## Testes

Para testar o módulo localmente:

```bash
# Compilar
go build

# Executar
./cli auth login --show
```

## Referências

- [RFC 6749 - OAuth 2.0](https://datatracker.ietf.org/doc/html/rfc6749)
- [RFC 7636 - PKCE](https://datatracker.ietf.org/doc/html/rfc7636)
- [Magalu Cloud OAuth](https://id.magalu.com)

## Melhorias Futuras

- [ ] Implementar device flow (login headless)
- [ ] Implementar login com QR code
- [ ] Adicionar armazenamento seguro de tokens (keyring)
- [ ] Implementar refresh de tokens
- [ ] Adicionar testes unitários
- [ ] Adicionar testes de integração
- [ ] Suporte a múltiplos tenants
- [ ] Cache de sessões

