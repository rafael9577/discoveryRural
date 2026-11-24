# API discoveryRural — guia para o front

## Convenções gerais

- **URL base:** a do servidor (ex.: `https://SEU-SERVICO.onrender.com` ou `http://localhost:8080`).
  Abaixo, os caminhos são relativos a ela.
- **Formato:** JSON em UTF-8 (`Content-Type: application/json`), exceto o cadastro de animal
  (`multipart/form-data`) e a imagem (`image/jpeg`).
- **Erros:** sempre `{"erro": "mensagem em português"}` com o status HTTP correspondente.
  A mensagem pode ser exibida ao usuário.
- **Dois tipos de credencial:**
  - **Usuário:** cabeçalho `Authorization: Bearer <token>` (token recebido no cadastro/login, vale 30 dias).
  - **Administrador:** cabeçalho `X-Admin-Key: <chave>`. É segredo do servidor: **nunca** coloque
    em app/site público; use só em ferramentas internas.
- **CORS:** a API **não** envia cabeçalhos CORS. Apps nativos/mobile funcionam normalmente; um front web em
  outro domínio será bloqueado pelo navegador até que o back-end adicione CORS (peça essa alteração).
- **Imagens:** `imagem_url` é um caminho relativo — monte `URL_BASE + imagem_url`.

### Status HTTP

| Status | Quando |
|---|---|
| 200 / 201 | Sucesso / criado |
| 400 | Dado inválido ou ausente (veja `erro`) |
| 401 | Token ausente/inválido/expirado, credenciais erradas ou chave de admin errada |
| 404 | Animal, avistamento, usuário ou rota inexistente |
| 405 | Método HTTP não permitido nessa rota |
| 409 | Conflito (e-mail, CPF ou nome de usuário já cadastrado) |
| 413 | Cadastro de animal maior que 12 MB |
| 429 | Muitas requisições (ver limites abaixo) |
| 500 | Erro interno (o detalhe fica só no log do servidor) |

### Limites de taxa (por IP)

| Ação | Limite |
|---|---|
| `POST /api/avistamento` | 1 a cada 10 s |
| `POST /api/users/register` e `POST /api/users/login` (somados) | 10 a cada 5 min |

Ao receber **429**, mostre a mensagem e tente de novo depois.

---

## Animais

### `GET /api/animais` — lista completa com avistamentos

```json
[
  {
    "id": 1,
    "nome_popular": "Coruja-buraqueira",
    "nome_cientifico": "Athene cunicularia",
    "breve_descricao": "Pequena coruja de hábitos terrestres.",
    "estado_conservacao": "Pouco preocupante",
    "habito": "Diurno/Crepuscular",
    "vocalizacao_postura": "Emite chamados agudos e balança a cabeça.",
    "peconhento": false,
    "aviso_seguranca": "Nunca tente manipular...",
    "protocolo_incidente": "Manter distância segura e isolar a área.",
    "contatos_resgate": "193 / Setor Ambiental",
    "raridade": 2,
    "bitmap_id": 18,
    "imagem_url": "/api/animal/imagem?id=1",
    "avistamentos": [ { "longitude": -35.0, "latitude": -8.0 } ]
  }
]
```

`avistamentos` é sempre um array (vazio se não houver). Atenção: aqui cada ponto é um **objeto**
`{longitude, latitude}`; em `GET /api/avistamentos` é um **array** `[longitude, latitude]`.
No máximo 1000 pontos no total são devolvidos (ordenados por animal e id).

### `GET /api/animais/ids` — só os IDs
`[1, 2, 3]`

### `GET /api/animais/resumo` — para listas/cartões
```json
[ { "id": 1, "nome_popular": "Coruja-buraqueira", "breve_descricao": "…", "imagem_url": "/api/animal/imagem?id=1" } ]
```

### `GET /api/animal?id=1` — detalhe de um animal
Mesmos campos de `/api/animais`, **sem** `avistamentos`.
Erros: `400` (id ausente/inválido), `404` ("Animal não encontrado").

### `GET /api/animal/imagem?id=1` — foto
Responde `image/jpeg` (450×600, até ~16 KB). Use direto em `<img src>` / componente de imagem.
`404` se não existir.

### `POST /api/animal/create` — cadastrar espécie *(admin)*
`multipart/form-data`. Chave de admin no cabeçalho `X-Admin-Key` (ou no campo `admin_key`, usado pela página `/`).

| Campo | Regra |
|---|---|
| `nome_popular` | obrigatório, até 100 caracteres |
| `nome_cientifico`, `estado_conservacao`, `habito` | opcionais, até 100 |
| `breve_descricao`, `vocalizacao_postura`, `aviso_seguranca`, `protocolo_incidente` | opcionais, até 500 |
| `contatos_resgate` | opcional, até 300 |
| `peconhento` | `on`, `true` ou `1` = sim; ausente = não |
| `raridade` | obrigatório, inteiro de 1 a 5 |
| `bitmap_id` | inteiro ≥ 0 (padrão 0) — veja a tabela de bits abaixo |
| `imagem` | arquivo JPEG ou PNG. Sem arquivo (ou arquivo inválido) usa uma imagem verde padrão |

Resposta de sucesso:
- Com cabeçalho `Accept: application/json` → `201` `{"id": 3}` (**recomendado para apps**).
- Sem ele → `303` redirecionando para `/` (comportamento da página administrativa).

Erros: `401`, `400` (mensagem do campo), `413`.

```js
const form = new FormData();
form.append("nome_popular", "Sapo-cururu");
form.append("raridade", "2");
form.append("bitmap_id", "10");
form.append("imagem", arquivo);
const r = await fetch(`${BASE}/api/animal/create`, {
  method: "POST",
  headers: { "X-Admin-Key": chave, Accept: "application/json" },
  body: form,
});
const { id } = await r.json();
```

### `PATCH /api/animal/update?id=1` (ou `PUT`) — editar *(admin)*
Atualização **parcial**: envie só os campos que mudam (JSON). Aceita:
`nome_popular, nome_cientifico, breve_descricao, estado_conservacao, habito, vocalizacao_postura,
peconhento (bool), aviso_seguranca, protocolo_incidente, contatos_resgate, raridade (1–5), bitmap_id (≥0)`
— mesmos limites de tamanho do cadastro. Não altera a imagem.

```js
await fetch(`${BASE}/api/animal/update?id=1`, {
  method: "PATCH",
  headers: { "Content-Type": "application/json", "X-Admin-Key": chave },
  body: JSON.stringify({ habito: "Noturno", raridade: 3 }),
});
```
Resposta `200`: o animal completo e atualizado (mesmo formato de `/api/animais`, com `avistamentos`).
Erros: `401`, `400` (id inválido, valor inválido, ou "nenhum campo para atualizar"), `404`.

### Tabela do `bitmap_id` (some os valores que se aplicam)

| Característica | Valor |
|---|---|
| No alto (árvores/telhados) | 1 |
| No chão (grama/asfalto) | 2 |
| Na água ou margem | 4 |
| Pássaro / ave | 8 |
| Coruja | 16 |
| Cágado / tartaruga | 32 |
| Lagarto | 64 |
| Cobra / serpente | 128 |
| Mamífero | 256 |

Exemplo: ave (8) que vive no chão (2) → `10`. Para testar um bit no front: `(a.bitmap_id & 8) !== 0`.

---

## Avistamentos

### `POST /api/avistamento` — registrar
```json
{ "animal_id": 1, "latitude": -8.0173, "longitude": -34.9521, "precisao": 12.5, "local_nome": "Praça central" }
```
`animal_id`, `latitude` (−90…90) e `longitude` (−180…180) são obrigatórios. `precisao` (metros, 0…100000) e
`local_nome` (cortado em 120 caracteres) são opcionais.

Resposta `201`: `{"id": 57, "animal_id": 1}`.
Erros: `400`, `404` ("Animal não encontrado"), `429` (1 a cada 10 s por IP).

### `GET /api/avistamento?id=57` — um avistamento
```json
{ "id": 57, "animal_id": 1, "latitude": -8.0173, "longitude": -34.9521, "precisao": 12.5, "local_nome": "Praça central", "data_hora": "2026-11-21 15:04:05" }
```
`precisao` é `null` quando não informada. `data_hora` em UTC. Erros: `400`, `404`.

### `GET /api/avistamentos` — pontos agrupados por animal (para mapas)
Opcional: `?animal_id=1`.
```json
[ { "animal_id": 1, "avistamentos": [[-34.9521, -8.0173], [-35.0, -8.0]] }, { "animal_id": 2, "avistamentos": [] } ]
```
Cada ponto é **`[longitude, latitude]`** (ordem GeoJSON; cuidado ao usar em bibliotecas que esperam lat/lng).
Animais sem relatos aparecem com lista vazia. Máximo de 1000 linhas. Erros: `400`, `404` (com `animal_id`).

---

## Usuários

O objeto `usuario` é sempre:
```json
{ "id": 1, "nome_usuario": "felipe_01", "cpf": "***.982.247-**", "email": "felipe@exemplo.com", "data_criacao": "2026-11-21 15:04:05", "reputacao": 5 }
```
O CPF volta sempre mascarado. `reputacao` vai de 0 a 10 (começa em 5).

### `POST /api/users/register`
```json
{ "nome_usuario": "felipe_01", "cpf": "529.982.247-25", "email": "felipe@exemplo.com", "senha": "minimo8caracteres" }
```
Regras: nome de usuário 3–20 caracteres (letras, números e `_`); CPF válido (com ou sem pontuação);
e-mail válido (guardado em minúsculas); senha de 8 a 72 caracteres.
Resposta `201`: `{"token": "…", "usuario": {…}}`.
Erros: `400` (regra violada), `409` ("Email já cadastrado", "CPF já cadastrado", "Nome de usuário já em uso"), `429`.

### `POST /api/users/login`
```json
{ "identificador": "felipe_01", "senha": "minimo8caracteres" }
```
`identificador` pode ser o nome de usuário **ou** o e-mail (se tiver `@`, é tratado como e-mail).
Resposta `200`: `{"token": "…", "usuario": {…}}`. Erros: `401` ("Credenciais inválidas"), `429`.

### `GET /api/users/me` — usuário logado
Cabeçalho `Authorization: Bearer <token>`. Resposta `200`: o objeto `usuario`. `401` se o token for inválido ou expirou.

### `POST /api/users/logout`
Cabeçalho `Authorization: Bearer <token>`. Resposta `200`: `{"ok": true}`. O token deixa de valer.

### `POST /api/users/reputacao` *(admin)*
```json
{ "id": 1, "delta": -0.5 }
```
Soma `delta` à reputação (resultado limitado entre 0 e 10). Resposta `200`: `{"id": 1, "reputacao": 4.5}`.
Erros: `401`, `404` ("Usuário não encontrado").

### Fluxo de sessão recomendado
1. Cadastro/login → guarde `token` em armazenamento seguro (Keychain/Keystore; no web, evite `localStorage` se possível).
2. Envie `Authorization: Bearer <token>` nas rotas de usuário.
3. Ao receber `401` em `/api/users/me`, apague o token e peça login novamente.
4. No logout, chame `/api/users/logout` e apague o token localmente.

---

## Página administrativa
`GET /` serve uma página HTML simples para listar e cadastrar espécies (usa os mesmos endpoints acima).
Qualquer rota inexistente responde `404` com `{"erro": "Rota não encontrada"}`.
