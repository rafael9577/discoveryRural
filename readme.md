# Discovery Rural – Guia da API para o Front-end

Este documento explica como o app (Expo / React Native) conversa com o servidor. Ele cobre os animais, os avistamentos (para o mapa de calor) e os usuários.

## 1. Conceitos básicos

| Item | Valor |
| --- | --- |
| Endereço base | Variável `EXPO_PUBLIC_API_URL` (ex.: `https://SEU-SERVICO.onrender.com`). Em testes locais, `http://IP-DO-COMPUTADOR:8080`. Nunca use `localhost` no celular. |
| Formato | JSON (as imagens são JPEG). |
| Imagens | A API devolve só o caminho (`/api/animal/imagem?id=1`). O app precisa juntar o endereço base: `API_URL + imagem_url`. |
| Coordenadas | Sempre `long` (longitude) e `lat` (latitude), em graus. |
| Textos sem valor | Chegam como `""` (texto vazio), nunca como `null`. |
| Autenticação | Só as rotas de usuário logado usam o cabeçalho `Authorization: Bearer <token>`. Animais e avistamentos são públicos. |
| Rotas de administrador | Cadastrar e atualizar animal exigem a chave de administrador (`X-Admin-Key`). **Não use no app.** |
| Primeira chamada lenta | O servidor gratuito do Render dorme após ~15 min sem uso. A primeira chamada pode levar cerca de 1 minuto. Use um indicador de carregamento. |

## 2. Formato de erro

- **Rotas de usuário, avistamento e animais padronizados:** JSON, com o campo `erro`.
  ```json
  { "erro": "Animal não encontrado" }
  ```
- **Rotas antigas de animal** (`/api/animal`, `/api/animais/resumo`, `/api/animal/imagem`, `/api/animal/create`) e de relato (`/api/relato`, `/api/relatos`): o erro vem em **texto simples**. No `api.js`, a função `chamar` já trata os dois casos.

| Status | Significado |
| --- | --- |
| 200 / 201 | Deu certo (201 = criado) |
| 400 | Dado inválido ou faltando |
| 401 | Sem login, token inválido ou chave de administrador errada |
| 404 | Animal ou avistamento não encontrado |
| 405 | Método HTTP não permitido nessa rota |
| 429 | Muitas tentativas seguidas; espere alguns segundos |
| 500 | Erro no servidor |

## 3. Objeto Animal (formato padrão)

É o formato da lista completa (`GET /api/animais`) e do detalhe (`GET /api/animal?id=N`).

| Campo | Tipo | Descrição |
| --- | --- | --- |
| `id` | número | Identificador do animal |
| `nome_popular` | texto | Nome popular |
| `nome_cientifico` | texto | Nome científico |
| `breve_descricao` | texto | Descrição resumida |
| `estado_conservacao` | texto | Ex.: "Pouco preocupante" |
| `habito` | texto | Ex.: "Diurno", "Noturno" |
| `vocalizacao_postura` | texto | Som e postura |
| `peconhento` | booleano | `true` se for peçonhento |
| `aviso_seguranca` | texto | Orientação de segurança |
| `protocolo_incidente` | texto | O que fazer em caso de incidente |
| `contatos_resgate` | texto | Contatos para resgate |
| `raridade` | número | De 1 a 5 |
| `bitmap_id` | número | Soma de bits das características (usada na identificação) |
| `imagem_url` | texto | Caminho da foto (juntar com o endereço base) |
| `avistamentos` | lista | Lista de `{ "long": número, "lat": número }`; vazia se não houver |

Exemplo:

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
  "imagem_url": "/api/animal/imagem?id=1",
  "avistamentos": [
    { "long": -34.9486, "lat": -8.0178 },
    { "long": -34.9503, "lat": -8.0152 }
  ]
}
```

## 4. Animais

### 4.1 Ids dos animais

`GET /api/animais/ids` → `200`

```json
[1, 2, 3]
```

### 4.2 Lista resumida (cards)

`GET /api/animais/resumo` → `200`

```json
[
  { "id": 1, "nome_popular": "Coruja-buraqueira", "breve_descricao": "Pequena coruja de hábitos terrestres.", "imagem_url": "/api/animal/imagem?id=1" }
]
```

### 4.3 Lista completa (objeto padrão, com avistamentos)

`GET /api/animais` → `200`: lista de objetos Animal (seção 3).

### 4.4 Detalhe de um animal

`GET /api/animal?id=N` → `200`: um objeto Animal (seção 3). `400` se faltar o `id` e `404` se não existir.

### 4.5 Foto do animal

`GET /api/animal/imagem?id=N` → `200`, `Content-Type: image/jpeg`. No app:

```jsx
<Image source={{ uri: API_URL + animal.imagem_url }} style={{ width: 120, height: 160 }} />
```

### 4.6 Cadastrar e atualizar (administrador)

Não fazem parte do app. Existem para a equipe e para a página de teste do servidor (`/`).

- `POST /api/animal/create`: formulário (`multipart/form-data`) com a chave de administrador.
- `PATCH /api/animal/update?id=N`: corpo JSON com **apenas os campos que mudam**, cabeçalho `X-Admin-Key`. Valida: `nome_popular` não pode ficar vazio, `raridade` de 1 a 5, `bitmap_id` maior ou igual a 0 e limites de tamanho dos textos. Devolve o animal atualizado (objeto padrão).

```json
{ "estado_conservacao": "Quase ameaçada", "raridade": 3 }
```

## 5. Avistamentos

### 5.1 Lista de avistamentos por animal (mapa de calor)

`GET /api/avistamentos` → `200`

```json
[
  { "animal_id": 1, "avistamentos": [ { "long": -34.9486, "lat": -8.0178 }, { "long": -34.9503, "lat": -8.0152 } ] },
  { "animal_id": 2, "avistamentos": [] }
]
```

- Todo animal aparece, mesmo sem avistamentos (lista vazia).
- Filtro: `GET /api/avistamentos?animal_id=N` (`400` se o id não for número, `404` se o animal não existir).
- Devolve no máximo 5.000 pontos por chamada.

Para o mapa de calor do `react-native-maps` (`Heatmap`), converta cada ponto:

```js
const pontos = grupos.flatMap(g => g.avistamentos)
  .map(p => ({ latitude: p.lat, longitude: p.long, weight: 1 }));
```

### 5.2 Buscar um avistamento

`GET /api/avistamento?id=N` → `200`

```json
{
  "id": 7,
  "animal_id": 1,
  "latitude": -8.0178,
  "longitude": -34.9486,
  "precisao": 10,
  "local_nome": "Perto da Reitoria",
  "data_hora": "2026-11-09T03:24:06Z"
}
```

`precisao` pode ser `null`. `400` se faltar o `id` e `404` se não existir.

### 5.3 Registrar um avistamento

`POST /api/avistamento` com `Content-Type: application/json` → `201`

```json
{ "animal_id": 1, "latitude": -8.0178, "longitude": -34.9486, "precisao": 10, "local_nome": "Perto da Reitoria" }
```

Resposta: `{ "id": 7, "animal_id": 1 }`.

| Campo | Obrigatório | Regra |
| --- | --- | --- |
| `animal_id` | sim | Número maior que 0, de um animal que existe (`404` se não existir) |
| `latitude` | sim | De -90 a 90 |
| `longitude` | sim | De -180 a 180 |
| `precisao` | não | De 0 a 100000 (metros) |
| `local_nome` | não | Até 120 caracteres |

Limite: **1 avistamento a cada 10 segundos por IP** (`429` se passar). Mostre uma mensagem amigável.

### 5.4 Rotas antigas (continuam funcionando)

`POST /api/relato` e `GET /api/relatos` (lista plana com latitude e longitude). Prefira as rotas novas, porque as antigas devem sair no futuro.

## 6. Usuários

| Rota | Método | O que faz |
| --- | --- | --- |
| `/api/users/register` | POST | Cadastra: `nome_usuario`, `cpf`, `email`, `senha`. Devolve `{ token, usuario }` |
| `/api/users/login` | POST | Entra: `identificador` (nome de usuário ou e-mail) e `senha`. Devolve `{ token, usuario }` |
| `/api/users/me` | GET | Perfil de quem está logado (com `Authorization: Bearer <token>`) |
| `/api/users/logout` | POST | Encerra a sessão (com `Authorization`) |

Regras do cadastro: nome de usuário com 3 a 20 caracteres (letras, números e `_`), CPF válido, e-mail válido e senha de 8 a 72 caracteres. O CPF volta sempre mascarado (`***.887.014-**`). Cadastro e login juntos aceitam 10 tentativas a cada 5 minutos por IP.

O `api.js` já guarda o token com o `expo-secure-store` e envia o cabeçalho `Authorization` nas chamadas com `auth: true`.

## 7. Funções sugeridas para o `api.js`

Seguem o estilo atual do arquivo (funções nomeadas, usando `chamar` e `FAUNA_URL`):

```js
export async function idsDosAnimais() {
  return chamar(FAUNA_URL, '/api/animais/ids');
}

export async function listarAnimaisCompletos() {
  const lista = await chamar(FAUNA_URL, '/api/animais');
  return lista.map(comUrlCompleta);
}

export async function listarAvistamentos(animalId) {
  const filtro = animalId ? `?animal_id=${animalId}` : '';
  return chamar(FAUNA_URL, `/api/avistamentos${filtro}`);
}

export async function buscarAvistamento(id) {
  return chamar(FAUNA_URL, `/api/avistamento?id=${id}`);
}

export async function registrarAvistamento({ animal_id, latitude, longitude, precisao, local_nome }) {
  return chamar(FAUNA_URL, '/api/avistamento', {
    method: 'POST',
    body: { animal_id, latitude, longitude, precisao, local_nome },
  });
}
```

## 8. Dados de exemplo (para o "avistamento mocado")

Enquanto a tela não chama a API, use este mesmo formato, para trocar depois sem refazer a tela:

```js
export const IDS_MOCK = [1, 2];

export const AVISTAMENTOS_MOCK = [
  { animal_id: 1, avistamentos: [{ long: -34.9486, lat: -8.0178 }, { long: -34.9503, lat: -8.0152 }] },
  { animal_id: 2, avistamentos: [] },
];
```

## 9. O que ainda não existe na API

- **Pontos de saúde (RF06):** não há tabela nem rota. Por enquanto, a lista com nome e coordenadas fica no próprio app.
- **Foto no avistamento:** a API grava a posição, a precisão e o nome do local, mas não recebe foto.
- **Privacidade:** as coordenadas voltam exatas e sem login. Se houver espécies raras ou perigosas, avise o back-end para arredondar.
