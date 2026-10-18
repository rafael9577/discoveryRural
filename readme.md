# Documentação da API – Discovery Rural

*Servidor em Go com banco de dados SQLite hospedado no Turso*

## 1. Visão geral

A API do Fauna Local é um servidor em Go que cadastra e consulta espécies da fauna. Os dados ficam em um banco SQLite hospedado no Turso (libSQL), acessado pela internet. Em desenvolvimento, o servidor responde em `http://localhost:8080`.

Além dos endpoints da API, a rota `/` serve uma página web simples de teste, com a listagem de espécies e um formulário de cadastro.

### Situação dos endpoints

| Nº | Endpoint | Situação |
| --- | --- | --- |
| 1 | `GET /api/animais/resumo` | Disponível |
| 2 | `GET /api/animal` | Disponível |
| 3 | `GET /api/animal/imagem` | Disponível |
| 4 | `GET /api/identificar` | **Em desenvolvimento** |
| 5 | `POST /api/animal/create` | Disponível (validação em desenvolvimento) |

## 2. Configuração

O servidor lê duas variáveis de ambiente e não inicia se alguma delas estiver ausente.

| Variável | Descrição |
| --- | --- |
| `TURSO_DATABASE_URL` | URL do banco no Turso (começa com `libsql://`). |
| `TURSO_AUTH_TOKEN` | Token de acesso ao banco. |

> **Atenção:** o token é uma senha de acesso ao banco. Nunca o coloque no código nem envie-o ao repositório. Mantenha-o em variáveis de ambiente ou em um arquivo `.env` listado no `.gitignore`.

## 3. Como executar

1. Instalar o Go e o Git.
2. Na pasta do back-end, baixar as dependências com `go mod tidy` (se a pasta ainda não tiver `go.mod`, criar antes com `go mod init <nome-do-modulo>`).
3. Definir as variáveis de ambiente, conforme o terminal:

   PowerShell (Windows):

   ```powershell
   $env:TURSO_DATABASE_URL = "libsql://..."
   $env:TURSO_AUTH_TOKEN = "seu-token"
   ```

   Bash (Linux, macOS, WSL):

   ```bash
   export TURSO_DATABASE_URL="libsql://..."
   export TURSO_AUTH_TOKEN="seu-token"
   ```

4. Executar o servidor com `go run .`
5. Acessar `http://localhost:8080`. Para encerrar o servidor, usar `Ctrl + C` no terminal.

## 4. Banco de dados

O banco `fauna` possui duas tabelas.

### Tabela `animais` (em uso)

| Coluna | Tipo | Descrição |
| --- | --- | --- |
| `id` | INTEGER, PK | Identificador automático. |
| `nome_popular` | TEXT, obrigatório | Nome popular da espécie. |
| `nome_cientifico` | TEXT | Nome científico. |
| `imagem` | BLOB | Foto em JPEG, corte 3:4 (ver endpoint 5). |
| `breve_descricao` | TEXT | Descrição resumida da espécie. |
| `estado_conservacao` | TEXT | Estado de conservação (ex.: Pouco preocupante). |
| `habito` | TEXT | Hábito (ex.: Diurno, Noturno). |
| `vocalizacao_postura` | TEXT | Vocalização e postura características. |
| `peconhento` | INTEGER, obrigatório | 1 se for peçonhento/venenoso; 0 se não for. |
| `aviso_seguranca` | TEXT | Orientações de segurança. |
| `protocolo_incidente` | TEXT | O que fazer em caso de incidente. |
| `contatos_resgate` | TEXT | Contatos para resgate. |
| `raridade` | INTEGER | Valor de 1 a 5 (restrição no banco). |
| `bitmap_id` | INTEGER, padrão 0 | Soma dos bits das características (ver seção 6). |

### Tabela `relatos` (reservada para uso futuro)

A tabela `relatos` (avistamentos com `animal_id`, `latitude`, `longitude`, `precisao`, `local_nome` e `data_hora`) existe no banco, mas está **reservada para uso futuro**: nenhum endpoint a utiliza por enquanto.

## 5. Endpoints

### 1. Listar Resumo dos Animais (Com Imagem)

Retorna uma lista contendo o ID, o nome popular e a URL da imagem de todas as espécies cadastradas (ideal para a listagem inicial e cards). Quando não há espécies cadastradas, a resposta é uma lista vazia (`[]`).

| | |
| --- | --- |
| **URL** | `/api/animais/resumo` |
| **Método** | `GET` |
| **Resposta** | `200 OK` (JSON) |

**Exemplo de resposta (200 OK):**

```json
[
  {
    "id": 1,
    "nome_popular": "Coruja-buraqueira",
    "imagem_url": "/api/animal/imagem?id=1"
  },
  {
    "id": 2,
    "nome_popular": "Cágado-de-barbicha",
    "imagem_url": "/api/animal/imagem?id=2"
  }
]
```

### 2. Obter Detalhes Completos de um Animal

Retorna todas as informações detalhadas, comportamentais, de segurança, raridade e bitmap de uma espécie específica. Campos de texto sem valor no banco são devolvidos como string vazia (`""`), e `peconhento` é devolvido como booleano.

| | |
| --- | --- |
| **URL** | `/api/animal` |
| **Método** | `GET` |
| **Parâmetro (query)** | `id` (inteiro, obrigatório): ID do animal desejado. |
| **Resposta** | `200 OK` (JSON) |

**Exemplo de requisição:** `GET /api/animal?id=1`

**Exemplo de resposta (200 OK):**

```json
{
  "id": 1,
  "nome_popular": "Coruja-buraqueira",
  "nome_cientifico": "Athene cunicularia",
  "breve_descricao": "Pequena coruja de hábitos terrestres.",
  "estado_conservacao": "Pouco preocupante",
  "habito": "Diurno/Crepuscular",
  "vocalizacao_postura": "Emite chamados agudos e balança a cabeça.",
  "peconhento": false,
  "aviso_seguranca": "Nunca tente manipular, alimentar ou se aproximar excessivamente de ninhos no chão.",
  "protocolo_incidente": "Manter distância segura e isolar a área.",
  "contatos_resgate": "193 / Setor Ambiental",
  "raridade": 2,
  "bitmap_id": 18,
  "imagem_url": "/api/animal/imagem?id=1"
}
```

### 3. Exibir Imagem do Animal

Endpoint binário que busca a imagem otimizada (corte 3:4, em torno de 16 KB) salva diretamente no banco de dados (BLOB).

| | |
| --- | --- |
| **URL** | `/api/animal/imagem` |
| **Método** | `GET` |
| **Parâmetro (query)** | `id` (inteiro): ID do animal. |
| **Header de resposta** | `Content-Type: image/jpeg` |

**Uso no front-end:**

```html
<img src="/api/animal/imagem?id=1" alt="Animal" />
```

### 4. Filtrar por Identificação (Bitmap ID)

> **EM DESENVOLVIMENTO.** Este endpoint ainda não está implementado no servidor. A descrição abaixo é a funcionalidade planejada e pode mudar até a implementação.

Endpoint a ser utilizado no fluxo de perguntas do aplicativo para filtrar as espécies compatíveis com as características selecionadas pelo usuário.

| | |
| --- | --- |
| **URL** | `/api/identificar` |
| **Método** | `GET` |
| **Parâmetro (query)** | `mask` (inteiro): valor da máscara resultante da soma dos bits correspondentes às respostas dadas pelo usuário (ver seção 6). |
| **Resposta prevista** | Array com os animais que casam com o filtro bit a bit: `(bitmap_id & mask) = mask`. |

**Exemplo de requisição:** `GET /api/identificar?mask=10`

### 5. Cadastrar Nova Espécie

Permite o cadastro de uma nova espécie enviando os dados via formulário multipart. A imagem enviada é processada e cortada automaticamente pelo back-end.

| | |
| --- | --- |
| **URL** | `/api/animal/create` |
| **Método** | `POST` |
| **Content-Type** | `multipart/form-data` |
| **Resposta de sucesso** | Redireciona para `/` (`303 See Other`), sem corpo JSON e sem devolver o ID do animal criado. |

> **VALIDAÇÃO EM DESENVOLVIMENTO.** Hoje o servidor não valida os campos enviados. A obrigatoriedade indicada abaixo é aplicada apenas pelo formulário HTML da rota `/`; a validação no back-end ainda será implementada.

#### Campos do formulário (Form Data)

| Campo | Tipo | Observação |
| --- | --- | --- |
| `nome_popular` | string | Obrigatório (no formulário). |
| `nome_cientifico` | string | |
| `breve_descricao` | string | Obrigatório (no formulário). |
| `estado_conservacao` | string | |
| `habito` | string | |
| `vocalizacao_postura` | string | |
| `peconhento` | checkbox / string | Enviar `"on"` se marcado. |
| `aviso_seguranca` | string | |
| `protocolo_incidente` | string | |
| `contatos_resgate` | string | |
| `raridade` | inteiro de 1 a 5 | O banco rejeita valores fora desse intervalo. |
| `bitmap_id` | inteiro | Obrigatório (no formulário). Ver seção 6. |
| `imagem` | arquivo binário | Foto da espécie. Obrigatória (no formulário). |

#### Comportamento atual

- **Validação:** o servidor não verifica campos obrigatórios. Uma requisição enviada diretamente à API, sem passar pelo formulário, é gravada mesmo com `nome_popular`, `breve_descricao`, `bitmap_id` ou `imagem` ausentes ou vazios.
- **Imagem:** o servidor aceita arquivos em JPEG ou PNG, recorta o centro da imagem na proporção 3:4, redimensiona para 450 x 600 pixels e a salva em JPEG. A qualidade começa em 85 e é reduzida em passos de 5 até o arquivo ficar com no máximo 16 KB; a redução para ao atingir a qualidade mínima (10), mesmo que o arquivo ainda ultrapasse esse tamanho.
- **Imagem ausente ou inválida:** se a imagem não for enviada ou não puder ser lida, o servidor não retorna erro. Ele grava uma imagem padrão (placeholder verde) no lugar.
- **Peçonhento:** o campo `peconhento` é gravado como verdadeiro apenas quando o valor enviado é `"on"`. Qualquer outro valor, ou a ausência do campo, é gravado como falso.
- **Raridade:** o banco só aceita valores de 1 a 5 no campo `raridade`.

## 6. Referência: Bitmap ID

O `bitmap_id` agrupa as características do animal em uma única soma baseada em potências de 2. Para cadastrar um animal, some apenas os valores das características que correspondem a ele. O mesmo esquema de bits forma a máscara (`mask`) do endpoint 4.

| Característica | Valor (bit) |
| --- | --- |
| No alto (árvores/telhados) | 1 |
| No chão (grama/asfalto) | 2 |
| Na água ou margem | 4 |
| É um pássaro / ave | 8 |
| É uma coruja | 16 |
| É um cágado / tartaruga | 32 |
| É um lagarto | 64 |
| É uma cobra / serpente | 128 |
| É um mamífero | 256 |

**Exemplo:** um pássaro (8) que vive no chão (2) tem `bitmap_id` = 8 + 2 = 10.

## 7. Tratamento de erros

> **EM DESENVOLVIMENTO.** O tratamento e a padronização das respostas de erro serão definidos em uma etapa futura e documentados nesta seção.

## 8. Segurança

> **Sem autenticação:** o endpoint `POST /api/animal/create` não exige autenticação. Quem tiver acesso à URL do servidor consegue cadastrar espécies no banco. Defina a proteção desse endpoint (por exemplo, uma chave de administrador) antes de publicar o servidor na internet.
