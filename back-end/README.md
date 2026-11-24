# discoveryRural — back-end (Go)

API REST para cadastro de animais, avistamentos e usuários. Banco: Turso (libsql).

## Como rodar

```bash
cp .env.example .env     # preencha TURSO_DATABASE_URL, TURSO_AUTH_TOKEN e ADMIN_KEY
go run ./cmd/servidor
go test ./...
```

As tabelas são criadas automaticamente na primeira execução (e duas espécies
padrão são cadastradas se a tabela de animais estiver vazia).

Documentação para quem consome a API: [`docs/API.md`](docs/API.md).

## Estrutura

```
cmd/servidor/main.go        raiz de composição: cria e liga todos os objetos
internal/
  config/                   leitura das variáveis de ambiente
  banco/                    conexão com o Turso e esquema das tabelas
  falha/                    erro de domínio tipado (validação, 404, 401, 409, 429, 500)
  httpx/                    JSON, middlewares, chave de admin, IP do cliente
  limitador/                limitador de taxa por IP (janela deslizante)
  animal/                   modelo, validação, imagem, repositório, serviço, handler
  avistamento/              modelo, repositório, serviço, handler
  usuario/                  modelo, CPF, senha/token, repositório, serviço, handler
  web/                      página administrativa embutida (pagina.html) e 404 em JSON
docs/API.md                 contrato da API para o front
```

## Arquitetura (cada módulo em 4 camadas)

| Camada | Responsabilidade | Conhece |
|---|---|---|
| `handler.go` | HTTP: lê a requisição, chama o serviço, escreve a resposta | serviço, `httpx` |
| `servico.go` | Regras de negócio e validação | interface do repositório |
| `repositorio.go` | SQL e acesso ao banco (interface + implementação `...SQL`) | `database/sql` |
| `modelo.go` | Entidades e DTOs | nada |

Princípios aplicados:

- **Responsabilidade única**: nenhum arquivo mistura HTTP, regra de negócio e SQL.
- **Inversão de dependência**: serviços dependem de *interfaces* (`Repositorio`,
  `HasherSenha`, `ProcessadorImagem`, `VerificadorAnimal`, `ProvedorPontos`),
  declaradas por quem as consome. Por isso os testes usam dublês, sem banco.
- **Injeção de dependências por construtor** (`NovoServico(...)`, `NovoHandler(...)`);
  só o `main.go` conhece as implementações concretas. Não há variáveis globais.
- **Encapsulamento**: campos privados, comportamento em métodos (`CPF.Valido()`,
  `Admin.Validar()`, `Limitador.Permitir()`).
- **Tratamento de erro único**: serviços devolvem `*falha.Falha`; `httpx.EscreverFalha`
  traduz para o status HTTP. Erros internos vão para o log, nunca para o cliente.
- **Uma conexão com o banco** compartilhada (antes eram três `sql.Open`).
- Servidor com timeouts e desligamento gracioso (SIGINT/SIGTERM).
