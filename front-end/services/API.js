import * as SecureStore from 'expo-secure-store';

const API_URL = 'https://discoveryrural.onrender.com/';
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
