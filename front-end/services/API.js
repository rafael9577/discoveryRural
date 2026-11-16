import * as SecureStore from 'expo-secure-store';

// Sem barra no final: os caminhos já começam com "/"
const API_URL = 'https://discoveryrural.onrender.com';
const USERS_URL = API_URL;
const FAUNA_URL = API_URL;

async function chamar(base, caminho, { method = 'GET', body, auth = false } = {}) {
  const headers = { 'Content-Type': 'application/json' };
  if (auth) {
    const token = await SecureStore.getItemAsync('token');
    if (token) headers.Authorization = `Bearer ${token}`;
  }
  const res = await fetch(`${base}${caminho}`, {
    method,
    headers,
    body: body ? JSON.stringify(body) : undefined,
  });
  const texto = await res.text();
  let data;
  try {
    data = JSON.parse(texto);
  } catch {
    data = { erro: texto };
  }
  if (!res.ok) throw new Error(data.erro || 'Erro de conexão');
  return data;
}

export async function cadastrar({ nome_usuario, cpf, email, senha }) {
  const data = await chamar(USERS_URL, '/api/users/register', {
    method: 'POST',
    body: { nome_usuario, cpf, email, senha },
  });
  await SecureStore.setItemAsync('token', data.token);
  return data.usuario;
}

export async function UserSingIn(identificador, senha) {
  const data = await chamar(USERS_URL, '/api/users/login', {
    method: 'POST',
    body: { identificador, senha },
  });
  await SecureStore.setItemAsync('token', data.token);
  return data.usuario;
}

export async function usuarioAtual() {
  const token = await SecureStore.getItemAsync('token');
  if (!token) return null;
  try {
    return await chamar(USERS_URL, '/api/users/me', { auth: true });
  } catch {
    await SecureStore.deleteItemAsync('token');
    return null;
  }
}

export async function sair() {
  try {
    await chamar(USERS_URL, '/api/users/logout', { method: 'POST', auth: true });
  } finally {
    await SecureStore.deleteItemAsync('token');
  }
}

function comUrlCompleta(animal) {
  return { ...animal, imagem_url: `${FAUNA_URL}${animal.imagem_url}` };
}

export async function listarAnimais() {
  const lista = await chamar(FAUNA_URL, '/api/animais/resumo');
  return lista.map(comUrlCompleta);
}

export async function detalheAnimal(id) {
  const animal = await chamar(FAUNA_URL, `/api/animal?id=${id}`);
  return comUrlCompleta(animal);
}

// ---------------------------------------------------------------------------
// Avistamentos
// ---------------------------------------------------------------------------

// GET /api/avistamentos devolve [{ animal_id, avistamentos: [...] }]. O back-end envia cada ponto
// como [longitude, latitude]; o guia da API descreve { long, lat }. Aceita os dois formatos e
// devolve uma lista plana: [{ id, animal_id, latitude, longitude }].
export function normalizarAvistamentos(grupos) {
  const lista = [];
  (Array.isArray(grupos) ? grupos : []).forEach((grupo) => {
    (grupo?.avistamentos || []).forEach((ponto, i) => {
      const longitude = Array.isArray(ponto) ? ponto[0] : ponto?.long ?? ponto?.longitude;
      const latitude = Array.isArray(ponto) ? ponto[1] : ponto?.lat ?? ponto?.latitude;
      if (Number.isFinite(latitude) && Number.isFinite(longitude)) {
        lista.push({ id: `${grupo.animal_id}-${i}`, animal_id: grupo.animal_id, latitude, longitude });
      }
    });
  });
  return lista;
}

// Todos os pontos do banco (ou só os de um animal, se animalId for informado)
export async function listarAvistamentos(animalId) {
  const filtro = animalId ? `?animal_id=${animalId}` : '';
  return normalizarAvistamentos(await chamar(FAUNA_URL, `/api/avistamentos${filtro}`));
}

// POST /api/avistamento. Limite do servidor: 1 avistamento a cada 10 segundos por IP (erro 429).
// local_nome aceita até 120 caracteres.
export async function registrarAvistamento({ animal_id, latitude, longitude, precisao, local_nome }) {
    return chamar(FAUNA_URL, '/api/avistamento', {
    method: 'POST',
    body: { animal_id, latitude, longitude, precisao, local_nome },
  });
}
