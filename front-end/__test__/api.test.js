o
jest.mock('expo-secure-store', () => ({
  getItemAsync: jest.fn(),
  setItemAsync: jest.fn(),
  deleteItemAsync: jest.fn(),
}));

import * as SecureStore from 'expo-secure-store';
import { UserSingIn, listarAnimais } from '../services/API';

const resposta = (ok, status, corpo) => ({
  ok,
  status,
  text: async () => (typeof corpo === 'string' ? corpo : JSON.stringify(corpo)),
});

beforeEach(() => {
  global.fetch = jest.fn();
  jest.clearAllMocks();
});

describe('CT20 - Falhas de comunicação com a API', () => {
  test('servidor indisponível: a chamada falha e nenhum token é salvo', async () => {
    global.fetch.mockRejectedValueOnce(new TypeError('Network request failed'));

    await expect(UserSingIn('a@b.com', '12345678')).rejects.toThrow('Network request failed');
    expect(SecureStore.setItemAsync).not.toHaveBeenCalled();
  });

  test('erro da API: a mensagem do servidor é repassada à interface', async () => {
    global.fetch.mockResolvedValueOnce(resposta(false, 401, { erro: 'Credenciais inválidas' }));

    await expect(UserSingIn('a@b.com', 'errada123')).rejects.toThrow('Credenciais inválidas');
  });

  test('resposta que não é JSON (ex.: 502 do Render) também vira erro', async () => {
    global.fetch.mockResolvedValueOnce(resposta(false, 502, 'Bad Gateway'));

    await expect(UserSingIn('a@b.com', '12345678')).rejects.toThrow('Bad Gateway');
  });

  test('recuperação: depois de uma falha, nova tentativa funciona e salva o token', async () => {
    global.fetch
      .mockRejectedValueOnce(new TypeError('Network request failed'))
      .mockResolvedValueOnce(
        resposta(true, 200, { token: 'abc123', usuario: { id: 1, nome_usuario: 'rafael' } })
      );

    await expect(UserSingIn('a@b.com', '12345678')).rejects.toThrow();
    const usuario = await UserSingIn('a@b.com', '12345678');

    expect(usuario.nome_usuario).toBe('rafael');
    expect(SecureStore.setItemAsync).toHaveBeenCalledWith('token', 'abc123');
  });
});

describe('Listagem de animais', () => {
  test('completa a URL das imagens', async () => {
    global.fetch.mockResolvedValueOnce(
      resposta(true, 200, [{ id: 1, nome_popular: 'Coruja-buraqueira', imagem_url: '/api/animal/imagem?id=1' }])
    );

    const lista = await listarAnimais();

    expect(lista).toHaveLength(1);
    expect(lista[0].imagem_url).toContain('/api/animal/imagem?id=1');
  });

  
  test.failing('[BUG] a URL chamada não deve ter barra dupla', async () => {
    global.fetch.mockResolvedValueOnce(resposta(true, 200, []));

    await listarAnimais();

    const url = global.fetch.mock.calls[0][0];
    expect(url.replace('https://', '')).not.toContain('//');
  });
});
