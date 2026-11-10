
import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react-native';
import Toast from 'react-native-toast-message';

import Busca from '../src/screens/userPage/informacao';
import Login from '../src/screens/login';
import * as Api from '../services/API';

jest.mock('react-native-safe-area-context', () =>
  require('react-native-safe-area-context/jest/mock').default
);

const mockNavigation = { navigate: jest.fn(), replace: jest.fn(), goBack: jest.fn() };
jest.mock('@react-navigation/native', () => ({
  useNavigation: () => mockNavigation,
}));

jest.mock('react-native-toast-message', () => ({
  __esModule: true,
  default: { show: jest.fn() },
}));

jest.mock('../services/API', () => ({
  listarAnimais: jest.fn(),
  UserSingIn: jest.fn(),
}));

beforeEach(() => {
  Api.listarAnimais.mockReset();
  Api.UserSingIn.mockReset();
  Toast.show.mockClear();
  mockNavigation.replace.mockClear();
  jest.spyOn(console, 'error').mockImplementation(() => {});
  jest.spyOn(console, 'log').mockImplementation(() => {});
});
afterEach(() => jest.restoreAllMocks());

describe('CT19/CT20 - Catálogo (aba busca) sem conexão', () => {
  test('mostra mensagem amigável e o botão de tentar novamente, sem detalhe técnico', async () => {
    Api.listarAnimais.mockRejectedValueOnce(new Error('Network request failed'));

    render(<Busca />);

    expect(await screen.findByText('Não foi possível carregar os animais.')).toBeTruthy();
    expect(screen.getByText('Tentar Novamente')).toBeTruthy();
    expect(screen.queryByText(/Network request failed/)).toBeNull();
  });

  test('recuperação: com a conexão de volta, "Tentar Novamente" carrega a lista', async () => {
    Api.listarAnimais
      .mockRejectedValueOnce(new Error('Network request failed'))
      .mockResolvedValueOnce([{ id: 1, nome_popular: 'Coruja-buraqueira', imagem_url: 'https://x/1.jpg' }]);

    render(<Busca />);
    fireEvent.press(await screen.findByText('Tentar Novamente'));

    expect(await screen.findByText('Coruja-buraqueira')).toBeTruthy();
  });
});

describe('CT20 - Login com falha', () => {
  const preencher = () => {
    fireEvent.changeText(screen.getByPlaceholderText('Digitel seu login/Email'), 'a@b.com');
    fireEvent.changeText(screen.getByPlaceholderText('Digite sua Senha'), '12345678');
  };

  test('erro da API: a mensagem aparece no toast e a tela continua utilizável', async () => {
    Api.UserSingIn.mockRejectedValueOnce(new Error('Credenciais inválidas'));

    render(<Login />);
    preencher();
    fireEvent.press(screen.getByText('Login'));

    await waitFor(() =>
      expect(Toast.show).toHaveBeenCalledWith({ type: 'error', text1: 'Credenciais inválidas' })
    );
    expect(screen.getByText('Login')).toBeTruthy();
    expect(mockNavigation.replace).not.toHaveBeenCalled();
  });

  test('recuperação: depois de uma falha, nova tentativa com sucesso abre o app', async () => {
    Api.UserSingIn
      .mockRejectedValueOnce(new Error('Credenciais inválidas'))
      .mockResolvedValueOnce({ id: 1, nome_usuario: 'rafael' });

    render(<Login />);
    preencher();

    fireEvent.press(screen.getByText('Login'));
    await waitFor(() => expect(Toast.show).toHaveBeenCalledTimes(1));

    fireEvent.press(screen.getByText('Login'));
    await waitFor(() => expect(mockNavigation.replace).toHaveBeenCalledWith('tabsRouts'));
  });

  test('campos vazios: avisa o usuário e não chama a API', () => {
    render(<Login />);
    fireEvent.changeText(screen.getByPlaceholderText('Digitel seu login/Email'), '');
    fireEvent.changeText(screen.getByPlaceholderText('Digite sua Senha'), '');
    fireEvent.press(screen.getByText('Login'));

    expect(Toast.show).toHaveBeenCalledWith({ type: 'info', text1: 'Preencha todos os campos' });
    expect(Api.UserSingIn).not.toHaveBeenCalled();
  });

  test.failing('[BUG] sem rede, a mensagem do toast não deve ser técnica', async () => {
    Api.UserSingIn.mockRejectedValueOnce(new TypeError('Network request failed'));

    render(<Login />);
    preencher();
    fireEvent.press(screen.getByText('Login'));

    await waitFor(() => expect(Toast.show).toHaveBeenCalled());
    expect(Toast.show).toHaveBeenCalledWith({
      type: 'error',
      text1: expect.not.stringMatching(/network request failed/i),
    });
  });
});

