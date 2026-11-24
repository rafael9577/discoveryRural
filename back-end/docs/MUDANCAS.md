# O que mudou na refatoração

## Contrato da API: o que continua igual
Todas as rotas, métodos, nomes de campos JSON, status de sucesso e mensagens de erro de
validação foram mantidos. O front existente continua funcionando.

## Pequenas diferenças de comportamento (intencionais)

| Antes | Agora |
|---|---|
| Erros de `/api/animais/resumo`, `/api/animal`, `/api/animal/create` e `/api/animal/imagem` vinham em **texto puro** | Todos os erros são `{"erro": "..."}` |
| `/api/animal/create` devolvia o texto cru do erro do banco em falhas 500 | Mensagem genérica; o detalhe vai só para o log |
| Cadastro de animal não validava nada (raridade inválida virava 500 do banco) | Mesmas regras da edição: nome obrigatório, limites de tamanho, raridade 1–5, `bitmap_id` ≥ 0 → 400 |
| Cadastro só entendia `peconhento=on` | Aceita também `true` e `1` |
| Cadastro sem limite de tamanho | Máximo de 12 MB (413) |
| Cadastro respondia só com redirect 303 | Com `Accept: application/json` responde `201 {"id": n}` (o redirect continua para a página) |
| `/api/animal` com id não numérico → 404 | → 400 ("id inválido ou ausente") |
| `/api/animais/resumo`, `/api/animal`, `/api/animal/imagem` aceitavam qualquer método | Só `GET` (405 nos demais) |
| Qualquer rota desconhecida devolvia a página HTML (200) | Só `/` devolve a página; o resto é 404 JSON |
| Falha de leitura de linha era ignorada (lista incompleta) | Vira erro 500 |
| Login com falha de banco respondia 401 | 500 (401 só para credenciais erradas) |
| Foto com origem não zero no `Bounds()` era recortada errado | Corrigido |
| Mapa do limitador de IPs crescia sem limite | Chaves ociosas são removidas |

## Pontos de atenção que NÃO foram alterados (decisão sua)
- **Limite de 1000 pontos é global**, não por animal: em `/api/animais` e `/api/avistamentos`,
  se houver mais de 1000 relatos no total, os animais de ID maior podem aparecer sem pontos.
- **IP atrás do proxy:** com `RENDER=true` usa-se o primeiro valor de `X-Forwarded-For`, que
  pode ser forjado pelo cliente; isso enfraquece os limites de taxa.
- **Sem CORS:** um front web em outro domínio será bloqueado pelo navegador.
- A chave de admin no campo `admin_key` do formulário (usado pela página `/`) trafega no corpo
  da requisição; use sempre HTTPS.
