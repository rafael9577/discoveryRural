import React, { useState, useEffect } from "react";
import { View, Text, TouchableOpacity, TextInput, Modal, ActivityIndicator, Alert } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";
import { useNavigation } from "@react-navigation/native";
import MapView from "react-native-maps";
import * as Location from "expo-location";
import { Ionicons } from "@expo/vector-icons";

import style from "./styled";
// placeholder da API:
// import { registrarOcorrencia } from "../../../../services/API";

export default () => {
  const navigation = useNavigation();

  const [animalId, setAnimalId] = useState("");
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
    setCoordenadas({
      latitude: regiaoMapa.latitude,
      longitude: regiaoMapa.longitude,
    });
    setMapaVisivel(false);
  };

  const salvarOcorrencia = async () => {
    if (!animalId || !coordenadas) {
      Alert.alert("Atenção", "Preencha a identificação e marque o local.");
      return;
    }

    try {
      setLoading(true);
      //id, descrição e as coordenadas
      const payload = { animalId, descricao, coordenadas };
      console.log("Dados prontos para envio:", payload);
      // registrarOcorrencia(payload);
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
      
      {/* CABEÇALHO */}
      <View style={style.header}>
        <View style={style.headerSpacer} />
        
        <View style={style.headerTitleBox}>
          <Text style={style.headerTitleText}>NOVA OCORRÊNCIA</Text>
        </View>

        <TouchableOpacity onPress={() => navigation.goBack()} style={style.backButton}>
          <Ionicons name="arrow-back" size={24} color="#FFFFFF" />
        </TouchableOpacity>
      </View>

      <View style={style.formContainer}>
        <TextInput
          style={style.input}
          placeholder="ID do Animal ou Brinco..."
          placeholderTextColor="#E0E0E0"
          value={animalId}
          onChangeText={setAnimalId}
        />

        <TextInput
          style={[style.input, style.textArea]}
          placeholder="Descrição detalhada..."
          placeholderTextColor="#E0E0E0"
          value={descricao}
          onChangeText={setDescricao}
          multiline
          numberOfLines={4}
        />

        {/* bloco de localização */}
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

      {/* Botão de submissão*/}
      <TouchableOpacity style={style.submitButton} onPress={salvarOcorrencia}>
        <Text style={style.submitButtonText}>Gravar Ocorrência</Text>
      </TouchableOpacity>

      {/* Loading overlay geral */}
      {loading && (
        <View style={style.centerContainer}>
          <ActivityIndicator size="large" color="#FFFFFF" />
          <Text style={style.loadingText}>Salvando registro...</Text>
        </View>
      )}

      {/* tela inteira para o Mapa */}
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