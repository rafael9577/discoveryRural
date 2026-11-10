
import React from 'react';
import { Linking } from 'react-native';
import { render, screen, fireEvent } from '@testing-library/react-native';

import Emergencia from '../src/screens/userPage/emergencia';

jest.mock('@react-navigation/native', () => ({
  useNavigation: () => ({ navigate: jest.fn(), goBack: jest.fn() }),
}));

describe('CT17 - Funções de emergência', () => {
  let openURL;

  beforeEach(() => {
    openURL = jest.spyOn(Linking, 'openURL').mockResolvedValue(true);
  });
  afterEach(() => jest.restoreAllMocks());

  test('o protocolo de primeiros socorros aparece ao abrir a tela, sem toques extras', () => {
    render(<Emergencia />);
    expect(screen.getByText('Guia de Emergência')).toBeTruthy();
    expect(screen.getByText('Primeiros Socorros e Protocolos')).toBeTruthy();
    expect(screen.getByText('Em caso de picada ou mordida')).toBeTruthy();
    expect(screen.getByText('O que fazer:')).toBeTruthy();
    expect(screen.getByText('O que NÃO fazer:')).toBeTruthy();
  });

  test('os contatos de resgate estão na mesma tela', () => {
    render(<Emergencia />);
    expect(screen.getByText('Contatos de Emergência (PE)')).toBeTruthy();
    expect(screen.getByText('CEATOX-PE (Assistência Toxicológica)')).toBeTruthy();
    expect(screen.getByText('CIPOMA / BPMA (Polícia Ambiental)')).toBeTruthy();
    expect(screen.getByText('Corpo de Bombeiros / SAMU')).toBeTruthy();
  });

  test.each([
    ['0800 722 6001', 'tel:08007226001'],
    ['193', 'tel:193'],
    ['192', 'tel:192'],
  ])('tocar em "%s" abre o discador com %s', (texto, esperado) => {
    render(<Emergencia />);
    fireEvent.press(screen.getByText(texto));
    expect(openURL).toHaveBeenCalledWith(esperado);
  });

  test('CIPOMA: hoje o toque liga para 190', () => {
    render(<Emergencia />);
    fireEvent.press(screen.getByText('(81) 3181-1700 / 190'));
    expect(openURL).toHaveBeenCalledWith('tel:190');
  });

 
  test.failing('[BUG] CIPOMA deveria discar o número exibido (81) 3181-1700', () => {
    render(<Emergencia />);
    fireEvent.press(screen.getByText('(81) 3181-1700 / 190'));
    expect(openURL).toHaveBeenCalledWith('tel:8131811700');
  });
});

describe('CT19 - Emergência sem conexão', () => {
  const fetchOriginal = global.fetch;
  afterEach(() => {
    global.fetch = fetchOriginal;
  });

  test('o conteúdo aparece completo sem nenhuma chamada de rede', () => {
    const semRede = jest.fn(() => Promise.reject(new Error('Network request failed')));
    global.fetch = semRede;

    render(<Emergencia />);

    expect(screen.getByText('Em caso de picada ou mordida')).toBeTruthy();
    expect(screen.getByText('Ao avistar um animal silvestre')).toBeTruthy();
    expect(screen.getByText('0800 722 6001')).toBeTruthy();
    expect(semRede).not.toHaveBeenCalled();
  });
});
