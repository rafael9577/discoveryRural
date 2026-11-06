import React, { useState } from 'react';
import { 
  View, 
  Text, 
  TouchableOpacity, 
  Switch, 
  //SafeAreaView 
} from 'react-native';

import MapView, { Marker, Heatmap, PROVIDER_GOOGLE } from 'react-native-maps';
import { Feather } from '@expo/vector-icons';
import { useNavigation } from '@react-navigation/native';
import { styles } from './styled';

import { SafeAreaView } from 'react-native-safe-area-context';

// placeholder da posição de animais e suas coordenadas, serve para o banco de dados
const ANIMAIS_MOCK = [
  { id: '1', nome: 'Cágado-de-barbicha', latitude: -8.0142, longitude: -34.9511 },
  { id: '2', nome: 'Bicho-Preguiça-comum', latitude: -8.0150, longitude: -34.9490 },
];

export default () => {
  const navigation = useNavigation();
  
  // controle da interface
  const [menuAberto, setMenuAberto] = useState(false); // controla se o menu lateral está visível ou não
  const [mapaDeCalorAtivo, setMapaDeCalorAtivo] = useState(false); // alterna entre os marcadores e o heatmap
  
  // estado do tipo dicionario/objeto para controlar quais animais estão selecionados no filtro
  const [animaisSelecionados, setAnimaisSelecionados] = useState({
    'Cágado-de-barbicha': true,
    'Bicho-Preguiça-comum': false,
  });

  // função que ativa ao clicar em um animal no menu, inverte o valor true/false do animal clicado
  const toggleAnimal = (nomeAnimal) => {
    setAnimaisSelecionados(prev => ({
      ...prev,
      [nomeAnimal]: !prev[nomeAnimal]
    }));
  };

  // cria uma nova lista contendo apenas os animais que estão marcados como true no estado "animaisSelecionados"
  const dadosFiltrados = ANIMAIS_MOCK.filter(animal => animaisSelecionados[animal.nome]);

  // formata os dados filtrados para o formato especifico que o componente <Heatmap> exige latitude, longitude, weight(o peso do negocio, 1 vale 1)
  const pontosHeatmap = dadosFiltrados.map(animal => ({
    latitude: animal.latitude,
    longitude: animal.longitude,
    weight: 1 // Peso padrão: 1 registro = peso 1
  }));

  // função auxiliar que retorna o design de uma linha do menu de filtro, texto + quadradinho do checkbox
  const renderCheckboxItem = (nomeAnimal) => (
    <TouchableOpacity 
      style={styles.checkboxItemContainer} 
      onPress={() => toggleAnimal(nomeAnimal)} // oo tocar na linha inteira, marca/desmarca o animal
      activeOpacity={0.7}
    >
      <Text style={styles.checkboxItemText}>{nomeAnimal}</Text>
      
      {/* caixa visual do Checkbox */}
      <View style={styles.checkboxBox}>
        {animaisSelecionados[nomeAnimal] && (
          <Feather name="check" size={16} color="#000" /> // mostra o ícone de "V" se estiver selecionado
        )}
      </View>
    </TouchableOpacity>
  );

  return (
    // SafeAreaView garante que o app não fique escondido atrás do entalhe da câmera (notch) em iPhones
    <SafeAreaView style={styles.container}>
      
      {/* HEADER/CABEÇALHO */}
      <View style={styles.header}>
        
        {/* botão de Voltar */}
        <TouchableOpacity style={styles.iconButton} onPress={() => navigation.goBack()}>
          <Feather name="arrow-left" size={28} color="#000" />
        </TouchableOpacity>

        {/* nome central da tela */}
        <View style={styles.headerTitleBox}>
          <Text style={styles.headerTitleText}>MAPA RURAL</Text>
        </View>

        {/* botão de abrir/fechar o fenu de filtros (o botão hamburguer) */}
        <TouchableOpacity 
          style={[styles.menuButton, menuAberto && styles.menuButtonAtivo]} 
          onPress={() => setMenuAberto(!menuAberto)}
        >
          <Feather 
            name="menu" 
            size={28} 
            color={menuAberto ? "#FFF" : "#000"} 
          />
        </TouchableOpacity>
      </View>

      {/* AREA PRINCIPAL DO MAPA */}
      <View style={styles.mapContainer}>
        <MapView
          provider={PROVIDER_GOOGLE} // força o uso do motor do Google maps
          style={styles.map}
          initialRegion={{
            latitude: -8.0142, // coordenada central inicial da UFRPE (placeholder)
            longitude: -34.9511, // (placeholder)
            latitudeDelta: 0.015, // define o nível de zoom vertical (placeholder)
            longitudeDelta: 0.012, // define o nível de zoom horizontal (placeholder)
          }}
        >
          {/* renderização condicional pt1: se o botão de calor estiver DESLIGADO, desenha os pinos vermelhos incompleto)*/} 
          {!mapaDeCalorAtivo && dadosFiltrados.map(animal => (
            <Marker
              key={animal.id}
              coordinate={{ latitude: animal.latitude, longitude: animal.longitude }}
              title={animal.nome}
            />
          ))}

          {/* renderização condicional pt2: se o botão de calor estiver LIGADO e houver dados, desenha a mancha colorida (incompleto) */}
          {mapaDeCalorAtivo && pontosHeatmap.length > 0 && (
            <Heatmap
              points={pontosHeatmap}
              radius={40} // tamanho da mancha de cada ponto (placeholder)
              opacity={0.7} // transparencia da camada de calor (placeholder)
              gradient={{
                colors: ["#00FF00", "#FF0000"], // do verde (poucos registros) ao vermelho (muitos registros)
                startPoints: [0.2, 1],
                colorMapSize: 256
              }}
            />
          )}
        </MapView>

        {/* MENU DE SOBREPOSIÇÃO (os filtros) */}
        {menuAberto && (
          <View style={styles.overlayMenu}>
            
            {/* header interno do menu (controlador do Heatmap) */}
            <View style={styles.overlayHeader}>
              <Text style={styles.overlayHeaderText}>MAPA DE CALOR (On/Off)</Text>
              
              {/* tnterruptor (toggle) do mapa de calor */}
              <Switch
                trackColor={{ false: "#767577", true: "#81b0ff" }}
                thumbColor={mapaDeCalorAtivo ? "#FFF" : "#f4f3f4"}
                ios_backgroundColor="#3e3e3e"
                onValueChange={() => setMapaDeCalorAtivo(previousState => !previousState)}
                value={mapaDeCalorAtivo}
              />
            </View>

            {/* corpo do menu com as opções de filtragem de animais */}
            <View style={styles.overlayBody}>
              <View style={styles.filterTitleContainer}>
                <Text style={styles.filterTitleText}>SELECIONE O ANIMAL</Text>
              </View>
              
              {/* itens dinamicos criados pela função auxiliar */}
              {renderCheckboxItem('Cágado-de-barbicha')}
              <View style={styles.separator} /> {/* Linha divisória fina */}
              
              {renderCheckboxItem('Bicho-Preguiça-comum')}
              <View style={styles.separator} />
              
              {/* itens desabilitados (placeholders para simular futuras opções visuais) */}
              <TouchableOpacity style={styles.checkboxItemContainer} disabled>
                <Text style={styles.checkboxItemText}></Text>
                <View style={styles.checkboxBox} />
              </TouchableOpacity>
              <View style={styles.separator} />
              
              <TouchableOpacity style={styles.checkboxItemContainer} disabled>
                <Text style={styles.checkboxItemText}></Text>
                <View style={styles.checkboxBox} />
              </TouchableOpacity>
              <View style={styles.separator} />
              
               <TouchableOpacity style={styles.checkboxItemContainer} disabled>
                <Text style={styles.checkboxItemText}></Text>
                <View style={styles.checkboxBox} />
              </TouchableOpacity>
            </View>
            
          </View>
        )}
      </View>
    </SafeAreaView>
  );
}