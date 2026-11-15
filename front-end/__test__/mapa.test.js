/**
 * Tela testada: src/screens/userPage/mapa/index.js
 *
 * Critérios de aceitação cobertos:
 *   CA1. Ao abrir a tela, o mapa de calor está desligado e são exibidos os marcadores
 *        apenas dos animais selecionados no filtro (inicialmente, só o Cágado-de-barbicha).
 *   CA2. O botão de menu abre e fecha o painel de filtros, que contém o interruptor
 *        "MAPA DE CALOR (On/Off)" e a lista de animais.
 *   CA3. Marcar ou desmarcar um animal no filtro mostra ou esconde o marcador dele.
 *   CA4. Ligar o mapa de calor troca os marcadores por uma camada de calor, com um ponto de
 *        peso 1 (e a coordenada do animal) para cada animal selecionado.
 *   CA5. Sem nenhum animal selecionado, a camada de calor não é desenhada.
 *   CA6. Desligar o mapa de calor devolve os marcadores.
 *   CA7. A caixa de seleção mostra o "check" apenas nos animais selecionados.
 *
 */

import { Switch } from 'react-native';
import { render, screen, fireEvent } from '@testing-library/react-native';

import Mapa from '../src/screens/userPage/mapa';

const mockNavigate = jest.fn();
const mockHeatmapProps = jest.fn();

jest.mock('@react-navigation/native', () => ({
  useNavigation: () => ({ navigate: mockNavigate }),
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

jest.mock('react-native-maps', () => {
  const React = require('react');
  const { View, Text } = require('react-native');
  const MapView = ({ children }) => React.createElement(View, null, children);
  const Marker = ({ title }) => React.createElement(Text, null, `marker:${title}`);
  const Heatmap = (props) => {
    mockHeatmapProps(props);
    return React.createElement(Text, null, `heatmap:${props.points.length}`);
  };
  return { __esModule: true, default: MapView, Marker, Heatmap, PROVIDER_GOOGLE: 'google' };
});

const CAGADO = 'Cágado-de-barbicha';
const PREGUICA = 'Bicho-Preguiça-comum';

const abrirMenu = () => fireEvent.press(screen.getByText('icon:menu'));
const ligarMapaDeCalor = () => fireEvent(screen.UNSAFE_getByType(Switch), 'valueChange', true);

beforeEach(() => {
  mockHeatmapProps.mockClear();
  mockNavigate.mockClear();
});

describe('RF04 - Mapa de Calor para Pesquisadores', () => {
  // TC-01 | CA1
  test('TC-01: ao abrir, mostra só o marcador do animal selecionado e nenhuma camada de calor', () => {
    render(<Mapa />);

    expect(screen.getByText('MAPA RURAL')).toBeTruthy();
    expect(screen.getByText(`marker:${CAGADO}`)).toBeTruthy();
    expect(screen.queryByText(`marker:${PREGUICA}`)).toBeNull();
    expect(screen.queryByText(/^heatmap:/)).toBeNull();
    expect(screen.queryByText('MAPA DE CALOR (On/Off)')).toBeNull(); // menu fechado
  });

  // TC-02 | CA2
  test('TC-02: o botão de menu abre e fecha o painel de filtros', () => {
    render(<Mapa />);

    abrirMenu();
    expect(screen.getByText('MAPA DE CALOR (On/Off)')).toBeTruthy();
    expect(screen.getByText('SELECIONE O ANIMAL')).toBeTruthy();
    expect(screen.getByText(CAGADO)).toBeTruthy();
    expect(screen.getByText(PREGUICA)).toBeTruthy();

    abrirMenu();
    expect(screen.queryByText('MAPA DE CALOR (On/Off)')).toBeNull();
    expect(screen.queryByText('SELECIONE O ANIMAL')).toBeNull();
  });

  // TC-03 | CA3
  test('TC-03: marcar a Preguiça no filtro faz o marcador dela aparecer', () => {
    render(<Mapa />);
    abrirMenu();

    fireEvent.press(screen.getByText(PREGUICA));

    expect(screen.getByText(`marker:${PREGUICA}`)).toBeTruthy();
    expect(screen.getByText(`marker:${CAGADO}`)).toBeTruthy();
  });

  test('TC-04: desmarcar o Cágado no filtro faz o marcador dele sumir', () => {
    render(<Mapa />);
    abrirMenu();

    fireEvent.press(screen.getByText(CAGADO));

    expect(screen.queryByText(`marker:${CAGADO}`)).toBeNull();
  });

  // TC-05 | CA4
  test('TC-05: ligar o mapa de calor troca os marcadores pela camada de calor', () => {
    render(<Mapa />);
    abrirMenu();

    ligarMapaDeCalor();

    expect(screen.queryByText(`marker:${CAGADO}`)).toBeNull();
    expect(screen.getByText('heatmap:1')).toBeTruthy();
  });

  test('TC-06: a camada de calor recebe um ponto de peso 1 com a coordenada do animal', () => {
    render(<Mapa />);
    abrirMenu();

    ligarMapaDeCalor();

    expect(mockHeatmapProps).toHaveBeenLastCalledWith(
      expect.objectContaining({
        points: [{ latitude: -8.0142, longitude: -34.9511, weight: 1 }],
      })
    );
  });

  test('TC-07: com o mapa de calor ligado, marcar outro animal adiciona outro ponto', () => {
    render(<Mapa />);
    abrirMenu();
    ligarMapaDeCalor();

    fireEvent.press(screen.getByText(PREGUICA));

    expect(screen.getByText('heatmap:2')).toBeTruthy();
    expect(mockHeatmapProps).toHaveBeenLastCalledWith(
      expect.objectContaining({
        points: [
          { latitude: -8.0142, longitude: -34.9511, weight: 1 },
          { latitude: -8.015, longitude: -34.949, weight: 1 },
        ],
      })
    );
  });

  // TC-08 | CA5
  test('TC-08: com o mapa de calor ligado e nenhum animal selecionado, nada é desenhado', () => {
    render(<Mapa />);
    abrirMenu();
    ligarMapaDeCalor();

    fireEvent.press(screen.getByText(CAGADO)); // desmarca o único selecionado

    expect(screen.queryByText(/^heatmap:/)).toBeNull();
    expect(screen.queryByText(/^marker:/)).toBeNull();
  });

  // TC-09 | CA6
  test('TC-09: desligar o mapa de calor devolve os marcadores', () => {
    render(<Mapa />);
    abrirMenu();
    ligarMapaDeCalor();
    expect(screen.queryByText(`marker:${CAGADO}`)).toBeNull();

    fireEvent(screen.UNSAFE_getByType(Switch), 'valueChange', false);

    expect(screen.queryByText(/^heatmap:/)).toBeNull();
    expect(screen.getByText(`marker:${CAGADO}`)).toBeTruthy();
  });

  // TC-10 | CA7
  test('TC-10: o "check" aparece só nos animais selecionados', () => {
    render(<Mapa />);
    abrirMenu();
    expect(screen.getAllByText('icon:check')).toHaveLength(1); // só o Cágado

    fireEvent.press(screen.getByText(PREGUICA));
    expect(screen.getAllByText('icon:check')).toHaveLength(2);

    fireEvent.press(screen.getByText(CAGADO));
    expect(screen.getAllByText('icon:check')).toHaveLength(1);
  });
});
