/**
 *
 * Critérios de aceitação cobertos (o que a tela já entrega do RF02):
 *   CA1. Quando o animal tem orientação de conduta, a tela mostra o bloco de alerta com
 *        "Aviso de Segurança", "Protocolo" e "Resgate".
 *   CA2. Só aparecem os itens de conduta que existem para o animal.
 *   CA3. Se o animal não tem nenhuma orientação de conduta, o bloco de alerta não é exibido.
 *   CA4. Enquanto os detalhes carregam, a tela já mostra o resumo recebido da lista.
 *   CA5. Se a API falha e não há resumo, a tela mostra a mensagem de erro e "Tentar novamente".
 *   CA6. "Tentar novamente" busca os dados de novo e exibe o animal.
 *  
 */

import { render, screen, fireEvent } from '@testing-library/react-native';

import DetalheAnimal from '../src/screens/userPage/detalheAnimal';
import { detalheAnimal } from '../services/API';

jest.mock('../services/API', () => ({
  detalheAnimal: jest.fn(),
}));

jest.mock('react-native-safe-area-context', () => {
  const { View } = require('react-native');
  return { SafeAreaView: View };
});

jest.mock('@expo/vector-icons', () => {
  const React = require('react');
  const { Text } = require('react-native');
  const Icone = ({ name }) => React.createElement(Text, null, `icon:${name}`);
  return { Feather: Icone, Ionicons: Icone };
});

const preguica = {
  id: 7,
  nome_popular: 'Bicho-preguiça',
  nome_cientifico: 'Bradypus variegatus',
  breve_descricao: 'Mamífero arborícola de movimentos lentos.',
  habito: 'Arborícola, vive nas copas das árvores',
  estado_conservacao: 'Pouco preocupante',
  raridade: 3,
  aviso_seguranca: 'Não toque nem tente pegar o animal.',
  protocolo_incidente: 'Em caso de arranhão, lave com água e sabão e procure atendimento.',
  contatos_resgate: 'CIPOMA/BPMA: (81) 3181-1700',
  imagem_url: 'https://exemplo.com/preguica.jpg',
};

const navigation = { goBack: jest.fn() };
const telaComRota = (params) => <DetalheAnimal navigation={navigation} route={{ params }} />;

beforeEach(() => {
  detalheAnimal.mockReset();
  jest.spyOn(console, 'error').mockImplementation(() => {});
});

afterEach(() => {
  jest.restoreAllMocks();
});

describe('RF02 - Guia de Conduta e Habitat', () => {

  // TC-01 | CA1
  test('TC-01: mostra aviso de segurança, protocolo e resgate quando existem', async () => {
    detalheAnimal.mockResolvedValue(preguica);

    render(telaComRota({ id: 7 }));

    expect(await screen.findByText(/Aviso de Segurança/)).toBeTruthy();
    expect(screen.getByText('Não toque nem tente pegar o animal.')).toBeTruthy();
    expect(screen.getByText(/Protocolo/)).toBeTruthy();
    expect(
      screen.getByText('Em caso de arranhão, lave com água e sabão e procure atendimento.')
    ).toBeTruthy();
    expect(screen.getByText(/Resgate/)).toBeTruthy();
    expect(screen.getByText('CIPOMA/BPMA: (81) 3181-1700')).toBeTruthy();
    expect(screen.getByText('icon:warning')).toBeTruthy();
  });

  // TC-02 | CA2
  test('TC-02: mostra apenas os itens de conduta que o animal possui', async () => {
    detalheAnimal.mockResolvedValue({
      ...preguica,
      protocolo_incidente: null,
      contatos_resgate: null,
    });

    render(telaComRota({ id: 7 }));

    expect(await screen.findByText(/Aviso de Segurança/)).toBeTruthy();
    expect(screen.queryByText(/Protocolo/)).toBeNull();
    expect(screen.queryByText(/Resgate/)).toBeNull();
  });

  // TC-03 | CA3
  test('TC-03: sem nenhuma orientação de conduta, o bloco de alerta não aparece', async () => {
    detalheAnimal.mockResolvedValue({
      ...preguica,
      aviso_seguranca: '',
      protocolo_incidente: '',
      contatos_resgate: '',
    });

    render(telaComRota({ id: 7 }));

    expect(await screen.findByText('Bicho-preguiça')).toBeTruthy();
    expect(screen.queryByText(/Aviso de Segurança/)).toBeNull();
    expect(screen.queryByText(/Protocolo/)).toBeNull();
    expect(screen.queryByText(/Resgate/)).toBeNull();
    expect(screen.queryByText('icon:warning')).toBeNull();
  });

  // TC-04 | CA4
  test('TC-04: enquanto carrega, já exibe o resumo recebido da lista', () => {
    detalheAnimal.mockReturnValue(new Promise(() => {})); // nunca resolve

    render(telaComRota({ id: 7, resumo: preguica }));

    expect(screen.getByText('Bicho-preguiça')).toBeTruthy();
    expect(screen.getByText('Arborícola, vive nas copas das árvores')).toBeTruthy();
  });

  test('TC-05: se a API falha mas há resumo, o resumo continua na tela (sem mensagem de erro)', async () => {
    detalheAnimal.mockRejectedValue(new Error('Erro de conexão'));

    render(telaComRota({ id: 7, resumo: preguica }));

    expect(await screen.findByText('Bicho-preguiça')).toBeTruthy();
    expect(screen.queryByText('Não foi possível carregar os detalhes.')).toBeNull();
  });

  // TC-06| CA5
  test('TC-06: se a API falha e não há resumo, mostra erro e "Tentar novamente"', async () => {
    detalheAnimal.mockRejectedValue(new Error('Erro de conexão'));

    render(telaComRota({ id: 7 }));

    expect(await screen.findByText('Não foi possível carregar os detalhes.')).toBeTruthy();
    expect(screen.getByText('Tentar novamente')).toBeTruthy();
  });

  // TC-7 | CA6
  test('TC-07: "Tentar novamente" busca os dados de novo e exibe o animal', async () => {
    detalheAnimal
      .mockRejectedValueOnce(new Error('Erro de conexão'))
      .mockResolvedValueOnce(preguica);

    render(telaComRota({ id: 7 }));
    fireEvent.press(await screen.findByText('Tentar novamente'));

    expect(await screen.findByText('Bicho-preguiça')).toBeTruthy();
    expect(screen.queryByText('Não foi possível carregar os detalhes.')).toBeNull();
    expect(detalheAnimal).toHaveBeenCalledTimes(2);
  });
});
