import React, { useState, useCallback, useRef } from 'react';
import { View, Text, TouchableOpacity, Switch, ScrollView, ActivityIndicator } from 'react-native';

import MapView, { Marker, Heatmap, PROVIDER_GOOGLE } from 'react-native-maps';
import { Feather } from '@expo/vector-icons';
import { useNavigation, useFocusEffect } from '@react-navigation/native';
import { styles } from './styled';
import { SafeAreaView } from 'react-native-safe-area-context';
import { listarAnimais, listarAvistamentos } from '../../../../services/API';

const HOSPITAIS_MOCK = [
  { id: 'h1', nome: 'Hospital Preencha 1', latitude: -8.029605, longitude: -34.957707 },
];

export default () => {
  const navigation = useNavigation();
  
  const avistamento = () => {
    navigation.navigate("avistamento");
  }

  const [menuAberto, setMenuAberto] = useState(false);
  const [mapaDeCalorAtivo, setMapaDeCalorAtivo] = useState(false);
  
  const [mostrarHospitais, setMostrarHospitais] = useState(false);
  const hospitaisRef = useRef(false); 

  const [animais, setAnimais] = useState([]);
  const [avistamentos, setAvistamentos] = useState([]);
  const [loading, setLoading] = useState(true);
  const [erro, setErro] = useState(null);

  const [animaisSelecionados, setAnimaisSelecionados] = useState({});

  const carregar = useCallback(async () => {
    try {
      setLoading(true);
      setErro(null);
      const [listaAnimais, listaAvistamentos] = await Promise.all([
        listarAnimais(),
        listarAvistamentos(),
      ]);
      setAnimais(listaAnimais);
      setAvistamentos(listaAvistamentos);
      setAnimaisSelecionados(prev => {
        const novo = {};
        listaAnimais.forEach(a => { novo[a.id] = prev[a.id] ?? !hospitaisRef.current; });
        return novo;
      });
    } catch (err) {
      console.error("Erro ao carregar o mapa:", err);
      setErro("Não foi possível carregar os avistamentos. Toque para tentar de novo.");
    } finally {
      setLoading(false);
    }
  }, []);

  useFocusEffect(useCallback(() => { carregar(); }, [carregar]));

  const definirMostrarHospitais = (valor) => {
    hospitaisRef.current = valor;
    setMostrarHospitais(valor);
  };

  const toggleHospitais = (valor) => {
    definirMostrarHospitais(valor);
    if (valor) {
      setAnimaisSelecionados(prev => {
        const resetados = {};
        Object.keys(prev).forEach(id => { resetados[id] = false; });
        return resetados;
      });
      setMapaDeCalorAtivo(false);
    }
  };

  const toggleAnimal = (id) => {
    if (!animaisSelecionados[id]) definirMostrarHospitais(false);
    setAnimaisSelecionados(prev => ({ ...prev, [id]: !prev[id] }));
  };

  const nomeDoAnimal = {};
  animais.forEach(a => { nomeDoAnimal[a.id] = a.nome_popular; });

  const dadosFiltrados = avistamentos.filter(a => animaisSelecionados[a.animal_id]);

  const pontosHeatmap = dadosFiltrados.map(a => ({
    latitude: a.latitude,
    longitude: a.longitude,
    weight: 1
  }));

  const statusMapa = mostrarHospitais
    ? null
    : loading
    ? 'Carregando avistamentos...'
    : erro
      ? erro
      : avistamentos.length === 0
        ? 'Nenhum avistamento registrado ainda.'
        : null;

  const renderCheckboxItem = (animal) => (
    <TouchableOpacity
      key={animal.id}
      style={styles.checkboxItemContainer}
      onPress={() => toggleAnimal(animal.id)}
      activeOpacity={0.7}
    >
      <Text style={styles.checkboxItemText}>{animal.nome_popular}</Text>

      <View style={styles.checkboxBox}>
        {animaisSelecionados[animal.id] && (
          <Feather name="check" size={16} color="#FFFFFF" />
        )}
      </View>
    </TouchableOpacity>
  );

  return (
    <SafeAreaView style={styles.container}>
      
      <View style={styles.header}>
        <View style={styles.headerSpacer} />

        <View style={styles.headerTitleBox}>
          <Text style={styles.headerTitleText}>MAPA RURAL</Text>
        </View>

        <TouchableOpacity 
          style={[styles.menuButton, menuAberto && styles.menuButtonAtivo]} 
          onPress={() => setMenuAberto(!menuAberto)}
        >
          <Feather name="menu" size={24} color="#FFFFFF" />
        </TouchableOpacity>
      </View>

      {/* AREA PRINCIPAL DO MAPA */}
      <View style={styles.mapContainer}>
        <MapView
          provider={PROVIDER_GOOGLE}
          style={styles.map}
          initialRegion={{
            latitude: -8.0142, 
            longitude: -34.9511, 
            latitudeDelta: 0.015,
            longitudeDelta: 0.012, 
          }}
        >

          {!mapaDeCalorAtivo && !mostrarHospitais && dadosFiltrados.map(a => (
            <Marker
              key={a.id}
              coordinate={{ latitude: a.latitude, longitude: a.longitude }}
              title={nomeDoAnimal[a.animal_id]}
            />
          ))}

          {mapaDeCalorAtivo && !mostrarHospitais && pontosHeatmap.length > 0 && (
            <Heatmap
              points={pontosHeatmap}
              radius={40} 
              opacity={0.7} 
              gradient={{
                colors: ["#00FF00", "#FF0000"],
                startPoints: [0.2, 1],
                colorMapSize: 256
              }}
            />
          )}

          {/* hospitais: posições fixas de HOSPITAIS_MOCK */}
          {mostrarHospitais && HOSPITAIS_MOCK.map(hospital => (
            <Marker
              key={hospital.id}
              coordinate={{ latitude: hospital.latitude, longitude: hospital.longitude }}
              title={hospital.nome}
              pinColor="blue" 
            />
          ))}

        </MapView>

        {statusMapa && (
          <TouchableOpacity
            disabled={!erro}
            onPress={carregar}
            activeOpacity={0.8}
            style={{
              position: 'absolute', bottom: 12, alignSelf: 'center',
              flexDirection: 'row', alignItems: 'center',
              backgroundColor: '#FFF', borderWidth: 1, borderColor: '#000',
              borderRadius: 12, paddingHorizontal: 12, paddingVertical: 8,
              maxWidth: '90%',
            }}
          >
            {loading && <ActivityIndicator size="small" color="#228B22" style={{ marginRight: 8 }} />}
            <Text style={{ color: '#000' }}>{statusMapa}</Text>
          </TouchableOpacity>
        )}

        {/* MENU DE SOBREPOSIÇÃO (Filtros) */}
        {menuAberto && (
          <View style={styles.overlayMenu}>
            
            <View style={styles.overlayHeader}>
              <Text style={styles.overlayHeaderText}>MAPA DE CALOR (On/Off)</Text>
              
              <Switch
                trackColor={{ false: "#767577", true: "#81b0ff" }}
                thumbColor={mapaDeCalorAtivo ? "#FFF" : "#f4f3f4"}
                ios_backgroundColor="#3e3e3e"
                onValueChange={(val) => {
                  setMapaDeCalorAtivo(val);
                  if (val) definirMostrarHospitais(false);
                }}
                value={mapaDeCalorAtivo}
              />
            </View>

            <View style={styles.overlayHeader}>
              <Text style={styles.overlayHeaderText}>UPA</Text>
              <Switch
                trackColor={{ false: "#767577", true: "#81b0ff" }}
                thumbColor={mostrarHospitais ? "#FFF" : "#f4f3f4"}
                ios_backgroundColor="#3e3e3e"
                onValueChange={toggleHospitais}
                value={mostrarHospitais}
              />
            </View>

            <View style={styles.overlayBody}>
              <View style={styles.filterTitleContainer}>
                <Text style={styles.filterTitleText}>SELECIONE O ANIMAL</Text>
              </View>
              
              <ScrollView style={{ maxHeight: 240 }}>
                {animais.map(animal => (
                  <View key={animal.id}>
                    {renderCheckboxItem(animal)}
                    <View style={styles.separator} />
                  </View>
                ))}
              </ScrollView>
            </View>
          </View>
        )}
      </View>

      <TouchableOpacity style={styles.avistamento} onPress={avistamento}>
        <Text style={styles.avistamentoText}>Adicionar Avistamento</Text>
      </TouchableOpacity>
    </SafeAreaView>
  );
}