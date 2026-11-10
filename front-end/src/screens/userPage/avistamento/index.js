import React, { useState, useEffect } from "react";
import { View, Text, TouchableOpacity, TextInput, Modal, ActivityIndicator, Alert } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useNavigation } from "@react-navigation/native";
import MapView from "react-native-maps";
import * as Location from "expo-location";

import style from "./styled";
// placeholder da API:
// import { registrarOcorrencia } from "../../../../services/API";

export default () => {
  const navigation = useNavigation();

  // para o formulário
  const [animalId, setAnimalId] = useState("");
  const [tipoOcorrencia, setTipoOcorrencia] = useState("");
  const [descricao, setDescricao] = useState("");
  
  // localização e mapa 
  const [coordenadas, setCoordenadas] = useState(null);
  const [mapaVisivel, setMapaVisivel] = useState(false);
  const [regiaoMapa, setRegiaoMapa] = useState(null);

  // estados de UI
  const [loading, setLoading] = useState(false);
  const [obtendoLocalizacao, setObtendoLocalizacao] = useState(false);

  useEffect(() => {
    solicitarPermissaoLocalizacao();
  }, []);

  const solicitarPermissaoLocalizacao = async () => {
    setObtendoLocalizacao(true);
    let { status } = await Location.requestForegroundPermissionsAsync();
    if (status !== 'granted') {
      Alert.alert("Permissão negada", "Precisamos de acesso à localização para marcar a ocorrência.");
      setObtendoLocalizacao(false);
      return;
    }

    let location = await Location.getCurrentPositionAsync({});
    const coordAtual = {
      latitude: location.coords.latitude,
      longitude: location.coords.longitude,
      latitudeDelta: 0.005,
      longitudeDelta: 0.005,
    };
    setRegiaoMapa(coordAtual);
    setObtendoLocalizacao(false);
  };

  const abrirMapa = () => {
    if (!regiaoMapa) {
      solicitarPermissaoLocalizacao();
    }
    setMapaVisivel(true);
  };

  const confirmarLocalizacao = () => {
    // guarda a coordenada exata do centro do mapa (onde o pino virtual está apontando)
    setCoordenadas({
      latitude: regiaoMapa.latitude,
      longitude: regiaoMapa.longitude,
    });
    setMapaVisivel(false);
  };

  const salvarOcorrencia = async () => {
    if (!animalId || !tipoOcorrencia || !coordenadas) {
      Alert.alert("Atenção", "Preencha a identificação, o tipo de ocorrência e marque o local.");
      return;
    }

    try {
      setLoading(true);
      const payload = { animalId, tipoOcorrencia, descricao, coordenadas };
      console.log("Dados prontos para envio:", payload);
      // await registrarOcorrencia(payload);
      Alert.alert("Sucesso", "Ocorrência registrada com sucesso!");
      navigation.goBack();
    } catch (err) {
      console.error("Erro ao salvar:", err);
      Alert.alert("Erro", "Não foi possível registrar a ocorrência.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <SafeAreaView style={style.SafeAreaView}>
      <View style={style.formContainer}>
        <Text style={style.title}>Nova Ocorrência</Text>

        <TextInput
          style={style.input}
          placeholder="ID do Animal ou Brinco..."
          placeholderTextColor="#ccc"
          value={animalId}
          onChangeText={setAnimalId}
        />

        <TextInput
          style={style.input}
          placeholder="Tipo (Ex: Doença, Fuga, Vacinado)..."
          placeholderTextColor="#ccc"
          value={tipoOcorrencia}
          onChangeText={setTipoOcorrencia}
        />

        <TextInput
          style={[style.input, style.textArea]}
          placeholder="Descrição detalhada..."
          placeholderTextColor="#ccc"
          value={descricao}
          onChangeText={setDescricao}
          multiline
          numberOfLines={4}
        />

        {/* Bloco de Localização */}
        <View style={style.locationCard}>
          <Text style={style.locationTitle}>Localização da Ocorrência</Text>
          {coordenadas ? (
            <Text style={style.locationCoords}>
              Lat: {coordenadas.latitude.toFixed(5)} {"\n"}Lon: {coordenadas.longitude.toFixed(5)}
            </Text>
          ) : (
            <Text style={style.locationCoords}>Nenhum local marcado</Text>
          )}
          
          <TouchableOpacity 
            style={style.mapButton} 
            onPress={abrirMapa}
            disabled={obtendoLocalizacao}
          >
            <Text style={style.mapButtonText}>
              {obtendoLocalizacao ? "Obtendo GPS..." : (coordenadas ? "Ajustar no Mapa" : "Marcar Localização")}
            </Text>
          </TouchableOpacity>
        </View>
      </View>

      {/* Botão Fixo de Submissão no Rodapé */}
      <TouchableOpacity style={style.submitButton} onPress={salvarOcorrencia}>
        <Text style={style.submitButtonText}>Gravar Ocorrência</Text>
      </TouchableOpacity>

      {/* Loading Overlay Geral */}
      {loading && (
        <View style={style.centerContainer}>
          <ActivityIndicator size="large" color="#FFFFFF" />
          <Text style={style.loadingText}>Salvando registro...</Text>
        </View>
      )}

      {/* Modal de Ecrã Inteiro para o Mapa */}
      <Modal visible={mapaVisivel} animationType="slide" transparent={false}>
        <View style={style.mapContainer}>
          {regiaoMapa && (
            <MapView
              style={style.map}
              initialRegion={regiaoMapa}
              onRegionChangeComplete={(region) => setRegiaoMapa(region)}
              showsUserLocation={true}
            />
          )}
          
          {/* Pino fixo no centro do ecrã simulando o marcador */}
          <View style={style.fixedMarker}>
            <View style={style.markerDot} />
            <View style={style.markerStick} />
          </View>

          <View style={style.mapFooter}>
            <TouchableOpacity style={style.cancelMapButton} onPress={() => setMapaVisivel(false)}>
              <Text style={style.cancelMapButtonText}>Cancelar</Text>
            </TouchableOpacity>
            
            <TouchableOpacity style={style.confirmMapButton} onPress={confirmarLocalizacao}>
              <Text style={style.confirmMapButtonText}>Confirmar Local</Text>
            </TouchableOpacity>
          </View>
        </View>
      </Modal>
    </SafeAreaView>
  );
};